package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mailMsg "github.com/TheDM-95/mail-hub/pkg/mail/msg"

	"github.com/TheDM-95/mail-hub/consumer/handler"
	"github.com/TheDM-95/mail-hub/pkg/mail"
	"github.com/TheDM-95/mail-hub/pkg/mail/services"
	"github.com/TheDM-95/mail-hub/pkg/mail/services/mailgun"
	"github.com/TheDM-95/mail-hub/pkg/mail/services/sendgrid"
	"github.com/TheDM-95/mail-hub/util/constant"
)

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Send mail consumer",
	Run:   runServeConsumerCommand,
}

func init() {
	serveCmd.AddCommand(consumerCmd)

	consumerCmd.Flags().StringP("kafkaBrokers", "", "kafka-brokers", "Kafka broker")
	consumerCmd.Flags().StringP("mailService", "", "sendgrid", "mail service")
	consumerCmd.Flags().StringP("sendgridApiKey", "", "sendgrid-api-key", "api key for Sendgrid")
	consumerCmd.Flags().StringP("sendgridDefaultSenderName", "", "sendgrid-default-sender-name", "default sender name for Sendgrid")
	consumerCmd.Flags().StringP("sendgridDefaultSenderMail", "", "sendgrid-default-sender-mail", "default sender mail for Sendgrid")
	consumerCmd.Flags().StringP("mailgunDomain", "", "mailgun-domain", "send domain for Mailgun")
	consumerCmd.Flags().StringP("mailgunApiKey", "", "mailgun-api-key", "api key for Mailgun")
	consumerCmd.Flags().StringP("mailgunDefaultSenderName", "", "mailgun-default-sender-name", "default sender name for Mailgun")
	consumerCmd.Flags().StringP("mailgunDefaultSenderMail", "", "mailgun-default-sender-mail", "default sender mail for Mailgun")

	_ = viper.BindPFlag("kafkaBrokers", consumerCmd.Flags().Lookup("kafkaBrokers"))
	_ = viper.BindPFlag("mailService", consumerCmd.Flags().Lookup("mailService"))
	_ = viper.BindPFlag("sendgridApiKey", consumerCmd.Flags().Lookup("sendgridApiKey"))
	_ = viper.BindPFlag("sendgridDefaultSenderName", consumerCmd.Flags().Lookup("sendgridDefaultSenderName"))
	_ = viper.BindPFlag("sendgridDefaultSenderMail", consumerCmd.Flags().Lookup("sendgridDefaultSenderMail"))
	_ = viper.BindPFlag("mailgunDomain", consumerCmd.Flags().Lookup("mailgunDomain"))
	_ = viper.BindPFlag("mailgunApiKey", consumerCmd.Flags().Lookup("mailgunApiKey"))
	_ = viper.BindPFlag("mailgunDefaultSenderName", consumerCmd.Flags().Lookup("mailgunDefaultSenderName"))
	_ = viper.BindPFlag("mailgunDefaultSenderMail", consumerCmd.Flags().Lookup("mailgunDefaultSenderMail"))
}

func runServeConsumerCommand(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	router, err := initConsumer()
	if err != nil {
		panic(err)
	}

	// Run Message Router
	go func() {
		err := router.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	<-c
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
}

func initConsumer() (*message.Router, error) {
	brokers := strings.Split(viper.GetString("kafkaBrokers"), ",")
	watermillLogger := watermill.NewStdLogger(false, false)

	// Create Router
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		return nil, err
	}

	// Add Middleware and Plugin
	middlewareRetry := middleware.Retry{MaxRetries: 2, InitialInterval: time.Second * 10, Logger: watermillLogger}
	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(middleware.CorrelationID, middlewareRetry.Middleware, middleware.Recoverer)

	// Initialize mailer
	var sendgridSvc, mailgunSvc services.MailService
	sendSvcConf := viper.GetString("mailService")
	switch sendSvcConf {
	case "sendgrid":
		sendgridSvc = initSendgridSvc()
	case "mailgun":
		mailgunSvc = initMailgunSvc()
	case "mixed":
		sendgridSvc = initSendgridSvc()
		mailgunSvc = initMailgunSvc()
	default:
		return nil, errors.New("unsupported send service")
	}

	if sendgridSvc != nil {
		err = initSendWorker(brokers, router, sendgridSvc, 30*time.Second)
	}

	if mailgunSvc != nil {
		err = initSendWorker(brokers, router, mailgunSvc, 30*time.Second)
	}

	return router, err
}

func initSendgridSvc() *sendgrid.MailService {
	dfSender := &mailMsg.EmailAddress{
		Name:  viper.GetString("sendgridDefaultSenderName"),
		Email: viper.GetString("sendgridDefaultSenderMail"),
	}
	sgApiKey := viper.GetString("sendgridApiKey")

	return sendgrid.NewMailService(
		sendgrid.WithApiKey(sgApiKey),
		sendgrid.WithDefaultSender(dfSender),
	)
}

func initMailgunSvc() *mailgun.MailService {
	dfSender := &mailMsg.EmailAddress{
		Name:  viper.GetString("mailgunDefaultSenderName"),
		Email: viper.GetString("mailgunDefaultSenderMail"),
	}
	mgDomain := viper.GetString("mailgunDomain")
	mgApiKey := viper.GetString("mailgunApiKey")

	return mailgun.NewMailService(
		mailgun.WithDefaultApiClient(mgDomain, mgApiKey),
		mailgun.WithDefaultSender(dfSender),
	)
}

func initSendWorker(brokers []string, router *message.Router, sendSvc services.MailService, timeout time.Duration) error {
	mailer := mail.NewMailer(
		mail.WithMailService(sendSvc),
		mail.WithTimeout(timeout),
	)
	sendMailHandler := handler.NewSendMailHandler(handler.WithMailer(mailer))

	handlerName := fmt.Sprintf("SendMailHandler%s", sendSvc.GetName())
	sendTopic := constant.QueueTopicSendMail

	// Declare Sarama Subscriber Config
	subscriberCfg := kafka.DefaultSaramaSubscriberConfig()
	subscriberCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create Subscriber
	watermillLogger := watermill.NewStdLogger(false, false)
	subscriberConfig := kafka.SubscriberConfig{
		Brokers:               brokers,
		Unmarshaler:           kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: subscriberCfg,
		ConsumerGroup:         constant.MailHubConsumerGroup,
	}
	subscriber, err := kafka.NewSubscriber(subscriberConfig, watermillLogger)
	if err != nil {
		return err
	}

	router.AddNoPublisherHandler(handlerName, sendTopic, subscriber, sendMailHandler.Handle)
	return nil
}

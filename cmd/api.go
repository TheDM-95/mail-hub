package cmd

import (
	"context"
	"fmt"
	"github.com/TheDM-95/mail-hub/util/publisher"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/TheDM-95/mail-hub/api/router"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "api",
	Short: "Serve the Shopify public app",
	Run:   runServeAppCmd,
}

func init() {
	serveCmd.AddCommand(appCmd)

	appCmd.Flags().StringP("address", "a", ":9092", "Serving address. Default is 9092")
	appCmd.Flags().StringP("kafkaBrokers", "q", "127.0.0.1:9092", "Kafka broker for messaging")

	_ = viper.BindPFlag("address", appCmd.Flags().Lookup("address"))
	_ = viper.BindPFlag("kafkaBrokers", appCmd.Flags().Lookup("kafkaBrokers"))
}

func runServeAppCmd(cmd *cobra.Command, args []string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	brokers := strings.Split(viper.GetString("kafkaBrokers"), ",")
	err := publisher.InitKafkaPublisher(brokers)
	if err != nil {
		panic(err)
	}

	address := viper.GetString("address")
	r := mux.NewRouter()
	router.ResolveRoute(r)

	s := &http.Server{
		Handler:      r,
		Addr:         address,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		fmt.Println("server started at " + address)
		err := s.ListenAndServe()
		if nil != err {
			log.Fatal(err)
		}
	}()

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	defer cancel()
}

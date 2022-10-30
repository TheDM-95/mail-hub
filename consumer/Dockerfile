FROM centos:7

EXPOSE 9090

ADD ./bin/send-mail-consumer /usr/bin/send-mail-consumer

CMD ["send-mail-consumer", "serve"]

FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/server
RUN mkdir /app/server/apps
RUN mkdir /app/server/apps/mail
WORKDIR /app/server/apps/mail

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD mail /app/server/apps/mail/mail

ENTRYPOINT [ "/app/server/apps/mail/mail" ]

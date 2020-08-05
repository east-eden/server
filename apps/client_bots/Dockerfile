FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/apps
RUN mkdir /app/apps/client_bots
WORKDIR /app/apps/client_bots

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD client_bots /app/apps/client_bots/client_bots

ENTRYPOINT [ "/app/apps/client_bots/client_bots" ]

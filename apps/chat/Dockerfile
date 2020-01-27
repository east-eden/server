FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/apps
RUN mkdir /app/apps/chat
WORKDIR /app/apps/chat

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD chat /app/apps/chat/chat

ENTRYPOINT [ "/app/apps/chat/chat" ]

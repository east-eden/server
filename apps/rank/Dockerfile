FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/server
RUN mkdir /app/server/apps
RUN mkdir /app/server/apps/rank
WORKDIR /app/server/apps/rank

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD rank /app/server/apps/rank/rank

ENTRYPOINT [ "/app/server/apps/rank/rank" ]

FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/server
RUN mkdir /app/server/apps
RUN mkdir /app/server/apps/game
WORKDIR /app/server/apps/game

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD game /app/server/apps/game/game

ENTRYPOINT [ "/app/server/apps/game/game" ]

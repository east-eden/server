FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/apps
RUN mkdir /app/apps/game
WORKDIR /app/apps/game

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD game /app/apps/game/game

ENTRYPOINT [ "/app/apps/game/game" ]

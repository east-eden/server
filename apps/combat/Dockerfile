FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/apps
RUN mkdir /app/apps/combat
WORKDIR /app/apps/combat

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD combat /app/apps/combat/combat

ENTRYPOINT [ "/app/apps/combat/combat" ]

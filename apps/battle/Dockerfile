FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/apps
RUN mkdir /app/apps/battle
WORKDIR /app/apps/battle

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD battle /app/apps/battle/battle

ENTRYPOINT [ "/app/apps/battle/battle" ]

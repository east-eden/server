FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/server
RUN mkdir /app/server/apps
RUN mkdir /app/server/apps/gate
WORKDIR /app/server/apps/gate

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD gate /app/server/apps/gate/gate

ENTRYPOINT [ "/app/server/apps/gate/gate" ]

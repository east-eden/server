FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/apps
RUN mkdir /app/apps/gate
WORKDIR /app/apps/gate

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD gate /app/apps/gate/gate

ENTRYPOINT [ "/app/apps/gate/gate" ]

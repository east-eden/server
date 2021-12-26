FROM alpine
RUN apk add tzdata

RUN mkdir /app
RUN mkdir /app/server
RUN mkdir /app/server/apps
RUN mkdir /app/server/apps/comment
WORKDIR /app/server/apps/comment

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ADD comment /app/server/apps/comment/comment

ENTRYPOINT [ "/app/server/apps/comment/comment" ]

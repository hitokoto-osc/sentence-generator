FROM alpine:3.18
ENV TZ Asia/Shanghai
RUN apk add alpine-conf tzdata && \
    /sbin/setup-timezone -z Asia/Shanghai && \
    apk del alpine-conf

ENV WORKDIR /app
VOLUME $WORKDIR/data
COPY data $WORKDIR/
COPY generator $WORKDIR/
ENTRYPOINT ["sh", "-c", "$WORKDIR/generator", "--", "start"]



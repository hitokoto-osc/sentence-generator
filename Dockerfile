FROM scratch
WORKDIR /src/app
USER hitokoto
VOLUME data
COPY data .
COPY generator .
ENTRYPOINT ["./generator", "start"]


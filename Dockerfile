FROM scratch
WORKDIR /src/app
RUN groupadd -r hitokoto && useradd -r -g hitokoto hitokoto
USER hitokoto
VOLUME data
COPY data .
COPY generator .
ENTRYPOINT ["./generator", "start"]


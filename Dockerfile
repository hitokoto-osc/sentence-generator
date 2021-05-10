FROM scratch
WORKDIR /
VOLUME data
COPY data /
COPY generator /
ENTRYPOINT ["/generator", "start"]


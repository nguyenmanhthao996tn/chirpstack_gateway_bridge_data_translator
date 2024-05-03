FROM golang:1.21.4-bookworm AS development

ENV PROJECT_PATH=/chirpstack-gw-protobuf-translator

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make

FROM debian:bookworm AS production

WORKDIR /root/
COPY --from=development /chirpstack-gw-protobuf-translator/build .
ENTRYPOINT ["./chirpstack-gw-protobuf-translator"]
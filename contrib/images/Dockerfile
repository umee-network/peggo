## image for docker Peggo release

FROM golang:1.19-bullseye AS builder
WORKDIR /src/peggo/
COPY . .
RUN make install &&  \
    wget https://raw.githubusercontent.com/CosmWasm/wasmvm/v1.1.1/internal/api/libwasmvm.x86_64.so

# prepare final image
FROM ubuntu:rolling
COPY --from=builder /go/bin/peggo /usr/local/bin/
COPY --from=builder /src/peggo/libwasmvm.x86_64.so /lib/
CMD ["peggo"]
STOPSIGNAL SIGTERM

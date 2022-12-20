## image for e2e tests
## docker build -t peggo-e2e -f ./contrib/images/e2e.dockerfile .

# Fetch base packages
FROM golang:1.19-bullseye AS builder

# Compile the peggo binary
WORKDIR /src/peggo/
COPY . .
RUN go mod download
RUN make install

# download umeed
WORKDIR /src/umee
RUN wget https://github.com/umee-network/umee/releases/download/v3.3.0-rc3/umeed-v3.3.0-rc3-linux-amd64 && \
  chmod +x umeed-v* && \
  cp umeed-v* umeed && \
  wget https://raw.githubusercontent.com/CosmWasm/wasmvm/v1.1.1/internal/api/libwasmvm.x86_64.so

# Prepare final image
# FROM gcr.io/distroless/cc:debug
FROM ubuntu:rolling
ARG IMG_TAG=latest
COPY --from=builder /go/bin/peggo /usr/local/bin/
COPY --from=builder /src/umee/umeed /usr/local/bin/
COPY --from=builder /src/umee/libwasmvm.x86_64.so /lib/
EXPOSE 26656 26657 1317 9090

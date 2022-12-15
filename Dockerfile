## image for docker Peggo release

ARG IMG_TAG=latest

# Fetch base packages
FROM golang:1.19-bullseye AS builder

# Compile the peggo binary
WORKDIR /src/peggo/
COPY . .
RUN go mod download
RUN make install

# Build umeed
WORKDIR /src/umee
# wget https://github.com/umee-network/umee/releases/download/v3.3.0-rc1/umeed-v3.3.0-rc1-linux-amd64
RUN https://github.com/umee-network/umee/releases/download/v3.0.3/umeed-v3.0.3-darwin-amd64.tar.gz && \
  tar -xvf umeed-v3.0.3-darwin-amd64.tar.gz  && \
  cd umeed-v3.0.3-darwin-amd64 && \
  chmod +x umeed* && \
  cp umeed* /usr/local/bin/umeed

# Add to a distroless container
FROM gcr.io/distroless/cc:$IMG_TAG
ARG IMG_TAG
COPY --from=builder /go/bin/peggo /usr/local/bin/
COPY --from=builder /usr/local/bin/umeed /usr/local/bin/
EXPOSE 26656 26657 1317 9090

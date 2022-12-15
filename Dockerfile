## image for docker Peggo release

ARG IMG_TAG=latest

# Fetch base packages
FROM golang:1.19-alpine AS base-builder
RUN apk add --no-cache make git libc-dev gcc linux-headers

# Compile the peggo binary
WORKDIR /src/peggo/
COPY . .
RUN go mod download
RUN make install

# Build umeed
WORKDIR /src/umee
ARG UMEE_VERSION=v3.3.0-rc1
RUN git clone https://github.com/umee-network/umee.git
RUN cd umee && git checkout ${UMEE_VERSION} && make build && cp ./build/umeed /usr/local/bin/

# Add to a distroless container
FROM gcr.io/distroless/cc:$IMG_TAG
ARG IMG_TAG
COPY --from=base-builder /go/bin/peggo /usr/local/bin/
COPY --from=base-builder /usr/local/bin/umeed /usr/local/bin/
EXPOSE 26656 26657 1317 9090

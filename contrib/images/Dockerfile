## image for e2e tests

# Fetch base packages
FROM golang:1.19-alpine AS builder
ENV PACKAGES make git libc-dev gcc linux-headers
RUN apk add --no-cache $PACKAGES
WORKDIR /src/app/
COPY . .
RUN make install 

FROM alpine:3.14
RUN apk add bash curl jq
COPY --from=builder /go/bin/peggo /usr/local/bin/
CMD ["peggo"]
STOPSIGNAL SIGTERM

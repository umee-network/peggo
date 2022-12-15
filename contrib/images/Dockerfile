## image for e2e tests

# Fetch base packages
FROM golang:1.19-bullseye AS builder
WORKDIR /src/peggo/
COPY . .
RUN make install 

FROM ubuntu:rolling
COPY --from=builder /go/bin/peggo /usr/local/bin/
CMD ["peggo"]
STOPSIGNAL SIGTERM

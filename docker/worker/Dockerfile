FROM golang:stretch as go-builder

RUN mkdir -p /go/src/github.com/nirnanaaa
COPY . /go/src/github.com/nirnanaaa/cloudive-mailer
WORKDIR /go/src/github.com/nirnanaaa/cloudive-mailer

RUN ./build.py

FROM alpine:3.7

RUN apk add --no-cache ca-certificates

# This is required for glibc compatibility
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY docker/master/optimistic-defaults.toml /config.toml
EXPOSE 9009
COPY --from=go-builder /go/src/github.com/nirnanaaa/cloudive-mailer/build/cloudive-mailer /usr/local/bin/mailer
ENTRYPOINT ["/usr/local/bin/mailer"]
CMD ["slave"]
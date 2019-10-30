FROM alpine:3.10

ADD placeholder .

ENTRYPOINT [ "./placeholder" ]

EXPOSE 80
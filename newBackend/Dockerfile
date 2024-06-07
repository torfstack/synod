FROM alpine:3.20.0

RUN mkdir /opt/kayvault

WORKDIR /opt/kayvault

COPY /bin/kayvault-backend kayvault-backend
COPY /sql/migrations sql/migrations

CMD ["./kayvault-backend"]
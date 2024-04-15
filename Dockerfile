FROM alpine

COPY .env .
ADD ./bin/app /app

RUN chmod +x ./app
ENTRYPOINT ["./app"]
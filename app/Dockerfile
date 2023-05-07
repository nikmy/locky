FROM alpine:latest

COPY locky_bot .
COPY db/sql sql

RUN chmod +x locky_bot

EXPOSE 80/tcp

CMD ["./locky_bot"]

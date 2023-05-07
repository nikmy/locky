FROM alpine:latest

COPY bot_linux /var/www/locky/locky_bot

WORKDIR /var/www/locky

RUN chmod +x locky_bot

EXPOSE 80/tcp

CMD ["./locky_bot"]

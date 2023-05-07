FROM alpine:latest

COPY locky_bot /var/www/locky

WORKDIR /var/www/locky

RUN chmod +x locky_bot

EXPOSE 80/tcp

CMD ["./locky_bot"]

FROM alpine:3.14
COPY . /app
CMD ["/bin/sh", "-c", "/app/thomas"]
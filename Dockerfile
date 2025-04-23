FROM ubuntu:24.04
COPY . /app
CMD ["/bin/sh", "-c", "/app/thomas"]
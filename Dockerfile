FROM alpine:3.14
COPY . /app
RUN ls -la /app # debug
WORKDIR /app
CMD /app/thomas
FROM golang:1.23.4

RUN apt-get update

ARG APP_NAME=app

RUN mkdir /$APP_NAME
COPY . /$APP_NAME
WORKDIR /$APP_NAME

RUN go mod download
RUN go build -o main .

EXPOSE 4730
CMD ["./main"]
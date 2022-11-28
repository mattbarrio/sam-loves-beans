FROM golang:1.19-alpine

RUN mkdir /app

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./beans

ARG API_KEY

COPY static ./static
COPY templates ./templates

EXPOSE 8080
CMD [ "/app/beans" ]

VOLUME /app/data
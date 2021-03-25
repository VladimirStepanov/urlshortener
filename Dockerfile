FROM golang:1.15.1

WORKDIR /app

COPY ./go.mod ./go.sum ./
COPY . .


RUN go mod download
RUN make

CMD ["./urlshortener"]
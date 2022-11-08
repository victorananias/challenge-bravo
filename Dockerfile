FROM golang:1.19

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go build -v

CMD ["challenge-bravo"]
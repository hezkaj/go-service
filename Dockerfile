FROM golang:alpine 

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build main.go

CMD [". /main"]
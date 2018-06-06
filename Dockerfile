FROM golang

WORKDIR /go/src/github.com/paulthom12345/route
COPY . . 
RUN go get
RUN go build 


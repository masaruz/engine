FROM golang:1.10

WORKDIR /go/src/engine
COPY . .

RUN go get -d -v   # "go get -d -v ./..."
RUN go install -v    # "go install -v ./..."

CMD ["engine"]
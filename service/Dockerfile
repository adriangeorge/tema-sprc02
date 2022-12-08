FROM golang:latest

ADD . /usr/local/go/src/service
WORKDIR /usr/local/go/src/service
# COPY go.mod ./
# COPY go.sum ./
# COPY api ./
# COPY database ./
# COPY models ./

RUN go mod download

COPY *.go ./

RUN go build -o ./service

EXPOSE 6000

CMD [ "./service" ]
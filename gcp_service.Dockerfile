FROM fedora
RUN dnf install git -y

FROM golang:1.17-alpine

RUN apk add --no-cache git

ENV GO111MODULE=on
ENV GOPATH /go
WORKDIR /gcp_service/go_app

COPY go.mod .

RUN go mod download

COPY . .

RUN  go build -o main .

EXPOSE 4001

CMD ["go", "run", "main.go"]




# ENTRYPOINT [ "main" ]

#CMD ["go", "run", "main.go"]

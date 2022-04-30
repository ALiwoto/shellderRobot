FROM golang:1.18-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./


COPY shellderRobot/ ./shellderRobot/
COPY *.go ./
COPY *.ini ./
COPY *.sh ./
COPY *.bat ./

RUN go mod tidy -e

RUN go build -o /botBinary

CMD [ "/botBinary" ]



FROM golang:1.20-bullseye as builder

WORKDIR /app

RUN apt update && apt install tzdata -y

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o sinarlog-scheduler

EXPOSE 8123
ENV TZ="Asia/Jakarta"

ENTRYPOINT ["/app/sinarlog-scheduler"]
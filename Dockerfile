FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o /envrep main.go
ENTRYPOINT ["/envrep"]

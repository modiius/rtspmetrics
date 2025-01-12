FROM golang:latest

COPY . .
RUN go build -o rtspmetrics .

CMD ["./rtspmetrics"]

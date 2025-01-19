FROM golang:latest

# FFmpeg libraries
RUN apt-get update && apt-get install -y \
    libavformat-dev \
    libswscale-dev \
    gcc \
    pkg-config \
    && apt-get clean && rm -rf /var/lib/apt/lists

WORKDIR /app

# cached if go.mod and go.sum are unchanged
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o rtspmetrics .

CMD ["./rtspmetrics"]

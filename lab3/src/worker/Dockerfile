FROM golang:latest

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o worker
ENV WORDLIST_PATH="/app/rockyou-50.txt"
ENTRYPOINT "/app/worker"

FROM golang:latest

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o server
ENTRYPOINT "/app/server"

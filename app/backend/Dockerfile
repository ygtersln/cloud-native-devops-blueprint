# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY *.go ./
# RUN go mod download
# For this demo, since we don't have go.sum, we run go mod tidy
RUN go env -w GOPROXY=direct && go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o blueprint-api .

# Final minimal stage
FROM gcr.io/distroless/static-debian12:latest

WORKDIR /
COPY --from=builder /app/blueprint-api .

EXPOSE 8080
USER nonroot:nonroot

CMD ["/blueprint-api"]

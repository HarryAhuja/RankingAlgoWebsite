FROM golang:1.22.4-alpine3.20 AS builder

# First stage -> copy all files from src to dst and build go binary
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./rankingalgo

# Second Stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app
# copy binary to this lightweight base image
COPY --from=builder /build/rankingalgo ./rankingalgo
CMD ["/app/rankingalgo"]
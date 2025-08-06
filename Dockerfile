FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Cross-compile for linux/amd64, as EKS nodes are likely amd64.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /octopipe .

FROM alpine:3.22
COPY --from=builder /octopipe /octopipe
EXPOSE 6652
ENTRYPOINT ["/octopipe"]
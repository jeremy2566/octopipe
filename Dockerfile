FROM golang:1.24 AS builder

# 为构建时变量添加参数
ARG VERSION=dev
ARG REVISION=unknown

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Cross-compile for linux/amd64, as EKS nodes are likely amd64.
# 使用 ldflags 注入版本和修订信息。
# 使用 -tags=purego 来禁用汇编，以实现稳健的交叉编译。
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=purego -ldflags="-X 'github.com/jeremy2566/octopipe/pkg/version.VERSION=${VERSION}' -X 'github.com/jeremy2566/octopipe/pkg/version.REVISION=${REVISION}'" -o /octopipe .

FROM alpine:3.22
COPY --from=builder /octopipe /octopipe
EXPOSE 6652
ENTRYPOINT ["/octopipe"]
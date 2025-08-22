# 第一阶段：在本地编译
FROM golang:1.24-alpine AS builder

# 安装 git（alpine 版本需要）
RUN apk add --no-cache git

# 为构建时变量添加参数
ARG VERSION=dev
ARG REVISION=unknown

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 使用 alpine 的原生架构编译，避免段错误
# 然后再为目标平台交叉编译
RUN CGO_ENABLED=0 go build -tags=purego -ldflags="-X 'github.com/jeremy2566/octopipe/pkg/version.VERSION=${VERSION}' -X 'github.com/jeremy2566/octopipe/pkg/version.REVISION=${REVISION}'" -o /octopipe .

FROM alpine:3.22
COPY --from=builder /octopipe /octopipe
EXPOSE 6652
ENTRYPOINT ["/octopipe"]
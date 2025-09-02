# 第一阶段：在本地编译
# 使用一个标准的 Debian-based Go 镜像来提高编译器的稳定性
# 这个镜像已经包含了 git，所以不需要再手动安装
FROM golang:1.24-bookworm AS builder

# 为构建时变量添加参数
ARG VERSION=dev
ARG REVISION=unknown

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 交叉编译为 linux 平台，并生成静态链接的二进制文件
# -a: 强制重新构建所有包
# -ldflags "-w -s": 减小二进制文件体积
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags=purego -ldflags="-w -s -X 'github.com/jeremy2566/octopipe/pkg/version.VERSION=${VERSION}' -X 'github.com/jeremy2566/octopipe/pkg/version.REVISION=${REVISION}'" -o /octopipe .

# 第二阶段：构建一个最小化的最终镜像
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /octopipe /octopipe
EXPOSE 6652
ENTRYPOINT ["/octopipe"]
# 第一阶段：用官方 Golang 镜像编译
FROM golang:1.24.3 AS builder

WORKDIR /app

# 只复制 go.mod 和 go.sum，先下载依赖，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 编译生成可执行文件 server
RUN go build -o server .

# 第二阶段：更小的运行镜像，拷贝编译好的二进制文件
FROM debian:bookworm-slim

# 拷贝二进制文件到根目录
COPY --from=builder /app/server /server

# 暴露端口
EXPOSE 8080

# 运行二进制
ENTRYPOINT ["/server"]


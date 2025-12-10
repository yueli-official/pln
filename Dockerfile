# ----------------------
# 前端构建阶段 (Vue)
# ----------------------
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend

# 接收构建参数（访问 GitHub Packages）
ARG GH_TOKEN=""

# 安装 pnpm
RUN npm install -g pnpm

# 复制依赖文件并配置认证
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN if [ -n "$GH_TOKEN" ]; then \
  echo "@yuelioi:registry=https://npm.pkg.github.com" > .npmrc && \
  echo "//npm.pkg.github.com/:_authToken=${GH_TOKEN}" >> .npmrc; \
  fi && \
  pnpm install --frozen-lockfile && pnpm store prune && \
  rm -f .npmrc

# 复制源码并构建
COPY frontend/ .
RUN pnpm build

# 清理开发文件，只保留 dist
RUN rm -rf src node_modules public package.json pnpm-lock.yaml vite.config.ts tsconfig*.json eslint.config.ts env.d.ts

# ----------------------
# 后端构建阶段 (Go)
# ----------------------
FROM golang:1.25.1-alpine AS backend-builder
WORKDIR /app

# 安装依赖工具
RUN apk add --no-cache git build-base

# 复制 go.mod / go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制后端源码
COPY . .

# 复制前端构建产物 dist
COPY --from=frontend-builder /frontend/dist ./frontend/dist

# 编译 Go 应用（main.go 在 cmd/main.go）
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd

# ----------------------
# 运行阶段
# ----------------------
FROM alpine:latest
WORKDIR /app

# 安装必要依赖
RUN apk add --no-cache libc6-compat

# 复制后端可执行文件
COPY --from=backend-builder /app/server ./

# 复制前端构建好的静态文件
COPY --from=backend-builder /app/frontend/dist ./frontend/dist

# 设置数据库环境变量
ENV DATABASE_URL=/app/index.db

# 暴露端口
EXPOSE 9000

# 启动服务
CMD ["./server"]
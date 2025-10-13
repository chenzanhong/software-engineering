FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install --registry https://registry.npmmirror.com

# 接收构建参数
ARG VUE_APP_API_BASE_URL
ENV VUE_APP_API_BASE_URL=$VUE_APP_API_BASE_URL

COPY . .
RUN chmod -R +x node_modules/.bin/

RUN npm run build  # 构建时会读取 .env 或 ARG 注入的变量

# 第二阶段：Nginx 托管
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

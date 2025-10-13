FROM node:16-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# 暴露开发端口（如 Vite 默认 3000）
EXPOSE 3000

# 启动开发服务器，监听 0.0.0.0
CMD ["npm", "run", "serve"]
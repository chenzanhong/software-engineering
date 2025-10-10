### Text2Image

Text2Image 是一个文本到图片的 Web 应用，实现了账号的注册与管理、文生图的前端抽象集成、个人信息管理以及文生图历史查看等功能。
（最新版本请查看czh_query分支）

## 前端启动

### 配置

1. **切换到前端项目目录**：
   ```bash
   cd frontend/text2image-vue
   ```

2. **创建 `.env.oss.local` 文件**：
   ```bash
   touch .env.oss.local
   vim .env.oss.local
   ```

3. **编辑 `.env.oss.local` 文件**：
   在 `.env.oss.local` 文件中，填入阿里云OSS的密钥。以下是一个示例配置：
   ```plaintext
   VUE_APP_OSS_REGION=oss-cn-shenzhen
   VUE_APP_OSS_ACCESS_KEY_ID=your_access_key_id
   VUE_APP_OSS_ACCESS_KEY_SECRET=your_access_key_secret
   VUE_APP_OSS_BUCKET=your_bucket_name
   ```

### 启动前端

1. **安装依赖**：
   ```bash
   npm install
   ```

2. **启动开发服务器**：
   ```bash
   npm run serve
   ```

## 后端启动

### 配置

1. **切换到后端项目目录**：
   ```bash
   cd backend/text-to-picture
   ```

2. **拉取依赖**：
   ```bash
   go mod tidy
   ```
   如果有拉取不完全的报错，可以使用 `go get + 依赖` 进行手动拉取。

3. **修改数据库配置**：
   - 进入 `backend/text-to-picture/config/DBconfig` 目录，修改 `DBconfig.yaml.example` 文件为自己的数据库配置,并重命名为 `DBconfig.yaml`。

4. **配置阿里云OSS**：
    在 `text-to-picture/config/oss` 目录下创建 `oss.env` 文件，写入阿里云OSS存储的配置信息。
   ```bash
   cd config/oss
   touch oss.env
   vim oss.env
   ```
   
   以下是一个示例配置：
     ```plaintext
     OSS_REGION=oss-cn-shenzhen
     OSS_ACCESS_KEY_ID=your_access_key_id
     OSS_ACCESS_KEY_SECRET=your_access_key_secret
     OSS_BUCKET=your_bucket_name
     ```


### 启动后端

1. **运行后端服务**：
   ```bash
   go run main.go
   ```

## 其他注意事项

- **暂无**

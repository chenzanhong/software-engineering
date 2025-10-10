# 文绘星河

文绘星河 是一个文本到图片的 Web 应用，实现了账号的注册与管理、文生图的前端抽象集成、个人信息管理以及文生图历史查看等功能。

![文绘星河](https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202412251121811.png)

## 项目结构

### 前端结构

1. **`frontend/text2image-vue/.env`**
   - 前端环境变量配置文件，用于头像的上传以及本地模型生成结果的存储，便于后续与后端对接。

2. **`frontend/text2image-vue/src/views`**
   - **`AboutView.vue`**: 项目介绍页
   - **`ExploreView.vue`**: 探索界面（含图片广场）
   - **`FavoritesView.vue`**: 收藏夹界面
   - **`GenerateView.vue`**: 文生图界面
   - **`HistoryView.vue`**: 文生图历史记录界面
   - **`InnerView.vue`**: 主导航界面（功能页为其子路由）
   - **`LogRegView.vue`**: 登录注册界面
   - **`SettingView.vue`**: 个人信息设置界面
   - **`UsageGuide.vue`**: 使用说明界面

   > 项目的所有页面视图定义，Vue-Router的路由定义基于此处的页面。

3. **`frontend/text2image-vue/src/components`**
   - 考虑到页面需要复用的组件不多，且定制性较强，所以没有特地进行组件拆分，此处无复用组件。

4. **`frontend/text2image-vue/src/assets`**
   - 存储静态资源，如项目介绍页的图片、按钮图标、首页的常驻图片等。

5. **`frontend/text2image-vue/src/router/index.js`**
   - 路由定义，路由跳转定义，路由守卫，页面懒加载等。

6. **`frontend/text2image-vue/src/api/index.js`**
   - API定义，封装了axios请求，自动在请求带上token，主要便于开发和前后端交互。

### 后端结构
-	API层（处理图片生成的API逻辑）
-	资源层（存放示例图片）
-	配置层（数据库和对象存储配置）
-	中间件层（JWT认证中间件）
-	模型层（数据模型及其数据库操作）
-	服务层（核心业务逻辑：用户认证、收藏、历史记录、图片生成等）
-	路由层（定义所有API路由）

   
## 前端启动

### 配置

1. **切换到前端项目目录**：
  
   ```bash
   cd frontend/text2image-vue
   ```
   
2. **创建 `.env` 文件**：
   ```bash
   touch .env
   vim .env
   ```

3. **编辑 `.env` 文件**：
   在 `.env` 文件中，填入阿里云OSS的密钥。以下是一个示例配置：
   ```plaintext
   VUE_APP_OSS_REGION=oss-cn-shenzhen.aliyuncs.com
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

3. **修改配置文件**：
   进入 `backend/text-to-picture/config/configs` 目录，修改 `config.yaml.example` 文件为自己的配置，并重命名为 `config.yaml`。

   以下是一个示例配置：
   ```yaml
   db:
     host: localhost # 数据库地址
     port: "5432"    # 数据库端口
     name: database  # 数据库名
     user: user      # 数据库用户名
     password: password # 数据库密码

   oss:
     OSS_REGION: region # oss区域
     OSS_ACCESS_KEY_ID: ... # oss key
     OSS_ACCESS_KEY_SECRET: ... # oss密钥
     OSS_BUCKET: bucket # oss bucket

   model:
     GEN_API_KEY: sk-6e79f5171c934d8fbbbdb0f4cd42d669 # api_key
     timeout: 30 # 轮询时间
   ```



## 其他注意事项

- **暂无**
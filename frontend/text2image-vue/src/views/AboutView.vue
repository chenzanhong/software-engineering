<template>
  <div class="project-intro">
    <!-- 头部 -->
    <header class="header">
      <div class="left-section">
        <div class="logo-container">
          <img src='../assets/button-icon/文生图.svg' alt="Logo" class="logo" @click="this.$router.push('log-reg')"/>
        </div>
        <h1 class="app-name" @click="this.$router.push('login')">{{ appName }}</h1>
      </div>
      <div class="right-section">
        <a href="/log-reg" class="link">登录</a>
        <a href="/log-reg" class="link">注册</a>
      </div>
    </header>

    <!-- 全屏轮播图 -->

      <el-carousel height="100vh">
      <el-carousel-item v-for="item in carouselItems" :key="item.id">
        <div class="carousel-content">
          <el-image
            style="width: 100%; height: auto"
            :src="item.image"
            :fit="scale-down"
            alt="轮播图"
            class="carousel-image"
          ></el-image>
        </div>
      </el-carousel-item>
    </el-carousel>
    

    <!-- 项目介绍 -->
    <div class="project-description">
      <h1>项目介绍</h1>
      <p>{{appName}}{{ projectDescription }}</p>
    </div>

    <!-- 关键词 -->
    <div class="keywords-container">
      <!-- 类别循环 -->
      <el-card v-for="item in keywords_features" :key="item.id" class="keyword-card">
        <div slot="header" class="card-header">
          <span class="card-title">{{ item.category }}</span>
          <el-tag type="info" effect="dark">{{ item.feature.length }} {{ item["desc-text"] }}</el-tag>
        </div>
        <el-row :gutter="20">
          <!-- 特性循环 -->
          <div class="feature-container"  :style="{ backgroundImage: `url(${item['bg-img']})`,backgroundSize: 'cover' }">
            <el-col :span="6" v-for="feature in item.feature" :key="feature">
              <div class="feature-item">
                <i class="el-icon-check" style="font-size: 24px; color: #1890ff;"></i>
                <span class="feature-text">{{ feature }}</span>
              </div>
            </el-col>
          </div>
        </el-row>
      </el-card>
    </div>

    <!-- 底部 -->
    <footer class="footer">
  <div class="github-container">
    <div class="github-icon">
      <a :href="githubLink" target="_blank">
        <img src="@/assets/relevent-brand-logo/github-fill.svg" alt="github" class="github-icon-img">
      </a>
    </div>
    <div class="github-link">
      <i icon="el-icon-coin"> <a :href="githubLink" target="_blank">GitHub 仓库链接</a></i>
    </div>
  </div>
  <div class="logo-bar">
    <img src="@/assets/relevent-brand-logo/1.png" alt="" class="logo-item" />
    <img src="@/assets/relevent-brand-logo/2.png" alt="" class="logo-item" />
    <img src="@/assets/relevent-brand-logo/3.png" alt="" class="logo-item" />
    <img src="@/assets/relevent-brand-logo/4.png" alt="" class="logo-item" />
    <img src="@/assets/relevent-brand-logo/5.png" alt="" class="logo-item" />
    <img src="@/assets/relevent-brand-logo/6.png" alt="" class="logo-item" />
  </div>
</footer>
  </div>
</template>

<script>
export default {
  data() {
    return {
      carouselItems: [],
      projectDescription: `是一个由深圳大学计算机与软件学院的学生组成的团队所开发的一个基于Flux.1模型的文本到图像生成Web平台。
      实现了账号注册,文本到图片的生成,图片收藏,图片下载等丰富功能。`,
      keywords_features: [
      {
    id: 1,
    category: '文本到图片',
    feature: [
      "自动化生成",
      "节省时间",
      "高质量图像",
      "多种风格选择",
      "易于定制",
      "支持多语言",
      "实时预览",
      "云存储集成"
    ],
    'bg-img': 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202411282246111.png', // 您可以根据实际情况替换这里的值
    'desc-text': '优势'  // 同上
  },
  {
    id: 2,
    category: 'Web应用',
    feature: [
      "方便访问",
      "快捷操作",
      "稳定性强",
      "一站直达",
      "跨平台兼容",
      "个性化设置",
      "数据安全",
      "持续更新"
    ],
    'bg-img': 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202411282346198.png',
    'desc-text': '特点'
  },
  {
    id: 3,
    category: '探索发现',
    feature: [
      "收藏功能",
      "评论互动",
      "批量导出",
      "高清预览",
      "标签分类",
      "趋势分析",
      "智能推荐",
      "社交分享"
    ],
    'bg-img': 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202411282351707.png',
    'desc-text': '功能'
  },
  {
    id: 4,
    category: '历史记录',
    feature: [
      "活动日志",
      "搜索历史",
      "版本回溯",
      "最近使用的列表",
      "过滤器",
      "删除选项",
      "备份与恢复",
      "同步功能",
      "安全加密"
    ],
    'bg-img': 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202411282345690.png',
    'desc-text': '特性'
  },
  {
    id: 5,
    category: '收藏界面',
    feature: [
      "个人收藏夹",
      "自定义排序",
      "分组管理",
      "备注功能",
      "隐私设置",
      "动态更新提醒",
      "视觉美化",
      "交互式元素",
      "分享选项"
    ],
    'bg-img': 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202411282349099.png',
    'desc-text': '功能'
  }
      ],
      githubLink: 'https://github.com/index-null/software-engineering'
    };
  },
  computed: {
    appName() {
      return this.$store.state.appName;
    },
  },
  methods: {
    goToRegister() {
      this.$router.push('/register');
    }
  },
  created() {
    for (let i = 1; i <= 9; i++) {
      this.carouselItems.push({
        id: i,
        image: require(`@/assets/about-images/carousel-imgs/${i}.png`)
      });
    }
  },
};
</script>

<style scoped>
.project-intro {
  background: #000000;
  color: white;
  padding: 20px;
  font-family: 'Arial', sans-serif;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  border-bottom: 2px solid #ffffff;
  margin-bottom: 20px;
}

.left-section {
  display: flex;
  align-items: center;
}

.logo-container {
  width: 50px;
  height: 50px;
  margin-right: 20px;
}

.logo {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.app-name {
  font-size: 3rem;
  margin: 0;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
  text-align: center;
  font-family: 'Arial', sans-serif;
  background: linear-gradient(270deg, rgb(255, 255, 255), rgba(150, 150, 150, 0.8), rgba(255, 255, 255));
  -webkit-background-clip: text;
  background-clip: text; /* 添加标准属性 */
  -webkit-text-fill-color: transparent;
  background-size: 200% 200%;
  animation: gradientAnimation 8s ease-in-out infinite;
  transform: translateY(0);
  animation: floatAnimation 4s ease-in-out infinite;
}
.right-section {
  display: flex;
  align-items: center;
}

.link {
  color: white;
  text-decoration: none;
  margin-left: 20px;
  font-size: 1.2rem;
  transition: color 0.3s ease;
}

.link:hover {
  color: #1890ff;
}

.carousel-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.project-description {
  margin-top: 40px;
  padding: 20px;
  border: 2px solid #ffffff;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.5);
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.9);
  background-image: url(https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202408311347062.jpg);
}

.project-description h1 {
  font-size: 2rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.project-description p {
  font-size: 1.2rem;
  line-height: 1.6;
}

.keywords-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-around;
  padding: 20px;
  background-color: #000000;
  color: #ffffff;
}

.keyword-card {
  width: 45%;
  margin-bottom: 20px;
  background-color: #000000;
  border: none;
}

.card-title {
  color: #ffffff;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-tag {
  margin-left: 10px;
}

.el-row {
  margin-top: 10px;
}

.el-col {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 10px;
}

.feature-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px; /* 间距 */
  border: 2px solid #d4d4d4;
  border-radius: 5px;
  padding: 10px; /* 内边距 */
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.feature-item:hover {
  background-color: #595959; /* 浅灰色背景 */
  color: #000000; /* 文字颜色变深 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影 */
  border-radius: 8px; /* 圆角 */
  transition: background-color 0.3s ease, box-shadow 0.3s ease, color 0.3s ease;
  cursor: pointer;
}

.feature-text {
  font-size: 1rem;
  margin-top: 5px;
  color: #ffffff;
}

.category-text {
  font-size: 0.8rem;
  color: #909399;
  margin-top: 5px;
}

.footer {
  display: flex;
  justify-content: space-between; /* 调整为 space-between */
  padding: 20px;
  border-top: 2px solid #ffffff;
  margin-top: 40px;
  background-color: #000000;
}

/* 将 github-icon 和 github-link 包含在一个容器中 */
.github-container {
  display: flex;
  align-items: center;
}

.github-icon {
  width: 24px;
  height: 24px;
  margin-right: 10px; /* 增加右侧间距 */
}

.github-icon-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.github-link {
  margin: 0;
  text-align: center;
}

.github-link a {
  color: white;
  text-decoration: none;
  font-size: 1.2rem;
  transition: transform 0.3s ease;
}

.github-link a:hover {
  transform: scale(1.1);
}

.logo-bar {
  display: flex;
  align-items: center;
  gap: 20px;
  border: 2px solid #ffffff;
  border-radius: 10px;
  background: linear-gradient(270deg, #585555, #ab9797, #946161); /* 初始渐变背景 */
  background-size: 200% 100%; /* 渐变背景的大小 */
  animation: gradientAnimation 5s ease infinite; /* 动画应用 */
  padding: 10px;
}

@keyframes gradientAnimation {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}
@keyframes floatAnimation {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-5px);
  }
}
.logo-item {
  width: 40px;
  height: 40px;
  object-fit: contain;
}
</style>
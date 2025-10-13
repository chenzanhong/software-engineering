<template>
  <div class="explore-container">
    <!-- 轮播图容器 -->
    <div class="carousel-img-container">
      <div class="carousel-img-bg" :style="{ backgroundImage: `url(${styles[currentStyle].img_url})` }">
        <!-- 样式选择器 -->
        <div class="selector">
          <div 
            v-for="(style, index) in styles" 
            :key="style.id" 
            :class="['selector-item', currentStyle === index ? '' : 'selector-item-unselected']"
            @click="selectStyle(index)"
            :style="{ backgroundColor: currentStyle === index ? style.buttonColor : '', color: currentStyle === index ? style.textColor : '' }"
          >
            {{ style.text1 }}
          </div>
        </div>
        <!-- 文本框 -->
        <div class="text-box">
          <h1 :style="{ color: styles[currentStyle].h1Color }">{{ styles[currentStyle].text1 }}</h1>
          <h3 :style="{ color: styles[currentStyle].h3Color }">{{ styles[currentStyle].text2 }}</h3>
        </div>
        <!-- 生成同款按钮 -->
        <div class="generate-same-style">
          <input 
            class="search-input" 
            type="text" 
            :placeholder="styles[currentStyle].prompt" 
            disabled
          >
          <div 
            class="generate-button" 
            @click="$router.push({ 
              name: 'generate', 
              query: { prompt: encodeURIComponent(styles[currentStyle].prompt) } 
            });"
            :style="{ backgroundColor: styles[currentStyle].buttonColor, color: styles[currentStyle].textColor }"
          >
            生成同款
          </div>
        </div>
      </div>
    </div>
    <!-- 生成入口容器 -->
    <div class="generate-entry-container">
      <div class="generate-entry">
        <div class="text1">文字作画</div>
        <div class="text2">
          文字描述生成画作
        </div>
        <div class="generate-button" @click="$router.push('/main/generate')">去生成</div>
      </div>
    </div>
    <!-- 图片广场容器 -->
    <div class="explore-ground-container">
      <h2 class="ground-title">图片广场</h2>
      <div class="image-grid">
        <div 
          class="image-item" 
          v-for="image in images" 
          :key="image.id"
          @dblclick="likeImage(image.id)"
        >
          <el-image
            style="width: 200px; height: 200px"
            :src="image.picture"
            fit="cover"
          />
          <div class="image-info">
            <span class="like-count" :class="{ 'liked': image.isliked }">
              <i class="el-icon-thumb" :class="{ 'liked': image.isliked }"></i> {{ image.likecount }}
            </span>
            <span class="username">生成用户: {{ image.username }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      images: [],
      styles: [
  {
    id: 1,
    img_url: require("@/assets/Picture-Square/1.png"),
    text1: "琼楼玉宇",
    text2: "中式色彩国风美学",
    prompt: "华丽的中国传统古代建筑，坐落在山腰上，山脚下有弯曲的河流流过，天空中飘着祥云",
    buttonColor: "#1f5a6b", // 保持不变
    textColor: "white", // 保持不变
    h1Color: "#ffffff", // 保持不变
    h3Color: "#ffffff" // 保持不变
  },
  {
    id: 2,
    img_url: require("@/assets/Picture-Square/2.png"),
    text1: "墨韵流芳",
    text2: "水墨留白古色古香",
    prompt: "中国古典女性，传统发型，回眸看镜头，中国传统服装",
    buttonColor: "#8B4513", // 棕色
    textColor: "white", // 棕色
    h1Color: "#8B4513", // 棕色
    h3Color: "#A0522D" // 浅棕色
  },
  {
    id: 3,
    img_url: require("@/assets/Picture-Square/3.png"),
    text1: "华丽时光旅者",
    text2: "厚涂风格原画细节",
    prompt: "一位黑色长发和黑色眼睛的年轻女性，黑色连衣裙，金玫瑰背景，柔和金黄色灯光",
    buttonColor: "#725F57", // 金色
    textColor: "white", // 棕色
    h1Color: "white", // 棕色
    h3Color: "white" // 浅棕色
  }
],
    currentStyle: 0 // 默认选中第一个样式
    };
  },
  mounted() {
    this.fetchImages();
  },
  methods: {
    fetchImages() {
      this.$axios.get('/auth/imageSquare')
        .then(response => {
          this.images = response.data.images;
        })
        .catch(error => {
          console.error("Error fetching images:", error);
        });
    },
    selectStyle(index) {
    this.currentStyle = index;
  },
  likeImage(imageId) {
    const image = this.images.find(img => img.id === imageId);
    if (!image) {
      console.error('Image not found');
      return;
    }

    // 先将本地的 isliked 字段设置为 true
    image.isliked = true;

    this.$axios.post('/auth/like', { url: image.picture })
      .then(response => {
        if (response.status === 200) {
          image.likecount = response.data.current_likes;
          this.$message.success('点赞成功');
        } else {
          // 如果点赞失败，恢复 isliked 状态
          image.isliked = false;
          this.$message.error(response.data.error);
        }
      })
      .catch(error => {
        // 如果请求失败，恢复 isliked 状态
        image.isliked = false;
        console.error('Error liking image:', error);
        this.$message.error(error.response ? error.response.data.error : '请求失败');
      });
  }
  }
};
</script>

<style scoped>
* {
    box-sizing: border-box; /* 确保所有元素的宽度和高度包括内边距和边框 */
    margin: 0;
    padding: 0;
}

.explore-container {
    display: flex;
    flex-direction: column;
    min-height: 120vh; /* 最小高度 */
    background-color: #F1F6FF;
}

.carousel-img-container {
    flex: 1.1; 
    border-radius: 1.5rem;
    position: relative;
    color: white;
}

.generate-entry-container {
    flex: 0.4; /* 减少 flex 值 */
    border-radius: 1.5rem;
    position: relative;
    color: white;
}
.explore-ground-container{
    flex: 8;
    border-radius: 1.5rem; /* 使用相对单位 */
    position: relative; /* 使子元素可以绝对定位 */
    color: white; /* 设置文本颜色为白色 */
}


.carousel-img-bg {
    background-image: url("@/assets/Picture-Square/1.png");
    background-size: cover;
    width: 100%;
    height: auto;
    top: 0;
    left: 0;
    border-radius: 1.5rem; /* 与容器圆角一致 */
}

.selector {
    display: flex;
    gap: 1rem; /* 使用相对单位 */
    margin-left: 2rem; /* 使用相对单位 */
}

.selector-item,
.selector-item-unselected {
    background-color: #BED7DA;
    border-radius: 3rem; /* 使用相对单位 */
    padding: 1rem 2rem; /* 使用相对单位 */
    color: white; /* 文本颜色设置为白色 */
    cursor: pointer; /* 添加鼠标悬停效果 */
    transition: background-color 0.3s; /* 添加过渡效果 */
}

.selector-item {
    background-color: #B9AFA8; /* 添加背景颜色 */
}

.selector-item:hover,
.selector-item-unselected:hover {
    background-color: #B9AFA8; /* 悬停时改变背景色 */
}

.generate-same-style {
    display: flex;
    margin-top: 2rem; /* 使用相对单位 */
    align-items: center; /* 垂直居中对齐 */
    margin-left: 6rem; /* 使用相对单位 */
    margin-right: 6rem; /* 使用相对单位 */
    margin-bottom: 6rem; /* 使用相对单位 */
    position: relative; /* 确保生成按钮和其他元素不重叠 */
    z-index: 1; /* 确保生成按钮在背景之上 */
    
}

.search-input {
    flex: 10;
    padding: 1rem 1.5rem;
    border: none;
    border-radius: 3rem; /* 圆角 */
    background-color: rgba(255, 255, 255, 0.5); /* 半透明背景 */
    color: #333;
    font-size: 1rem;
    outline: none; /* 去掉聚焦时的默认边框 */
    text-align: left; /* 确保placeholder左侧对齐 */
    transition: background-color 0.3s; /* 添加过渡效果 */
    margin-bottom: 2rem;
}

.search-input:focus {
    background-color: rgba(255, 255, 255, 1); /* 聚焦时背景变白 */
}

.generate-button {
    flex: 1;
    background-color: #ff6347; /* 添加按钮背景色 */
    color: white;
    border: none;
    border-radius: 3rem; /* 圆角 */
    padding: 1rem 2rem;
    cursor: pointer;
    transition: background-color 0.3s; /* 添加过渡效果 */
    margin-bottom: 2rem;
}

.generate-button:hover {
    background-color: #ff4500; /* 悬停时改变背景色 */
}

.generate-entry-container {
    flex: 1;
    display: flex;
    align-items: flex-end;
    margin-bottom: 2rem; /* 使用相对单位 */
    margin-top: 0.5rem; /* 使用相对单位 */
}

.generate-entry {
    width: 50%;
    height: 100%;
    background-image: url("@/assets/Picture-Square/it.gif");
    background-size: cover; /* 保持图片完整显示 */
    background-repeat: no-repeat; /* 防止图片重复 */
    background-position: center; /* 图片居中显示 */
    border-radius: 2.5rem;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    align-items: flex-start;
    padding: 2rem; /* 添加内边距 */
    box-sizing: border-box;
}

.text1 {
    color: black; /* 设置文本颜色为黑色 */
    margin-bottom: 1rem; /* 使用相对单位 */
}

.text2 {
    color: rgb(97, 90, 90);
    margin-bottom: 1rem; /* 使用相对单位 */
    font-size: 0.875rem; /* 使用相对单位 */
}

.generate-button {
    background-color: #007BFF; /* 蓝色背景 */
    color: white; /* 白色文字 */
    border: none;
    border-radius: 3rem;
    padding: 1rem 2rem;
    cursor: pointer;
    transition: background-color 0.3s; /* 添加过渡效果 */
}

.generate-button:hover {
    background-color: #0056b3; /* 悬停时改变背景色 */
}

.explore-ground-container {
  flex: 8;
  padding: 2rem;
  box-sizing: border-box;
}

.ground-title {
  font-size: 2rem;
  margin-bottom: 2rem;
  text-align: left;
  color: black;
  font-family: Arial, Helvetica, sans-serif;
}

.image-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: space-between;
}

.image-item {
  flex: 0 0 calc(25% - 1rem); /* 一行展示4张图片 */
  position: relative;
  overflow: hidden;
  border-radius: 1rem;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
  transition: transform 0.3s; /* 添加过渡效果 */
}

.image-item:hover {
  transform: scale(1.05); /* 悬停时放大图片 */
}
.image-info {
  display: flex;
  justify-content: space-around;
  align-items: center;
  margin-top: 0.5rem;
}

.like-count {
  color: #606266;
}

.username {
  color: #409EFF;
}

.el-image__inner {
  border-radius: 1rem; /* 保持图片圆角 */
}
.like-count.liked,
.el-icon-thumb.liked {
  color: red; /* 设置点赞图标和点赞数为红色 */
}
.text-box {
    margin-top: 2rem; /* 使用相对单位 */
    display: flex;
    flex-direction: column;
    text-align: left;
    margin-left: 3rem; /* 使用相对单位 */
    margin-top: 9rem; /* 使用相对单位 */
}

.text-box h1 {
    font-size: 2rem; /* 使用相对单位 */
    margin-bottom: 0.5rem; /* 使用相对单位 */
}

.text-box h3 {
    font-size: 1.2rem; /* 使用相对单位 */
}
</style>
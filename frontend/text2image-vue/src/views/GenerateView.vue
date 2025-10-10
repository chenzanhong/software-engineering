<template>
  <div class="main-container">
    <!-- 主容器，包含表单和结果展示区域 -->
    <div class="form-container">
      <!-- 表单容器，包含输入字段和提交按钮 -->
      <div class="form-header">
        <!-- 表单头部，包含标题和教程链接 -->
        <div class="form-title">文字作画</div>
        <div class="tutorial" @click="$router.push('/usage')">使用指南</div>
      </div>
      <div class="form-appname">文绘星河</div>
      <!-- 应用名称 -->
      <div class="form-body">
        <!-- 表单主体，包含各个输入项 -->
        <div class="form-item">
          <!-- 输入项：画面描述 -->
          <div class="form-label">画面描述</div>
          <el-input
            type="textarea"
            autosize
            :rows="8"
            placeholder="试试输入你心中的画面,尽量描述具体,可以尝试一些风格修饰词辅助你的表达"
            v-model="form.prompt"
            maxlength="100"
            show-word-limit
          />
        </div>
        <div class="form-item">
          <!-- 输入项：API调用方式 -->
          <div class="form-label">API调用方式</div>
          <div class="api-mode">
            <el-button :type="apiMode === 'remote' ? 'primary' : 'default'" @click="setApiMode('remote')">远程API</el-button>
            <el-button :type="apiMode === 'local' ? 'primary' : 'default'" @click="setApiMode('local')">本地API</el-button>
          </div>
        </div>
        <div class="form-item">
          <!-- 输入项：尺寸 -->
          <div class="form-label">尺寸</div>
          <el-select v-model="selectedSize" placeholder="请选择尺寸" @change="updateSize">
            <el-option
              v-for="item in sizeOptions"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </div>
        <div class="form-item">
          <!-- 输入项：步数 -->
          <div class="form-label">步数</div>
          <el-select v-model="form.steps" placeholder="请选择步数">
            <el-option
              v-for="item in stepsOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </div>
        <div class="form-item">
          <!-- 输入项：种子 -->
          <div class="form-label">种子</div>
          <el-input v-model="form.seed" placeholder="请输入种子值">
            <template #suffix>
              <i class="el-input__icon el-icon-refresh" @click="generateRandomSeed"></i>
            </template>
          </el-input>
        </div>
        <div class="form-submit">
          <!-- 提交按钮 -->
          <el-button type="primary" native-type="submit" @click="handleSubmit">生成 (消耗20积分)</el-button>
        </div>
      </div>
    </div>
    <div class="result-container">
      <!-- 结果展示容器，包含生成的图片和相关信息 -->
      <div class="result-header">
        <!-- 结果头部，包含应用名称和重新生成按钮 -->
        <div class="appName">文绘星河</div>
        <div class="regenerate">
          <button @click="regenerateImage">再次生成</button>
        </div>
      </div>
      <div class="result-content" v-loading="loading">
        <!-- 结果内容，包含生成的图片卡片 -->
        <div class="image-card" v-for="(img, index) in temp_generatedImg_results" :key="index">
          <el-image
            style="width: 100%; height: 100%; border-radius: 8px;"
            :src="img.img_url"
            fit="contain"
            lazy
          />
          <div class="overlay">
            <!-- 图片卡片上的覆盖层，包含图片信息和操作按钮 -->
            <div class="image-info">
              <p>Prompt: {{ img.prompt }}</p>
              <p>Width: {{ img.width }}</p>
              <p>Height: {{ img.height }}</p>
              <p>Seed: {{ img.seed }}</p>
              <p>Steps: {{ img.steps }}</p>
            </div>
            <div class="image-buttons">
              <el-button type="primary" icon="el-icon-edit" circle @click="reuseParameters(img)"></el-button>
              <el-button type="warning" icon="el-icon-star-off" circle @click="favoriteImage(img)"></el-button>
              <el-button type="danger" icon="el-icon-delete" circle @click="deleteImage(index)"></el-button>
            </div>
          </div>
        </div>
      </div>
      <div v-if="temp_generatedImg_results.length === 0" class="placeholder">生成的图片将在这里显示</div>
      <!-- 占位符，当没有生成图片时显示 -->
    </div>
  </div>
</template>
<script>
import OSS from 'ali-oss';
export default {
  data() {
    return {
      form: {
        prompt: '',
        width: 1024,
        height: 1024,
        steps: 10,
        seed: 1024
      },
      stepsOptions: [
        { value: 10, label: '10' },
        { value: 15, label: '15' },
        { value: 20, label: '20' },
        { value: 25, label: '25' },
        { value: 30, label: '30' },
        { value: 35, label: '35' },
        { value: 40, label: '40' }
      ],
      temp_generatedImg_results: [ {
            "prompt": "陈绮贞的专辑图片<<Pussy>>",
            "width": 1024,
            "height": 1024,
            "seed": 1024,
            "steps": 10,
            "img_url": `https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202408311347062.jpg`,
          }],
      loading: false,
      apiMode: 'remote',
      sizeOptions: [
    { label: '512 x 1024px', value: '512x1024' },
    { label: '768 x 1024px', value: '768x1024' },
    { label: '576 x 1024px', value: '576x1024' },
    { label: '768 x 512px', value: '768x512' },
    { label: '1024 x 576px', value: '1024x576' },
    { label: '1024 x 1024px', value: '1024x1024' }
],
    selectedSize: '1024x1024'
    };
  },
  methods: {
    async initOSSClient() {
      this.client = new OSS({
        region: process.env.VUE_APP_OSS_REGION,
        accessKeyId: process.env.VUE_APP_OSS_ACCESS_KEY_ID,
        accessKeySecret: process.env.VUE_APP_OSS_ACCESS_KEY_SECRET,
        bucket: process.env.VUE_APP_OSS_BUCKET,
      });
    },
    async uploadImageToOSS(img_url) {
      try {
        await this.initOSSClient();
        const base64Data = img_url.split(',')[1];
        const blob = this.b64toBlob(base64Data, 'image/png');
        const fileName = `generated_image_${Date.now()}.png`;
        const result = await this.client.put(`generated_images/${fileName}`, blob);
        const imageUrl = result.url;

        // 更新图片 URL
        const lastIdx = this.temp_generatedImg_results.length - 1;
        this.$set(this.temp_generatedImg_results[lastIdx], 'img_url', imageUrl);
         // 发送历史记录请求
    await this.$axios.post('http://localhost:8080/auth/generate/addhistory', {
      prompt: this.form.prompt,
      width: this.form.width,
      height: this.form.height,
      seed: this.form.seed,
      steps: this.form.steps,
      pictureURL: imageUrl
    });
        this.$message.success('图片上传成功');
      } catch (error) {
        console.error('Error uploading image:', error);
        this.$message.error
      }
    },
    b64toBlob(b64Data, contentType = '', sliceSize = 512) {
      const byteCharacters = atob(b64Data);
      const byteArrays = [];

      for (let offset = 0; offset < byteCharacters.length; offset += sliceSize) {
        const slice = byteCharacters.slice(offset, offset + sliceSize);

        const byteNumbers = new Array(slice.length);
        for (let i = 0; i < slice.length; i++) {
          byteNumbers[i] = slice.charCodeAt(i);
        }

        const byteArray = new Uint8Array(byteNumbers);
        byteArrays.push(byteArray);
      }

      const blob = new Blob(byteArrays, {type: contentType});
      return blob;
    },
    setApiMode(mode) {
      this.apiMode = mode;
    },
    handleSubmit() {
      if (this.apiMode === 'local') {
        this.handleLocalSubmit();
      } else {
        this.handleRemoteSubmit();
      }
    },
    reuseParameters(img) {
    this.form.prompt = img.prompt;
    this.form.width = img.width;
    this.form.height = img.height;
    this.form.seed = img.seed;
    this.form.steps = img.steps;
    this.selectedSize = `${img.width}x${img.height}`;
  },
  favoriteImage(img) {
    console.log('Favorite image:', img);
    this.$axios.post('http://localhost:8080/auth/addFavoritedImage', {url: img.img_url}, {
                        headers: {
                            'Content-Type': 'application/json', // 设置请求头
                        },                      
                    })
      .then(response => {
        if (response.status === 200) {
          this.$message.success('收藏成功');
        } else {
          this.$message.error(response.data.message);
        }
      })
      .catch(error => {
        this.$message.error(error.response ? error.response.data.message : '请求失败');
      });
  },
  deleteImage(index) {
    this.temp_generatedImg_results.splice(index, 1);
  },
    updateSize() {
  const selected = this.sizeOptions.find(option => option.value === this.selectedSize);
  if (selected) {
    console.log('Old width:', this.form.width, 'Old height:', this.form.height);
    this.form.width = parseInt(selected.value.split('x')[0], 10);
    this.form.height = parseInt(selected.value.split('x')[1], 10);
    console.log('New width:', this.form.width, 'New height:', this.form.height);
  }
},
    regenerateImage() {
      this.loading = true;
      this.handleSubmit().finally(() => {
        this.loading = false;
      });
    },
    handleLocalSubmit() {
      // 强制将 seed 转换为整数
      this.form.seed = parseInt(this.form.seed, 10) || 0;

      // 确保转换后的 seed 是有效的整数
      if (isNaN(this.form.seed)) {
        this.$message.error('种子值必须是有效的整数');
        return;
      }
      this.$message.success('提交成功,正在生成图片...');
      const currentScore = parseInt(localStorage.getItem("score"), 10) || 0;
      localStorage.setItem("score", currentScore - 20);

      // 添加占位图片
      const placeholderImg = {
        prompt: this.form.prompt,
        width: this.form.width,
        height: this.form.height,
        seed: this.form.seed,
        steps: this.form.steps,
        img_url: 'https://via.placeholder.com/150'}
      this.temp_generatedImg_results.push(placeholderImg);

      this.loading = true;

      const url = "http://172.26.224.111:7860"; // 使用WSL中的Windows主机IP地址
      const testPayload = {
        "prompt": this.form.prompt,
        "seed": this.form.seed,
        "sampler_name": "Euler",
        "scheduler": "Simple",
        "batch_size": 1,
        "steps": this.form.steps,
        "cfg_scale": 1,
        "distilled_cfg_scale": 3.5,
        "width": this.form.width,
        "height": this.form.height
      };

      return this.$axios.post(`${url}/sdapi/v1/txt2img`, testPayload, {
        timeout: 300000 // 设置超时时间为300秒
      }).then(response => {
        if (response && response.data && 'images' in response.data) {
          console.log(response.data);
          let img_item = {
            "prompt": this.form.prompt,
            "width": this.form.width,
            "height": this.form.height,
            "seed": this.form.seed,
            "steps": this.form.steps,
            "img_url": `data:image/png;base64,${response.data.images[0]}`,
          };
          console.log(img_item);

          // 替换占位图片
          const lastIdx = this.temp_generatedImg_results.length - 1;
          this.$set(this.temp_generatedImg_results, lastIdx, img_item);

          this.$message.success(response.data.message);
          this.uploadImageToOSS(img_item.img_url);
        } else {
          this.$message.error('服务器返回数据异常');
        }
      }).catch(error => {
        this.$message.error(error.response ? error.response.data.message : '请求失败');
      }).finally(() => {
        this.loading = false;
      });
    },
    handleRemoteSubmit() {
      // 强制将 seed 转换为整数
      this.form.seed = parseInt(this.form.seed, 10) || 0;
  
      // 确保转换后的 seed 是有效的整数
      if (isNaN(this.form.seed)) {
        this.$message.error('种子值必须是有效的整数');
        return;
      }
      this.$message.success('提交成功,正在生成图片...');
      const currentScore = parseInt(localStorage.getItem("score"), 10) || 0;
      localStorage.setItem("score", currentScore - 20);
      // 添加占位图片
      const placeholderImg = {
        prompt: this.form.prompt,
        width: this.form.width,
        height: this.form.height,
        seed: this.form.seed,
        steps: this.form.steps,
        img_url: 'https://via.placeholder.com/150'}
      this.temp_generatedImg_results.push(placeholderImg);
  
      this.loading = true;

      return this.$axios.post('http://localhost:8080/auth/generate', this.form, {
        timeout: 300000 // 设置超时时间为300秒
      }).then(response => {
        if (response && response.data) {
          console.log(response.data);
          let img_item = {
            "prompt": this.form.prompt,
            "width": this.form.width,
            "height": this.form.height,
            "seed": this.form.seed,
            "steps": this.form.steps,
            "img_url": response.data.image_url,
          };
          console.log(img_item);

          // 替换占位图片
          const lastIdx = this.temp_generatedImg_results.length - 1;
          this.$set(this.temp_generatedImg_results, lastIdx, img_item);

          this.$message.success(response.data.message);
        } else {
          this.$message.error('服务器返回数据异常');
        }
      }).catch(error => {
        this.$message.error(error.response ? error.response.data.message : '请求失败');
      }).finally(() => {
        this.loading = false;
      });
    },
    generateRandomSeed() {
      this.form.seed = Math.floor(Math.random() * 4369000);
    },
    
  },
  watch: {
    '$route.query.prompt': {
      handler(newPrompt) {
        if (newPrompt) {
          try {
            this.form.prompt = decodeURIComponent(newPrompt);
            console.log("Decoded Prompt:", this.form.prompt);
          } catch (e) {
            console.error("解码失败:", e);
            this.form.prompt = ''; // 或者其他的错误处理逻辑
          }
        }
      },
      immediate: true // 立即执行一次，确保首次加载时也生效
    }
  },

};
</script>
<style scoped>
/* 主容器样式 */
.main-container {
  height: 98vh;
  width: 100%;
  display: flex;
  justify-content: space-between;
  gap: 1vw;
  background-color: #f0f2f5;
  font-family: Arial, Helvetica, sans-serif;
}

/* 表单容器样式 */
.form-container {
  flex: 1;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow-y: auto;
}

/* 表单头部样式 */
.form-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

/* 表单标题样式 */
.form-title {
  font-size: 22px;
  font-weight: bold;
}

/* 教程样式 */
.tutorial {
  font-size: 12px;
  font-weight: bolder;
  color: #5f6163;
  border: 1px gray solid;
  border-radius: 30px;
  padding: 6px 20px;
}

/* 表单主体样式 */
.form-body {
  display: flex;
  flex-direction: column;
  height: calc(100% - 60px); /* Adjust height to accommodate header and footer */
}
.form-appname{
  text-align: left;
  font-size: 12px;
  font-weight: bolder;
  color: #000000;
  background-color: #F7F8FC;
  border-radius: 30px;
  padding: 6px 20px;
  margin-bottom: 5vh;
}
/* 表单项样式 */
.form-item {
  display: flex;
  flex-direction: column;
  margin-bottom: 15px;
}

/* 表单标签样式 */
.form-label {
  font-size: 18px;
  margin-bottom: 5px;
  text-align: left;
}

/* 子标签样式 */
.form-sub-label {
  font-size: 14px;
  margin-bottom: 5px;
}

/* 输入框样式 */
.el-input,
.el-input-number,
.el-select {
  width: 100%;
}

/* 数字输入框样式 */
.el-input-number {
  width: 100%;
}

/* 选择框样式 */
.el-select {
  width: 100%;
}

/* 选项样式 */
.el-option {
  width: 100%;
}

/* 提交按钮样式 */
.form-submit {
  display: flex;
  justify-content: center;
  width: 90%;
  margin: 0 auto;
}


/* 结果容器样式 */
.result-container {
  flex: 3.5;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

/* 结果头部样式 */
.result-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

/* 应用名称样式 */
.appName {
  font-size: 22px;
  font-weight: bold;
}

/* 重新生成按钮样式 */
.regenerate button {
  font-size: 14px;
  font-weight: bold;
  color: #5f6163;
  border: 1px gray solid;
  border-radius: 30px;
  padding: 6px 20px;
  background-color: transparent;
}

.result-content {
  display: flex;
  flex-wrap: wrap; /* 允许换行 */
  justify-content: flex-start; /* 从左侧开始 */
  align-items: flex-start;
  padding: 20px;
  gap: 10px; /* 设置间距 */
}

/* 生成的图片样式 */
.generated-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

/* 占位符样式 */
.placeholder {
  font-size: 16px;
  color: #909399;
  text-align: center;
  padding: 20px;
}

/* 其他未指定样式的组件类 */
.el-row,
.el-col {
  width: 100%;
}

.el-button {
  width: 100%;
}

.image-card {
  position: relative; /* 添加相对定位 */
  overflow: hidden;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease, opacity 0.3s ease; /* 添加 opacity 过渡 */
  height: 280px;
  width: 280px;
}

.image-card:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
  opacity: 0.9; /* 可选：调整悬停时的透明度 */
}

.overlay {
  position: absolute; /* 设置为绝对定位 */
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5); /* 半透明黑色背景 */
  display: flex;
  justify-content: space-around;
  align-items: center;
  opacity: 0; /* 初始透明度为0 */
  transition: opacity 0.3s ease; /* 添加过渡效果 */
  flex-direction: column; /* 确保内容垂直居中 */
}

.image-card:hover .overlay {
  opacity: 1; /* 鼠标悬停时透明度为1 */
}

.image-info {
  color: white;
  text-align: center;
  margin-bottom: 10px;
  overflow-y: auto; /* 添加垂直滚动 */
  max-height: 200px; /* 设置最大高度 */
}

.image-buttons {
  display: flex;
  gap: 10px;
}
/* 美化滚动条 */
.image-info::-webkit-scrollbar {
  width: 8px;
}

.image-info::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

.image-info::-webkit-scrollbar-thumb {
  background: rgba(100, 94, 94, 0.6);
  border-radius: 4px;
}

.image-info::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 1);
}

/* 结果内容样式 */
.result-content {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
}

/* 刷新图标样式 */
.el-input__icon.el-icon-refresh {
  cursor: pointer;
}

::v-deep .el-textarea__inner {
  background-color: #F7F8FC; /* 浅灰色背景 */
  color: #5951f2; /* 深灰色文本 */
  font-weight: bold; /* 加粗 */
  border-radius: 6px;
  font-family: Arial, Helvetica, sans-serif;
}

.api-mode {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}
</style>
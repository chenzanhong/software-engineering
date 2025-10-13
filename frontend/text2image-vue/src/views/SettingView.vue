<template>
  <div class="setting-container">
    <!-- 顶部闪动条 -->
    <div class="shinning-bar"></div>

    <!-- 用户资料卡片 -->
    <el-card class="user-profile">
      <!-- 用户资料头部 -->
      <div class="profile-header">
        <!-- 头像容器，点击时触发头像点击事件 -->
        <div class="avatar-container" @click="handleAvatarClick">
          <img :src="user.avatar" alt="User Avatar" class="avatar">
          <!-- 编辑模式下显示的头像遮罩层 -->
          <div class="avatar-overlay" v-show="isEditing"></div>
        </div>

        <!-- 用户基本信息展示，编辑模式下隐藏 -->
        <div class="profile-info" v-if="!isEditing">
          <div class="info-item">
            <span class="info-label">昵称:</span>
            <span class="info-value">{{ user.username }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">邮箱:</span>
            <span class="info-value">{{ user.email }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">个签:</span>
            <span class="info-value">{{ user.personalSignature }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">性别:</span>
            <span class="info-value">{{ user.gender }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">爱好:</span>
            <span class="info-value">{{ user.hobbies }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">居住地:</span>
            <span class="info-value">{{ user.location }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">积分:</span>
            <span class="info-value">{{ user.collectedPictures }}</span>
          </div>
        </div>
      </div>

      <!-- 签到按钮，已签到则禁用 -->
      <el-button type="success" @click="signIn" class="sign-in-button" :disabled="hasSignedInToday" v-show="!isEditing">
        {{ hasSignedInToday ? '已签到' : '签到' }}
      </el-button>

      <!-- 编辑按钮，点击后进入编辑模式 -->
      <el-button type="primary" @click="toggleEdit" class="edit-button" v-show="!isEditing">
        编辑
      </el-button>

      <!-- 注销按钮，点击后确认注销 -->
      <el-button type="danger" @click="confirmLogout" class="logout-button" v-show="!isEditing">
        注销
      </el-button>

      <!-- 编辑表单，编辑模式下显示 -->
      <el-form v-if="isEditing" label-width="80px" class="edit-form">
        <el-form-item label="昵称">
          <el-input v-model="user.username" disabled class="form-item-input"></el-input>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="user.email" class="form-item-input"></el-input>
        </el-form-item>
        <el-form-item label="个签">
          <el-input v-model="user.personalSignature" class="form-item-input"></el-input>
        </el-form-item>
        <el-form-item label="性别">
          <el-select v-model="user.gender" placeholder="请选择性别" class="form-item-input">
            <el-option label="男" value="男"></el-option>
            <el-option label="女" value="女"></el-option>
            <el-option label="保密" value="保密"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="爱好">
          <el-input v-model="user.hobbies" class="form-item-input"></el-input>
        </el-form-item>
        <el-form-item label="居住地">
          <el-select v-model="user.location" placeholder="请选择省份" class="form-item-input">
            <el-option v-for="province in provinces" :key="province" :label="province" :value="province"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="积分">
          <el-input v-model="user.collectedPictures" disabled class="form-item-input"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="saveChanges" class="save-button">保存</el-button>
        </el-form-item>
      </el-form>

      <!-- 隐藏的文件输入控件，用于上传头像 -->
      <input type="file" ref="fileInput" @change="handleFileChange" style="display: none;">
    </el-card>
  </div>
</template>

<script>
import OSS from 'ali-oss';
export default {
  data() {
    return {
      user: {
        avatar: localStorage.getItem('avatarUrl') || 'https://whxh-czh.oss-cn-shenzhen.aliyuncs.com/generate/main.jpg',
        username: localStorage.getItem('username') || '未知用户',
        email: localStorage.getItem('email') || '未知邮箱',
        collectedPictures: localStorage.getItem('score') || 0,
        personalSignature: localStorage.getItem('personalSignature') || '快创建一个个性签名吧!',
        gender: localStorage.getItem('gender') || '未知', 
        hobbies: localStorage.getItem('hobbies') || '未知',
        location: localStorage.getItem('location') || '未知',
      },
      isEditing: false,
      hasSignedInToday: false, // 添加此属性以跟踪签到状态
      client: null,
      provinces: [
        '北京市', '天津市', '上海市', '重庆市', '河北省', '山西省', '辽宁省', '吉林省', '黑龙江省', '江苏省', '浙江省', '安徽省', '福建省', '江西省', '山东省', '河南省', '湖北省', '湖南省', '广东省', '海南省', '四川省', '贵州省', '云南省', '陕西省', '甘肃省', '青海省', '台湾省', '内蒙古自治区', '广西壮族自治区', '西藏自治区', '宁夏回族自治区', '新疆维吾尔自治区', '香港特别行政区', '澳门特别行政区'
      ]
    };
  },
  methods: {
    handleAvatarClick() {
      if (this.isEditing) {
        this.deleteOldAvatar().then(() => {
          this.$refs.fileInput.click();
        });
      }
    },
    async deleteOldAvatar() {
      // 📢 打印配置信息（用于调试）
       const region = process.env.VUE_APP_OSS_REGION;
  const accessKeyId = process.env.VUE_APP_OSS_ACCESS_KEY_ID;
  const accessKeySecret = process.env.VUE_APP_OSS_ACCESS_KEY_SECRET;
  const bucket = process.env.VUE_APP_OSS_BUCKET;
      console.log('OSS Config:', {
        region,
        accessKeyId: accessKeyId ? '✔️ 已设置' : '❌ 未设置',
        accessKeySecret: accessKeySecret ? '✔️ 已设置' : '❌ 未设置',
        bucket,
      });
      if (this.user.avatar) {
        const urlParts = this.user.avatar.split('/');
        const fileName = urlParts[urlParts.length - 1];
        try {
          await this.initOSSClient();
          await this.client.delete(`avator/${fileName}`);
          this.$message.success('旧头像删除成功');
        } catch (error) {
          this.$message.error('删除旧头像失败');
        }
      }
    },
    async initOSSClient() {
      this.client = new OSS({
        region: process.env.VUE_APP_OSS_REGION,
        accessKeyId: process.env.VUE_APP_OSS_ACCESS_KEY_ID,
        accessKeySecret: process.env.VUE_APP_OSS_ACCESS_KEY_SECRET,
        bucket: process.env.VUE_APP_OSS_BUCKET,
      });
    },
    async handleFileChange(event) {
      const file = event.target.files[0];
      console.log('Selected file:', file); // 添加这行来调试
      if (!file) return;

      try {
        await this.initOSSClient();
        const result = await this.client.put(`avator/${file.name}`, file);
        this.user.avatar = result.url;
        localStorage.setItem('avatarUrl', this.user.avatar);
        this.$message.success('头像上传成功');
      } catch (error) {
        console.error('Upload failed:', error); // More detailed error logging
        this.$message.error('上传失败');
      }
    },
    toggleEdit() {
      this.isEditing = !this.isEditing;
    },
    // 签到
    checkSignInStatus() {
      const lastSignIn = localStorage.getItem('lastSignInDate');
      if (!lastSignIn) {
        this.hasSignedInToday = false;
        return;
      }

      const todayStart = new Date();
      todayStart.setHours(0, 0, 0, 0);
      
      const lastSignInDate = new Date(lastSignIn);
      lastSignInDate.setHours(0, 0, 0, 0);

      this.hasSignedInToday = todayStart.getTime() === lastSignInDate.getTime();
    },
    setupResetTimer() {
      const now = new Date();
      const tomorrow = new Date(now);
      tomorrow.setDate(now.getDate() + 1);
      tomorrow.setHours(0, 0, 0, 0);

      const timeToTomorrow = tomorrow.getTime() - now.getTime();

      setTimeout(() => {
        this.hasSignedInToday = false;
        localStorage.removeItem('lastSignInDate'); // 清除本地存储中的最后签到日期
      }, timeToTomorrow);
    },
    async signIn() {
      try {
        const response = await this.$axios.get('/auth/score', {
          headers: {
            Authorization: localStorage.getItem('token') || ''
          }
        });

        if (response.status === 200) {
          this.hasSignedInToday = true; // 更新签到状态
          this.user.collectedPictures = response.data.score; // 更新积分显示
          this.$message.success(response.data.message);
          // 更新本地存储中的最后签到日期
          localStorage.setItem('lastSignInDate', new Date().toISOString());
        } else {
          this.$message.error('签到失败');
        }
      } catch (error) {
        console.error('签到失败:', error.response?.data || error.message);
        this.$message.error('签到失败，请稍后再试');
      }
    },

    // 账号注销
    confirmLogout() {
      const isConfirmed = confirm('确认注销吗？注销后本系统将清除关于你的所有信息记录。');
      if (isConfirmed) {
        this.performLogout();
      } else {
        this.$message({
          type: 'info',
          message: '已取消注销'
        });
      }
    },

    performLogout() {
      this.$axios.delete('/auth/root/deleteOneUser', {
        params: {
          isOwn: true
        },
        headers: {
          Authorization: localStorage.getItem('token') || ''
        }
      }).then(response => {
        if (response.status === 200) {
          this.$message.success('注销成功');
          // 清除本地存储的数据和执行其他清理工作
          localStorage.clear();
          // 可能还需要重定向到登录页面或主页等
          this.$router.push({ name: 'log-reg' }); // 重定向到登录/注册页面
        } else {
          this.$message.error('服务器返回异常');
        }
      }).catch(error => {
        console.error('注销失败:', error.response?.data || error.message);
        this.$message.error('注销失败，请稍后再试');
      });
    },

    saveChanges() {
    localStorage.setItem('personalSignature', this.user.personalSignature);
    localStorage.setItem('gender', this.user.gender);
    localStorage.setItem('hobbies', this.user.hobbies);
    localStorage.setItem('location', this.user.location);
    let updatedUser = {
      "email": this.user.email,
      "avatar_url": this.user.avatar
    };

    this.$axios.put('/auth/user/update', updatedUser)
      .then(response => {
        if (response && response.data) {
          this.$message.success(response.data.message);
        } else {
          this.$message.error('服务器返回数据异常');
        }
      })
      .catch(error => {
        this.$message.error(error.response ? error.response.data.message : '请求失败');
      });

    this.isEditing = false;
    this.$nextTick(this.$forceUpdate);
  }
  },
  created() {
  },
  mounted() {
    this.checkSignInStatus();
    this.setupResetTimer();
  }
};
</script>

<style scoped>
.setting-container {
  background-color: #F1F6FF;
  min-height: 100vh;
  width: 90vw;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.shinning-bar {
  height: 18vh;
  width: 90vw;
  margin-bottom: 6vh;
  background: linear-gradient(to right, #ff7e5f, #feb47b,#6a89c1,#F1F6FF); /* 彩色渐变效果 */
}

.user-profile {
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
  text-align: center;
  background-color: #f7f7f7; /* 苹果风格的背景色 */
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.profile-header {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  margin-bottom: 20px;
}

.avatar-container {
  position: relative;
  cursor: pointer;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  transition: box-shadow 0.3s ease;
}

.avatar-container:hover {
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 18px;
}

.profile-info {
  margin-top: 10px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  animation: fadeIn 0.5s ease-in-out; /* 添加简单的动画效果 */
}

.info-label {
  font-size: 16px;
  color: #555; /* 苹果风格的标签颜色 */
}

.info-value {
  font-size: 16px;
  color: #333; /* 苹果风格的值颜色 */
}

.sign-in-button {
  margin-left: 10px;
  animation: fadeIn 0.5s ease-in-out; /* 添加简单的动画效果 */
}

.edit-button {
  margin-top: 20px;
  animation: fadeIn 0.5s ease-in-out; /* 添加简单的动画效果 */
}
.save-button {
  margin-left: -10%;
  animation: fadeIn 0.5s ease-in-out; /* 添加简单的动画效果 */
}

.logout-button {
  margin-left: 10px;
  animation: fadeIn 0.5s ease-in-out; /* 添加简单的动画效果 */
}

.edit-form {
  margin-top: 20px;
  margin-left: -10%;
  margin-bottom: 10px;
}

.form-item-input {
  width: 100%; /* 设置宽度为100%，或者根据需要设置固定宽度 */
}

/* 定义动画 */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
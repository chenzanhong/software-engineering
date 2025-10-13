<template>
  <div class="container">
    <!-- 整体容器 -->
    <div class="log-reg-container">
      <!-- 左侧容器，包含欢迎文案 -->
      <div class="left-side-container">
        <div class="text-bold">用简单的文案</div>
        <div class="text-bold-smaller">创作精彩的图片!</div>
      </div>
      
      <!-- 右侧容器，包含登录和注册表单 -->
      <div class="right-side-container">
        <!-- 标题区域，包含应用图标和名称 -->
        <div class="title">
          <img src="@/assets/button-icon/文生图-gray.svg" alt="">
          <div class="title-text">{{ appName }}</div>
        </div>
        
        <!-- 标签页组件，包含登录和注册两个标签页 -->
        <el-tabs type="border-card">
          <!-- 登录标签页 -->
          <el-tab-pane label="登录">
            <h2>登录</h2>
            <!-- 登录表单 -->
            <el-form
              :model="loginForm"
              status-icon
              :rules="loginRules"
              ref="loginForm"
              label-position="left"
              label-width="70px"
              class="login-form"
            >
              <el-form-item label="用户名" prop="username">
                <el-input v-model="loginForm.username"></el-input>
              </el-form-item>
              <el-form-item label="密码" prop="password">
                <el-input
                  type="password"
                  v-model="loginForm.password"
                  autocomplete="off"
                ></el-input>
              </el-form-item>
            </el-form>
            <!-- 登录表单按钮组 -->
            <div class="btn-group">
              <el-button type="primary" @click="submitLoginForm('loginForm')">登录</el-button>
              <el-button @click="resetLoginForm('loginForm')">重置</el-button>
            </div>
          </el-tab-pane>
          
          <!-- 注册标签页 -->
          <el-tab-pane label="注册">
            <h2>注册</h2>
            <!-- 注册表单 -->
            <el-form
              :model="registerForm"
              status-icon
              :rules="registerRules"
              ref="registerForm"
              label-position="left"
              label-width="80px"
              class="register-form"
            >
              <el-form-item label="用户名" prop="uname">
                <el-input v-model="registerForm.uname"></el-input>
              </el-form-item>
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="registerForm.email"></el-input>
              </el-form-item>
              <el-form-item label="密码" prop="pass">
                <el-input
                  type="password"
                  v-model="registerForm.pass"
                  autocomplete="off"
                ></el-input>
              </el-form-item>
              <el-form-item label="确认密码" prop="password">
                <el-input
                  type="password"
                  v-model="registerForm.password"
                  autocomplete="off"
                ></el-input>
              </el-form-item>
            </el-form>
            <!-- 注册表单按钮组 -->
            <div class="btn-group">
              <el-button type="primary" @click="submitRegisterForm('registerForm')">提交</el-button>
              <el-button @click="resetRegisterForm('registerForm')">重置</el-button>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import apiClient from '@/api';

export default {
  computed: {
    appName() {
      return this.$store.state.appName;
    },
  },

  data() {
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'));
      } else if (value.length < 6 || value.length > 16) {
        callback(new Error('密码长度必须在6到16个字符之间'));
      } else {
        if (this.registerForm.password !== '') {
          this.$refs.registerForm.validateField('password');
        }
        callback();
      }
    };
    var validatePass2 = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请再次输入密码"));
      } else if (value !== this.registerForm.pass) {
        callback(new Error("两次输入密码不一致!"));
      } else {
        callback();
      }
    };
    return {
      loginForm: {
        username: '',
        password: '',
      },
      loginRules: {
        username: [
          { required: true, message: "用户名不能为空！", trigger: "blur" },
        ],
        password: [
          { required: true, message: "密码不能为空！", trigger: "blur" },
          { min: 6, max: 16, message: '密码长度必须在6到16个字符之间', trigger: 'change' }
        ],
      },
      registerForm: {
        uname: "",
        email: "",
        pass: "",
        password: "",
      },
      registerRules: {
        uname: [{ required: true, message: "请输入用户名", trigger: "blur" }],
        email: [
          { required: true, message: "请输入邮箱地址", trigger: "blur" },
          { type: "email", message: "请输入正确的邮箱地址", trigger: ["blur", "change"] },
        ],
        pass: [{ required: true, validator: validatePass, trigger: "blur" }],
        password: [
          { required: true, validator: validatePass2, trigger: "blur" },
        ],
      },
    };
  },
  methods: {
  async hashPassword(password) {
    const encoder = new TextEncoder();
    const data = encoder.encode(password);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    return hashHex;
  },
  async submitLoginForm(formName) {
    this.$refs[formName].validate(async (valid) => {
      if (valid) {
        let hashedPassword = await this.hashPassword(this.loginForm.password);
        let formData = {
          username: this.loginForm.username,
          password: hashedPassword
        };
        try {
          const response = await apiClient.post('/login', formData);
          if (response.data.code === 200) {
            localStorage.setItem('token', response.data.token);
            localStorage.setItem('username', this.loginForm.username);
            localStorage.setItem('avatarUrl', response.data.avatar);
            this.$message.success('登录成功');
            this.$router.push('/main');
          } 
        } catch (error) {
          this.$message.error(error.response ? error.response.data.message : '请求失败');
        }
      } else {
        console.log("error submit!!");
        return false;
      }
    });
  },
  resetLoginForm(formName) {
    this.$refs[formName].resetFields();
  },
  async submitRegisterForm(formName) {
    this.$refs[formName].validate(async (valid) => {
      if (valid) {
        try {
          const hashedPassword = await this.hashPassword(this.registerForm.password);
          const formattedFormData = {
            email: this.registerForm.email,
            username: this.registerForm.uname,
            password: hashedPassword
          };
          const response = await this.$axios.post('/register', formattedFormData);
          if (response.data.code === 200) {
            this.$message.success('注册成功');
          } 
        } catch (error) {
          this.$message.error(error.response ? error.response.data.message : '请求失败');
        }
      }
    });
  },
    resetRegisterForm(formName) {
      this.$refs[formName].resetFields();
    },
  },
};
</script>

<style scoped>
.log-reg-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  background-color: 0F131C;
}

.left-side-container {
  flex: 1;
  background-image: url(https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202412091935661.jpg);
  background-size: cover;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  border-radius: 10px; /* 添加圆角 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影 */
}
.text-bold {
  font-weight: bold;
  font-size: 40px;
  color: white;
  margin-top: 10vh;
  margin-bottom: 10px; /* 调整两段文字之间的间距 */
  transform: translateX(-12vw); /* 向左偏移 */
}

.text-bold-smaller {
  font-weight: bold;
  font-size: 30px;
  color: white;
  transform: translateX(-5vw); /* 向右偏移 */
}

.right-side-container {
  flex: 1;
  display: flex;
  gap: 0px;
  align-items: center;
  flex-direction: column;
  background-color: #0F131C; /* 添加背景色 */
  padding: 20px; /* 添加内边距 */
  border-radius: 10px; /* 添加圆角 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影 */
}

.title {
  font-size: 40px;
  display: flex;
  justify-content: center;
  align-items: center; /* 修改对齐方式 */
  gap: 0;
  height: 10vh;
  margin-bottom: 10vh;
}

.title img {
  width: 70px;
  height: 70px;
  margin-right: 10px;
}

.title-text {
  font-family: 'Roboto', sans-serif; /* 使用现代字体 */

  color: #333333; /* 文字颜色 */
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1); /* 添加文字阴影 */
  background: linear-gradient(to right, #6a11cb, #2575fc); /* 渐变背景色 */
  -webkit-background-clip: text; /* 将背景应用于文字 */
  -webkit-text-fill-color: transparent; /* 使文字透明以显示背景 */
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333333;
}

.login-form, .register-form {
  margin: 0 auto;
  width: 100%;
  max-width: 400px; /* 设置最大宽度 */
}

.btn-group {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.el-button {
  flex: 1;
  margin: 0 5px;
}

.el-button + .el-button {
  margin-left: 10px;
}

.el-input {
  width: 100%;
  margin-bottom: 10px;
  border-radius: 4px;
  transition: border-color 0.3s ease;
}

/* 新增样式 */
.el-form-item__label {
  color: #555555; /* 调整标签颜色 */
}

.el-input__inner {
  border: 1px solid #dcdcdc; /* 调整输入框边框颜色 */
  transition: border-color 0.3s ease;
}

.el-input__inner:focus {
  border-color: #409EFF; /* 输入框聚焦时边框颜色 */
}

.el-button--primary {
  background-color: #409EFF; /* 主按钮背景颜色 */
  border-color: #409EFF; /* 主按钮边框颜色 */
}

.el-button--primary:hover {
  background-color: #66b1ff; /* 主按钮悬停时背景颜色 */
  border-color: #66b1ff; /* 主按钮悬停时边框颜色 */
}

.el-button--default {
  background-color: #ffffff; /* 默认按钮背景颜色 */
  border-color: #dcdcdc; /* 默认按钮边框颜色 */
}

.el-button--default:hover {
  background-color: #ebeef5; /* 默认按钮悬停时背景颜色 */
  border-color: #c6e2ff; /* 默认按钮悬停时边框颜色 */
}
</style>
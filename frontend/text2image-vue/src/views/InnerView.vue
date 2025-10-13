<template>
  <div class="inner-view">
    <!-- 左侧容器 -->
    <div class="left-side-container">
      <!-- 上侧容器，包含logo和导航栏 -->
      <div class="upside-container">
        <!-- Logo 和 标题 -->
        <div class="logo-container">
          <img :src="svgPaths.logo" alt="Logo" class="logo"/>
          <span class="title">{{ appName }}</span>
        </div>
        <!-- 导航栏 -->
        <nav class="nav-container">
          <!-- 探索发现 -->
          <div class="nav-item" :class="{ 'active': activeRoute === '/main/explore' }" @click="setActiveRoute('/main/explore')">
            <img :src="svgPaths.explore" alt="Explore" class="nav-icon"/>
            <router-link to="/main/explore" class="nav-link">探索发现</router-link>
          </div>
          <!-- 文字作画 -->
          <div class="nav-item" :class="{ 'active': activeRoute === '/main/generate' }" @click="setActiveRoute('/main/generate')">
            <img :src="svgPaths.generate" alt="generate" class="nav-icon"/>
            <router-link to="/main/generate" class="nav-link">文字作画</router-link>
          </div>
          <!-- 我的收藏 -->
          <div class="nav-item" :class="{ 'active': activeRoute === '/main/favorites' }" @click="setActiveRoute('/main/favorites')">
            <img :src="svgPaths.favorite" alt="favourite" class="nav-icon"/>
            <router-link to="/main/favorites" class="nav-link">我的收藏</router-link>
          </div>
          <!-- 历史记录 -->
          <div class="nav-item" :class="{ 'active': activeRoute === '/main/history' }" @click="setActiveRoute('/main/history')">
            <img :src="svgPaths.history" alt="History" class="nav-icon"/>
            <router-link to="/main/history" class="nav-link">历史记录</router-link>
          </div>
        </nav>
      </div>
      <!-- 下侧容器，包含用户信息和设置 -->
      <div class="downside-container">
        <!-- 用户头像 -->
        <div class="avatar-container">
          <el-avatar :size="100" :src="avatarUrl"></el-avatar>
        </div>
        <!-- 退出登录确认 -->
        <el-popconfirm
          confirm-button-text='是'
          cancel-button-text='取消'
          icon="el-icon-info"
          icon-color="red"
          title="确定要退出登录吗？"
          @confirm="handleLogout"
        >
          <div class="nav-item" slot="reference">
            <img :src="svgPaths.avatar" alt="info" class="nav-icon"/>
            <router-link to="#" class="nav-link">{{ username }}</router-link>
          </div>
        </el-popconfirm>
        <!-- 账户信息 -->
        <div class="nav-item" :class="{ 'active': activeRoute === '/main/settings' }" @click="setActiveRoute('/main/settings')">
          <img :src="svgPaths.setting" alt="info" class="nav-icon"/>
          <router-link to="/main/setting" class="nav-link">账户信息</router-link>
        </div>
      </div>
    </div>
    <!-- 视图容器，用于显示路由视图 -->
    <div class="view-container">
      <router-view></router-view>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      activeRoute: this.$route.path,
      svgPaths: {
      logo: require('../assets/button-icon/文生图-gray.svg'),
      explore: require('../assets/button-icon/explore-gray.svg'),
      generate: require('../assets/button-icon/img-gray.svg'),
      favorite: require('../assets/button-icon/favourite-gray.svg'),
      history: require('../assets/button-icon/history-gray.svg'),
      avatar: require('../assets/button-icon/avator-gray.svg'),
      setting: require('../assets/button-icon/setting-gray.svg')
      }
    };
  },
  computed: {
    appName() {
      return this.$store.state.appName;
    },
    avatarUrl() {
      return localStorage.getItem('avatarUrl') || 'https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202412101217874.png';
    },
    username() {
      return localStorage.getItem('username') || '未知用户';
    }
  },
  methods: {
    setActiveRoute(route) {
      this.activeRoute = route;
    },
    handleLogout() {
      // 清空 token
      localStorage.removeItem('token');
      localStorage.removeItem('lastSignInDate');
      // 跳转到登录注册页面
      this.$router.push('/log-reg');
    }
  },
  watch: {
    $route(to) {
      this.activeRoute = to.path;
    }
  },
  created() {
//     {
//     "user": {
//         "id": 1,
//         "email": "root@example.com",
//         "username": "root",
//         "password": "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a",
//         "avatar_url": "https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202407272335307.png",
//         "score": 10000,
//         "token": "",
//         "create_time": "0001-01-01T00:00:00Z"
//     }
// }
    this.$axios.get('/auth/user/info').then(response => {
      localStorage.setItem('avatarUrl', response.data.user.avatar_url);
      localStorage.setItem('username', response.data.user.username);
      localStorage.setItem('email',response.data.user.email)
      localStorage.setItem('score',response.data.user.score)
      localStorage.setItem('createTime',response.data.user.create_time)
      localStorage.setItem('score', response.data.user.score);
      this.$message.success('用户信息获取成功');
      this.$nextTick(() => {
        this.$forceUpdate(); // 强制更新组件
      });
    }).catch(error => {
      this.$message.error('获取用户信息失败');
      console.error('获取用户信息失败:', error);
    });
  }
}
</script>
<style scoped>
.inner-view, .inner-view * {
  font-family: 'PingFang SC', sans-serif;
  font-weight: 700;
}

.inner-view {
  display: flex;
  height: 100vh;
  background-color: #f0f2f5;
}

.left-side-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 200px;
  height: 100vh;
  background-color: #F1F6FF;
  color: #ffffff;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
}

.upside-container {
  flex: 5;
  display: flex;
  flex-direction: column;
  gap: 10%;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-content: center;
}

.logo-container img {
  max-width: 100%;
  height: auto;
  max-height: 30px;
}

.logo-container .title {
  font-size: 20px;
  color: #4944e8;
}

.nav-container {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 5px;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-item:hover {
  background-color: #fff;
  border-radius: 20px;
}

.nav-item.active {
  background-color: #fff;
  border-radius: 20px;
  padding: 10px;
}

.nav-item.active .nav-link {
  color: #615CED;
}

.nav-icon {
  max-width: 100%;
  height: auto;
  max-height: 30px;
}

.nav-link {
  color: #707070;
  text-decoration: none;
  font-size: 16px;
  padding: 10px;
  border-radius: 4px;
  transition: background-color 0.3s ease;
}
.downside-container {
  padding-bottom: 40px;
}

.view-container {
  flex: 10;
  background-color: #ffffff;
  padding: 20px;
  margin-left: 200px;
  box-shadow: -2px 0 10px rgba(0, 0, 0, 0.1);
}
</style>
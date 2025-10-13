<template>
  <div>
    <!-- 标题区域 -->
    <h2 class="history-title">历史记录</h2> 

    <!-- 日期选择器和查询按钮 -->
    <div class="block" style="margin-left:50px;margin-top: 50px;">
      <span class="demonstration"></span>
      <span class="data-picker-lable">
        选择日期：
      </span>
      <el-date-picker
        v-model="value1"
        type="datetimerange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
      ></el-date-picker>
      <el-button icon="el-icon-check" circle @click="handleDateChange()"></el-button>
    </div>

    <!-- 批量管理按钮 -->
    <el-button 
      type="primary" 
      icon="el-icon-delete" 
      class="delete-button" 
      @click="removemoreRecord"
    >
      批量管理
    </el-button>

    <!-- 搜索框 -->
    <div style="margin-top: 50px;">
      <el-input 
        placeholder="请输入内容" 
        v-model="input" 
        class="input-with-select"
      >
        <el-button 
          slot="append" 
          icon="el-icon-search" 
          @click="searchimage"
        >
          搜索
        </el-button>
      </el-input>
    </div>

    <!-- 历史记录展示区域 -->
    <div class="history-container">
      <!-- 如果有历史记录则显示图片卡片 -->
      <div v-if="historyRecords && historyRecords.length" class="image-gallery-container">
        <div 
          v-for="image in historyRecords" 
          :key="image.id" 
          class="image-card" 
          @mouseover="hoveredImage = image.id" 
          @mouseleave="hoveredImage = null"
        >
          <el-checkbox v-model="checked[image.id]"></el-checkbox><br />
          <img :src="image.url" :alt="image.name" class="image">
          
          <!-- 悬停时显示的操作按钮 -->
          <div class="overlay" v-if="hoveredImage === image.id">
            <button circle @click="downloadImage(image)">下载图像</button>
            <button round @click="deleteRecord(image)">删除</button> 
            <button round @click="addFavoriteImage(image)">收藏</button> 
          </div>
        </div>
      </div>

      <!-- 如果没有历史记录则显示提示信息 -->
      <div v-else>
        <img :src="require('@/assets/nofavorites.png')" alt="暂无历史记录">
        <h1>暂无历史记录</h1>
      </div>
    </div>
  </div>
</template>

<script>
// import axios from 'axios';
import apiClient from '@/api';
export default {
  data() {
    return {
      historyRecords: [],
      token: localStorage.getItem('token') || '',
      hoveredImage: null,  // 用于追踪当前悬停的图片
      pickerOptions: {
        shortcuts: [
          {
            text: '最近一周',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit('pick', [start, end]);
            }
          },
          {
            text: '最近一个月',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              picker.$emit('pick', [start, end]);
            }
          },
          {
            text: '最近三个月',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
              picker.$emit('pick', [start, end]);
            }
          }
        ]
      },
      value1: [new Date(2000, 10, 10, 10, 10), new Date(2030, 10, 11, 10, 10)],
      value2: '',
      checked: {},
      input: '',
    };
  },
  methods: {
    // 按关键字搜索图像
    async searchimage(){
      try {
        console.log(this.input);
        // 分割输入的字符串为多个feature，使用正则表达式匹配空白字符和中文、英文逗号，并过滤掉空字符串
        const features = this.input.split(/[\s，,]+/).filter(Boolean);

        // 创建URLSearchParams 对象用于构建查询字符串
        const params = new URLSearchParams();
        features.forEach(feature => params.append('feature', feature));

        const response = await apiClient.get(
            '/auth/image/feature',
            {
                headers: {
                    'Authorization': this.token,  // 携带 token
                },
                params: params,
            }
        );
        const data = response.data;
        console.log(data);
        
        const images = data.images;
        this.historyRecords = [];//置为空
        images.forEach(item => {
          if (!this.historyRecords.some(record => record.id === item.id)) {
            this.historyRecords.push({
              id: item.id,
              name: item.username,
              params: item.params,
              url: item.picture,
              likecount: item.likecount,
              createtime: item.create_time
            });
          } else {
            console.log(`记录 id: ${item.id} 已经存在，跳过添加`);
          }
        });
      } catch (error) {
        console.error('获取图片失败:', error.response?.data || error.message);
      }
    },
    async addFavoriteImage(image) {
      try {
                
                const response = await apiClient.post(
                    '/auth/addFavoritedImage',
                    {url:image.url},
                    {
                        headers: {
                            'Authorization': localStorage.getItem('token'),  // 携带 token
                            'Content-Type': 'application/json', // 设置请求头
                        },                      
                    }
                );
                if (response.status === 200) {
                    // this.getHistoryImages();  // 收藏成功后重新获取收藏列表
                    this.$message.success('收藏图像成功');
                }
            } catch (error) {
              if (error.response) {
                // 请求成功发出但服务器返回了非2xx的状态码
                if (error.response.status === 409) {
                  this.$message.warning("收藏失败：您之前已经收藏了该图像");
                } else {
                  console.error('收藏图像失败:', error.response.data || error.message);
                  this.$message.error('收藏图像失败，请稍后再试或联系管理员');
                }
              } else {
                // 对于网络错误或其他情况，直接使用 error.message
                console.error('收藏图像失败:', error.message);
                this.$message.error('收藏图像失败，请检查网络连接后重试');
              }
            }
        },
    async getHistoryImages() {
      try {
        const response = await apiClient.get("/auth/user/images", {
          headers: {
            Authorization: this.token,  // 携带 token
          },
        });

        const data = await response.data;
        //console.log(data);
        console.log(localStorage.getItem('token')? '有token':'没有token');
        console.log(1234567);

        const images = data.images;
        images.forEach(item => {
          if (!this.historyRecords.some(record => record.id === item.id)) {
            this.historyRecords.push({
              id: item.id,
              name: item.username,
              params: item.params,
              url: item.picture,
              likecount: item.likecount,
              createtime: item.create_time
            });
          } else {
            console.log(`记录 id: ${item.id} 已经存在，跳过添加`);
          }
        });
      } catch (error) {
        console.error('获取收藏的图片失败:', error.response?.data || error.message);
      }
    },
    formatDateToISO(date) {
    // 将传入的日期字符串转换为 Date 对象
    const dateObj = new Date(date);
    
    // 获取各个时间部分
    const year = dateObj.getFullYear();
    const month = String(dateObj.getMonth() + 1).padStart(2, '0'); // 月份从0开始，所以需要加1
    const day = String(dateObj.getDate()).padStart(2, '0');
    const hours = String(dateObj.getHours()).padStart(2, '0');
    const minutes = String(dateObj.getMinutes()).padStart(2, '0');
    const seconds = String(dateObj.getSeconds()).padStart(2, '0');
    const milliseconds = String(dateObj.getMilliseconds()).padStart(3, '0');
    
    // 拼接成目标格式: 2006-01-02T15:04:05.000000Z
    return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}Z`;
    },
    handleDateChange() {
      const { value1 } = this;
      if (value1 && value1.length === 2) {
        const [startDate, endDate] =value1;
        console.log('开始:', startDate);
        console.log('结束:', endDate);
        const start = this.formatDateToISO(startDate);  // 转换格式
        const end =  this.formatDateToISO(endDate);
        console.log('开始日期:', start);
        console.log('结束日期:', end);
        this.getPeriodImages(start, end);
      }
    },
    async getPeriodImages(start,end) {
      try {
        const response = await apiClient.get("/auth/user/images/timeRange?start_time="+start+"&end_time="+end, {
          headers: {
            Authorization: this.token,  // 携带 token
          },
          //params:{startDate,endDate},//传递开始时间和结束时间
        });
        console.log('get开始日期:', start);
        console.log('get结束日期:', end);

        const data = await response.data;
        console.log(data);

        const images = data.images;
        this.historyRecords = [];//置为空
        images.forEach(item => {
          if (!this.historyRecords.some(record => record.id === item.id)) {
            
            this.historyRecords.push({
              id: item.id,
              name: item.username,
              params: item.params,
              url: item.picture,
              likecount: item.likecount,
              createtime: item.create_time
            });
          } else {
            console.log(`记录 id: ${item.id} 已经存在，跳过添加`);
          }
        });
      } catch (error) {
        console.error('获取图片失败:', error.response?.data || error.message);
      }
    },
    async removemoreRecord() {
            try {
                const checkedKeys = Object.keys(this.checked).filter(key => this.checked[key]);
                if (!checkedKeys.length) {
                this.$message.warning('请选择要删除的图片');
                return;
                }

                // 使用 Promise.all 并发处理所有选定图片的取消收藏请求
                const promises = checkedKeys.map(async (key) => {
                const image = this.historyRecords.find(img => img.id === parseInt(key));
                if (image) {
                    await this.deleteRecord(image);
                }
                });

                await Promise.all(promises);

                this.$message.success(`成功删除了 ${checkedKeys.length} 张图片的历史记录`);
            } catch (error) {
                console.error('多张图片删除失败:', error.message);
                this.$message.error('批量删除时发生了错误');
            }
        },
    async deleteRecord(image) {
      try {
                const response = await apiClient.post(
                    '/auth/user/deleteImages', {ids:[image.id]},{
                    headers: {
                        Authorization: this.token,
                        'Content-Type':'application/json'
                    },
                    //params: { url: image.url },  // 传递图像的收藏表url
                }
                );
                
                if (response.status === 200) {
                    this.historyRecords = this.historyRecords.filter(i => i.url !== image.url);  // 从历史列表中移除已删除的图像
                    this.$message.success('删除成功');
                }
          } catch (error) {
                console.error('删除失败:', error.response?.data || error.message);
          }
    
    },
    viewImage(record) {
      window.open(record.url, '_blank');
    },
    downloadImage(image) {
            const link = document.createElement('a');
            link.href = image.url;
            link.download = image.name;
            link.click();
    },
  },
  mounted() {
    this.getHistoryImages();
  }
};
</script>

<style scoped>
.date-button{
  margin-left: 20px;
}
 .el-input {
    width: 500px;
  }
  .el-input-group__prepend {
    background-color: #fff;
    width: 200px;
  }
.history-title {
    position: fixed;
    /* 固定位置 */
    top: 30px;
    /* 离页面顶部20px */
    left: 300px;
    /* 离页面右边20px */
    z-index: 900;
    /* 确保按钮显示在页面最上面 */
    border-radius: 30px;
}
.delete-button {
    position: fixed;
    /* 固定位置 */
    top: 50px;
    /* 离页面顶部20px */
    right: 50px;
    /* 离页面右边20px */
    z-index: 900;
    /* 确保按钮显示在页面最上面 */
    border-radius: 30px;
}
.overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 50%;
    top:50%;
    height: 50%;
    background: whitesmoke;
    display: flex;
    justify-content: center;
    align-items: center;
    opacity: 0;
    transition: opacity 0.3s ease;
}

.image-gallery-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(365px, 1fr));
  gap: 80px;
  padding: 20px;
  margin-left: 100px;
  margin-top: 50px;
  width: 80%;
}

.image-card {
  position: relative;
  overflow: hidden;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  height:600px;
  width:400px;
}

.image-card:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

.image {
  width: 90%;
  height: 90%;
  object-fit: cover;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.image-card:hover .overlay {
  opacity: 1;
}

.overlay button {
  background: white;
  border: none;
  padding: 10px 20px;
  margin: 10px;
  cursor: pointer;
  border-radius: 5px;
}

.history-container {
  padding: 20px;
  margin: 0 auto;
  text-align: center;
}

.history-records {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 20px;
  
}

.record {
  border: 1px solid #ccc;
  padding: 10px;
  border-radius: 4px;
  width: calc(33.333% - 20px);
  text-align: center;
}

.history-image {
  display: block;
  max-width: 100%;
  height: auto;
  margin: 0 auto;
}

.record-details {
  margin-top: 10px;
}
</style>
<template>
    <!-- 主容器 -->
    <div>
        <!-- 批量管理按钮 -->
        <el-button 
            type="primary" 
            icon="el-icon-delete" 
            class="delete-button" 
            @click="removemoreFavorite"
        >
            批量管理
        </el-button>

        <!-- 图片展示区域 -->
        <div v-if="images && images.length" class="image-gallery-container">
            <!-- 每张图片卡片 -->
            <div 
                v-for="image in images" 
                :key="image.id" 
                class="image-card"
                @mouseover="hoveredImage = image.id"
                @mouseleave="hoveredImage = null"
            >
                <!-- 复选框用于批量选择 -->
                <el-checkbox v-model="checked[image.id]"></el-checkbox><br />
                
                <!-- 图片预览 -->
                <img :src="image.url" :alt="image.name" class="image">
                
                <!-- 悬停显示的操作按钮 -->
                <div class="overlay" v-if="hoveredImage === image.id">
                    <button @click="toggleFavorite(image)">
                        {{ image.isFavorite ? '取消收藏' : '收藏' }}
                    </button>
                    <button @click="downloadImage(image)">下载</button>
                    <!-- <i class="el-icon-delete"></i> -->
                </div>
            </div>
        </div>

        <!-- 暂无收藏提示 -->
        <div v-else>
            <img :src="require('@/assets/nofavorites.png')" alt="暂无收藏">
            <h1>暂无收藏</h1>
        </div>
    </div>
</template>

<script>
// import axios from 'axios';
import apiClient from '@/api';

export default {
    data() {
        return {
            images: [],  // 存储用户的收藏图片
            hoveredImage: null,  // 用于追踪当前悬停的图片
            token: localStorage.getItem('token') || '',  // 获取用户的 token
            checked: {},
        };
    },
    mounted() {
        this.getFavoritedImages();  // 组件挂载时获取用户的收藏图片
    },
    methods: {
        // 获取用户收藏的图片
        async getFavoritedImages() {
            //console.log(this.token);
            try {
                 const response = await apiClient.get('/auth/user/favoritedimages', {
                    headers: {
                        Authorization: this.token,
                    },
                });
                
                const data = await response.json();
                this.images = data.map(item => ({
                id: item.id,
                url: item.picture,
                name: item.username,
                isFavorite: true,
                }));

                // 初始化 checked 对象以匹配新的 images 数组
                this.checked = {}; // 清空旧的 checked
                this.images.forEach(image => {
                    this.$set(this.checked, image.id, false); // 使用 $set 确保反应性
                });
            } catch (error) {
                console.error('获取收藏的图片失败:', error.response?.data || error.message);
            }
        },

        // 收藏图像(此功能无法在收藏界面实现，一旦取消收藏图片会立即消失)
        // async addFavoriteImage(imageId, imageUrl) {
        //     try {
        //         const response = await axios.post(
        //             'http://113.44.170.52:8080/auth/addFavoritedImage',
        //             { id: imageId, url: imageUrl },  // 只需传递 id 或 url
        //             {
        //                 headers: {
        //                     Authorization: this.token,  // 携带 token
        //                 },
        //             }
        //         );
        //         if (response.status === 200) {
        //             this.getFavoritedImages();  // 收藏成功后重新获取收藏列表
        //             this.$message.success('收藏图像成功');
        //         }
        //     } catch (error) {
        //         console.error('收藏图像失败:', error.response?.data || error.message);
        //     }
        // },
        async removemoreFavorite() {
            try {
                const checkedKeys = Object.keys(this.checked).filter(key => this.checked[key]);
                if (!checkedKeys.length) {
                this.$message.warning('请选择要取消收藏的图片');
                return;
                }

                // 使用 Promise.all 并发处理所有选定图片的取消收藏请求
                const promises = checkedKeys.map(async (key) => {
                const image = this.images.find(img => img.id === parseInt(key));
                if (image) {
                    await this.removeFavorite(image);
                }
                });

                await Promise.all(promises);

                this.$message.success(`成功取消了 ${checkedKeys.length} 张图片的收藏`);
            } catch (error) {
                console.error('取消多张图片收藏失败:', error.message);
                this.$message.error('批量取消收藏时发生了错误');
            }
        },
        // 取消收藏图像
        async removeFavorite(image) {
            try {
                const response = await apiClient.delete(
                    '/auth/deleteFavoritedImage', {
                    headers: {
                        Authorization: this.token,
                    },

                    params: { id: image.id },  // 传递图像的收藏表id或url

                    }
                );
                
                if (response.status === 200) { 
                    // 移除已取消收藏的图像之前，先检查是否存在于 images 中
                    if (this.images.some(i => i.id === image.id)) {
                        this.images = this.images.filter(i => i.id !== image.id);  // 从收藏列表中移除已取消收藏的图像
                        this.$message.success('取消收藏成功');
                        delete this.checked[image.id]; // 移除对应的 checked 元素
                    }
                }
            } catch (error) {
                console.error('取消收藏失败:', error.response?.data || error.message);
            }
        },

        // 切换收藏状态
        toggleFavorite(image) {
            if (image.isFavorite) {
                // 如果当前是收藏状态，调用取消收藏
                this.removeFavorite(image);
            } else {
                // 如果当前不是收藏状态，调用添加收藏
                this.addFavoriteImage(image.id, image.url);
            }
            image.isFavorite = !image.isFavorite;  // 切换收藏状态
        },

        // 下载图像
        downloadImage(image) {
            const link = document.createElement('a');
            link.href = image.url;
            link.download = image.name;
            link.click();
        },
    },
};
</script>

<style scoped>
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

.image-gallery-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 80px;
    padding: 20px;
    margin-left: 130px;
    margin-top: 100px;
    width: 80%;
    /* 容器宽度设置为页面宽度的 100% */
}

.image-card {
    position: relative;
    overflow: hidden;
    border-radius: 10px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
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
</style>

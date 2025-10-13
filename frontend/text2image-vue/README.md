# text2image-vue

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Lints and fixes files
```
npm run lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).


## 问题
目前的url是直接访问的，需要在oss上开放“公共读”权限
或者：
```js
const { OSS } = require('@alicloud/oss');

const client = new OSS({
  endpoint: 'oss-cn-shenzhen.aliyuncs.com',
  accessKeyId: 'your-access-key-id',
  accessKeySecret: 'your-access-key-secret',
  bucket: 'whxh-czh'
});

async function getSignedUrl(fileName) {
  const url = await client.signatureUrl(`generate/images/${fileName}`, {
    method: 'GET',
    expires: 3600 // 1小时有效
  });
  return url;
}

// 使用
getSignedUrl('image_2024-12-22_11-20-37.png').then(url => {
  console.log(url); // 输出带签名的 URL
});
```
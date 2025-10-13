const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer:{
    allowedHosts:["caohaitong.xyz"],
    host:"0.0.0.0",
  }
})

<template>
  <div id="app">
    <router-view v-if="!$route.meta.keepAlive" :key="$route.fullPath">
      <!-- 这里是不被缓存的视图组件，比如详情B页面-->
    </router-view>

  </div>
</template>

<script>

  export default {
    name: 'App',
  }

</script>

<style>
  #app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    margin-top: 1px;
  }
</style>

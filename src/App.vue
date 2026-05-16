<template>
  <router-view />
</template>

<style lang="scss">
html,
body,
#app {
  height: 100%;
  padding: 0;
  margin: 0;
  overflow: hidden;
  background-color: rgb(248 248 248);
}

*,
*::before,
*::after {
  box-sizing: border-box;
}

li,
ol {
  list-style: none;
}

// 移动端全局适配
@media (width <= 768px) {
  html,
  body {
    // 防止 iOS Safari 双击缩放
    touch-action: manipulation;

    // 平滑滚动
    -webkit-overflow-scrolling: touch;
  }
}
</style>

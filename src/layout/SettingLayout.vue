<script setup lang="ts">
import { ref } from 'vue';
import {
  HomeFilled,
  Avatar,
  Management,
  Expand,
  Fold,
} from '@element-plus/icons-vue';

// 控制展开状态，false 为展开，true 为收缩
const isCollapse = ref(false);
</script>

<template>
  <div class="setting-layout">
    <el-container>
      <el-header>Header</el-header>
      <el-container>
        <el-aside :class="['admin-aside', isCollapse ? 'is-collapsed' : '']">
          <el-scrollbar>
            <el-menu class="admin-menu" :collapse="isCollapse">
              <el-menu-item index="0">
                <el-icon><HomeFilled /></el-icon>
                <template #title>总览</template>
              </el-menu-item>
              <el-sub-menu index="1">
                <template #title>
                  <el-icon><Avatar /></el-icon>
                  <span>用户管理</span>
                </template>
                <el-menu-item-group>
                  <el-menu-item index="1-1">用户列表</el-menu-item>
                </el-menu-item-group>
              </el-sub-menu>
              <el-sub-menu index="2">
                <template #title>
                  <el-icon><Management /></el-icon>
                  <span>课程管理</span>
                </template>
                <el-menu-item-group>
                  <el-menu-item index="2-1">课程列表</el-menu-item>
                </el-menu-item-group>
              </el-sub-menu>
            </el-menu>
          </el-scrollbar>

          <div class="collapse-trigger" @click="isCollapse = !isCollapse">
            <el-icon :size="20">
              <Expand v-if="isCollapse" />
              <Fold v-else />
            </el-icon>
          </div>
        </el-aside>
        <el-main><router-view /></el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style scoped lang="scss">
.admin-header {
  line-height: 60px;
  border-bottom: 1px solid var(--el-border-color-light);
}

.admin-aside {
  position: relative; // 为右下角按钮提供定位基准
  width: 200px;
  height: calc(100vh - 60px);
  overflow: hidden;
  background-color: rgb(255 255 255);
  border-right: 1px solid var(--el-border-color-light);
  transition: width 0.3s ease;

  .admin-menu {
    border-right: none; // 去掉菜单默认边框

    // 解决折叠时文字闪烁的细节
    &:not(.el-menu--collapse) {
      width: 200px;
    }
  }

  .collapse-trigger {
    position: absolute;
    right: 0;
    bottom: 20px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    width: 100%; // 初始全宽，居中
    height: 40px;
    padding-right: 20px;
    color: var(--el-text-color-secondary);
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      color: var(--el-color-primary);
      background-color: var(--el-fill-color-light);
    }
  }

  &.is-collapsed {
    width: 64px;

    .collapse-trigger {
      justify-content: center;
      padding-right: 0;
    }
  }
}

// 隐藏滚动条背景
:deep(.el-scrollbar__view) {
  height: 100%;
}
</style>

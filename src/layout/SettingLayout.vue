<script setup lang="ts">
import {
  HomeFilled,
  Avatar,
  Management,
  Expand,
  Fold,
  Menu as IconMenu,
} from "@element-plus/icons-vue";
import { ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";

// 控制展开状态，false 为展开，true 为收缩
const isCollapse = ref(false);
// 移动端侧边栏开关
const mobileOpen = ref(false);
const router = useRouter();
const route = useRoute();

// 移动端路由切换后自动关闭侧边栏
watch(
  () => route.path,
  () => {
    mobileOpen.value = false;
  },
);

function goToIndex() {
  router.push({ name: "Index" });
}

function toggleMobileSidebar() {
  mobileOpen.value = !mobileOpen.value;
  // 移动端打开时确保侧边栏是展开状态
  if (mobileOpen.value) {
    isCollapse.value = false;
  }
}
</script>

<template>
  <div class="setting-layout">
    <el-container>
      <el-header class="admin-header">
        <el-row>
          <el-col
            class="link-index"
            :xs="18"
            :sm="8"
            :lg="6"
            @click="goToIndex"
          >
            <el-image class="logo" src="/favicon.svg" />
            <span class="logo-txt">家庭学坊后台管理</span>
          </el-col>
          <el-col :xs="6" :sm="0" class="mobile-menu-col">
            <el-button
              class="mobile-menu-btn"
              :icon="IconMenu"
              text
              @click="toggleMobileSidebar"
            />
          </el-col>
        </el-row>
      </el-header>
      <el-container>
        <el-aside
          :class="[
            'admin-aside',
            isCollapse ? 'is-collapsed' : '',
            mobileOpen ? 'is-mobile-open' : '',
          ]"
        >
          <el-scrollbar>
            <el-menu
              class="admin-menu"
              default-active="/setting"
              :collapse="isCollapse"
              :router="true"
              :default-openeds="['1', '2']"
            >
              <el-menu-item index="/setting">
                <el-icon><HomeFilled /></el-icon>
                <template #title>总览</template>
              </el-menu-item>
              <el-sub-menu index="1">
                <template #title>
                  <el-icon><Avatar /></el-icon>
                  <span>用户管理</span>
                </template>
                <el-menu-item-group>
                  <el-menu-item index="/setting/users">用户列表</el-menu-item>
                </el-menu-item-group>
              </el-sub-menu>
              <el-sub-menu index="2">
                <template #title>
                  <el-icon><Management /></el-icon>
                  <span>课程管理</span>
                </template>
                <el-menu-item-group>
                  <el-menu-item index="/setting/categories"
                    >课程分类</el-menu-item
                  >
                  <el-menu-item index="/setting/courses">课程列表</el-menu-item>
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
        <!-- 移动端遮罩 -->
        <div
          :class="['mobile-overlay', mobileOpen ? 'is-visible' : '']"
          @click="mobileOpen = false"
        />
        <el-main><router-view /></el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style lang="scss">
.admin-header {
  background-color: rgb(255 255 255);
  border-bottom: 1px solid var(--el-border-color-light);

  .link-index {
    cursor: pointer;

    .logo {
      width: 50px;
      height: 50px;
      margin-top: 6px;
    }

    .logo-txt {
      margin-left: 6px;
      font-size: 24px;
      font-weight: bold;
      vertical-align: 15px;
    }
  }

  .mobile-menu-col {
    display: none;
    text-align: right;

    .mobile-menu-btn {
      margin-top: 10px;
      font-size: 24px;
    }
  }
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

/* 表格拖拽时的影子样式保持 */
.sortable-ghost {
  outline: 2px dashed rgb(64 158 255); /* 增加虚线框增强交互感 */
  background-color: rgb(245 247 250 / 60%) !important;
}

/* 表格操作栏统一样式 */
.header-panel {
  margin-bottom: 10px;

  .btns {
    text-align: right;
  }
}

// 移动端适配
@media (width <= 768px) {
  .admin-header {
    .link-index {
      .logo-txt {
        font-size: 16px;
        vertical-align: 12px;
      }

      .logo {
        width: 36px;
        height: 36px;
      }
    }

    .mobile-menu-col {
      display: block;
    }
  }

  .admin-aside {
    position: fixed;
    top: 56px;
    left: 0;
    z-index: 100;
    height: calc(100vh - 56px);
    transform: translateX(-100%);
    transition: transform 0.3s ease;

    // 展开时从左边滑入
    &.is-mobile-open {
      box-shadow: 2px 0 16px rgb(0 0 0 / 15%);
      transform: translateX(0);
    }

    &.is-collapsed {
      width: 64px;
      transform: translateX(0);
    }

    &.is-collapsed.is-mobile-open {
      width: 200px;
    }
  }

  // 移动端遮罩层
  .mobile-overlay {
    position: fixed;
    inset: 0;
    z-index: 99;
    pointer-events: none;
    background: rgb(0 0 0 / 50%);
    opacity: 0;
    transition: opacity 0.3s ease;

    &.is-visible {
      pointer-events: auto;
      opacity: 1;
    }
  }
}

// 表格操作栏移动端适配
@media (width <= 768px) {
  .header-panel {
    // 移动端按钮占满宽度，左对齐
    .el-col {
      flex: 1 1 100%;
      max-width: 100%;
    }

    .btns {
      text-align: left;
    }
  }

  // 表格横向滚动
  :deep(.el-table) {
    // 移动端表格自适应
    .el-table__body-wrapper {
      overflow-x: auto;
    }
  }

  // 对话框宽度适配
  :deep(.el-dialog) {
    --el-dialog-width: 90%;
  }

  // el-main 内边距减小
  .el-main {
    padding: 10px;
  }
}
</style>

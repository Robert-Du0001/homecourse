<script setup lang="ts">
import {
  ArrowLeft,
  ArrowRight,
  Document,
  Fold,
  Expand,
  FullScreen,
  ArrowRight as ArrowRightIcon,
} from "@element-plus/icons-vue";
import DPlayer from "dplayer";
import { ElMessageBox } from "element-plus";
import { onMounted, ref, watch, onBeforeUnmount } from "vue";
import { useRoute, useRouter } from "vue-router";

import type { CatchData } from "@/lib/js/api";
import type { EpisodesResource } from "@/types/episode";
import type { DPlayerDanmaku, DPlayerEvents } from "dplayer";

import { request } from "@/lib/js/api";

const videoRef = ref(null);
const route = useRoute();
const router = useRouter();
const episode = ref<EpisodesResource>();
const isCollapse = ref(false);
const isLandscape = ref(false);
let dp: DPlayer | null = null;

/** 1. 获取并初始化数据 */
async function loadEpisodeData(id: number) {
  // 如果当前已经在看这一集了，就不重复加载
  if (episode.value?.id === id) return;

  try {
    const res = await request<EpisodesResource>("GET", `/episodes/${id}`);
    episode.value = res.data;

    // 初始化或更新播放器
    initOrUpdatePlayer(id);
  } catch (e) {
    console.error(e);
    const rd = e as CatchData;
    ElMessageBox.alert(rd.msg || "剧集加载失败", "温馨提示");
  }
}

/** 2. 播放器逻辑：支持无刷新切换视频源 */
function initOrUpdatePlayer(id: number) {
  const videoUrl = `/videos/${id}`;

  if (!videoRef.value) {
    console.warn("Video container not found, retrying...");
    return;
  }

  if (dp) {
    // 如果播放器已存在，直接切换视频源，体验更顺滑
    dp.switchVideo(
      { url: videoUrl, type: "auto" },
      undefined as unknown as DPlayerDanmaku,
    );
    dp.play();
  } else {
    // 第一次加载时创建实例
    dp = new DPlayer({
      container: videoRef.value,
      autoplay: true,
      screenshot: true,
      hotkey: true,
      playbackSpeed: [0.5, 0.75, 1, 1.25, 1.5, 2, 2.5, 2.75, 3],
      video: { url: videoUrl, type: "auto" },
    });

    // 自动播放下一集逻辑
    dp.on("ended" as DPlayerEvents, () => {
      if (episode.value?.navigation?.next) {
        goToEpisode(episode.value.navigation.next.id);
      }
    });
  }
}

/** 3. 页面跳转逻辑 */
function goToEpisode(id: number) {
  router.push(`/episodes/${id}`);
}

/** 4. 附件相关逻辑（保留你的原始代码） */
function handleViewAttachment(id: number) {
  window.open(`/attachments/${id}`, "_blank");
}

function toggleSidebar() {
  isCollapse.value = !isCollapse.value;
}

/** 5. 横屏播放 */
async function toggleLandscape() {
  const videoEl = videoRef.value as HTMLElement | null;
  if (!videoEl) return;

  if (isLandscape.value) {
    // 退出横屏
    try {
      await document.exitFullscreen();
    } catch {
      // 忽略退出全屏失败
    }
  } else {
    // 进入横屏全屏
    try {
      await videoEl.requestFullscreen();
      // 锁定屏幕方向为横屏
      const orientation = screen.orientation as ScreenOrientation & {
        lock?: (o: string) => Promise<void>;
      };
      if (orientation?.lock) {
        try {
          await orientation.lock("landscape");
        } catch {
          // 部分浏览器不支持 orientation.lock
        }
      }
    } catch (e) {
      console.error("进入全屏失败:", e);
    }
  }
}

function onFullscreenChange() {
  if (!document.fullscreenElement) {
    // 退出全屏时恢复竖屏
    isLandscape.value = false;
    if (screen.orientation) {
      (
        screen.orientation as ScreenOrientation & { unlock?: () => void }
      ).unlock?.();
    }
  } else {
    isLandscape.value = true;
  }
}

// 生命周期：组件挂载时加载
onMounted(() => {
  loadEpisodeData(Number(route.params.id));
  document.addEventListener("fullscreenchange", onFullscreenChange);
});

// 核心：监听路由变化，点击“下一集”时 URL 变了，这里会触发重新请求
watch(
  () => route.params.id,
  (newId) => {
    if (newId) loadEpisodeData(Number(newId));
  },
);

// 组件卸载前销毁播放器，释放内存
onBeforeUnmount(() => {
  if (dp) dp.destroy();
  document.removeEventListener("fullscreenchange", onFullscreenChange);
});
</script>

<template>
  <div
    class="page-container"
    :class="{ 'is-collapsed': isCollapse, 'is-fullscreen': isLandscape }"
  >
    <aside class="sidebar">
      <div class="sidebar-header">
        <span v-show="!isCollapse">附件资料</span>
        <el-icon class="toggle-btn" @click="toggleSidebar">
          <Fold v-if="!isCollapse" />
          <Expand v-else />
        </el-icon>
      </div>

      <div class="attachment-list">
        <template v-if="episode?.attachments?.length">
          <el-tooltip
            v-for="file in episode.attachments"
            :key="file.id"
            :content="file.name"
            placement="right"
            :disabled="!isCollapse"
          >
            <div class="attachment-item" @click="handleViewAttachment(file.id)">
              <el-icon class="file-icon"><Document /></el-icon>
              <span v-show="!isCollapse" class="file-name">{{
                file.name
              }}</span>
            </div>
          </el-tooltip>
        </template>
        <el-empty
          v-else
          v-show="!isCollapse"
          description="暂无附件"
          :image-size="60"
        />
      </div>
    </aside>

    <main class="main-content">
      <!-- 移动端侧边栏开关 -->
      <div class="mobile-attachment-toggle">
        <el-button
          v-if="episode?.attachments?.length"
          size="small"
          @click="isCollapse = !isCollapse"
        >
          <el-icon><Document /></el-icon>
          {{ isCollapse ? "展开附件" : "收起附件" }}
          ({{ episode.attachments.length }})
        </el-button>
      </div>

      <el-breadcrumb
        v-if="episode"
        :separator-icon="ArrowRight"
        class="breadcrumb"
      >
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item
          v-if="episode.group?.course"
          :to="{ path: '/courses/' + episode.group.course_id }"
        >
          {{ episode.group.course.title }}
        </el-breadcrumb-item>
        <el-breadcrumb-item>{{ episode.title }}</el-breadcrumb-item>
      </el-breadcrumb>

      <div class="video-toolbar">
        <el-button
          class="landscape-btn"
          :icon="FullScreen"
          size="small"
          @click="toggleLandscape"
        >
          {{ isLandscape ? "退出横屏" : "横屏播放" }}
        </el-button>
      </div>

      <div ref="videoRef" class="video-player" />

      <div v-if="episode?.navigation" class="video-navigation">
        <el-button
          :disabled="!episode.navigation.prev"
          :icon="ArrowLeft"
          @click="goToEpisode(episode.navigation.prev!.id)"
        >
          上一集：{{ episode.navigation.prev?.title || "没有了" }}
        </el-button>

        <el-button
          :disabled="!episode.navigation.next"
          type="primary"
          @click="goToEpisode(episode.navigation.next!.id)"
        >
          下一集：{{ episode.navigation.next?.title || "没有了" }}
          <el-icon class="el-icon--right"><ArrowRightIcon /></el-icon>
        </el-button>
      </div>
    </main>
  </div>
</template>

<style scoped lang="scss">
/* 此处保留你原有的全部 CSS 样式，并添加以下导航样式 */
.video-navigation {
  display: flex;
  gap: 20px;
  justify-content: space-between;
  width: 100%;
  max-width: 1000px;
  margin-top: 24px;

  .el-button {
    flex: 1;
    height: 48px;
    font-size: 14px;

    /* 核心：处理标题过长的情况 */
    span {
      display: inline-block;
      max-width: 200px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.page-container {
  display: flex;
  min-height: calc(100vh - 100px);
  background-color: var(--el-bg-color-page);
  border-radius: 10px;
  transition: all 0.3s ease;

  .sidebar {
    display: flex;
    flex-direction: column;
    width: 260px;
    background-color: var(--el-bg-color);
    border-right: 1px solid var(--el-border-color-light);
    transition: width 0.3s ease;

    .sidebar-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 20px;
      font-weight: bold;
      border-bottom: 1px solid rgb(235 238 245);

      .toggle-btn {
        font-size: 20px;
        cursor: pointer;

        &:hover {
          color: rgb(102 177 255);
        }
      }
    }

    .attachment-list {
      flex: 1;
      padding: 10px;
      overflow-y: auto;

      .attachment-item {
        display: flex;
        gap: 10px;
        align-items: center;
        padding: 12px;
        margin-bottom: 8px;
        font-size: 14px;
        color: rgb(96 98 102);
        cursor: pointer;
        border-radius: 6px;

        &:hover {
          color: var(--el-color-primary);
          background-color: rgb(64 158 255 / 10%);
        }

        .file-name {
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }

  &.is-collapsed {
    .sidebar {
      width: 60px;

      .sidebar-header {
        justify-content: center;
        padding: 20px 0;
      }
    }
  }

  .main-content {
    display: flex;
    flex: 1;
    flex-direction: column;
    align-items: center;
    padding: 20px;

    .breadcrumb {
      align-self: flex-start;
      margin-bottom: 20px;
    }

    .video-toolbar {
      display: none; // 默认隐藏，仅移动端显示
      justify-content: flex-end;
      width: 100%;
      max-width: 1000px;
      margin-bottom: 8px;

      .landscape-btn {
        font-size: 13px;
      }
    }

    .mobile-attachment-toggle {
      display: none;
      align-self: stretch;
      margin-bottom: 10px;
    }

    .video-player {
      width: 100%;
      max-width: 1000px;
      aspect-ratio: 16 / 9;
      background-color: rgb(0 0 0);
      border-radius: 12px;
      box-shadow: 0 8px 24px rgb(0 0 0 / 15%);
    }
  }
}

// 移动端适配
@media (width <= 768px) {
  .page-container {
    flex-direction: column;

    .sidebar {
      position: fixed;
      top: 56px;
      right: 0;
      z-index: 50;
      width: 260px;
      height: calc(100vh - 56px);
      transform: translateX(100%);
      transition: transform 0.3s ease;

      &.is-mobile-visible {
        box-shadow: -2px 0 16px rgb(0 0 0 / 15%);
        transform: translateX(0);
      }
    }

    &.is-collapsed {
      // 移动端: 使用 is-collapsed 控制 sidebar 显隐
      .sidebar {
        position: static;
        width: 260px;
        height: auto;
        max-height: 60vh;
        border-right: none;
        border-bottom: 1px solid var(--el-border-color-light);
        border-radius: 0;
        transform: none;

        .sidebar-header {
          padding: 12px 16px;
        }

        .attachment-list {
          max-height: 200px;
        }
      }
    }

    .main-content {
      padding: 10px;

      .video-toolbar {
        display: flex;
      }

      .mobile-attachment-toggle {
        display: block;
      }

      .breadcrumb {
        margin-bottom: 10px;
        font-size: 12px;
      }

      .video-navigation {
        flex-direction: column;
        gap: 10px;
        margin-top: 16px;

        .el-button {
          height: 40px;
          font-size: 13px;

          span {
            max-width: 120px;
          }
        }
      }
    }
  }
}

// ========== 全屏横屏模式 ==========
.page-container.is-fullscreen {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgb(0 0 0);
  border-radius: 0;

  .sidebar,
  .breadcrumb,
  .video-toolbar,
  .video-navigation,
  .mobile-attachment-toggle {
    display: none !important;
  }

  .main-content {
    justify-content: center;
    padding: 0;

    .video-player {
      width: 100%;
      max-width: 100%;
      height: 100vh;
      aspect-ratio: auto;
      border-radius: 0;
      box-shadow: none;
    }
  }
}
</style>

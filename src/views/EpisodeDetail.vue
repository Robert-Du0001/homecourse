<script setup lang="ts">
import { ArrowRight, Document, Fold, Expand } from "@element-plus/icons-vue";
import DPlayer from "dplayer";
import { ElMessageBox } from "element-plus";
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";

import type { CatchData } from "@/lib/js/api";
import type { EpisodesResource } from "@/types/episode";

import { request } from "@/lib/js/api";

const videoRef = ref(null);
const route = useRoute();
const episodeId = route.params.id;
const episode = ref<EpisodesResource>();
/** 控制附件侧边栏展开收起 */
const isCollapse = ref(false);

/**
 * 在新窗口打开附件
 * @param id 附件 ID 或完整路径
 */
function handleViewAttachment(id: number) {
  window.open(`/attachments/${id}`, "_blank");
}

/**
 * 切换附件侧边栏展开收起
 */
function toggleSidebar() {
  isCollapse.value = !isCollapse.value;
}

onMounted(async () => {
  try {
    const episodeRes = await request<EpisodesResource>(
      "GET",
      `/episodes/${episodeId}`,
    );
    episode.value = episodeRes.data;
  } catch (e) {
    const rd = e as CatchData;
    ElMessageBox.alert(rd.msg, "温馨提示", {
      confirmButtonText: "确认",
    });
  }

  new DPlayer({
    container: videoRef.value,
    autoplay: true,
    screenshot: true,
    video: {
      url: `/videos/${episodeId}`,
      type: "auto",
    },
  });
});
</script>

<template>
  <div class="page-container" :class="{ 'is-collapsed': isCollapse }">
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

      <div ref="videoRef" class="video-player" />
    </main>
  </div>
</template>

<style scoped lang="scss">
.page-container {
  display: flex;
  min-height: calc(100vh - 80px);
  background-color: var(--el-bg-color-page);
  border-radius: 10px;
  transition: all 0.3s ease;

  // 侧边栏样式
  .sidebar {
    display: flex;
    flex-direction: column;
    width: 260px;
    overflow: hidden;
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
        transition: background 0.2s;

        &:hover {
          color: var(--el-color-primary);
          background-color: rgb(64 158 255 / 10%); // 10% 透明度蓝
        }

        .file-name {
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }

  // 收起状态
  &.is-collapsed {
    .sidebar {
      width: 60px;

      .sidebar-header {
        justify-content: center;
        padding: 20px 0;
      }
    }
  }

  // 主内容区
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

    .video-player {
      width: 100%;
      max-width: 1000px;
      aspect-ratio: 16 / 9; // 自动保持 16:9 比例
      background-color: rgb(0 0 0);
      border-radius: 12px;
      box-shadow: 0 8px 24px rgb(0 0 0 / 15%);
    }
  }
}
</style>

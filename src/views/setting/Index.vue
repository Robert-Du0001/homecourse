<script setup lang="ts">
import { TransitionPresets, useTransition } from "@vueuse/core";
import { onMounted, ref } from "vue";

import { request } from "@/lib/js/api";

type Stat = {
  total: number;
};

const totalCourse = ref(0);
const totalEpisode = ref(0);
const totalAttachment = ref(0);
const totalCourseTran = useTransition(totalCourse, {
  duration: 500,
  transition: TransitionPresets.easeOutExpo,
});
const totalEpisodeTran = useTransition(totalEpisode, {
  duration: 500,
  transition: TransitionPresets.easeOutExpo,
});
const totalAttachmentTran = useTransition(totalAttachment, {
  duration: 500,
  transition: TransitionPresets.easeOutExpo,
});

/**
 * 获取课程总数
 */
function getTotalCourse() {
  request<Stat>("GET", "/admin/courses/statistic").then(function (res) {
    totalCourse.value = res.data.total;
  });
}

/**
 * 获取剧集总数
 */
function getTotalEpisode() {
  request<Stat>("GET", "/admin/episodes/statistic").then(function (res) {
    totalEpisode.value = res.data.total;
  });
}

/**
 * 获取附件总数
 */
function getTotalAttachment() {
  request<Stat>("GET", "/admin/attachments/statistic").then(function (res) {
    totalAttachment.value = res.data.total;
  });
}

onMounted(function () {
  getTotalCourse();

  getTotalEpisode();

  getTotalAttachment();
});
</script>

<template>
  <el-main>
    <el-row :gutter="40">
      <el-col :xs="24" :sm="12" :md="8" class="mb-4">
        <div class="statistic-card">
          <el-statistic :value="totalCourseTran">
            <template #title>
              <div style="display: inline-flex; align-items: center">
                添加课程数
              </div>
            </template>
          </el-statistic>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8" class="mb-4">
        <div class="statistic-card">
          <el-statistic :value="totalEpisodeTran">
            <template #title>
              <div style="display: inline-flex; align-items: center">
                添加剧集数
              </div>
            </template>
          </el-statistic>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8" class="mb-4">
        <div class="statistic-card">
          <el-statistic :value="totalAttachmentTran">
            <template #title>
              <div style="display: inline-flex; align-items: center">
                添加附件数
              </div>
            </template>
          </el-statistic>
        </div>
      </el-col>
    </el-row>
  </el-main>
</template>

<style scoped lang="scss">
.statistic-card {
  height: 100%;
  padding: 20px;
  background-color: var(--el-bg-color-overlay);
  border-radius: 4px;
}
</style>

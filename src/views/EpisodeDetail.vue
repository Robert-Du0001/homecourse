<script setup lang="ts">
import type { EpisodesResource } from '@/types/episode';
import type { CatchData } from '@/lib/js/api';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { request } from '@/lib/js/api';
import { ArrowRight } from '@element-plus/icons-vue';
import DPlayer from 'dplayer';
import { ElMessageBox } from 'element-plus';

const videoRef = ref(null);
const route = useRoute();
const episodeId = route.params.id;
const episode = ref<EpisodesResource>();

onMounted(async () => {
  try {
    const episodeRes = await request<EpisodesResource>(
      'GET',
      `/episodes/${episodeId}`,
    );
    episode.value = episodeRes.data;
  } catch (e) {
    const rd = e as CatchData;
    ElMessageBox.alert(rd.msg, '温馨提示', {
      confirmButtonText: '确认',
    });
  }

  new DPlayer({
    container: videoRef.value,
    autoplay: true,
    screenshot: true,
    video: {
      url: `/media/${episodeId}`,
      type: 'auto',
    },
  });
});
</script>

<template>
  <div class="content">
    <el-breadcrumb v-if="episode" :separator-icon="ArrowRight">
      <el-breadcrumb-item :to="{ path: '/' }"> 首页 </el-breadcrumb-item>
      <el-breadcrumb-item
        v-if="episode.course"
        :to="{ path: '/courses/' + episode.course_id }"
      >
        {{ episode.course.title }}
      </el-breadcrumb-item>
      <el-breadcrumb-item>{{ episode.title }}</el-breadcrumb-item>
    </el-breadcrumb>

    <div ref="videoRef" class="video" />
  </div>
</template>

<style scoped lang="scss">
.content {
  width: 1200px;
  margin: 0 auto;

  .video {
    width: 800px;
    height: 500px;
    margin: 20px auto;
  }
}
</style>

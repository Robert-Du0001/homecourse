<script setup>
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { request } from '@/lib/js/api';
import DPlayer from 'dplayer';

const videoRef = ref(null);
const route = useRoute();
const episodeId = route.params.id;
const episode = ref({});

onMounted(async () => {
  const episodeRes = await request('get', `/episodes/${episodeId}`);
  episode.value = episodeRes.data;

  const dp = new DPlayer({
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
  <el-breadcrumb :separator-icon="ArrowRight">
    <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
    <el-breadcrumb-item 
      v-if="episode.course" 
      :to="{path: '/courses/'+episode.course_id}"
    >
      {{ episode.course.title }}
    </el-breadcrumb-item>
    <el-breadcrumb-item>{{ episode.title }}</el-breadcrumb-item>
  </el-breadcrumb>

  <div ref="videoRef" class="video"></div>
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
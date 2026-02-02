<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { request } from '@/lib/js/api';
import { ArrowRight } from '@element-plus/icons-vue';

const defaultBgs = ['/img/bg-course-01.png', '/img/bg-course-02.png'];

const route = useRoute();
const courseId = route.params.id;
const course = ref({});
const episodes = ref([]);

onMounted(async function() {
  const courseRes = await request('get', `/courses/${courseId}`);
  course.value = courseRes.data;

  // 根据 ID 请求后端的 episodes 接口
  const episodesRes = await request('get', `/episodes?course_id=${courseId}`);
  episodes.value = episodesRes.data;
});
</script>

<template>
<div class="content">
  <div class="course-panel">
    <el-breadcrumb :separator-icon="ArrowRight">
      <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item>{{ course.title }}</el-breadcrumb-item>
    </el-breadcrumb>

    <div class="course-info">
      <div class="cover">
        <img
          :src="course.cover_path || defaultBgs[course.id % 2]" 
          alt="封面"
        >
      </div>
      <div class="intro">
        <div class="title">{{ course.title }}</div>
        <div class="description">{{ course.description || '暂无简介'}}</div>
      </div>
    </div>
  </div>
  <div class="episodes-panel">
    <div class="catalogue-title">
      <img src="/img/catalogue-title.png" alt="课程目录">
    </div>
    <div v-if="!episodes.length" class="episodes-empty">
      暂无课程内容
    </div>
    <ul v-else class="episodes">
      <li v-for="(e, i) in episodes" :key="i" class="episode-item">
        <el-link type="primary" :href="'/episodes/' + e.id">
          {{ e.title.split('.')[0] }}
        </el-link>
      </li>
    </ul>
  </div>
</div>
<el-backtop :right="100" :bottom="100" />
</template>

<style scoped lang="scss">
.content {
  width: 1200px;
  margin: 0 auto;

  .course-panel {
    background: #ffffff;
    min-height: 368px;
    padding: 16px 20px;
    border-radius: 12px;

    .course-info {
      display: flex;
      gap: 30px;
      margin-top: 30px;

      .cover {
        width: 510px;
        margin-left: 0px;
        position: relative;
        overflow: hidden;
        border-radius: 12px;
        height: 288px;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
      }

      .intro {
        flex: 1;

        .title {
          font-size: 20px;
          font-weight: bold;
        }

        .description {
          margin-top: 10px;
          min-height: 60px;
  
          display: -webkit-box;           /* 必须结合使用的属性，将对象作为弹性伸缩盒子模型显示 */
          -webkit-box-orient: vertical;   /* 必须结合使用的属性，设置或检索伸缩盒对象的子元素的排列方式 */
          -webkit-line-clamp: 11;          /* 用来限制在一个块元素显示的文本的行数 */
          overflow: hidden;               /* 必须结合使用的属性，隐藏超出范围的内容 */
          text-overflow: ellipsis;        /* 文本溢出包含元素时显示省略符号 */
        }
      }
    }
  }

  .episodes-panel {
    margin-top: 30px;
    background: #ffffff;
    padding: 12px 0;
    border-radius: 12px;

    .catalogue-title {
      margin-left: 30px;

      img {
        height: 40px;
      }
    }

    .episode-item {
      padding: 6px;
      font-size: 16px;
    }
  }
}
</style>
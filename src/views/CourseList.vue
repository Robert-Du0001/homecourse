<script setup lang="ts">
import type { CourseResource } from '@/types/course';
import { ref, onMounted } from 'vue'; 
import { useRouter } from 'vue-router';
import { request } from '@/lib/js/api';
import { ElMessage } from 'element-plus';

const router = useRouter();
const defaultBgs = ['/img/bg-course-01.png', '/img/bg-course-02.png'];
const courses = ref<CourseResource[]>([]);
const loading = ref(false);

onMounted(async function() {
  loading.value = true;
  try {
    await loadCourses();
  } finally {
    loading.value = false;
  }
});

async function loadCourses() {
  const { data } = await request<CourseResource[]>('GET', '/courses?category_id=-1');
  courses.value = data;
}

function goToDetail(courseId: number) {
  router.push({
    name: 'CourseDetail',
    params: { id: courseId },
  });
}

/**
 * 扫描文件
 */
async function scan() {
  const { msg } = await request('PUT', '/episodes/scan');
  ElMessage.success(msg);
  
  await loadCourses();
}
</script>

<template>
  <div class="operate-btn">
    <el-button
      type="primary"
      @click="scan"
    >
      扫描文件
    </el-button>
  </div>
  <div
    v-loading="loading"
    class="course-list"
  >
    <el-card
      v-for="(course, i) in courses"
      :key="i"
      class="course-card"
    >
      <template #header>
        <div class="card-header">
          <b style="font-size: 16px;">{{ course.title }}</b>
        </div>
      </template>
    
      <div
        class="content"
        @click="goToDetail(course.id)"
      >
        <img
          :src="course.cover_path || defaultBgs[course.id % 2]" 
          alt="封面"
        >
        <div class="description">
          {{ course.description || '暂无课程描述...' }}
        </div>
      </div>

      <template #footer>
        <span style="font-size: 12px; color: rgb(153 153 153);">
          {{ course.created_at.split(' ')[0] }} </span>
      </template>
    </el-card>
  </div>
</template>

<style scoped lang="scss">
.operate-btn {
  text-align: right;
}

.course-list {
  display: grid;

  // 核心：创建 5 列，每列等分
  grid-template-columns: repeat(5, 1fr); 
  gap: 20px;
  padding: 20px;

  // 响应式：根据屏幕宽度调整列数
  @media (width <= 1400px) { grid-template-columns: repeat(4, 1fr); }

  @media (width <= 1100px) { grid-template-columns: repeat(3, 1fr); }

  @media (width <= 768px)  { grid-template-columns: repeat(2, 1fr); }

  @media (width <= 480px)  { grid-template-columns: repeat(1, 1fr); }

  .course-card {
    display: flex;
    flex-direction: column;
    height: 100%; // 让卡片充满 Grid 单元格高度

    :deep(.el-card__body) {
      display: flex;
      flex: 1; // 让 card body 自动撑开
      flex-direction: column;
    }

    .content {
      display: flex;
      flex: 1; // 关键：让内容区占据所有剩余空间，从而实现对齐
      flex-direction: column;
      cursor: pointer;

      img {
        width: 100%;
        aspect-ratio: 16 / 9; // 强制图片比例统一，防止高度抖动
        object-fit: cover;
      }

      .description {
        flex: 1; // 描述文字部分向下撑开
        margin-top: 12px;
      }
    }
  }
}  
</style>

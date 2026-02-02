<script setup>
import { ref, onMounted } from 'vue'; 
import { useRouter } from 'vue-router';
import { request } from '@/lib/js/api';
import { ElMessage } from 'element-plus';

const router = useRouter();
const defaultBgs = ['/img/bg-course-01.png', '/img/bg-course-02.png'];
const courses = ref([]);
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
  const { data } = await request('get', '/courses?category_id=-1');
  courses.value = data;
}

function goToDetail(courseId) {
  router.push({
    name: 'CourseDetail',
    params: { id: courseId }
  });
}

/**
 * 扫描文件
 */
async function scan() {
  const { msg } = await request('put', '/episodes/scan');
  ElMessage.success(msg);
  
  await loadCourses();
}
</script>

<template>
<div class="operate-btn">
   <el-button type="primary" @click="scan">扫描文件</el-button>
</div>
<div class="course-list" v-loading="loading">
  <el-card class="course-card" v-for="(course, i) in courses" :key="i">
    <template #header>
      <div class="card-header">
        <b style="font-size: 16px;">{{ course.title }}</b>
      </div>
    </template>
    
    <div class="content">
      <img
       :src="course.cover_path || defaultBgs[course.id % 2]" 
       alt="封面"
       @click="goToDetail(course.id)"
      >
      <div class="description">
        {{ course.description || '暂无课程描述...' }}
      </div>
    </div>

    <template #footer>
      <span style="font-size: 12px; color: #999;">
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
  @media (max-width: 1400px) { grid-template-columns: repeat(4, 1fr); }
  @media (max-width: 1100px) { grid-template-columns: repeat(3, 1fr); }
  @media (max-width: 768px)  { grid-template-columns: repeat(2, 1fr); }
  @media (max-width: 480px)  { grid-template-columns: repeat(1, 1fr); }

  .course-card {
    height: 100%; // 让卡片充满 Grid 单元格高度
    display: flex;
    flex-direction: column;

    :deep(.el-card__body) {
      flex: 1; // 让 card body 自动撑开
      display: flex;
      flex-direction: column;
    }

    .content {
      flex: 1; // 关键：让内容区占据所有剩余空间，从而实现对齐
      display: flex;
      flex-direction: column;

      img {
        width: 100%;
        aspect-ratio: 16 / 9; // 强制图片比例统一，防止高度抖动
        object-fit: cover;
        cursor: pointer;
      }

      .description {
        margin-top: 12px;
        flex: 1; // 描述文字部分向下撑开
      }
    }
  }
}  
</style>
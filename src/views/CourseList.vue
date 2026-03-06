<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

import type { CatchData } from '@/lib/js/api';
import type { CategoryResource } from '@/types/category';
import type { CourseResource } from '@/types/course';

import { request } from '@/lib/js/api';
import { getDefaultBgImg } from '@/lib/js/helper';

const router = useRouter();

/** 分类数据 */
const categories = ref<CategoryResource[]>();
/** 课程列表 */
const courses = ref<CourseResource[]>([]);
/** 加载课程列表状态 */
const loading = ref(false);
/** 当前选中的课程分类 */
const activeCategoryId = ref(0);

/**
 * 获取课程分类数据
 */
async function loadCategories() {
  try {
    const { data } = await request<CategoryResource[]>('GET', `/categories`);
    categories.value = data;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 获取课程数据
 */
async function loadCourses(categoryId: number) {
  try {
    loading.value = true;
    const { data } = await request<CourseResource[]>(
      'GET',
      `/courses?category_id=${categoryId}`,
    );
    courses.value = data;
    activeCategoryId.value = categoryId;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  } finally {
    loading.value = false;
  }
}

/**
 * 跳转到课程详情
 * @param courseId 课程ID
 */
function goToDetail(courseId: number) {
  router.push({
    name: 'CourseDetail',
    params: { id: courseId },
  });
}

onMounted(async function () {
  await loadCategories();
  await loadCourses(0);
});
</script>

<template>
  <!-- 分类筛选 -->
  <div class="capsule-container">
    <el-check-tag
      :checked="activeCategoryId === 0"
      class="capsule-tag"
      @change="loadCourses(0)"
    >
      全部
    </el-check-tag>

    <el-check-tag
      v-for="item in categories"
      :key="item.id"
      :checked="activeCategoryId === item.id"
      class="capsule-tag"
      @change="loadCourses(item.id)"
    >
      {{ item.name }}
    </el-check-tag>
  </div>

  <div v-loading="loading" class="course-list">
    <el-card v-for="(course, i) in courses" :key="i" class="course-card">
      <template #header>
        <div class="card-header">
          <b style="font-size: 16px">{{ course.title }}</b>
        </div>
      </template>

      <div class="content" @click="goToDetail(course.id)">
        <img
          :src="course.cover_path || getDefaultBgImg(course.id)"
          alt="封面"
        />
        <div class="description">
          {{ course.description || '暂无课程描述...' }}
        </div>
      </div>

      <template #footer>
        <span style="font-size: 12px; color: rgb(153 153 153)">
          {{ course.created_at.split(' ')[0] }}
        </span>
      </template>
    </el-card>
  </div>
</template>

<style scoped lang="scss">
.operate-btn {
  text-align: right;
}

.capsule-container {
  display: flex;
  flex-wrap: wrap; /* 自动换行 */
  gap: 12px; /* 标签之间的间距 */
  margin: 20px 0;

  .capsule-tag {
    /* 基础样式 */
    padding: 8px 20px;
    font-size: 14px;
    font-weight: 500;
    color: rgb(144 147 153);
    cursor: pointer;

    /* 未选中状态：浅灰色背景 */
    background-color: rgb(244 244 245);
    border: none;
    border-radius: 20px; /* 足够大的圆角形成胶囊状 */
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      color: rgb(64 158 255);
      background-color: rgb(233 233 235);
    }

    /* 选中状态：Element 主色调 */
    &.is-checked {
      color: rgb(255 255 255);
      background-color: rgb(64 158 255);
      box-shadow: 0 4px 12px rgb(64 158 255 / 30%); /* 增加一点发光投影 */
    }
  }
}

.course-list {
  display: grid;

  // 核心：创建 5 列，每列等分
  grid-template-columns: repeat(5, 1fr);
  gap: 20px;
  padding: 20px;

  // 响应式：根据屏幕宽度调整列数
  @media (width <= 1400px) {
    grid-template-columns: repeat(4, 1fr);
  }

  @media (width <= 1100px) {
    grid-template-columns: repeat(3, 1fr);
  }

  @media (width <= 768px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (width <= 480px) {
    grid-template-columns: repeat(1, 1fr);
  }

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

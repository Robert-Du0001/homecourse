<script setup lang="ts">
import { ElMessage } from "element-plus";
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";

import type { CatchData } from "@/lib/js/api";
import type { CategoryResource } from "@/types/category";
import type { CourseResource } from "@/types/course";

import { request } from "@/lib/js/api";
import { getDefaultBgImg } from "@/lib/js/helper";

/** 分类选项类型 */
type Option = {
  label: string;
  value: number;
};

const router = useRouter();

/** 分类数据 */
const categories = ref<CategoryResource[]>();
/** 课程列表 */
const courses = ref<CourseResource[]>([]);
/** 当前选中的课程分类 */
const activeCategoryId = ref(0);

/** 分类选项 */
const segmentedOptions = computed(() => {
  let baseOptions: Option[] = [];
  if (categories.value) {
    baseOptions = categories.value.map((item) => ({
      label: item.name,
      value: item.id,
    }));
  }

  return [{ label: "全部", value: 0 }, ...baseOptions];
});

/**
 * 获取课程分类数据
 */
async function loadCategories() {
  try {
    const { data } = await request<CategoryResource[]>("GET", `/categories`);
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
    const { data } = await request<CourseResource[]>(
      "GET",
      `/courses?category_id=${categoryId}`,
    );
    courses.value = data;
    activeCategoryId.value = categoryId;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 跳转到课程详情
 * @param courseId 课程ID
 */
function goToDetail(courseId: number) {
  router.push({
    name: "CourseDetail",
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
  <div v-if="segmentedOptions.length > 1" class="filter-container">
    <el-segmented
      v-model="activeCategoryId"
      :options="segmentedOptions"
      size="large"
      @change="loadCourses"
    />
  </div>

  <div v-if="courses.length" class="course-list">
    <el-card v-for="(course, i) in courses" :key="i" class="course-card">
      <template #header>
        <div class="card-header">
          <b style="font-size: 16px">{{ course.title }}</b>
        </div>
      </template>

      <div class="content" @click="goToDetail(course.id)">
        <el-image
          class="cover"
          :src="course.cover_path || getDefaultBgImg(course.id)"
          alt="封面"
          fit="contain"
        />
        <div class="description">
          {{ course.description || "暂无课程描述..." }}
        </div>
      </div>

      <template #footer>
        <span style="font-size: 12px; color: rgb(153 153 153)">
          {{ course.created_at.split(" ")[0] }}
        </span>
      </template>
    </el-card>
  </div>
  <el-empty v-else description="暂无课程" />
</template>

<style scoped lang="scss">
.operate-btn {
  text-align: right;
}

.filter-container {
  margin: 20px 0;

  .el-segmented {
    /* 组件整体圆角，让它更圆润，像个胶囊容器 */
    --el-border-radius-base: 24px;

    /* 未选中状态的背景色 */
    --el-segmented-bg-color: rgb(240 242 245);

    /* --- 关键改变：选中的按钮样式 --- */

    /* 1. 选中的滑块背景色（好看的活力橙） */
    --el-segmented-item-selected-bg-color: rgb(255 149 0);

    /* 2. 选中的文字颜色（改为白色，与橙色更搭，比默认黑色好看） */
    --el-segmented-item-selected-color: rgb(255 255 255);

    /* 3. 给选中的胶囊增加一点轻微的发光阴影，增强质感 */
    --el-segmented-item-selected-box-shadow: 0 4px 10px rgb(255 120 45 / 20%);

    :deep(.el-segmented__item + .el-segmented__item) {
      margin-left: 10px;
    }

    :deep(
      .el-segmented__item:not(.is-selected):hover .el-segmented__item-label
    ) {
      color: rgb(255 149 0); /* 亮橙色 */
    }
  }
}

.course-list {
  display: grid;

  // 核心：创建 5 列，每列等分
  grid-template-columns: repeat(5, 1fr);
  gap: 20px;
  height: calc(100vh - 180px);
  padding: 20px;
  overflow-y: auto;

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
    height: 360px; // 让卡片充满 Grid 单元格高度

    :deep(.el-card__body) {
      display: flex;
      flex: 1; // 让 card body 自动撑开
      flex-direction: column;
      overflow: hidden;
    }

    .content {
      display: flex;
      flex: 1; // 关键：让内容区占据所有剩余空间，从而实现对齐
      flex-direction: column;
      cursor: pointer;

      .cover {
        width: 100%;
        aspect-ratio: 16 / 9; // 强制图片比例统一，防止高度抖动
      }

      .description {
        flex: 1; // 描述文字部分向下撑开
        margin-top: 12px;
      }
    }
  }
}

.el-empty {
  height: calc(100vh - 180px);
}
</style>

<script setup lang="ts">
import { Edit, Delete, Plus, Rank, Files } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { object, string } from "yup";

import type { CatchData } from "@/lib/js/api";
import type { CategoryResource } from "@/types/category";
import type { CourseResource } from "@/types/course";
import type { UploadFile } from "element-plus";
import type { ValidationError } from "yup";

import { request } from "@/lib/js/api";
import { getDefaultBgImg } from "@/lib/js/helper";
import { TableSortable } from "@/lib/js/tableSortable";

/** 表格排序实例 */
let sortable: TableSortable | null = null;
/** 分类数据 */
const categories = ref<CategoryResource[]>();
/** 课程数据 */
const courses = ref<CourseResource[]>();
/** 是否显示添加/编辑课程对话框 */
const dialogVisible = ref(false);
/**
 * 对话框模式
 * @default
 * 'add' 添加模式
 * 'edit' 编辑模式
 */
const dialogMode = ref("add");
/** 添加/编辑课程表单 */
const courseForm = ref({
  /** 课程分类ID */
  id: 0,
  /** 课程分类ID */
  category_id: 0,
  /** 课程标题 */
  title: "",
  /** 课程简介 */
  description: "",
  /** 课程封面 */
  cover_file: null as File | null,
  /** 课程封面路径，用于预览 */
  cover_preview: "",
});
/** 课程表单验证规则 */
const courseSchema = object({
  /** 课程标题 */
  title: string()
    .required("请输入课程标题")
    .max(20, "课程标题不能超过20个字符"),
  /** 课程简介 */
  description: string().max(200, "课程简介不能超过200个字符"),
});
/** 路由实例 */
const router = useRouter();

/** 分类查找映射表 { id: name } */
const categoryMap = computed(() => {
  const map: Record<number, string> = {};
  if (categories.value) {
    categories.value.forEach((item) => {
      map[item.id] = item.name;
    });
  }
  return map;
});

/**
 * 扫描文件
 */
async function scan() {
  const { msg } = await request("PUT", "/admin/courses/scan");
  ElMessage.success(msg);

  await loadCourses();
}

/**
 * 获取课程数据
 */
async function loadCourses() {
  try {
    const { data } = await request<CategoryResource[]>("GET", `/categories`);

    categories.value = data;

    // 在最前面加个未分类
    categories.value.unshift({
      id: 0,
      name: "未分类",
      is_default: false,
      sort: 0,
      created_at: "",
    });
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }

  // 获取课程数据
  try {
    const { data } = await request<CourseResource[]>("GET", `/admin/courses`);
    courses.value = data;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 删除课程
 * @param id 删除的课程ID
 */
function delCourse(id: number) {
  ElMessageBox.confirm("此操作将永久删除该课程, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    request("DELETE", `/admin/courses/${id}`)
      .then(async ({ msg }) => {
        ElMessage.success(msg);

        await loadCourses();
      })
      .catch(({ msg }) => {
        ElMessage.error(msg);
      });
  });
}

/**
 * 添加课程
 */
function addCourse() {
  dialogMode.value = "add";
  dialogVisible.value = true;
  courseForm.value = {
    id: 0,
    category_id: categories.value![0]!.id,
    title: "",
    description: "",
    cover_file: null,
    cover_preview: "",
  };
}

/**
 * 设置课程
 */
async function setCourse() {
  try {
    await courseSchema.validate(courseForm.value);
  } catch (e) {
    const { message } = e as ValidationError;
    ElMessage.error(message);
    return;
  }

  try {
    const apiMethod = dialogMode.value === "add" ? "POST" : "PUT";
    const api =
      dialogMode.value === "add"
        ? "/admin/courses"
        : `/admin/courses/${courseForm.value.id}`;

    // 处理成 FormData
    const formData = new FormData();
    if (courseForm.value.cover_file) {
      formData.append("cover_file", courseForm.value.cover_file);
    }
    formData.append("title", courseForm.value.title);
    formData.append("category_id", String(courseForm.value.category_id));
    formData.append("description", courseForm.value.description);

    const { msg } = await request(apiMethod, api, formData);
    ElMessage.success(msg);
    dialogVisible.value = false;
    await loadCourses();
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 编辑课程
 * @param course 课程数据
 */
async function editCourse(course: CourseResource) {
  dialogMode.value = "edit";
  courseForm.value = {
    id: course.id,
    category_id: course.category_id,
    title: course.title,
    description: course.description,
    cover_file: null,
    cover_preview: course.cover_path,
  };
  dialogVisible.value = true;
}

/**
 * 处理封面选择
 */
function handleCoverChange(file: UploadFile) {
  if (!file.raw) {
    ElMessage.error("请选择封面文件!");
    return;
  }

  // 简单的文件类型校验
  const isValid =
    file.raw.type === "image/jpeg" || file.raw.type === "image/png";
  if (!isValid) {
    ElMessage.error("封面只能是 JPG 或 PNG 格式!");
    return;
  }

  courseForm.value.cover_file = file.raw;

  // 生成本地预览图 URL
  if (courseForm.value.cover_preview) {
    URL.revokeObjectURL(courseForm.value.cover_preview); // 释放旧内存
  }
  courseForm.value.cover_preview = URL.createObjectURL(file.raw);
}

/**
 * 定义排序后的逻辑
 * @param newIndex
 * @param oldIndex
 */
async function handleSort(newIndex: number, oldIndex: number) {
  if (!courses.value || courses.value?.length === 0) return;

  // 1. 内存同步
  const targetRow = courses.value.splice(oldIndex, 1)[0];
  courses.value.splice(newIndex, 0, targetRow!);

  // 2. 持久化
  try {
    const ids = courses.value.map((item) => item.id);
    const { msg } = await request("PUT", "/admin/courses/sort", { ids });
    ElMessage.success(msg);
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
    await loadCourses(); // 失败回滚
  }
}

onMounted(function () {
  // 初始化排序
  sortable = new TableSortable(".el-table__body-wrapper tbody", handleSort);
  sortable.init();

  loadCourses();
});
</script>

<template>
  <el-row class="header-panel" justify="end">
    <el-col class="btns" :span="6" justify="end">
      <el-button type="primary" @click="addCourse">添加课程</el-button>
      <el-button type="primary" @click="scan">扫描文件</el-button>
    </el-col>
  </el-row>

  <!-- 添加/编辑课程对话框 -->
  <el-dialog
    v-model="dialogVisible"
    :title="dialogMode === 'add' ? '添加课程' : '编辑课程'"
    width="400"
    :center="true"
    class="course-dialog"
  >
    <el-form :model="courseForm">
      <el-form-item label="课程封面" label-width="80px">
        <el-upload
          class="cover-uploader"
          action="#"
          :auto-upload="false"
          :show-file-list="false"
          :on-change="handleCoverChange"
        >
          <el-image
            v-if="courseForm.cover_preview"
            style="width: 100%; height: 100%"
            :src="courseForm.cover_preview"
            fit="contain"
          />
          <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
          <template #tip>
            <div class="el-upload__tip">建议比例 16:9，支持 jpg/png</div>
          </template>
        </el-upload>
      </el-form-item>

      <el-form-item label="课程分类" label-width="80px">
        <el-select
          v-model="courseForm.category_id"
          placeholder="Select"
          style="width: 240px"
        >
          <el-option
            v-for="item in categories"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="课程标题" required label-width="80px">
        <el-input
          v-model="courseForm.title"
          autocomplete="off"
          placeholder="请输入课程标题"
          maxlength="20"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="课程简介" label-width="80px">
        <el-input
          v-model="courseForm.description"
          type="textarea"
          :rows="4"
          placeholder="请简要介绍课程内容..."
          maxlength="200"
          show-word-limit
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="setCourse"> 确认 </el-button>
      </div>
    </template>
  </el-dialog>

  <el-table
    :data="courses"
    :stripe="true"
    :border="true"
    height="600"
    style="width: 100%"
    row-key="id"
    empty-text="暂无课程"
  >
    <el-table-column width="50" align="center">
      <template #default>
        <div class="drag-handler">
          <el-icon><Rank /></el-icon>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="封面" width="100">
      <template #default="{ row }: { row: CourseResource }">
        <el-image
          style="width: 50px; height: 50px"
          :src="row.cover_path || getDefaultBgImg(row.id)"
          fit="cover"
        />
      </template>
    </el-table-column>
    <el-table-column prop="title" label="课程分类" width="120">
      <template #default="{ row }: { row: CourseResource }">
        <el-tag>{{ categoryMap[row.category_id] || "未分类" }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="title" label="课程标题" width="280" />
    <el-table-column prop="description" label="简介" width="580" />
    <el-table-column prop="created_at" label="创建日期" />
    <el-table-column fixed="right" label="操作" min-width="120">
      <template #default="{ row }: { row: CourseResource }">
        <el-tooltip content="编辑" placement="top">
          <el-button
            type="primary"
            :icon="Edit"
            circle
            @click="editCourse(row)"
          />
        </el-tooltip>
        <el-tooltip content="剧集分组管理" placement="top">
          <el-button
            type="warning"
            :icon="Files"
            circle
            @click="router.push(`/setting/courses/${row.id}/groups`)"
          />
        </el-tooltip>
        <el-tooltip content="删除" placement="top">
          <el-button
            type="danger"
            :icon="Delete"
            circle
            @click="delCourse(row.id)"
          />
        </el-tooltip>
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped lang="scss">
// 拖拽样式
:deep(.el-table__row) {
  .el-table__cell:has(.drag-handler) {
    transition: background-color 0.2s ease;

    &:hover {
      cursor: move;
      background-color: rgb(64 158 255 / 10%) !important;

      .drag-handler {
        color: rgb(64 158 255);
      }
    }
  }
}

.course-dialog {
  .cover-uploader {
    :deep(.el-upload) {
      position: relative;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      aspect-ratio: 16 / 9;
      overflow: hidden;
      cursor: pointer;
      border: 1px dashed var(--el-border-color);
      border-radius: 8px;
      transition: var(--el-transition-duration-fast);

      &:hover {
        border-color: var(--el-color-primary);
      }
    }
  }

  .cover-uploader-icon {
    font-size: 28px;
    color: rgb(140 147 157);
  }
}
</style>

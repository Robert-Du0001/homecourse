<script setup lang="ts">
import { Edit, Delete, Plus } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { onMounted, ref } from 'vue';
import { object, string } from 'yup';

import type { CatchData } from '@/lib/js/api';
import type { CourseResource } from '@/types/course';
import type { UploadFile } from 'element-plus';
import type { ValidationError } from 'yup';

import { request } from '@/lib/js/api';
import { getDefaultBgImg } from '@/lib/js/helper';

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
const dialogMode = ref('add');

/** 添加/编辑课程表单 */
const courseForm = ref({
  /** 课程分类ID */
  id: 0,
  /** 课程标题 */
  title: '',
  /** 课程简介 */
  description: '',
  /** 课程封面 */
  cover_path: '',
});

/** 课程表单验证规则 */
const courseSchema = object({
  /** 课程标题 */
  title: string()
    .required('请输入课程标题')
    .max(20, '课程标题不能超过20个字符'),
  /** 课程简介 */
  description: string().max(200, '课程简介不能超过200个字符'),
});

/**
 * 扫描文件
 */
async function scan() {
  const { msg } = await request('PUT', '/admin/episodes/scan');
  ElMessage.success(msg);

  await loadCourses();
}

/**
 * 获取课程数据
 */
async function loadCourses() {
  try {
    const { data } = await request<{
      courses: CourseResource[];
      total: number;
    }>('GET', `/admin/courses?page=1&limit=10`);
    courses.value = data.courses;
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
  ElMessageBox.confirm('此操作将永久删除该课程, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    request('DELETE', `/admin/courses/${id}`)
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
  dialogMode.value = 'add';
  dialogVisible.value = true;
  courseForm.value = {
    id: 0,
    title: '',
    description: '',
    cover_path: '',
  };
}

/**
 * 设置课程
 */
async function setCourse() {
  let course;
  try {
    course = await courseSchema.validate(courseForm.value);
  } catch (e) {
    const { message } = e as ValidationError;
    ElMessage.error(message);
    return;
  }

  try {
    const apiMethod = dialogMode.value === 'add' ? 'POST' : 'PUT';
    const api =
      dialogMode.value === 'add'
        ? '/admin/courses'
        : `/admin/courses/${courseForm.value.id}`;
    const { msg } = await request(apiMethod, api, course);
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
  dialogMode.value = 'edit';
  courseForm.value = {
    id: course.id,
    title: course.title,
    description: course.description,
    cover_path: course.cover_path,
  };
  dialogVisible.value = true;
}

/**
 * 处理封面选择
 */
function handleCoverChange(file: UploadFile) {
  console.log(file);

  if (!file.raw) {
    ElMessage.error('请选择封面文件!');
    return;
  }

  // 1. 简单的文件类型校验
  const isValid =
    file.raw.type === 'image/jpeg' || file.raw.type === 'image/png';
  if (!isValid) {
    ElMessage.error('封面只能是 JPG 或 PNG 格式!');
    return;
  }

  // 3. 存储文件对象，准备上传
  // courseForm.value.cover_path = file.raw;
}

onMounted(loadCourses);
</script>

<template>
  <el-row class="opration-panel" justify="end">
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
          <img
            v-if="courseForm.cover_path"
            :src="courseForm.cover_path"
            class="cover-img"
          />
          <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
          <template #tip>
            <div class="el-upload__tip">建议比例 16:9，支持 jpg/png</div>
          </template>
        </el-upload>
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
  >
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column label="封面" width="100">
      <template #default="{ row }: { row: CourseResource }">
        <el-image
          style="width: 50px; height: 50px"
          :src="row.cover_path || getDefaultBgImg(row.id)"
          fit="cover"
        />
      </template>
    </el-table-column>
    <el-table-column prop="title" label="课程标题" width="280" />
    <el-table-column prop="description" label="简介" width="580" />
    <el-table-column prop="created_at" label="创建日期" />
    <el-table-column fixed="right" label="操作" min-width="120">
      <template #default="{ row }: { row: CourseResource }">
        <el-button
          type="primary"
          :icon="Edit"
          circle
          @click="editCourse(row)"
        />
        <el-button
          type="danger"
          :icon="Delete"
          circle
          @click="delCourse(row.id)"
        ></el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped lang="scss">
.opration-panel {
  margin-bottom: 10px;

  .btns {
    text-align: right;
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

  .cover-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .cover-uploader-icon {
    font-size: 28px;
    color: rgb(140 147 157);
  }
}
</style>

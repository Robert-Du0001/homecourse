<script setup lang="ts">
import type { CourseResource } from '@/types/course';
import type { CatchData } from '@/lib/js/api';
import { request } from '@/lib/js/api';
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';

const courses = ref<CourseResource[]>();

/**
 * 获取课程数据
 */
async function loadCourses() {
  try {
    const { data } = await request<{
      courses: CourseResource[];
      total: number;
    }>('GET', `/courses?page=1&limit=10`);
    courses.value = data.courses;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 删除课程
 * @param id 删除课程的ID
 */
function delCourse(id: number) {
  request('DELETE', `/users/${id}`)
    .then(async ({ msg }) => {
      ElMessage.success(msg);

      await loadCourses();
    })
    .catch(({ msg }) => {
      ElMessage.error(msg);
    });
}

onMounted(loadCourses);
</script>

<template>
  <el-table :data="courses" :stripe="true" :border="true" style="width: 100%">
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column prop="title" label="课程标题" width="180" />
    <el-table-column prop="description" label="简介" width="180" />
    <!-- <el-table-column prop="role" label="简介" width="180">
      <template #default="{ row }: { row: UserResource }">
        {{ row.role === 1 ? '管理员' : '普通用户' }}
      </template>
    </el-table-column> -->
    <el-table-column prop="created_at" label="创建日期" />
    <el-table-column fixed="right" label="操作" min-width="120">
      <template #default="{ row }: { row: CourseResource }">
        <el-button link type="danger" size="small" @click="delCourse(row.id)"
          >删除</el-button
        >
      </template>
    </el-table-column>
  </el-table>
</template>

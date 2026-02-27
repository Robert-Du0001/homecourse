<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { onMounted, ref } from 'vue';

import type { CatchData } from '@/lib/js/api';
import type { CourseResource } from '@/types/course';

import { request } from '@/lib/js/api';
import { getDefaultBgImg } from '@/lib/js/helper';

const courses = ref<CourseResource[]>();

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
  request('DELETE', `/admin/courses/${id}`)
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
  <el-table
    :data="courses"
    :stripe="true"
    :border="true"
    height="600"
    style="width: 100%"
  >
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column label="封面" width="180">
      <template #default="{ row }: { row: CourseResource }">
        <el-image
          style="width: 50px; height: 50px"
          :src="row.cover_path || getDefaultBgImg(row.id)"
          fit="cover"
        />
      </template>
    </el-table-column>
    <el-table-column prop="title" label="课程标题" width="180" />
    <el-table-column prop="description" label="简介" width="180" />
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

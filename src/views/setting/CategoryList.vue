<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { onMounted, ref } from 'vue';

import type { CatchData } from '@/lib/js/api';
import type { CategoryResource } from '@/types/category';

import { request } from '@/lib/js/api';

const categories = ref<CategoryResource[]>();

/**
 * 获取课程分类数据
 */
async function loadCategories() {
  try {
    const { data } = await request<{
      categories: CategoryResource[];
      total: number;
    }>('GET', `/admin/categories?page=1&limit=10`);
    categories.value = data.categories;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 删除课程分类
 * @param id 删除的课程分类ID
 */
function delCategory(id: number) {
  request('DELETE', `/admin/categories/${id}`)
    .then(async ({ msg }) => {
      ElMessage.success(msg);

      await loadCategories();
    })
    .catch(({ msg }) => {
      ElMessage.error(msg);
    });
}

/**
 * 设置默认分类
 * @param val 设置值
 */
function setDefault(id: number, val: boolean) {
  if (val === true) {
    request('PUT', `/admin/categories/${id}/default`)
      .then(async ({ msg }) => {
        ElMessage.success(msg);

        await loadCategories();
      })
      .catch(({ msg }) => {
        ElMessage.error(msg);
      });
  }
}

onMounted(loadCategories);
</script>

<template>
  <el-table
    :data="categories"
    :stripe="true"
    :border="true"
    height="600"
    style="width: 100%"
  >
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column prop="name" label="分类名" width="180" />
    <el-table-column label="默认分类" width="180">
      <template #default="{ row }: { row: CategoryResource }">
        <el-switch
          v-model="row.is_default"
          :disabled="row.is_default"
          @change="setDefault(row.id, $event)"
        />
      </template>
    </el-table-column>
    <el-table-column prop="created_at" label="创建日期" />
    <el-table-column fixed="right" label="操作" min-width="120">
      <template #default="{ row }: { row: CategoryResource }">
        <el-button link type="danger" size="small" @click="delCategory(row.id)"
          >删除</el-button
        >
      </template>
    </el-table-column>
  </el-table>
</template>

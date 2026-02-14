<script setup lang="ts">
import type { UserResource } from '@/types/user';
import type { CatchData } from '@/lib/js/api';
import { request } from '@/lib/js/api';
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';

const users = ref<UserResource[]>();

onMounted(async function () {
  try {
    const { data } = await request<{ users: UserResource[]; total: number }>(
      'GET',
      `/users?page=1&limit=10`,
    );
    users.value = data.users;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
});
</script>

<template>
  <el-table :data="users" style="width: 100%">
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column prop="name" label="用户名" width="180" />
    <el-table-column prop="role" label="角色" width="180" />
    <el-table-column prop="created_at" label="注册时间" />
  </el-table>
</template>

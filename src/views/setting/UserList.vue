<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { onMounted, ref } from 'vue';

import type { CatchData } from '@/lib/js/api';
import type { UserResource } from '@/types/user';

import { request } from '@/lib/js/api';

const users = ref<UserResource[]>();

/**
 * 获取用户数据
 */
async function loadUsers() {
  try {
    const { data } = await request<{ users: UserResource[]; total: number }>(
      'GET',
      `/admin/users?page=1&limit=10`,
    );
    users.value = data.users;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 删除用户
 * @param id 删除用户的ID
 */
function delUser(id: number) {
  request('DELETE', `/admin/users/${id}`)
    .then(async ({ msg }) => {
      ElMessage.success(msg);

      await loadUsers();
    })
    .catch(({ msg }) => {
      ElMessage.error(msg);
    });
}

onMounted(loadUsers);
</script>

<template>
  <el-table :data="users" :stripe="true" :border="true" style="width: 100%">
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column prop="name" label="用户名" width="180" />
    <el-table-column prop="role" label="角色" width="180">
      <template #default="{ row }: { row: UserResource }">
        {{ row.role === 1 ? '管理员' : '普通用户' }}
      </template>
    </el-table-column>
    <el-table-column prop="created_at" label="注册时间" />
    <el-table-column fixed="right" label="操作" min-width="120">
      <template #default="{ row }: { row: UserResource }">
        <el-button
          v-if="row.role !== 1"
          link
          type="danger"
          size="small"
          @click="delUser(row.id)"
          >删除</el-button
        >
      </template>
    </el-table-column>
  </el-table>
</template>

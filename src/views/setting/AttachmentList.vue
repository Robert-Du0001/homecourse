<script setup lang="ts">
import { ArrowRight, Delete, View } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";

import type { AttachmentResource } from "@/types/attachment";
import type { UploadUserFile, UploadFile } from "element-plus";

import { request, type CatchData } from "@/lib/js/api";
import { useUserStore } from "@/store/user";

/** 当前路由对象 */
const route = useRoute();
/** 附件列表 */
const attachments = ref<AttachmentResource[]>([]);
const userStore = useUserStore();
/** 选择的文件列表 */
const fileList = ref<UploadUserFile[]>([]);
/** 文件计数器 */
const fileCounter = ref(0);

/**
 * 获取附件列表
 */
async function loadAttachments() {
  try {
    const { data } = await request<AttachmentResource[]>(
      "GET",
      `/episodes/${route.params.episode_id}/attachments`,
    );
    attachments.value = data;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 上传失败处理
 * 这里不用计数，失败后会自动删除fileList.value的值
 */
async function handleError(err: Error, uploadFile: UploadFile) {
  ElMessage.error(`文件 ${uploadFile.name} 上传失败：${err.message}`);
}

watch(fileCounter, (val) => {
  // 保险起见这里使用大于等于，因为失败后会自动删除fileList中的元素
  if (val && val >= fileList.value.length) {
    loadAttachments();
    // 清空文件列表
    fileList.value = [];
    fileCounter.value = 0;
  }
});

/**
 * 删除附件
 * @param {number} id 附件ID
 */
async function delAttachment(id: number) {
  ElMessageBox.confirm("此操作将永久删除文件, 是否继续?", "error", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    request("DELETE", `/admin/attachments/${id}`)
      .then(async ({ msg }) => {
        ElMessage.success(msg);

        await loadAttachments();
      })
      .catch(({ msg }) => {
        ElMessage.error(msg);
      });
  });
}

/**
 * 在新窗口打开附件
 * @param id 附件 ID 或完整路径
 */
const handleViewAttachment = (id: number) => {
  window.open(`/attachments/${id}`, "_blank");
};

onMounted(function () {
  loadAttachments();
});
</script>

<template>
  <el-row class="header-panel" justify="space-between">
    <el-col :span="10">
      <el-breadcrumb :separator-icon="ArrowRight">
        <el-breadcrumb-item :to="{ path: '/setting/courses' }"
          >课程列表</el-breadcrumb-item
        >
        <el-breadcrumb-item
          :to="{ path: `/setting/courses/${route.params.course_id}/groups` }"
          >剧集分组管理</el-breadcrumb-item
        >
        <el-breadcrumb-item
          :to="{
            path: `/setting/courses/${route.params.course_id}/groups/${route.params.group_id}/episodes`,
          }"
          >剧集管理</el-breadcrumb-item
        >
        <el-breadcrumb-item>附件管理</el-breadcrumb-item>
      </el-breadcrumb>
    </el-col>
    <el-col class="btns" :span="6" justify="end">
      <el-upload
        v-model:file-list="fileList"
        class="upload-demo"
        :action="`/api/admin/episodes/${route.params.episode_id}/attachments`"
        name="attachment_file"
        :show-file-list="false"
        :headers="{ Authorization: `Bearer ${userStore.token}` }"
        :multiple="true"
        :on-success="() => fileCounter++"
        :on-error="handleError"
      >
        <el-button type="primary">上传附件</el-button>
      </el-upload>
    </el-col>
  </el-row>

  <el-table
    :data="attachments"
    :stripe="true"
    :border="true"
    height="600"
    style="width: 100%"
    row-key="id"
    empty-text="暂无附件"
  >
    <el-table-column prop="name" label="文件名" width="460" />
    <el-table-column prop="file_path" label="文件路径" width="520" />
    <el-table-column prop="created_at" label="创建日期" />
    <el-table-column fixed="right" label="操作" width="260">
      <template #default="{ row }: { row: AttachmentResource }">
        <el-tooltip content="查看" placement="top">
          <el-button
            type="primary"
            :icon="View"
            circle
            @click="handleViewAttachment(row.id)"
          />
        </el-tooltip>
        <el-tooltip content="删除" placement="top">
          <el-button
            type="danger"
            :icon="Delete"
            circle
            @click="delAttachment(row.id)"
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
</style>

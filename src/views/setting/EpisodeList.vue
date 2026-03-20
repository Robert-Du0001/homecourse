<script setup lang="ts">
import {
  ArrowRight,
  Delete,
  Edit,
  Folder,
  Rank,
} from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { number, object, string, ValidationError } from "yup";

import type { EpisodesItemResource } from "@/types/episode";

import { request, type CatchData } from "@/lib/js/api";
import { TableSortable } from "@/lib/js/tableSortable";

/** 当前路由对象 */
const route = useRoute();
/** 路由实例 */
const router = useRouter();
/** 剧集分组列表 */
const episodes = ref<EpisodesItemResource[]>([]);
/** 是否显示添加/编辑课程分组对话框 */
const dialogVisible = ref(false);
/**
 * 对话框模式
 * @default
 * 'add' 添加模式
 * 'edit' 编辑模式
 */
const dialogMode = ref("add");
/** 添加/编辑分组表单 */
const episodeForm = ref({
  /** 剧集ID */
  id: 0,
  /** 所属的剧集分组ID */
  group_id: route.params.group_id,
  /** 剧集标题 */
  title: "",
  /** 剧集文件路径 */
  file_path: "",
});
/** 分组表单验证规则 */
const episodeSchema = object({
  /** 分组名 */
  title: string()
    .required("请输入剧集标题")
    .max(20, "剧集标题不能超过20个字符"),
  /** 剧集文件路径 */
  file_path: string().required("请填写剧集文件的路径"),
  group_id: number().required("请选择所属的剧集分组"),
});
/** 表格排序实例 */
let sortable: TableSortable | null = null;

/**
 * 获取课程分组数据
 */
async function loadEpisodes() {
  try {
    const { data } = await request<EpisodesItemResource[]>(
      "GET",
      `/groups/${route.params.group_id}/episodes`,
    );
    episodes.value = data;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 设置课程分类
 */
async function setEpisode() {
  let episode;
  try {
    episode = await episodeSchema.validate(episodeForm.value);
  } catch (e) {
    const { message } = e as ValidationError;
    ElMessage.error(message);
    return;
  }

  try {
    const apiMethod = dialogMode.value === "add" ? "POST" : "PUT";
    const api =
      dialogMode.value === "add"
        ? `/admin/episodes`
        : `/admin/episodes/${episodeForm.value.id}`;
    const { msg } = await request(apiMethod, api, episode);
    ElMessage.success(msg);
    dialogVisible.value = false;
    await loadEpisodes();
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 添加分组
 */
function addEpisode() {
  dialogMode.value = "add";
  episodeForm.value.id = 0;
  episodeForm.value.title = "";
  episodeForm.value.file_path = "";
  dialogVisible.value = true;
}

/**
 * 编辑分组
 * @param {EpisodesItemResource} episode 待编辑的分组
 */
function editEpisode(episode: EpisodesItemResource) {
  dialogMode.value = "edit";
  episodeForm.value.id = episode.id;
  episodeForm.value.title = episode.title;
  episodeForm.value.file_path = episode.file_path;
  dialogVisible.value = true;
}

/**
 * 删除分组
 * @param {number} id 分组ID
 */
async function delEpisode(id: number) {
  ElMessageBox.confirm(
    "此操作将永久删除该剧集，并且会删除关联的附件，但不会删除对应剧集文件, 是否继续?",
    "提示",
    {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    },
  ).then(() => {
    request("DELETE", `/admin/episodes/${id}`)
      .then(async ({ msg }) => {
        ElMessage.success(msg);

        await loadEpisodes();
      })
      .catch(({ msg }) => {
        ElMessage.error(msg);
      });
  });
}

/**
 * 定义排序后的逻辑
 * @param newIndex
 * @param oldIndex
 */
async function handleSort(newIndex: number, oldIndex: number) {
  if (!episodes.value || episodes.value?.length === 0) return;

  // 1. 内存同步
  const targetRow = episodes.value.splice(oldIndex, 1)[0];
  episodes.value.splice(newIndex, 0, targetRow!);

  // 2. 持久化
  try {
    const ids = episodes.value.map((item) => item.id);
    const { msg } = await request("PUT", "/admin/episodes/sort", { ids });
    ElMessage.success(msg);
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
    await loadEpisodes(); // 失败回滚
  }
}

onMounted(function () {
  // 初始化排序
  sortable = new TableSortable(".el-table__body-wrapper tbody", handleSort);
  sortable.init();

  loadEpisodes();
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
        <el-breadcrumb-item>剧集管理</el-breadcrumb-item>
      </el-breadcrumb>
    </el-col>
    <el-col class="btns" :span="6" justify="end">
      <el-button type="primary" @click="addEpisode">添加剧集</el-button>
    </el-col>
  </el-row>

  <!-- 添加/编辑剧集对话框 -->
  <el-dialog
    v-model="dialogVisible"
    :title="dialogMode === 'add' ? '添加剧集' : '编辑剧集'"
    width="400"
    :center="true"
  >
    <el-form :model="episodeForm">
      <el-form-item label="剧集标题" required label-width="80px">
        <el-input
          v-model="episodeForm.title"
          placeholder="请输入标题名"
          autocomplete="off"
          maxlength="20"
          show-word-limit
        />
      </el-form-item>
      <el-form-item label="文件路径" required label-width="80px">
        <el-input
          v-model="episodeForm.file_path"
          placeholder="请输入文件路径"
          autocomplete="off"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="setEpisode"> 确认 </el-button>
      </div>
    </template>
  </el-dialog>

  <el-table
    :data="episodes"
    :stripe="true"
    :border="true"
    height="600"
    style="width: 100%"
    row-key="id"
    empty-text="暂无剧集分组"
  >
    <el-table-column width="50" align="center">
      <template #default>
        <div class="drag-handler">
          <el-icon><Rank /></el-icon>
        </div>
      </template>
    </el-table-column>

    <el-table-column prop="title" label="标题" width="280" />
    <el-table-column prop="file_path" label="文件路径" width="380" />
    <el-table-column prop="created_at" label="创建日期" />
    <el-table-column fixed="right" label="操作" min-width="120">
      <template #default="{ row }: { row: EpisodesItemResource }">
        <el-tooltip content="编辑" placement="top">
          <el-button
            type="primary"
            :icon="Edit"
            circle
            @click="editEpisode(row)"
          />
        </el-tooltip>
        <el-tooltip content="附件管理" placement="top">
          <el-button
            type="warning"
            :icon="Folder"
            circle
            @click="
              router.push(
                `/setting/courses/${route.params.course_id}/groups/${row.group_id}/episodes/${row.id}/attachments`,
              )
            "
          />
        </el-tooltip>
        <el-tooltip content="删除" placement="top">
          <el-button
            type="danger"
            :icon="Delete"
            circle
            @click="delEpisode(row.id)"
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

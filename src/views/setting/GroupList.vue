<script setup lang="ts">
import { ArrowRight, Delete, Edit, Rank } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { object, string, ValidationError } from "yup";

import type { GroupResource } from "@/types/group";

import { request, type CatchData } from "@/lib/js/api";
import { TableSortable } from "@/lib/js/tableSortable";

/** 当前路由对象 */
const route = useRoute();
/** 剧集分组列表 */
const groups = ref<GroupResource[]>([]);
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
const groupForm = ref({
  /** 课程分组ID */
  id: 0,
  /** 所属的课程ID */
  course_id: route.params.id,
  /** 分组名 */
  name: "",
});
/** 分组表单验证规则 */
const groupSchema = object({
  /** 分组名 */
  name: string().required("请输入分组名").max(10, "分组名不能超过10个字符"),
});
/** 表格排序实例 */
let sortable: TableSortable | null = null;

/**
 * 获取课程分组数据
 */
async function loadGroups() {
  try {
    const { data } = await request<GroupResource[]>(
      "GET",
      `/courses/${route.params.id}/groups`,
    );
    groups.value = data;
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 设置课程分类
 */
async function setGroup() {
  let category;
  try {
    category = await groupSchema.validate(groupForm.value);
  } catch (e) {
    const { message } = e as ValidationError;
    ElMessage.error(message);
    return;
  }

  try {
    const apiMethod = dialogMode.value === "add" ? "POST" : "PUT";
    const api =
      dialogMode.value === "add"
        ? `/admin/groups`
        : `/admin/groups/${groupForm.value.id}`;
    const { msg } = await request(apiMethod, api, category);
    ElMessage.success(msg);
    dialogVisible.value = false;
    await loadGroups();
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 添加分组
 */
function addGroup() {
  dialogMode.value = "add";
  groupForm.value.id = 0;
  groupForm.value.name = "";
  dialogVisible.value = true;
}

/**
 * 编辑分组
 * @param {GroupResource} group 待编辑的分组
 */
function editGroup(group: GroupResource) {
  dialogMode.value = "edit";
  groupForm.value.id = group.id;
  groupForm.value.name = group.name;
  dialogVisible.value = true;
}

/**
 * 删除分组
 * @param {number} id 分组ID
 */
async function delGroup(id: number) {
  ElMessageBox.confirm("此操作将永久删除该分组, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    request("DELETE", `/admin/groups/${id}`)
      .then(async ({ msg }) => {
        ElMessage.success(msg);

        await loadGroups();
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
  if (!groups.value || groups.value?.length === 0) return;

  // 1. 内存同步
  const targetRow = groups.value.splice(oldIndex, 1)[0];
  groups.value.splice(newIndex, 0, targetRow!);

  // 2. 持久化
  try {
    const ids = groups.value.map((item) => item.id);
    const { msg } = await request("PUT", "/admin/groups/sort", { ids });
    ElMessage.success(msg);
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
    await loadGroups(); // 失败回滚
  }
}

onMounted(function () {
  // 初始化排序
  sortable = new TableSortable(".el-table__body-wrapper tbody", handleSort);
  sortable.init();

  loadGroups();
});
</script>

<template>
  <el-row class="header-panel" justify="space-between">
    <el-col :span="10">
      <el-breadcrumb :separator-icon="ArrowRight">
        <el-breadcrumb-item :to="{ path: '/setting/courses' }"
          >课程列表</el-breadcrumb-item
        >
        <el-breadcrumb-item>剧集分组管理</el-breadcrumb-item>
      </el-breadcrumb>
    </el-col>
    <el-col class="btns" :span="6" justify="end">
      <el-button type="primary" @click="addGroup">添加剧集分组</el-button>
    </el-col>
  </el-row>

  <!-- 添加/编辑分组对话框 -->
  <el-dialog
    v-model="dialogVisible"
    :title="dialogMode === 'add' ? '添加分组' : '编辑分组'"
    width="400"
    :center="true"
  >
    <el-form :model="groupForm">
      <el-form-item label="分组名" required label-width="80px">
        <el-input
          v-model="groupForm.name"
          placeholder="请输入分组名"
          autocomplete="off"
          maxlength="10"
          show-word-limit
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="setGroup"> 确认 </el-button>
      </div>
    </template>
  </el-dialog>

  <el-table
    :data="groups"
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

    <el-table-column prop="name" label="分组名" width="580" />
    <el-table-column prop="created_at" label="创建日期" />
    <el-table-column fixed="right" label="操作" min-width="120">
      <template #default="{ row }: { row: GroupResource }">
        <el-tooltip content="编辑" placement="top">
          <el-button
            type="primary"
            :icon="Edit"
            circle
            @click="editGroup(row)"
          />
        </el-tooltip>
        <el-tooltip content="删除" placement="top">
          <el-button
            type="danger"
            :icon="Delete"
            circle
            @click="delGroup(row.id)"
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

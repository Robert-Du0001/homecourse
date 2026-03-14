<script setup lang="ts">
import { Edit, Delete, Rank, QuestionFilled } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, ref } from "vue";
import { object, string } from "yup";

import type { CatchData } from "@/lib/js/api";
import type { CategoryResource } from "@/types/category";
import type { ValidationError } from "yup";

import { request } from "@/lib/js/api";
import { TableSortable } from "@/lib/js/tableSortable";

/** 表格排序实例 */
let sortable: TableSortable | null = null;
/** 分类数据 */
const categories = ref<CategoryResource[]>();
/** 是否显示添加/编辑课程分类对话框 */
const dialogVisible = ref(false);
/**
 * 对话框模式
 * @default
 * 'add' 添加模式
 * 'edit' 编辑模式
 */
const dialogMode = ref("add");
/** 添加/编辑分类表单 */
const categoryForm = ref({
  /** 课程分类ID */
  id: 0,
  /** 分类名 */
  name: "",
});

/** 分类表单验证规则 */
const categorySchema = object({
  /** 分类名 */
  name: string().required("请输入分类名").max(10, "分类名不能超过10个字符"),
});

/**
 * 定义排序后的逻辑
 * @param newIndex
 * @param oldIndex
 */
async function handleSort(newIndex: number, oldIndex: number) {
  if (!categories.value || categories.value?.length === 0) return;

  // 1. 内存同步
  const targetRow = categories.value.splice(oldIndex, 1)[0];
  categories.value.splice(newIndex, 0, targetRow!);

  // 2. 持久化
  try {
    const ids = categories.value.map((item) => item.id);
    const { msg } = await request("PUT", "/admin/categories/sort", { ids });
    ElMessage.success(msg);
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
    await loadCategories(); // 失败回滚
  }
}

/**
 * 获取课程分类数据
 */
async function loadCategories() {
  try {
    const { data } = await request<CategoryResource[]>("GET", `/categories`);
    categories.value = data;
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
  ElMessageBox.confirm("此操作将永久删除该分类, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    request("DELETE", `/admin/categories/${id}`)
      .then(async ({ msg }) => {
        ElMessage.success(msg);

        await loadCategories();
      })
      .catch(({ msg }) => {
        ElMessage.error(msg);
      });
  });
}

/**
 * 设置默认分类
 * @param val 设置值
 */
function setDefault(id: number, val: boolean) {
  if (val === true) {
    request("PUT", `/admin/categories/${id}/default`)
      .then(async ({ msg }) => {
        ElMessage.success(msg);

        await loadCategories();
      })
      .catch(({ msg }) => {
        ElMessage.error(msg);
      });
  }
}

/**
 * 设置课程分类
 */
async function setCategory() {
  let category;
  try {
    category = await categorySchema.validate(categoryForm.value);
  } catch (e) {
    const { message } = e as ValidationError;
    ElMessage.error(message);
    return;
  }

  try {
    const apiMethod = dialogMode.value === "add" ? "POST" : "PUT";
    const api =
      dialogMode.value === "add"
        ? "/admin/categories"
        : `/admin/categories/${categoryForm.value.id}`;
    const { msg } = await request(apiMethod, api, category);
    ElMessage.success(msg);
    dialogVisible.value = false;
    await loadCategories();
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
}

/**
 * 添加分类
 */
async function addCategory() {
  dialogMode.value = "add";
  dialogVisible.value = true;
  categoryForm.value = {
    id: 0,
    name: "",
  };
}

/**
 * 编辑分类
 * @param category 分类数据
 */
async function editCategory(category: CategoryResource) {
  dialogMode.value = "edit";
  categoryForm.value = {
    id: category.id,
    name: category.name,
  };
  dialogVisible.value = true;
}

onMounted(() => {
  // 初始化排序
  sortable = new TableSortable(".el-table__body-wrapper tbody", handleSort);
  sortable.init();

  loadCategories();
});
</script>

<template>
  <el-row class="opration-panel" justify="end">
    <el-col class="btns" :span="6" justify="end">
      <el-button type="primary" @click="addCategory">添加分类</el-button>
    </el-col>
  </el-row>

  <!-- 添加/编辑分类对话框 -->
  <el-dialog
    v-model="dialogVisible"
    :title="dialogMode === 'add' ? '添加分类' : '编辑分类'"
    width="400"
    :center="true"
  >
    <el-form :model="categoryForm">
      <el-form-item label="分类名" required label-width="80px">
        <el-input
          v-model="categoryForm.name"
          placeholder="请输入分类名"
          autocomplete="off"
          maxlength="10"
          show-word-limit
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="setCategory"> 确认 </el-button>
      </div>
    </template>
  </el-dialog>

  <el-table
    ref="tableRef"
    :data="categories"
    :stripe="true"
    :border="true"
    height="600"
    style="width: 100%"
    row-key="id"
  >
    <el-table-column width="50" align="center">
      <template #default>
        <div class="drag-handler">
          <el-icon><Rank /></el-icon>
        </div>
      </template>
    </el-table-column>

    <el-table-column prop="name" label="分类名" width="580" />
    <el-table-column width="180">
      <template #header>
        默认分类
        <el-tooltip content="扫描课程时，将自动添加到此分类中" placement="top">
          <el-icon><QuestionFilled /></el-icon>
        </el-tooltip>
      </template>
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
        <el-tooltip content="编辑" placement="top">
          <el-button
            type="primary"
            :icon="Edit"
            circle
            @click="editCategory(row)"
          />
        </el-tooltip>
        <el-tooltip content="删除" placement="top">
          <el-button
            type="danger"
            :icon="Delete"
            circle
            @click="delCategory(row.id)"
          />
        </el-tooltip>
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

/* 拖拽时的影子样式保持 */
.sortable-ghost {
  outline: 2px dashed rgb(64 158 255); /* 增加虚线框增强交互感 */
  background-color: rgb(245 247 250 / 60%) !important;
}
</style>

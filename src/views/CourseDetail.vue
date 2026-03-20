<script setup lang="ts">
import { ArrowRight } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { ref, onMounted, computed } from "vue";
import { useRoute } from "vue-router";

import type { CatchData } from "@/lib/js/api";
import type { CourseResource } from "@/types/course";
import type { EpisodesItemResource } from "@/types/episode";
import type { GroupResource } from "@/types/group";

import { request } from "@/lib/js/api";
import { getDefaultBgImg } from "@/lib/js/helper";

/** 分类选项类型 */
type Option = {
  label: string;
  value: number;
};

const route = useRoute();
const courseId = route.params.id;
const course = ref<CourseResource>();
const groups = ref<GroupResource[]>([]);
const episodes = ref<EpisodesItemResource[]>([]);
/** 当前选中的剧集分组 */
const activeGroupId = ref(0);

/** 分组选项 */
const segmentedOptions = computed(() => {
  let baseOptions: Option[] = [];
  if (groups.value) {
    baseOptions = groups.value.map((item) => ({
      label: item.name,
      value: item.id,
    }));
  }

  return baseOptions;
});

/**
 * 获取剧集数据
 */
async function loadEpisodes() {
  const { data: episodesData } = await request<EpisodesItemResource[]>(
    "GET",
    `/groups/${activeGroupId.value}/episodes`,
  );
  episodes.value = episodesData;
}

onMounted(async function () {
  try {
    // 获取课程数据
    const { data: courseData } = await request<CourseResource>(
      "GET",
      `/courses/${courseId}`,
    );
    course.value = courseData;

    // 获取剧集分组数据
    const { data: groupsData } = await request<GroupResource[]>(
      "GET",
      `/courses/${courseId}/groups`,
    );
    groups.value = groupsData;

    activeGroupId.value = groups.value[0]!.id;

    await loadEpisodes();
  } catch (e) {
    const { msg } = e as CatchData;
    ElMessage.error(msg);
  }
});
</script>

<template>
  <div class="content">
    <div v-if="course" class="course-panel">
      <el-breadcrumb :separator-icon="ArrowRight">
        <el-breadcrumb-item :to="{ path: '/' }"> 首页 </el-breadcrumb-item>
        <el-breadcrumb-item>{{ course.title }}</el-breadcrumb-item>
      </el-breadcrumb>

      <div class="course-info">
        <div class="cover">
          <img
            :src="course.cover_path || getDefaultBgImg(course.id)"
            alt="封面"
          />
        </div>
        <div class="intro">
          <div class="title">
            {{ course.title }}
          </div>
          <div class="description">
            {{ course.description || "暂无简介" }}
          </div>
        </div>
      </div>
    </div>
    <div class="episodes-panel">
      <div class="catalogue-title">
        <img
          class="catalogue-title-img"
          src="/img/catalogue-title.png"
          alt="课程目录"
        />
      </div>

      <!-- 分组筛选 -->
      <div v-if="segmentedOptions.length > 1" class="filter-container">
        <el-segmented
          v-model="activeGroupId"
          :options="segmentedOptions"
          @change="loadEpisodes"
        />
      </div>

      <ul v-if="episodes.length" class="episodes">
        <li v-for="(e, i) in episodes" :key="i" class="episode-item">
          <el-link type="primary" :href="'/episodes/' + e.id">
            {{ e.title.split(".")[0] }}
          </el-link>
        </li>
      </ul>
      <el-empty v-else description="暂无课程内容" />
    </div>
  </div>
  <el-backtop :right="100" :bottom="100" />
</template>

<style scoped lang="scss">
.content {
  width: 1200px;
  height: calc(100vh - 100px);
  margin: 0 auto;

  .course-panel {
    min-height: 368px;
    padding: 16px 20px;
    background: rgb(255 255 255);
    border-radius: 12px;

    .course-info {
      display: flex;
      gap: 30px;
      margin-top: 30px;

      .cover {
        position: relative;
        width: 510px;
        height: 288px;
        margin-left: 0;
        overflow: hidden;
        border-radius: 12px;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
      }

      .intro {
        flex: 1;

        .title {
          font-size: 20px;
          font-weight: bold;
        }

        .description {
          display: -webkit-box; /* 必须结合使用的属性，将对象作为弹性伸缩盒子模型显示 */
          min-height: 60px;
          margin-top: 10px;
          overflow: hidden; /* 必须结合使用的属性，隐藏超出范围的内容 */
          text-overflow: ellipsis; /* 文本溢出包含元素时显示省略符号 */
          -webkit-line-clamp: 11; /* 用来限制在一个块元素显示的文本的行数 */
          -webkit-box-orient: vertical; /* 必须结合使用的属性，设置或检索伸缩盒对象的子元素的排列方式 */
        }
      }
    }
  }

  .filter-container {
    margin: 20px 40px;

    .el-segmented {
      /* 组件整体圆角，让它更圆润，像个胶囊容器 */
      --el-border-radius-base: 24px;

      /* 未选中状态的背景色 */
      --el-segmented-bg-color: rgb(240 242 245);

      /* --- 关键改变：选中的按钮样式 --- */

      /* 1. 选中的滑块背景色（好看的活力橙） */
      --el-segmented-item-selected-bg-color: rgb(255 149 0);

      /* 2. 选中的文字颜色（改为白色，与橙色更搭，比默认黑色好看） */
      --el-segmented-item-selected-color: rgb(255 255 255);

      /* 3. 给选中的胶囊增加一点轻微的发光阴影，增强质感 */
      --el-segmented-item-selected-box-shadow: 0 4px 10px rgb(255 120 45 / 20%);

      :deep(.el-segmented__item + .el-segmented__item) {
        margin-left: 10px;
      }

      :deep(
        .el-segmented__item:not(.is-selected):hover .el-segmented__item-label
      ) {
        color: rgb(255 149 0); /* 亮橙色 */
      }
    }
  }

  .episodes-panel {
    padding: 12px 0;
    margin-top: 30px;
    background: rgb(255 255 255);
    border-radius: 12px;

    .catalogue-title {
      margin-left: 30px;

      .catalogue-title-img {
        height: 40px;
      }
    }

    .episode-item {
      padding: 6px;
      font-size: 16px;
    }
  }
}
</style>

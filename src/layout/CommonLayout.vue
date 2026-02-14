<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';

const userStore = useUserStore();
const router = useRouter();

function goToIndex() {
  router.push({ name: 'Index' });
}

function goToSetting() {
  router.push({ name: 'Setting' });
}

function logout() {
  userStore.logout();
}
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-header class="header">
        <el-row>
          <el-col
            class="link-index"
            :xs="10"
            :sm="6"
            :lg="6"
            @click="goToIndex"
          >
            <el-image class="logo" src="/favicon.svg" />
            <span class="logo-txt">家庭学坊</span>
          </el-col>
          <el-col
            :xs="{ span: 6, offset: 7 }"
            :sm="{ span: 6, offset: 12 }"
            :lg="{ span: 6, offset: 12 }"
          >
            <el-row justify="end">
              <el-col :span="6" :lg="{ push: 3 }">
                <el-dropdown>
                  <div class="avatar">
                    <el-avatar :size="50" src="/img/avatar.png" />
                  </div>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="goToSetting">
                        设置
                      </el-dropdown-item>
                      <el-dropdown-item divided @click="logout">
                        登出
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-col>
            </el-row>
          </el-col>
        </el-row>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<style scoped lang="scss">
.header {
  background-color: rgb(255 255 255);
  border-bottom: 1px solid var(--el-border-color);
  box-shadow: 2px 2px 12px 0 rgb(0 0 0 / 10%);

  .link-index {
    cursor: pointer;

    .logo {
      width: 50px;
      height: 50px;
      margin-top: 6px;
    }

    .logo-txt {
      margin-left: 6px;
      font-size: 24px;
      font-weight: bold;
      vertical-align: 15px;
    }
  }

  .avatar {
    margin-top: 5px;
    text-align: right;
    outline: none;
  }
}
</style>

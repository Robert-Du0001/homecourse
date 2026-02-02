<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import { request } from '@/lib/js/api';
import { ElMessage } from 'element-plus';

const router = useRouter();
const userStore = useUserStore();

let labelPosition = ref('right');

onMounted(() => {
  // 相当于Element+规定的xs尺寸
  if (window.innerWidth < 768) {
    labelPosition.value = 'top';
  }else {
    labelPosition.value = 'right';
  }

  window.addEventListener('resize', () => {
    // 相当于Element+规定的xs尺寸
    if (window.innerWidth < 768) {
      labelPosition.value = 'top';
    }else {
      labelPosition.value = 'right';
    }
  });
});

const formRef = ref();
const btnDisabled = ref(false);

const ruleForm = ref({
  name: '',
  password: '',
  remember: false,
});

const rules = ref({
  name: [
    { required: true, message: '账号不能为空', trigger: 'blur' },
    { max: 10, message: '账号不能超过10个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '密码不能为空', trigger: 'blur' },
    { min: 8, max: 20, message: '密码长度需要8位到20位', trigger: 'blur' },
  ],
});

const submitForm = async (formEl) => {
  if (!formEl) return;

  await formEl.validate(async (valid) => {
    if (valid) {
      btnDisabled.value = true;
      try {
        const { msg, data } = await request('post', '/users/token', ruleForm.value);
        ElMessage.success(msg);

        userStore.$patch(data);
        userStore.setToken(data.token);

        router.replace('/');
      }catch(e) {
        console.error(e);
        btnDisabled.value = false;
      }
    }
  });
};
</script>

<template>
  <el-row>
    <el-col
      :xs="16"
      :sm="12"
      class="login-panel"
    >
      <div class="login-title">
        欢迎登录家庭学坊
      </div>
      <el-form
        ref="formRef"
        :model="ruleForm"
        :rules="rules"
        label-width="120px"
        :label-position="labelPosition"
        status-icon
      >
        <el-form-item
          label="账号"
          prop="name"
        >
          <el-input
            v-model="ruleForm.name"
            type="text"
            placeholder="请输入您的账号"
          />
        </el-form-item>
        <el-form-item
          label="密码"
          prop="password"
        >
          <el-input
            v-model="ruleForm.password"
            type="password"
            placeholder="请输入您的密码"
          />
        </el-form-item>
        <!-- <el-form-item
          style="display: flex; align-items: center;"
          label="记住我"
          prop="remember"
        >
          <el-switch
            v-model="ruleForm.remember"
            size="large"
          />
        </el-form-item> -->
        <el-form-item>
          <el-button
            type="primary"
            plain
            :disabled="btnDisabled"
            @click="submitForm(formRef)"
          >
            登录
          </el-button>
          <el-link
            class="to-register"
            type="info"
            href="/register"
          >
            没有账号？点击注册
          </el-link>
        </el-form-item>
      </el-form>
    </el-col>
  </el-row>
</template>

<style scoped lang="scss">
.el-row {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-panel {
  border: 1px solid var(--el-border-color);
  box-shadow: var(--el-box-shadow-light);
  padding: 20px;
  border-radius: 10px;
  background-color: var(--el-bg-color-overlay);

  .login-title {
    font-size: 22px;
    text-align: center;
    font-weight: bold;
    margin-top: 20px;
    margin-bottom: 30px;
  }

  .to-register {
    margin-left: 10px;
  }
}
</style>

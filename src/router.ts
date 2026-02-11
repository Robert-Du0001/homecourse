import type { UserResource } from '@/types/user';
import type { CatchData } from '@/lib/js/api';
import { createRouter, createWebHistory } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/store/user';
import { UserRole } from '@/types/user';
import { request } from '@/lib/js/api';
import { ElMessageBox } from 'element-plus';

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
  },
  {
    path: '/',
    component: () => import('@/layout/CommonLayout.vue'), // 通用布局
    children: [
      {
        path: '',
        name: 'Index',
        component: () => import('@/views/CourseList.vue'),
      },
      {
        path: '/courses/:id',
        name: 'CourseDetail',
        component: () => import('@/views/CourseDetail.vue'),
      },
      {
        path: '/episodes/:id',
        component: () => import('@/views/EpisodeDetail.vue'),
      },
    ],
  },
  {
    path: '/setting',
    component: () => import('@/layout/SettingLayout.vue'), // 后台设置布局
    meta: { requiresAdmin: true },
    children: [
      {
        path: '',
        name: 'Setting',
        component: () => import('@/views/setting/Index.vue'),
      },
    ],
  },
];

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  // 滚动行为，跳转后页面自动滚回顶部
  scrollBehavior() {
    return { top: 0 };
  },
});

// 拦截器
router.beforeEach(async (to) => {
  const userStore = useUserStore();

  const { token, role } = storeToRefs(userStore);

  if (!token.value && to.name !== 'Register' && to.name !== 'Login') {
    return { name: 'Login' };
  } else if (token.value && (to.name === 'Register' || to.name === 'Login')) {
    return { name: 'Index' };
  }

  if (to.name === 'Index') {
    if (!userStore.name) {
      try {
        const resData = await request<UserResource>('GET', '/users');
        userStore.$patch(resData.data);
      } catch (e) {
        const rd = e as CatchData;
        ElMessageBox.alert(rd.msg, '温馨提示', {
          confirmButtonText: '确认',
        });
      }
    }
  }

  const needsAdmin = to.matched.some((record) => record.meta.requiresAdmin);
  if (needsAdmin && role.value !== UserRole.ADMIN) {
    return { name: 'Index' };
  }

  return true;
});

export default router;

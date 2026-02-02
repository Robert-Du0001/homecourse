import { createRouter, createWebHistory } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/store/user';

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
    component: () => import('@/layout/CommonLayout.vue'), // 这里是你的通用布局
    children: [
      { 
        path: '', 
        name: 'Index',
        component: () => import('@/views/CourseList.vue') 
      },
      { 
        path: '/courses/:id', 
        name: 'CourseDetail', 
        component: () => import('@/views/CourseDetail.vue'),
      },
      {
        path: '/episodes/:id',
        component: () => import('@/views/EpisodeDetail.vue'),
      }
    ]
  },
];

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  // 滚动行为，跳转后页面自动滚回顶部
  scrollBehavior() { return { top: 0 } }
});

// 拦截器
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore();

  const { token } = storeToRefs(userStore);

  if (!token.value && (to.name !== 'Register' && to.name !== 'Login')) {
    next({name: 'Login'});
  }else if (token.value && (to.name === 'Register' || to.name === 'Login')) {
    next({name: 'Index'});
  }else {
    next();
  }
});

export default router;
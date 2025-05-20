import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 布局组件
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: DefaultLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/pages/DashboardPage.vue'),
          meta: { title: '仪表盘', requiresAuth: true }
        },
        {
          path: 'errors',
          name: 'errors',
          component: () => import('@/pages/ErrorsPage.vue'),
          meta: { title: '错误监控', requiresAuth: true }
        },
        {
          path: 'errors/:id',
          name: 'error-detail',
          component: () => import('@/pages/ErrorDetailPage.vue'),
          meta: { title: '错误详情', requiresAuth: true, hiddenInSidebar: true }
        },
        {
          path: 'performance',
          name: 'performance',
          component: () => import('@/pages/PerformancePage.vue'),
          meta: { title: '性能监控', requiresAuth: true }
        },
        {
          path: 'behavior',
          name: 'behavior',
          component: () => import('@/pages/BehaviorPage.vue'),
          meta: { title: '用户行为', requiresAuth: true }
        },
        {
          path: 'projects',
          name: 'projects',
          component: () => import('@/pages/ProjectsPage.vue'),
          meta: { title: '项目管理', requiresAuth: true }
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/pages/SettingsPage.vue'),
          meta: { title: '设置', requiresAuth: true }
        }
      ]
    },
    {
      path: '/auth',
      component: AuthLayout,
      meta: { requiresAuth: false },
      children: [
        {
          path: 'login',
          name: 'login',
          component: () => import('@/pages/LoginPage.vue'),
          meta: { title: '登录', requiresAuth: false }
        },
        {
          path: 'register',
          name: 'register',
          component: () => import('@/pages/RegisterPage.vue'),
          meta: { title: '注册', requiresAuth: false }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/pages/NotFoundPage.vue'),
      meta: { title: '页面不存在' }
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  // 设置页面标题
  document.title = `${to.meta.title || '错误监控后台'} - Web Tracing Admin`

  if (requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login' })
  } else if (to.path === '/auth/login' && authStore.isAuthenticated) {
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

export default router

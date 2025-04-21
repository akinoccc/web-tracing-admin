import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/pages/dashboard/index.vue'),
          meta: { title: '仪表盘' }
        },
        {
          path: 'errors',
          name: 'errors',
          component: () => import('@/pages/errors/index.vue'),
          meta: { title: '错误监控' }
        },
        {
          path: 'errors/:id',
          name: 'error-detail',
          component: () => import('@/pages/errors/detail.vue'),
          meta: {
            title: '错误详情',
            hiddenInSidebar: true
          }
        },
        {
          path: 'performance',
          name: 'performance',
          component: () => import('@/pages/performance/index.vue'),
          meta: { title: '性能监控' }
        },
        {
          path: 'requests',
          name: 'requests',
          component: () => import('@/pages/requests/index.vue'),
          meta: { title: '请求监控' }
        },
        {
          path: 'resources',
          name: 'resources',
          component: () => import('@/pages/resources/index.vue'),
          meta: { title: '资源监控' }
        },
        {
          path: 'routes',
          name: 'routes',
          component: () => import('@/pages/routes/index.vue'),
          meta: { title: '路由监控' }
        },
        {
          path: 'events',
          name: 'events',
          component: () => import('@/pages/events/index.vue'),
          meta: { title: '事件监控' }
        },
        {
          path: 'projects',
          name: 'projects',
          component: () => import('@/pages/projects/index.vue'),
          meta: { title: '项目管理' }
        },
        {
          path: 'projects/:id/settings',
          name: 'project-settings',
          component: () => import('@/pages/projects/detail.vue'),
          meta: { title: '项目设置', hiddenInSidebar: true }
        }
      ]
    },
    {
      path: '/auth',
      component: () => import('@/layouts/AuthLayout.vue'),
      meta: { requiresAuth: false },
      children: [
        {
          path: 'login',
          name: 'login',
          component: () => import('@/pages/auth/login.vue'),
          meta: { title: '登录' }
        },
        {
          path: 'register',
          name: 'register',
          component: () => import('@/pages/auth/register.vue'),
          meta: { title: '注册' }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/pages/error/404.vue'),
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
  } else if (!requiresAuth && authStore.isAuthenticated) {
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

export default router

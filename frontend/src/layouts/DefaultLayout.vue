<template>
  <div class="min-h-screen bg-background">
    <header class="sticky top-0 z-40 w-full border-b bg-background">
      <div class="container flex h-16 items-center space-x-4 sm:justify-between sm:space-x-0">
        <div class="flex gap-6 md:gap-10">
          <router-link to="/" class="flex items-center space-x-2">
            <span class="inline-block font-bold">Web Tracing Admin</span>
          </router-link>
          <nav class="flex gap-6">
            <router-link
              v-for="item in navItems"
              :key="item.href"
              :to="item.href"
              class="flex items-center text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
              :class="{ 'text-primary': isActive(item.href) }"
            >
              {{ item.title }}
            </router-link>
          </nav>
        </div>
        <div class="flex flex-1 items-center justify-end space-x-4">
          <div class="flex items-center space-x-1">
            <div class="relative">
              <div class="flex items-center space-x-2">
                <span v-if="projectStore.currentProject" class="text-sm font-medium">
                  {{ projectStore.currentProject.name }}
                </span>
                <Button variant="ghost" size="sm" @click="logout">
                  退出登录
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
    <main class="container py-6">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/project'
import {Button} from '@/components/ui/button'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const projectStore = useProjectStore()

const navItems = ref([
  {
    title: '仪表盘',
    href: '/'
  },
  {
    title: '错误监控',
    href: '/errors'
  },
  {
    title: '性能监控',
    href: '/performance'
  },
  {
    title: '请求监控',
    href: '/requests'
  },
  {
    title: '资源监控',
    href: '/resources'
  },
  {
    title: '路由监控',
    href: '/routes'
  },
  {
    title: '事件监控',
    href: '/events'
  },
  {
    title: '设置',
    href: '/settings'
  }
])

// 判断当前路由是否激活
const isActive = (href: string) => {
  if (href === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(href)
}

// 退出登录
const logout = () => {
  authStore.logout()
}

onMounted(() => {
  // 初始化项目数据
  projectStore.init()
})
</script>

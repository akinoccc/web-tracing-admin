<template>
  <div class="min-h-screen bg-background flex">
    <!-- 侧边栏 -->
    <aside class="w-64 border-r bg-background flex-shrink-0 h-screen sticky top-0">
      <div class="p-4 border-b">
        <router-link to="/" class="flex items-center space-x-2">
          <span class="inline-block font-bold">Web Tracing Admin</span>
        </router-link>
      </div>
      <nav class="flex flex-col p-4 space-y-2">
        <router-link
          v-for="item in navItems"
          :key="item.href"
          :to="item.href"
          class="flex items-center p-2 rounded-md text-sm font-medium text-muted-foreground transition-colors hover:text-primary hover:bg-accent"
          :class="{ 'text-primary bg-accent': isActive(item.href) }"
        >
          {{ item.title }}
        </router-link>
      </nav>
    </aside>
    
    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-col">
      <header class="sticky top-0 z-40 w-full border-b bg-background">
        <div class="flex h-16 items-center justify-end px-4">
          <div class="flex items-center space-x-4">
            <span v-if="projectStore.currentProject" class="text-sm font-medium">
              {{ projectStore.currentProject.name }}
            </span>
            <Button variant="ghost" size="sm" @click="logout">
              退出登录
            </Button>
          </div>
        </div>
      </header>
      <main class="flex-1 p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/project'
import {Button} from '@/components/ui/button'
import router from '@/router'

const route = useRoute()
const authStore = useAuthStore()
const projectStore = useProjectStore()

// 从路由配置中获取导航菜单
const navItems = computed(() => {
  const mainRoute = router.options.routes.find(route => route.path === '/')
  if (!mainRoute || !mainRoute.children) return []
  
  return mainRoute.children
    .filter(route => route.meta?.title && !route.meta?.hiddenInSidebar)
    .map(route => ({
      title: route.meta?.title as string,
      href: route.path === '' ? '/' : `/${route.path}`
    }))
})

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

<template>
  <div class="flex flex-col items-center justify-center min-h-screen">
    <h1 class="text-6xl font-bold">404</h1>
    <p class="text-xl mt-4">页面不存在</p>
    <Button class="mt-8" @click="goHome">返回首页</Button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-[80vh] p-4">
    <h1 class="text-6xl font-bold text-primary">404</h1>
    <p class="text-xl mt-4 mb-8 text-center">页面不存在或已被移除</p>
    <Button @click="router.push('/')">返回首页</Button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'

const router = useRouter()
</script>

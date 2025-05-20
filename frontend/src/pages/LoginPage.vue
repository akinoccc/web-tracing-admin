<template>
  <div class="space-y-6">
    <div class="space-y-2 text-center">
      <h1 class="text-2xl font-semibold tracking-tight">登录</h1>
      <p class="text-sm text-muted-foreground">
        输入您的账号和密码登录系统
      </p>
    </div>
    <div class="space-y-4">
      <div class="space-y-2">
        <label for="username" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">用户名</label>
        <Input
          id="username"
          v-model="username"
          placeholder="请输入用户名"
          :disabled="authStore.loading"
          @keyup.enter="handleLogin"
        />
      </div>
      <div class="space-y-2">
        <label for="password" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">密码</label>
        <Input
          id="password"
          v-model="password"
          type="password"
          placeholder="请输入密码"
          :disabled="authStore.loading"
          @keyup.enter="handleLogin"
        />
      </div>
      <div v-if="authStore.error" class="text-sm text-red-500">
        {{ authStore.error }}
      </div>
      <Button
        class="w-full"
        :disabled="authStore.loading || !username || !password"
        @click="handleLogin"
      >
        {{ authStore.loading ? '登录中...' : '登录' }}
      </Button>
      <div class="text-center text-sm">
        还没有账号？
        <router-link to="/auth/register" class="text-primary hover:underline">
          注册
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { toast } from '@/utils/toast'

const authStore = useAuthStore()
const username = ref('')
const password = ref('')

const handleLogin = async () => {
  if (!username.value || !password.value) {
    toast.error('请输入用户名和密码')
    return
  }

  try {
    await authStore.login(username.value, password.value)
    toast.success('登录成功')
  } catch (error) {
    // 错误已在 store 中处理
  }
}
</script>

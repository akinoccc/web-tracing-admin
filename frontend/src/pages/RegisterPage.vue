<template>
  <div class="space-y-6">
    <div class="space-y-2 text-center">
      <h1 class="text-2xl font-semibold tracking-tight">注册</h1>
      <p class="text-sm text-muted-foreground">
        创建您的账号
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
        />
      </div>
      <div class="space-y-2">
        <label for="email" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">邮箱</label>
        <Input
          id="email"
          v-model="email"
          type="email"
          placeholder="请输入邮箱"
          :disabled="authStore.loading"
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
        />
      </div>
      <div class="space-y-2">
        <label for="confirmPassword" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">确认密码</label>
        <Input
          id="confirmPassword"
          v-model="confirmPassword"
          type="password"
          placeholder="请再次输入密码"
          :disabled="authStore.loading"
        />
      </div>
      <div v-if="authStore.error" class="text-sm text-red-500">
        {{ authStore.error }}
      </div>
      <Button
        class="w-full"
        :disabled="authStore.loading || !isFormValid"
        @click="handleRegister"
      >
        {{ authStore.loading ? '注册中...' : '注册' }}
      </Button>
      <div class="text-center text-sm">
        已有账号？
        <router-link to="/auth/login" class="text-primary hover:underline">
          登录
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { toast } from '@/utils/toast'

const authStore = useAuthStore()
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const isFormValid = computed(() => {
  return (
    username.value.length >= 3 &&
    email.value.includes('@') &&
    password.value.length >= 6 &&
    password.value === confirmPassword.value
  )
})

const handleRegister = async () => {
  if (!isFormValid.value) {
    if (username.value.length < 3) {
      toast.error('用户名至少需要3个字符')
      return
    }
    if (!email.value.includes('@')) {
      toast.error('请输入有效的邮箱地址')
      return
    }
    if (password.value.length < 6) {
      toast.error('密码至少需要6个字符')
      return
    }
    if (password.value !== confirmPassword.value) {
      toast.error('两次输入的密码不一致')
      return
    }
    return
  }

  try {
    await authStore.register(username.value, password.value, email.value)
    toast.success('注册成功，请登录')
  } catch (error) {
    // 错误已在 store 中处理
  }
}
</script>

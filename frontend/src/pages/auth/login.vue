<template>
  <Card>
    <CardHeader>
      <CardTitle>登录</CardTitle>
      <CardDescription>
        输入您的用户名和密码登录系统
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div class="space-y-2">
          <label for="username" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">用户名</label>
          <Input
            id="username"
            v-model="form.username"
            placeholder="请输入用户名"
            required
          />
        </div>
        <div class="space-y-2">
          <label for="password" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">密码</label>
          <Input
            id="password"
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            required
          />
        </div>
        <div v-if="authStore.error" class="text-sm text-destructive">
          {{ authStore.error }}
        </div>
        <Button type="submit" class="w-full" :loading="authStore.loading">
          登录
        </Button>
      </form>
    </CardContent>
    <CardFooter>
      <div class="text-sm text-center w-full">
        还没有账号？
        <router-link to="/auth/register" class="underline text-primary">
          注册
        </router-link>
      </div>
    </CardFooter>
  </Card>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

const authStore = useAuthStore()

const form = reactive({
  username: '',
  password: ''
})

const handleSubmit = async () => {
  await authStore.login(form.username, form.password)
}
</script>

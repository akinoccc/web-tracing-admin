import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthService } from '@/services'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()
  const authService = useAuthService()

  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<any | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)

  // 加载用户信息
  const loadUser = async () => {
    if (!token.value) return

    try {
      // 这里应该有一个获取用户信息的 API
      // 暂时从 localStorage 获取
      const userStr = localStorage.getItem('user')
      if (userStr) {
        user.value = JSON.parse(userStr)
      }
    } catch (err: any) {
      error.value = err.message
    }
  }

  // 登录
  const login = async (username: string, password: string) => {
    loading.value = true
    error.value = null

    try {
      const response = await authService.login({ username, password })
      token.value = response.data.token
      user.value = response.data.user

      // 保存到 localStorage
      localStorage.setItem('token', response.data.token)
      localStorage.setItem('user', JSON.stringify(response.data.user))

      router.push('/')
    } catch (err: any) {
      error.value = err.response?.data?.message || '登录失败'
    } finally {
      loading.value = false
    }
  }

  // 注册
  const register = async (username: string, password: string, email: string) => {
    loading.value = true
    error.value = null

    try {
      await authService.register({ username, password, email })
      router.push('/auth/login')
    } catch (err: any) {
      error.value = err.response?.data?.message || '注册失败'
    } finally {
      loading.value = false
    }
  }

  // 登出
  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/auth/login')
  }

  // 初始化
  loadUser()

  return {
    token,
    user,
    loading,
    error,
    isAuthenticated,
    login,
    register,
    logout
  }
})

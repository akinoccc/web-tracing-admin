import { createApiService } from './api'
import type { 
  PostApiAuthLoginMutationRequest, 
  PostApiAuthLoginMutationResponse,
  PostApiAuthRegisterMutationRequest
} from '@/types'

// 认证服务
export function useAuthService() {
  const api = createApiService()

  return {
    // 登录
    login: (data: PostApiAuthLoginMutationRequest) => {
      return api.post<PostApiAuthLoginMutationResponse>('/api/auth/login', data)
    },

    // 注册
    register: (data: PostApiAuthRegisterMutationRequest) => {
      return api.post('/api/auth/register', data)
    }
  }
}

import { useApi } from '@/composables/useApi'

// 创建基础API服务
export function createApiService() {
  const { get, post, put, delete: del } = useApi()

  return {
    // 通用请求方法
    get,
    post,
    put,
    delete: del
  }
}

import { createApiService } from './api'
import type { 
  GetApiEventsStatsRequestErrorQueryParams,
  GetApiEventsStatsRequestError200
} from '@/types'

// 请求监控服务
export function useRequestService() {
  const api = createApiService()

  return {
    // 获取请求错误统计数据
    getRequestErrorStats: (params: GetApiEventsStatsRequestErrorQueryParams) => {
      return api.get<GetApiEventsStatsRequestError200>('/api/events/stats/request-error', params)
    }
  }
}

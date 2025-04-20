import { createApiService } from './api'
import type { 
  GetApiEventsStatsPerformanceQueryParams,
  GetApiEventsStatsPerformance200
} from '@/types'

// 性能监控服务
export function usePerformanceService() {
  const api = createApiService()

  return {
    // 获取性能统计数据
    getPerformanceStats: (params: GetApiEventsStatsPerformanceQueryParams) => {
      return api.get<GetApiEventsStatsPerformance200>('/api/events/stats/performance', params)
    }
  }
}

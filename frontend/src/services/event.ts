import { createApiService } from './api'
import type { 
  GetApiEventsQueryParams,
  GetApiEvents200,
  GetApiEventsIdPathParams,
  GetApiEventsId200,
  GetApiEventsStatsQueryParams,
  GetApiEventsStats200,
  GetApiEventsStatsBrowserQueryParams,
  GetApiEventsStatsBrowser200,
  GetApiEventsStatsOsQueryParams,
  GetApiEventsStatsOs200,
  GetApiEventsStatsDeviceQueryParams,
  GetApiEventsStatsDevice200,
  GetApiEventsStatsErrorTypeQueryParams,
  GetApiEventsStatsErrorType200
} from '@/types'

// 事件服务
export function useEventService() {
  const api = createApiService()

  return {
    // 获取事件列表
    getEvents: (params: GetApiEventsQueryParams) => {
      return api.get<GetApiEvents200>('/api/events', params)
    },

    // 获取事件详情
    getEventDetail: (id: number) => {
      return api.get<GetApiEventsId200>(`/api/events/${id}`)
    },

    // 获取事件统计数据
    getEventStats: (params: GetApiEventsStatsQueryParams) => {
      return api.get<GetApiEventsStats200>('/api/events/stats', params)
    },

    // 获取浏览器分布统计
    getBrowserStats: (params: GetApiEventsStatsBrowserQueryParams) => {
      return api.get<GetApiEventsStatsBrowser200>('/api/events/stats/browser', params)
    },

    // 获取操作系统分布统计
    getOSStats: (params: GetApiEventsStatsOsQueryParams) => {
      return api.get<GetApiEventsStatsOs200>('/api/events/stats/os', params)
    },

    // 获取设备分布统计
    getDeviceStats: (params: GetApiEventsStatsDeviceQueryParams) => {
      return api.get<GetApiEventsStatsDevice200>('/api/events/stats/device', params)
    },

    // 获取错误类型分布统计
    getErrorTypeStats: (params: GetApiEventsStatsErrorTypeQueryParams) => {
      return api.get<GetApiEventsStatsErrorType200>('/api/events/stats/error-type', params)
    }
  }
}

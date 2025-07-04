/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { ApiErrorResponse } from './api/ErrorResponse.ts'
import type { ServicePerformanceStatsResponse } from './service/PerformanceStatsResponse.ts'

export type GetApiEventsStatsPerformanceQueryParams = {
  /**
   * @description 项目ID
   * @type integer
   */
  projectId: number
  /**
   * @description 开始时间
   * @type string
   */
  startTime: string
  /**
   * @description 结束时间
   * @type string
   */
  endTime: string
}

/**
 * @description 统计数据
 */
export type GetApiEventsStatsPerformance200 = ServicePerformanceStatsResponse

/**
 * @description 请求错误
 */
export type GetApiEventsStatsPerformance400 = ApiErrorResponse

/**
 * @description 未授权
 */
export type GetApiEventsStatsPerformance401 = ApiErrorResponse

/**
 * @description 内部错误
 */
export type GetApiEventsStatsPerformance500 = ApiErrorResponse

export type GetApiEventsStatsPerformanceQueryResponse = GetApiEventsStatsPerformance200

export type GetApiEventsStatsPerformanceQuery = {
  Response: GetApiEventsStatsPerformance200
  QueryParams: GetApiEventsStatsPerformanceQueryParams
  Errors: GetApiEventsStatsPerformance400 | GetApiEventsStatsPerformance401 | GetApiEventsStatsPerformance500
}
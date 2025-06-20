/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { ApiErrorResponse } from './api/ErrorResponse.ts'
import type { ServicePerformanceStatsResponse } from './service/PerformanceStatsResponse.ts'

export type GetApiPerformanceStatsQueryParams = {
  /**
   * @description 项目ID
   * @type integer
   */
  projectId: number
  /**
   * @description 开始时间戳
   * @type integer | undefined
   */
  startTime?: number
  /**
   * @description 结束时间戳
   * @type integer | undefined
   */
  endTime?: number
}

/**
 * @description 性能统计信息
 */
export type GetApiPerformanceStats200 = ServicePerformanceStatsResponse

/**
 * @description 请求错误
 */
export type GetApiPerformanceStats400 = ApiErrorResponse

/**
 * @description 未授权
 */
export type GetApiPerformanceStats401 = ApiErrorResponse

/**
 * @description 内部错误
 */
export type GetApiPerformanceStats500 = ApiErrorResponse

export type GetApiPerformanceStatsQueryResponse = GetApiPerformanceStats200

export type GetApiPerformanceStatsQuery = {
  Response: GetApiPerformanceStats200
  QueryParams: GetApiPerformanceStatsQueryParams
  Errors: GetApiPerformanceStats400 | GetApiPerformanceStats401 | GetApiPerformanceStats500
}
/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { ApiErrorResponse } from './api/ErrorResponse.ts'
import type { ServiceDistributionStatsResponse } from './service/DistributionStatsResponse.ts'

export type GetApiEventsStatsBrowserQueryParams = {
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
export type GetApiEventsStatsBrowser200 = ServiceDistributionStatsResponse

/**
 * @description 请求错误
 */
export type GetApiEventsStatsBrowser400 = ApiErrorResponse

/**
 * @description 未授权
 */
export type GetApiEventsStatsBrowser401 = ApiErrorResponse

/**
 * @description 内部错误
 */
export type GetApiEventsStatsBrowser500 = ApiErrorResponse

export type GetApiEventsStatsBrowserQueryResponse = GetApiEventsStatsBrowser200

export type GetApiEventsStatsBrowserQuery = {
  Response: GetApiEventsStatsBrowser200
  QueryParams: GetApiEventsStatsBrowserQueryParams
  Errors: GetApiEventsStatsBrowser400 | GetApiEventsStatsBrowser401 | GetApiEventsStatsBrowser500
}
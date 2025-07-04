/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { ApiErrorResponse } from './api/ErrorResponse.ts'
import type { ServiceEventListResponse } from './service/EventListResponse.ts'

export type GetApiEventsQueryParams = {
  /**
   * @description 项目ID
   * @type integer
   */
  projectId: number
  /**
   * @description 事件类型
   * @type string | undefined
   */
  eventType?: string
  /**
   * @description 页码
   * @type integer | undefined
   */
  page?: number
  /**
   * @description 每页数量
   * @type integer | undefined
   */
  pageSize?: number
}

/**
 * @description 事件列表
 */
export type GetApiEvents200 = ServiceEventListResponse

/**
 * @description 请求错误
 */
export type GetApiEvents400 = ApiErrorResponse

/**
 * @description 未授权
 */
export type GetApiEvents401 = ApiErrorResponse

/**
 * @description 内部错误
 */
export type GetApiEvents500 = ApiErrorResponse

export type GetApiEventsQueryResponse = GetApiEvents200

export type GetApiEventsQuery = {
  Response: GetApiEvents200
  QueryParams: GetApiEventsQueryParams
  Errors: GetApiEvents400 | GetApiEvents401 | GetApiEvents500
}
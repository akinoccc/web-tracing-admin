/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { ApiErrorResponse } from './api/ErrorResponse.ts'
import type { ModelProject } from './model/Project.ts'

export type GetApiProjectsIdPathParams = {
  /**
   * @description 项目ID
   * @type integer
   */
  id: number
}

/**
 * @description 项目详情
 */
export type GetApiProjectsId200 = ModelProject

/**
 * @description 请求错误
 */
export type GetApiProjectsId400 = ApiErrorResponse

/**
 * @description 未授权
 */
export type GetApiProjectsId401 = ApiErrorResponse

/**
 * @description 项目不存在
 */
export type GetApiProjectsId404 = ApiErrorResponse

/**
 * @description 内部错误
 */
export type GetApiProjectsId500 = ApiErrorResponse

export type GetApiProjectsIdQueryResponse = GetApiProjectsId200

export type GetApiProjectsIdQuery = {
  Response: GetApiProjectsId200
  PathParams: GetApiProjectsIdPathParams
  Errors: GetApiProjectsId400 | GetApiProjectsId401 | GetApiProjectsId404 | GetApiProjectsId500
}
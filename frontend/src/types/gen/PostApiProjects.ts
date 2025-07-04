/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { ApiErrorResponse } from './api/ErrorResponse.ts'
import type { ModelProject } from './model/Project.ts'
import type { ServiceCreateProjectRequest } from './service/CreateProjectRequest.ts'

/**
 * @description 创建成功
 */
export type PostApiProjects200 = ModelProject

/**
 * @description 请求错误
 */
export type PostApiProjects400 = ApiErrorResponse

/**
 * @description 未授权
 */
export type PostApiProjects401 = ApiErrorResponse

/**
 * @description 内部错误
 */
export type PostApiProjects500 = ApiErrorResponse

/**
 * @description 项目信息
 */
export type PostApiProjectsMutationRequest = ServiceCreateProjectRequest

export type PostApiProjectsMutationResponse = PostApiProjects200

export type PostApiProjectsMutation = {
  Response: PostApiProjects200
  Request: PostApiProjectsMutationRequest
  Errors: PostApiProjects400 | PostApiProjects401 | PostApiProjects500
}
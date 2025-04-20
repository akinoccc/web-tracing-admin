import { createApiService } from './api'
import type { 
  PostApiProjectsMutationRequest,
  PutApiProjectsIdMutationRequest,
  ModelProject
} from '@/types'

// 项目服务
export function useProjectService() {
  const api = createApiService()

  return {
    // 获取项目列表
    getProjects: () => {
      return api.get<ModelProject[]>('/api/projects')
    },

    // 获取项目详情
    getProject: (id: number) => {
      return api.get<ModelProject>(`/api/projects/${id}`)
    },

    // 创建项目
    createProject: (data: PostApiProjectsMutationRequest) => {
      return api.post<ModelProject>('/api/projects', data)
    },

    // 更新项目
    updateProject: (id: number, data: PutApiProjectsIdMutationRequest) => {
      return api.put<ModelProject>(`/api/projects/${id}`, data)
    },

    // 删除项目
    deleteProject: (id: number) => {
      return api.delete(`/api/projects/${id}`)
    }
  }
}

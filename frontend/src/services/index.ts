import { createApiService } from './api'
import type { ServiceLoginRequest } from '@/types/gen/service/LoginRequest'
import type { ServiceRegisterRequest } from '@/types/gen/service/RegisterRequest'
import type { ServiceCreateProjectRequest } from '@/types/gen/service/CreateProjectRequest'
import type { ServiceUpdateProjectRequest } from '@/types/gen/service/UpdateProjectRequest'

// 认证服务
export function useAuthService() {
  const api = createApiService()

  return {
    login: (data: ServiceLoginRequest) => api.post('/api/auth/login', data),
    register: (data: ServiceRegisterRequest) => api.post('/api/auth/register', data)
  }
}

// 项目服务
export function useProjectService() {
  const api = createApiService()

  return {
    getProjects: () => api.get('/api/projects'),
    getProject: (id: number) => api.get(`/api/projects/${id}`),
    createProject: (data: ServiceCreateProjectRequest) => api.post('/api/projects', data),
    updateProject: (id: number, data: ServiceUpdateProjectRequest) => api.put(`/api/projects/${id}`, data),
    deleteProject: (id: number) => api.delete(`/api/projects/${id}`)
  }
}

// 错误监控服务
export function useErrorService() {
  const api = createApiService()

  return {
    getErrors: (params: { 
      projectId: number, 
      page?: number, 
      pageSize?: number,
      startTime?: number,
      endTime?: number,
      errorType?: string,
      severity?: string
    }) => api.get('/api/errors', params),
    getErrorDetail: (id: number) => api.get(`/api/errors/${id}`),
    getErrorStats: (params: { 
      projectId: number, 
      startTime?: number, 
      endTime?: number 
    }) => api.get('/api/errors/stats', params)
  }
}

// 性能监控服务
export function usePerformanceService() {
  const api = createApiService()

  return {
    getPerformance: (params: { 
      projectId: number, 
      page?: number, 
      pageSize?: number,
      startTime?: number,
      endTime?: number,
      type?: string
    }) => api.get('/api/performance', params),
    getPerformanceStats: (params: { 
      projectId: number, 
      startTime?: number, 
      endTime?: number 
    }) => api.get('/api/performance/stats', params),
    getResourcePerformance: (params: { 
      projectId: number, 
      page?: number, 
      pageSize?: number,
      startTime?: number,
      endTime?: number,
      resourceType?: string
    }) => api.get('/api/performance/resources', params)
  }
}

// 用户行为服务
export function useBehaviorService() {
  const api = createApiService()

  return {
    getPageViews: (params: { 
      projectId: number, 
      page?: number, 
      pageSize?: number,
      startTime?: number,
      endTime?: number
    }) => api.get('/api/behavior/pv', params),
    getClicks: (params: { 
      projectId: number, 
      page?: number, 
      pageSize?: number,
      startTime?: number,
      endTime?: number
    }) => api.get('/api/behavior/clicks', params),
    getBehaviorStats: (params: { 
      projectId: number, 
      startTime?: number, 
      endTime?: number 
    }) => api.get('/api/behavior/stats', params)
  }
}

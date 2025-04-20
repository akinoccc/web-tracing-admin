import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useProjectService } from '@/services'

export const useProjectStore = defineStore('project', () => {
  const projectService = useProjectService()

  // 状态
  const projects = ref<any[]>([])
  const currentProject = ref<any | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const hasProjects = computed(() => projects.value.length > 0)

  // 获取所有项目
  const fetchProjects = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await projectService.getProjects()
      projects.value = response.data

      // 如果有项目且没有选择当前项目，则选择第一个
      if (projects.value.length > 0 && !currentProject.value) {
        currentProject.value = projects.value[0]
        localStorage.setItem('currentProject', JSON.stringify(currentProject.value))
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取项目列表失败'
    } finally {
      loading.value = false
    }
  }

  // 获取项目详情
  const fetchProject = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await projectService.getProject(id)
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取项目详情失败'
      return null
    } finally {
      loading.value = false
    }
  }

  // 创建项目
  const createProject = async (name: string, description: string) => {
    loading.value = true
    error.value = null

    try {
      const response = await projectService.createProject({ name, description })
      projects.value.push(response.data)
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '创建项目失败'
      return null
    } finally {
      loading.value = false
    }
  }

  // 更新项目
  const updateProject = async (id: number, name: string, description: string) => {
    loading.value = true
    error.value = null

    try {
      const response = await projectService.updateProject(id, { name, description })

      // 更新项目列表
      const index = projects.value.findIndex(p => p.id === id)
      if (index !== -1) {
        projects.value[index] = response.data
      }

      // 如果更新的是当前项目，也更新当前项目
      if (currentProject.value && currentProject.value.id === id) {
        currentProject.value = response.data
        localStorage.setItem('currentProject', JSON.stringify(currentProject.value))
      }

      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '更新项目失败'
      return null
    } finally {
      loading.value = false
    }
  }

  // 删除项目
  const deleteProject = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      await projectService.deleteProject(id)

      // 从项目列表中移除
      projects.value = projects.value.filter(p => p.id !== id)

      // 如果删除的是当前项目，重置当前项目
      if (currentProject.value && currentProject.value.id === id) {
        currentProject.value = projects.value.length > 0 ? projects.value[0] : null
        if (currentProject.value) {
          localStorage.setItem('currentProject', JSON.stringify(currentProject.value))
        } else {
          localStorage.removeItem('currentProject')
        }
      }

      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || '删除项目失败'
      return false
    } finally {
      loading.value = false
    }
  }

  // 设置当前项目
  const setCurrentProject = (project: any) => {
    currentProject.value = project
    localStorage.setItem('currentProject', JSON.stringify(project))
  }

  // 初始化
  const init = () => {
    // 尝试从 localStorage 恢复当前项目
    const savedProject = localStorage.getItem('currentProject')
    if (savedProject) {
      try {
        currentProject.value = JSON.parse(savedProject)
      } catch (e) {
        localStorage.removeItem('currentProject')
      }
    }

    // 获取项目列表
    fetchProjects()
  }

  return {
    projects,
    currentProject,
    loading,
    error,
    hasProjects,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject,
    setCurrentProject,
    init
  }
})

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useErrorService } from '@/services'
import { useProjectStore } from './project'
import type { ServiceErrorListResponse } from '@/types/gen/service/ErrorListResponse'
import type { ServiceErrorDetailResponse } from '@/types/gen/service/ErrorDetailResponse'
import type { ServiceErrorStatsResponse } from '@/types/gen/service/ErrorStatsResponse'

export const useErrorStore = defineStore('error', () => {
  const errorService = useErrorService()
  const projectStore = useProjectStore()

  // 状态
  const errors = ref<ServiceErrorListResponse | null>(null)
  const errorDetail = ref<ServiceErrorDetailResponse | null>(null)
  const errorStats = ref<ServiceErrorStatsResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const filters = ref({
    page: 1,
    pageSize: 10,
    startTime: undefined as number | undefined,
    endTime: undefined as number | undefined,
    errorType: undefined as string | undefined,
    severity: undefined as string | undefined
  })

  // 计算属性
  const hasErrors = computed(() => errors.value?.list && errors.value.list.length > 0)
  const totalErrors = computed(() => errors.value?.total || 0)
  const totalPages = computed(() => Math.ceil((errors.value?.total || 0) / filters.value.pageSize))

  // 获取错误列表
  const fetchErrors = async () => {
    if (!projectStore.currentProject) return

    loading.value = true
    error.value = null

    try {
      const response = await errorService.getErrors({
        projectId: projectStore.currentProject.id,
        page: filters.value.page,
        pageSize: filters.value.pageSize,
        startTime: filters.value.startTime,
        endTime: filters.value.endTime,
        errorType: filters.value.errorType,
        severity: filters.value.severity
      })
      errors.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取错误列表失败'
    } finally {
      loading.value = false
    }
  }

  // 获取错误详情
  const fetchErrorDetail = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await errorService.getErrorDetail(id)
      errorDetail.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取错误详情失败'
    } finally {
      loading.value = false
    }
  }

  // 获取错误统计信息
  const fetchErrorStats = async (startTime?: number, endTime?: number) => {
    if (!projectStore.currentProject) return

    loading.value = true
    error.value = null

    try {
      const response = await errorService.getErrorStats({
        projectId: projectStore.currentProject.id,
        startTime,
        endTime
      })
      errorStats.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取错误统计信息失败'
    } finally {
      loading.value = false
    }
  }

  // 设置过滤条件
  const setFilters = (newFilters: Partial<typeof filters.value>) => {
    filters.value = { ...filters.value, ...newFilters }
    // 如果修改了非分页相关的过滤条件，重置页码
    if (Object.keys(newFilters).some(key => key !== 'page' && key !== 'pageSize')) {
      filters.value.page = 1
    }
    fetchErrors()
  }

  // 重置过滤条件
  const resetFilters = () => {
    filters.value = {
      page: 1,
      pageSize: 10,
      startTime: undefined,
      endTime: undefined,
      errorType: undefined,
      severity: undefined
    }
    fetchErrors()
  }

  return {
    errors,
    errorDetail,
    errorStats,
    loading,
    error,
    filters,
    hasErrors,
    totalErrors,
    totalPages,
    fetchErrors,
    fetchErrorDetail,
    fetchErrorStats,
    setFilters,
    resetFilters
  }
})

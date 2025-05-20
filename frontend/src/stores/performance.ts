import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { usePerformanceService } from '@/services'
import { useProjectStore } from './project'
import type { ServicePerformanceListResponse } from '@/types/gen/service/PerformanceListResponse'
import type { ServicePerformanceStatsResponse } from '@/types/gen/service/PerformanceStatsResponse'
import type { ServiceResourcePerformanceListResponse } from '@/types/gen/service/ResourcePerformanceListResponse'

export const usePerformanceStore = defineStore('performance', () => {
  const performanceService = usePerformanceService()
  const projectStore = useProjectStore()

  // 状态
  const performance = ref<ServicePerformanceListResponse | null>(null)
  const performanceStats = ref<ServicePerformanceStatsResponse | null>(null)
  const resourcePerformance = ref<ServiceResourcePerformanceListResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const filters = ref({
    page: 1,
    pageSize: 10,
    startTime: undefined as number | undefined,
    endTime: undefined as number | undefined,
    type: undefined as string | undefined
  })
  const resourceFilters = ref({
    page: 1,
    pageSize: 10,
    startTime: undefined as number | undefined,
    endTime: undefined as number | undefined,
    resourceType: undefined as string | undefined
  })

  // 计算属性
  const hasPerformanceData = computed(() => performance.value?.list && performance.value.list.length > 0)
  const totalPerformanceItems = computed(() => performance.value?.total || 0)
  const totalPerformancePages = computed(() => Math.ceil((performance.value?.total || 0) / filters.value.pageSize))
  
  const hasResourcePerformanceData = computed(() => resourcePerformance.value?.list && resourcePerformance.value.list.length > 0)
  const totalResourceItems = computed(() => resourcePerformance.value?.total || 0)
  const totalResourcePages = computed(() => Math.ceil((resourcePerformance.value?.total || 0) / resourceFilters.value.pageSize))

  // 获取性能数据
  const fetchPerformance = async () => {
    if (!projectStore.currentProject) return

    loading.value = true
    error.value = null

    try {
      const response = await performanceService.getPerformance({
        projectId: projectStore.currentProject.id,
        page: filters.value.page,
        pageSize: filters.value.pageSize,
        startTime: filters.value.startTime,
        endTime: filters.value.endTime,
        type: filters.value.type
      })
      performance.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取性能数据失败'
    } finally {
      loading.value = false
    }
  }

  // 获取性能统计信息
  const fetchPerformanceStats = async (startTime?: number, endTime?: number) => {
    if (!projectStore.currentProject) return

    loading.value = true
    error.value = null

    try {
      const response = await performanceService.getPerformanceStats({
        projectId: projectStore.currentProject.id,
        startTime,
        endTime
      })
      performanceStats.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取性能统计信息失败'
    } finally {
      loading.value = false
    }
  }

  // 获取资源性能数据
  const fetchResourcePerformance = async () => {
    if (!projectStore.currentProject) return

    loading.value = true
    error.value = null

    try {
      const response = await performanceService.getResourcePerformance({
        projectId: projectStore.currentProject.id,
        page: resourceFilters.value.page,
        pageSize: resourceFilters.value.pageSize,
        startTime: resourceFilters.value.startTime,
        endTime: resourceFilters.value.endTime,
        resourceType: resourceFilters.value.resourceType
      })
      resourcePerformance.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取资源性能数据失败'
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
    fetchPerformance()
  }

  // 设置资源过滤条件
  const setResourceFilters = (newFilters: Partial<typeof resourceFilters.value>) => {
    resourceFilters.value = { ...resourceFilters.value, ...newFilters }
    // 如果修改了非分页相关的过滤条件，重置页码
    if (Object.keys(newFilters).some(key => key !== 'page' && key !== 'pageSize')) {
      resourceFilters.value.page = 1
    }
    fetchResourcePerformance()
  }

  // 重置过滤条件
  const resetFilters = () => {
    filters.value = {
      page: 1,
      pageSize: 10,
      startTime: undefined,
      endTime: undefined,
      type: undefined
    }
    fetchPerformance()
  }

  // 重置资源过滤条件
  const resetResourceFilters = () => {
    resourceFilters.value = {
      page: 1,
      pageSize: 10,
      startTime: undefined,
      endTime: undefined,
      resourceType: undefined
    }
    fetchResourcePerformance()
  }

  return {
    performance,
    performanceStats,
    resourcePerformance,
    loading,
    error,
    filters,
    resourceFilters,
    hasPerformanceData,
    totalPerformanceItems,
    totalPerformancePages,
    hasResourcePerformanceData,
    totalResourceItems,
    totalResourcePages,
    fetchPerformance,
    fetchPerformanceStats,
    fetchResourcePerformance,
    setFilters,
    setResourceFilters,
    resetFilters,
    resetResourceFilters
  }
})

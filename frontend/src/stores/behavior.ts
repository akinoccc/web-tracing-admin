import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useBehaviorService } from '@/services'
import { useProjectStore } from './project'
import type { ServicePVListResponse } from '@/types/gen/service/PVListResponse'
import type { ServiceClickListResponse } from '@/types/gen/service/ClickListResponse'
import type { ServiceBehaviorStatsResponse } from '@/types/gen/service/BehaviorStatsResponse'

export const useBehaviorStore = defineStore('behavior', () => {
  const behaviorService = useBehaviorService()
  const projectStore = useProjectStore()

  // 状态
  const pageViews = ref<ServicePVListResponse | null>(null)
  const clicks = ref<ServiceClickListResponse | null>(null)
  const behaviorStats = ref<ServiceBehaviorStatsResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const pvFilters = ref({
    page: 1,
    pageSize: 10,
    startTime: undefined as number | undefined,
    endTime: undefined as number | undefined
  })
  const clickFilters = ref({
    page: 1,
    pageSize: 10,
    startTime: undefined as number | undefined,
    endTime: undefined as number | undefined
  })

  // 计算属性
  const hasPageViews = computed(() => pageViews.value?.list && pageViews.value.list.length > 0)
  const totalPageViews = computed(() => pageViews.value?.total || 0)
  const totalPVPages = computed(() => Math.ceil((pageViews.value?.total || 0) / pvFilters.value.pageSize))
  
  const hasClicks = computed(() => clicks.value?.list && clicks.value.list.length > 0)
  const totalClicks = computed(() => clicks.value?.total || 0)
  const totalClickPages = computed(() => Math.ceil((clicks.value?.total || 0) / clickFilters.value.pageSize))

  // 获取页面访问数据
  const fetchPageViews = async () => {
    if (!projectStore.currentProject) return

    loading.value = true
    error.value = null

    try {
      const response = await behaviorService.getPageViews({
        projectId: projectStore.currentProject.id,
        page: pvFilters.value.page,
        pageSize: pvFilters.value.pageSize,
        startTime: pvFilters.value.startTime,
        endTime: pvFilters.value.endTime
      })
      pageViews.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取页面访问数据失败'
    } finally {
      loading.value = false
    }
  }

  // 获取点击数据
  const fetchClicks = async () => {
    if (!projectStore.currentProject) return

    loading.value = true
    error.value = null

    try {
      const response = await behaviorService.getClicks({
        projectId: projectStore.currentProject.id,
        page: clickFilters.value.page,
        pageSize: clickFilters.value.pageSize,
        startTime: clickFilters.value.startTime,
        endTime: clickFilters.value.endTime
      })
      clicks.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取点击数据失败'
    } finally {
      loading.value = false
    }
  }

  // 获取行为统计信息
  const fetchBehaviorStats = async (startTime?: number, endTime?: number) => {
    if (!projectStore.currentProject) return

    loading.value = true
    error.value = null

    try {
      const response = await behaviorService.getBehaviorStats({
        projectId: projectStore.currentProject.id,
        startTime,
        endTime
      })
      behaviorStats.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取行为统计信息失败'
    } finally {
      loading.value = false
    }
  }

  // 设置页面访问过滤条件
  const setPVFilters = (newFilters: Partial<typeof pvFilters.value>) => {
    pvFilters.value = { ...pvFilters.value, ...newFilters }
    // 如果修改了非分页相关的过滤条件，重置页码
    if (Object.keys(newFilters).some(key => key !== 'page' && key !== 'pageSize')) {
      pvFilters.value.page = 1
    }
    fetchPageViews()
  }

  // 设置点击过滤条件
  const setClickFilters = (newFilters: Partial<typeof clickFilters.value>) => {
    clickFilters.value = { ...clickFilters.value, ...newFilters }
    // 如果修改了非分页相关的过滤条件，重置页码
    if (Object.keys(newFilters).some(key => key !== 'page' && key !== 'pageSize')) {
      clickFilters.value.page = 1
    }
    fetchClicks()
  }

  // 重置页面访问过滤条件
  const resetPVFilters = () => {
    pvFilters.value = {
      page: 1,
      pageSize: 10,
      startTime: undefined,
      endTime: undefined
    }
    fetchPageViews()
  }

  // 重置点击过滤条件
  const resetClickFilters = () => {
    clickFilters.value = {
      page: 1,
      pageSize: 10,
      startTime: undefined,
      endTime: undefined
    }
    fetchClicks()
  }

  return {
    pageViews,
    clicks,
    behaviorStats,
    loading,
    error,
    pvFilters,
    clickFilters,
    hasPageViews,
    totalPageViews,
    totalPVPages,
    hasClicks,
    totalClicks,
    totalClickPages,
    fetchPageViews,
    fetchClicks,
    fetchBehaviorStats,
    setPVFilters,
    setClickFilters,
    resetPVFilters,
    resetClickFilters
  }
})

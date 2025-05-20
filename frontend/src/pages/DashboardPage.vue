<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold tracking-tight">仪表盘</h1>
      <div class="flex items-center space-x-2">
        <Button variant="outline" size="sm" @click="refreshData">
          刷新
        </Button>
      </div>
    </div>

    <div v-if="!projectStore.currentProject" class="flex flex-col items-center justify-center p-8 border rounded-lg">
      <p class="text-lg text-center text-muted-foreground mb-4">
        请先创建或选择一个项目
      </p>
      <Button @click="router.push('/projects')">
        前往项目管理
      </Button>
    </div>

    <div v-else>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex justify-center p-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>

      <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <!-- 错误统计卡片 -->
        <Card>
          <CardHeader>
            <CardTitle>错误统计</CardTitle>
            <CardDescription>最近错误情况</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">总错误数</span>
                <span class="text-2xl font-bold">{{ errorStats?.stats?.totalErrors || 0 }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">今日错误</span>
                <span class="text-2xl font-bold">{{ errorStats?.stats?.errorsToday || 0 }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">昨日错误</span>
                <span class="text-2xl font-bold">{{ errorStats?.stats?.errorsYesterday || 0 }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">影响用户数</span>
                <span class="text-2xl font-bold">{{ errorStats?.stats?.affectedUsers || 0 }}</span>
              </div>
            </div>
          </CardContent>
          <CardFooter>
            <Button variant="outline" class="w-full" @click="router.push('/errors')">
              查看详情
            </Button>
          </CardFooter>
        </Card>

        <!-- 性能统计卡片 -->
        <Card>
          <CardHeader>
            <CardTitle>性能统计</CardTitle>
            <CardDescription>页面性能指标</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">平均FCP</span>
                <span class="text-2xl font-bold">{{ formatTime(performanceStats?.stats?.avgFCP) }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">平均LCP</span>
                <span class="text-2xl font-bold">{{ formatTime(performanceStats?.stats?.avgLCP) }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">平均TTFB</span>
                <span class="text-2xl font-bold">{{ formatTime(performanceStats?.stats?.avgTTFB) }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">平均加载时间</span>
                <span class="text-2xl font-bold">{{ formatTime(performanceStats?.stats?.avgLoad) }}</span>
              </div>
            </div>
          </CardContent>
          <CardFooter>
            <Button variant="outline" class="w-full" @click="router.push('/performance')">
              查看详情
            </Button>
          </CardFooter>
        </Card>

        <!-- 用户行为卡片 -->
        <Card>
          <CardHeader>
            <CardTitle>用户行为</CardTitle>
            <CardDescription>用户访问情况</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">总PV</span>
                <span class="text-2xl font-bold">{{ behaviorStats?.pvStats?.totalPV || 0 }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">总UV</span>
                <span class="text-2xl font-bold">{{ behaviorStats?.pvStats?.totalUV || 0 }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">今日PV</span>
                <span class="text-2xl font-bold">{{ behaviorStats?.pvStats?.pvToday || 0 }}</span>
              </div>
              <div class="flex flex-col">
                <span class="text-sm text-muted-foreground">今日UV</span>
                <span class="text-2xl font-bold">{{ behaviorStats?.pvStats?.uvToday || 0 }}</span>
              </div>
            </div>
          </CardContent>
          <CardFooter>
            <Button variant="outline" class="w-full" @click="router.push('/behavior')">
              查看详情
            </Button>
          </CardFooter>
        </Card>
      </div>

      <!-- 错误趋势图 -->
      <Card class="mt-6">
        <CardHeader>
          <CardTitle>错误趋势</CardTitle>
          <CardDescription>最近7天错误趋势</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="h-80">
            <!-- 这里可以使用图表库，如 Chart.js 或 ECharts -->
            <div v-if="errorStats?.trend && errorStats.trend.length > 0" class="h-full">
              <!-- 简单的柱状图示例 -->
              <div class="flex h-full items-end space-x-2">
                <div 
                  v-for="item in errorStats.trend" 
                  :key="item.date" 
                  class="flex-1 bg-primary/20 hover:bg-primary/30 rounded-t-md relative group"
                  :style="{ height: `${getBarHeight(item.count)}%` }"
                >
                  <div class="absolute bottom-0 left-0 right-0 text-xs text-center -mb-6 opacity-70">
                    {{ formatDate(item.date) }}
                  </div>
                  <div class="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-background border px-2 py-1 rounded text-xs hidden group-hover:block whitespace-nowrap">
                    {{ item.date }}: {{ item.count }} 个错误
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="h-full flex items-center justify-center">
              <p class="text-muted-foreground">暂无数据</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { useErrorStore } from '@/stores/error'
import { usePerformanceStore } from '@/stores/performance'
import { useBehaviorStore } from '@/stores/behavior'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import type { ServiceErrorStatsResponse } from '@/types/gen/service/ErrorStatsResponse'
import type { ServicePerformanceStatsResponse } from '@/types/gen/service/PerformanceStatsResponse'
import type { ServiceBehaviorStatsResponse } from '@/types/gen/service/BehaviorStatsResponse'

const router = useRouter()
const projectStore = useProjectStore()
const errorStore = useErrorStore()
const performanceStore = usePerformanceStore()
const behaviorStore = useBehaviorStore()

const loading = ref(true)
const errorStats = ref<ServiceErrorStatsResponse | null>(null)
const performanceStats = ref<ServicePerformanceStatsResponse | null>(null)
const behaviorStats = ref<ServiceBehaviorStatsResponse | null>(null)

// 获取数据
const fetchData = async () => {
  if (!projectStore.currentProject) return

  loading.value = true

  try {
    // 获取错误统计
    const errorResponse = await errorStore.fetchErrorStats()
    errorStats.value = errorStore.errorStats

    // 获取性能统计
    const perfResponse = await performanceStore.fetchPerformanceStats()
    performanceStats.value = performanceStore.performanceStats

    // 获取行为统计
    const behaviorResponse = await behaviorStore.fetchBehaviorStats()
    behaviorStats.value = behaviorStore.behaviorStats
  } catch (error) {
    console.error('获取仪表盘数据失败', error)
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 格式化时间（毫秒）
const formatTime = (ms?: number) => {
  if (!ms) return '0ms'
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

// 格式化日期
const formatDate = (dateStr: string) => {
  return dateStr.split('-').slice(1).join('/')
}

// 计算柱状图高度百分比
const getBarHeight = (count: number) => {
  const maxCount = Math.max(...(errorStats.value?.trend?.map(item => item.count) || [0]))
  if (maxCount === 0) return 0
  return (count / maxCount) * 90 // 最高90%，留出空间显示数值
}

onMounted(() => {
  fetchData()
})
</script>

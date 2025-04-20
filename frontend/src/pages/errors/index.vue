<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">错误监控</h1>
      <div class="flex items-center space-x-2">
        <Button variant="outline" @click="refreshData">
          刷新
        </Button>
      </div>
    </div>

    <div v-if="!projectStore.currentProject" class="flex flex-col items-center justify-center py-12">
      <h2 class="text-xl font-semibold mb-4">请先选择一个项目</h2>
      <Button @click="router.push('/')">返回仪表盘</Button>
    </div>

    <div v-else>
      <!-- 错误统计卡片 -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">今日错误</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.todayErrors }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.todayTrend > 0 ? '+' : '' }}{{ stats.todayTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">本周错误</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.weekErrors }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.weekTrend > 0 ? '+' : '' }}{{ stats.weekTrend }}% 相比上周
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">影响用户数</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.affectedUsers }}</div>
            <p class="text-xs text-muted-foreground">
              占总用户的 {{ stats.affectedPercentage }}%
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">解决率</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.resolvedRate }}%</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.resolvedTrend > 0 ? '+' : '' }}{{ stats.resolvedTrend }}% 相比上周
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- 错误趋势图 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>错误趋势</CardTitle>
          <CardDescription>过去 7 天的错误数量趋势</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="h-[300px]">
            <!-- 这里将来放置图表组件 -->
            <div class="flex items-center justify-center h-full text-muted-foreground">
              图表加载中...
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 错误列表 -->
      <Card>
        <CardHeader>
          <CardTitle>错误列表</CardTitle>
          <CardDescription>项目中发生的所有错误</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center space-x-2 mb-4">
              <Input
                v-model="searchQuery"
                placeholder="搜索错误..."
                class="max-w-sm"
              />
              <select
                v-model="filterType"
                class="h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="all">所有类型</option>
                <option value="TypeError">TypeError</option>
                <option value="ReferenceError">ReferenceError</option>
                <option value="SyntaxError">SyntaxError</option>
                <option value="RangeError">RangeError</option>
                <option value="URIError">URIError</option>
                <option value="EvalError">EvalError</option>
                <option value="InternalError">InternalError</option>
              </select>
            </div>

            <div v-if="loading" class="flex justify-center py-8">
              <div class="animate-spin">&#8635;</div>
            </div>

            <div v-else-if="filteredErrors.length === 0" class="text-center py-8 text-muted-foreground">
              暂无错误数据
            </div>

            <div v-else class="border rounded-md">
              <div class="grid grid-cols-12 gap-4 p-4 font-medium border-b">
                <div class="col-span-5">错误信息</div>
                <div class="col-span-2">类型</div>
                <div class="col-span-2">浏览器</div>
                <div class="col-span-2">时间</div>
                <div class="col-span-1">操作</div>
              </div>
              <div
                v-for="error in filteredErrors"
                :key="error.id"
                class="grid grid-cols-12 gap-4 p-4 border-b last:border-0 hover:bg-muted/50"
              >
                <div class="col-span-5 truncate">
                  {{ error.message }}
                </div>
                <div class="col-span-2">{{ error.type }}</div>
                <div class="col-span-2">{{ error.browser }}</div>
                <div class="col-span-2">{{ formatDate(error.time) }}</div>
                <div class="col-span-1">
                  <Button variant="ghost" size="sm" @click="viewErrorDetail(error.id)">
                    详情
                  </Button>
                </div>
              </div>
            </div>

            <!-- 分页 -->
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-foreground">
                显示 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize, totalErrors) }} 条，共 {{ totalErrors }} 条
              </div>
              <div class="flex items-center space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  :disabled="currentPage === 1"
                  @click="currentPage--"
                >
                  上一页
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  :disabled="currentPage * pageSize >= totalErrors"
                  @click="currentPage++"
                >
                  下一页
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { format } from 'date-fns'
import { useProjectStore } from '@/stores/project'
import { useEventService } from '@/services'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const router = useRouter()
const projectStore = useProjectStore()
const eventService = useEventService()

// 统计数据
const stats = reactive({
  todayErrors: 0,
  todayTrend: 0,
  weekErrors: 0,
  weekTrend: 0,
  affectedUsers: 0,
  affectedPercentage: 0,
  resolvedRate: 0,
  resolvedTrend: 0
})

// 错误列表
const errors = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')
const filterType = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const totalErrors = ref(0)

// 过滤后的错误列表
const filteredErrors = computed(() => {
  let result = errors.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(error =>
      error.message.toLowerCase().includes(query) ||
      error.type.toLowerCase().includes(query)
    )
  }

  // 类型过滤
  if (filterType.value !== 'all') {
    result = result.filter(error => error.type === filterType.value)
  }

  return result
})

// 格式化日期
const formatDate = (dateString: string) => {
  try {
    const date = new Date(dateString)
    return format(date, 'yyyy-MM-dd HH:mm')
  } catch (e) {
    return dateString
  }
}

// 查看错误详情
const viewErrorDetail = (id: number) => {
  router.push(`/errors/${id}`)
}

// 刷新数据
const refreshData = () => {
  fetchErrors()
  fetchStats()
}

// 获取错误列表
const fetchErrors = async () => {
  if (!projectStore.currentProject) return

  loading.value = true

  try {
    const response = await eventService.getEvents({
      projectId: projectStore.currentProject.id,
      eventType: 'error',
      page: currentPage.value,
      pageSize: pageSize.value
    })

    errors.value = response.data.items.map((item: any) => {
      const errorEvent = item.errorEvent || {}
      const baseInfo = item.event?.baseInfo || {}

      return {
        id: item.id,
        message: errorEvent.message || '',
        type: errorEvent.type || '',
        browser: `${baseInfo.browser || ''} ${baseInfo.browserVersion || ''}`.trim(),
        time: item.createdAt,
        stack: errorEvent.stack || ''
      }
    })

    totalErrors.value = response.data.total
  } catch (error) {
    console.error('Failed to fetch errors:', error)
  } finally {
    loading.value = false
  }
}

// 获取统计数据
const fetchStats = async () => {
  if (!projectStore.currentProject) return

  try {
    // 获取当前日期和一周前的日期
    const now = new Date()
    const oneWeekAgo = new Date()
    oneWeekAgo.setDate(oneWeekAgo.getDate() - 7)

    const startTime = format(oneWeekAgo, 'yyyy-MM-dd')
    const endTime = format(now, 'yyyy-MM-dd')

    const response = await eventService.getErrorStats({
      projectId: projectStore.currentProject.id,
      startTime,
      endTime
    })

    // 处理返回的统计数据
    const data = response.data

    stats.todayErrors = data.todayCount || 0
    stats.todayTrend = data.todayTrend || 0
    stats.weekErrors = data.weekCount || 0
    stats.weekTrend = data.weekTrend || 0
    stats.affectedUsers = data.affectedUsers || 0
    stats.affectedPercentage = data.affectedPercentage || 0
    stats.resolvedRate = data.resolvedRate || 0
    stats.resolvedTrend = data.resolvedTrend || 0
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

// 监听项目变化
watch(() => projectStore.currentProject, (newProject) => {
  if (newProject) {
    refreshData()
  }
})

// 监听搜索和过滤条件变化
watch([searchQuery, filterType], () => {
  currentPage.value = 1
})

onMounted(() => {
  if (projectStore.currentProject) {
    refreshData()
  }
})
</script>

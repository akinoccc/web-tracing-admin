<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">事件监控</h1>
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
      <!-- 事件统计卡片 -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">今日点击事件</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.todayClicks }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.todayClicksTrend > 0 ? '+' : '' }}{{ stats.todayClicksTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">本周点击事件</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.weekClicks }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.weekClicksTrend > 0 ? '+' : '' }}{{ stats.weekClicksTrend }}% 相比上周
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">今日曝光事件</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.todayExposures }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.todayExposuresTrend > 0 ? '+' : '' }}{{ stats.todayExposuresTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">本周曝光事件</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.weekExposures }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.weekExposuresTrend > 0 ? '+' : '' }}{{ stats.weekExposuresTrend }}% 相比上周
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- 事件趋势图 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>事件趋势</CardTitle>
          <CardDescription>过去 7 天的事件数量趋势</CardDescription>
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

      <!-- 事件列表 -->
      <Card>
        <CardHeader>
          <CardTitle>事件列表</CardTitle>
          <CardDescription>用户交互事件记录</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center space-x-2 mb-4">
              <Input
                v-model="searchQuery"
                placeholder="搜索事件..."
                class="max-w-sm"
              />
              <select
                v-model="filterType"
                class="h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option v-for="eventType in eventTypes" :key="eventType.value" :value="eventType.value">
                  {{ eventType.label }}
                </option>
              </select>
            </div>

            <div v-if="loading" class="flex justify-center py-8">
              <div class="animate-spin">&#8635;</div>
            </div>

            <div v-else-if="filteredEvents.length === 0" class="text-center py-8 text-muted-foreground">
              暂无事件数据
            </div>

            <div v-else>
              <Table>
                <TableCaption>事件列表</TableCaption>
                <TableHeader>
                  <TableRow>
                    <TableHead>ID</TableHead>
                    <TableHead>事件ID</TableHead>
                    <TableHead>事件类型</TableHead>
                    <TableHead>标题</TableHead>
                    <TableHead>页面URL</TableHead>
                    <TableHead>Referer</TableHead>
                    <TableHead>浏览器</TableHead>
                    <TableHead>设备类型</TableHead>
                    <TableHead>操作系统</TableHead>
                    <TableHead>触发时间</TableHead>
                    <TableHead>创建时间</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow 
                    v-for="event in filteredEvents" 
                    :key="event.id" 
                    class="hover:bg-muted/50"
                  >
                    <TableCell>{{ event.id || '-' }}</TableCell>
                    <TableCell>{{ event.eventId || '-' }}</TableCell>
                    <TableCell>{{ eventTypes.find(type => type.value === event.eventType)?.label || '未知' }}</TableCell>
                    <TableCell class="max-w-[150px] truncate" :title="event.title || ''">
                      {{ event.title || '-' }}
                    </TableCell>
                    <TableCell class="max-w-[200px] truncate" :title="event.triggerPageUrl || ''">
                      {{ event.triggerPageUrl || '-' }}
                    </TableCell>
                    <TableCell class="max-w-[150px] truncate" :title="event.referer || ''">
                      {{ event.referer || '-' }}
                    </TableCell>
                    <TableCell>
                      {{ event.baseInfo?.browser || '-' }}
                      {{ event.baseInfo?.browserVersion ? `(${event.baseInfo.browserVersion})` : '' }}
                    </TableCell>
                    <TableCell>{{ event.baseInfo?.deviceType || '-' }}</TableCell>
                    <TableCell>{{ event.baseInfo?.os || '-' }}</TableCell>
                    <TableCell>{{ event.triggerTime ? formatTimestamp(event.triggerTime) : '-' }}</TableCell>
                    <TableCell>{{ formatDate(event.createdAt || '') }}</TableCell>
                  </TableRow>
                </TableBody>
              </Table>

              <!-- 分页 -->
              <div class="flex items-center justify-between mt-4">
                <div class="text-sm text-muted-foreground">
                  显示 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize, totalEvents) }} 条，共 {{ totalEvents }} 条
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
                    :disabled="currentPage * pageSize >= totalEvents"
                    @click="currentPage++"
                  >
                    下一页
                  </Button>
                </div>
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
import { 
  Table, 
  TableBody, 
  TableCaption, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
import { ModelEvent } from '@/types'

const router = useRouter()
const projectStore = useProjectStore()
const eventService = useEventService()

const eventTypes = [
  { label: '所有类型', value: 'all' },
  { label: '点击事件', value: 'click' },
  { label: '曝光事件', value: 'intersection' },
  { label: 'PV事件', value: 'pv' },
  { label: 'PV时长事件', value: 'pv-duration' },
  { label: '错误事件', value: 'error' },
  { label: '性能事件', value: 'performance' },
  { label: '请求错误事件', value: 'request' }
]

// 统计数据
const stats = reactive({
  todayClicks: 0,
  todayClicksTrend: 0,
  weekClicks: 0,
  weekClicksTrend: 0,
  todayExposures: 0,
  todayExposuresTrend: 0,
  weekExposures: 0,
  weekExposuresTrend: 0
})

// 事件列表
const events = ref<ModelEvent[]>([])
const loading = ref(false)
const searchQuery = ref('')
const filterType = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const totalEvents = ref(0)

// 过滤后的事件列表
const filteredEvents = computed(() => {
  let result = events.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(event =>
      (event.triggerPageUrl && event.triggerPageUrl.toLowerCase().includes(query)) ||
      (event.title && event.title.toLowerCase().includes(query)) ||
      (event.eventType && event.eventType.toLowerCase().includes(query)) ||
      (event.baseInfo?.browser && event.baseInfo.browser.toLowerCase().includes(query)) ||
      (event.baseInfo?.os && event.baseInfo.os.toLowerCase().includes(query)) ||
      (event.baseInfo?.deviceType && event.baseInfo.deviceType.toLowerCase().includes(query))
    )
  }

  // 类型过滤
  if (filterType.value !== 'all') {
    result = result.filter(event => event.eventType === filterType.value)
  }

  return result
})

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  try {
    const date = new Date(dateString)
    return format(date, 'yyyy-MM-dd HH:mm:ss')
  } catch (e) {
    return dateString
  }
}

// 格式化时间戳
const formatTimestamp = (timestamp: number) => {
  try {
    const date = new Date(timestamp)
    return format(date, 'yyyy-MM-dd HH:mm:ss')
  } catch (e) {
    return `${timestamp}`
  }
}

// 刷新数据
const refreshData = () => {
  fetchEvents()
  // 暂时注释掉统计数据的获取，API 响应格式不匹配
  // fetchStats()
}

// 获取事件列表
const fetchEvents = async () => {
  if (!projectStore.currentProject) return

  loading.value = true

  try {
    const response = await eventService.getEvents({
      projectId: projectStore.currentProject.id,
      eventType: filterType.value === 'all' ? undefined : filterType.value,
      page: currentPage.value,
      pageSize: pageSize.value
    })

    console.log(response.data)

    events.value = response.data?.list ?? []

    totalEvents.value = response.data?.total ?? 0
  } catch (error) {
    console.error('Failed to fetch events:', error)
  } finally {
    loading.value = false
  }
}

// 获取统计数据 - 暂时注释掉，API 响应格式不匹配
/*
const fetchStats = async () => {
  if (!projectStore.currentProject) return

  try {
    // 获取当前日期和一周前的日期
    const now = new Date()
    const oneWeekAgo = new Date()
    oneWeekAgo.setDate(oneWeekAgo.getDate() - 7)

    const startTime = format(oneWeekAgo, 'yyyy-MM-dd')
    const endTime = format(now, 'yyyy-MM-dd')

    // 获取点击和曝光事件的统计数据 - 需要等待接口文档更新
    const clickResponse = await eventService.getEventStats({
      projectId: projectStore.currentProject.id,
      startTime,
      endTime
    })

    const exposureResponse = await eventService.getEventStats({
      projectId: projectStore.currentProject.id,
      startTime,
      endTime
    })

    // 临时使用，等待接口文档更新
    stats.todayClicks = 0
    stats.todayClicksTrend = 0
    stats.weekClicks = 0
    stats.weekClicksTrend = 0
    stats.todayExposures = 0
    stats.todayExposuresTrend = 0
    stats.weekExposures = 0
    stats.weekExposuresTrend = 0
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}
*/

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

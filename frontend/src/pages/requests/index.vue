<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">请求监控</h1>
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
      <!-- 请求统计卡片 -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">请求总数</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.totalRequests }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.totalRequestsTrend > 0 ? '+' : '' }}{{ stats.totalRequestsTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">成功率</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.successRate }}%</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.successRateTrend > 0 ? '+' : '' }}{{ stats.successRateTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">平均响应时间</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.avgResponseTime }}ms</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.avgResponseTimeTrend > 0 ? '+' : '' }}{{ stats.avgResponseTimeTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">错误请求</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.errorRequests }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.errorRequestsTrend > 0 ? '+' : '' }}{{ stats.errorRequestsTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- 请求趋势图 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>请求趋势</CardTitle>
          <CardDescription>过去 7 天的请求数量和成功率趋势</CardDescription>
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

      <!-- 请求列表 -->
      <Card>
        <CardHeader>
          <CardTitle>请求列表</CardTitle>
          <CardDescription>HTTP 请求记录</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center space-x-2 mb-4">
              <Input
                v-model="searchQuery"
                placeholder="搜索请求..."
                class="max-w-sm"
              />
              <select
                v-model="filterStatus"
                class="h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="all">所有状态</option>
                <option value="success">成功</option>
                <option value="error">失败</option>
              </select>
              <select
                v-model="filterMethod"
                class="h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="all">所有方法</option>
                <option value="GET">GET</option>
                <option value="POST">POST</option>
                <option value="PUT">PUT</option>
                <option value="DELETE">DELETE</option>
              </select>
            </div>

            <div v-if="loading" class="flex justify-center py-8">
              <div class="animate-spin">&#8635;</div>
            </div>

            <div v-else-if="filteredRequests.length === 0" class="text-center py-8 text-muted-foreground">
              暂无请求数据
            </div>

            <div v-else class="border rounded-md">
              <div class="grid grid-cols-12 gap-4 p-4 font-medium border-b">
                <div class="col-span-1">方法</div>
                <div class="col-span-4">URL</div>
                <div class="col-span-1">状态</div>
                <div class="col-span-2">响应时间</div>
                <div class="col-span-2">时间</div>
                <div class="col-span-2">操作</div>
              </div>
              <div
                v-for="request in filteredRequests"
                :key="request.id"
                class="grid grid-cols-12 gap-4 p-4 border-b last:border-0 hover:bg-muted/50"
              >
                <div class="col-span-1">{{ request.method }}</div>
                <div class="col-span-4 truncate" :title="request.url">
                  {{ request.url }}
                </div>
                <div class="col-span-1">
                  <span
                    :class="{
                      'text-green-500': request.status >= 200 && request.status < 300,
                      'text-red-500': request.status >= 400
                    }"
                  >
                    {{ request.status }}
                  </span>
                </div>
                <div class="col-span-2">{{ request.duration }}ms</div>
                <div class="col-span-2">{{ formatDate(request.time) }}</div>
                <div class="col-span-2">
                  <Button variant="ghost" size="sm" @click="showRequestDetail(request)">
                    详情
                  </Button>
                </div>
              </div>
            </div>

            <!-- 分页 -->
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-foreground">
                显示 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize, totalRequests) }} 条，共 {{ totalRequests }} 条
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
                  :disabled="currentPage * pageSize >= totalRequests"
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

    <!-- 请求详情对话框 -->
    <div v-if="showDetailDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <Card class="w-full max-w-3xl max-h-[80vh] overflow-auto">
        <CardHeader>
          <CardTitle>请求详情</CardTitle>
          <CardDescription>{{ selectedRequest?.method }} {{ selectedRequest?.url }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-6">
            <div>
              <h3 class="text-lg font-medium mb-2">基本信息</h3>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <p class="text-sm text-muted-foreground">请求方法</p>
                  <p class="font-medium">{{ selectedRequest?.method }}</p>
                </div>
                <div>
                  <p class="text-sm text-muted-foreground">状态码</p>
                  <p class="font-medium" :class="{
                    'text-green-500': selectedRequest?.status >= 200 && selectedRequest?.status < 300,
                    'text-red-500': selectedRequest?.status >= 400
                  }">{{ selectedRequest?.status }}</p>
                </div>
                <div>
                  <p class="text-sm text-muted-foreground">响应时间</p>
                  <p class="font-medium">{{ selectedRequest?.duration }}ms</p>
                </div>
                <div>
                  <p class="text-sm text-muted-foreground">请求时间</p>
                  <p class="font-medium">{{ formatDate(selectedRequest?.time) }}</p>
                </div>
              </div>
            </div>

            <div>
              <h3 class="text-lg font-medium mb-2">请求 URL</h3>
              <pre class="bg-muted p-4 rounded-md overflow-x-auto text-sm">{{ selectedRequest?.url }}</pre>
            </div>

            <div>
              <h3 class="text-lg font-medium mb-2">请求参数</h3>
              <pre class="bg-muted p-4 rounded-md overflow-x-auto text-sm">{{ selectedRequest?.params || '无' }}</pre>
            </div>

            <div>
              <h3 class="text-lg font-medium mb-2">响应数据</h3>
              <pre class="bg-muted p-4 rounded-md overflow-x-auto text-sm">{{ selectedRequest?.response || '无' }}</pre>
            </div>
          </div>
        </CardContent>
        <CardFooter>
          <Button @click="showDetailDialog = false">关闭</Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { format } from 'date-fns'
import { useProjectStore } from '@/stores/project'
import { useEventService, useRequestService } from '@/services'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const router = useRouter()
const projectStore = useProjectStore()
const eventService = useEventService()
const requestService = useRequestService()

// 统计数据
const stats = reactive({
  totalRequests: 0,
  totalRequestsTrend: 0,
  successRate: 0,
  successRateTrend: 0,
  avgResponseTime: 0,
  avgResponseTimeTrend: 0,
  errorRequests: 0,
  errorRequestsTrend: 0
})

// 请求列表
const requests = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')
const filterStatus = ref('all')
const filterMethod = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const totalRequests = ref(0)

// 请求详情对话框
const showDetailDialog = ref(false)
const selectedRequest = ref<any>(null)

// 过滤后的请求列表
const filteredRequests = computed(() => {
  let result = requests.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(request =>
      request.url.toLowerCase().includes(query)
    )
  }

  // 状态过滤
  if (filterStatus.value !== 'all') {
    if (filterStatus.value === 'success') {
      result = result.filter(request => request.status >= 200 && request.status < 300)
    } else if (filterStatus.value === 'error') {
      result = result.filter(request => request.status >= 400)
    }
  }

  // 方法过滤
  if (filterMethod.value !== 'all') {
    result = result.filter(request => request.method === filterMethod.value)
  }

  return result
})

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''

  try {
    const date = new Date(dateString)
    return format(date, 'yyyy-MM-dd HH:mm:ss')
  } catch (e) {
    return dateString
  }
}

// 显示请求详情
const showRequestDetail = (request: any) => {
  selectedRequest.value = request
  showDetailDialog.value = true
}

// 刷新数据
const refreshData = () => {
  fetchRequests()
  fetchStats()
}

// 获取请求列表
const fetchRequests = async () => {
  if (!projectStore.currentProject) return

  loading.value = true

  try {
    const response = await eventService.getEvents({
      projectId: projectStore.currentProject.id,
      eventType: 'request',
      page: currentPage.value,
      pageSize: pageSize.value
    })

    requests.value = response.data.items.map((item: any) => {
      const requestEvent = item.requestEvent || {}

      return {
        id: item.id,
        method: requestEvent.method || 'GET',
        url: requestEvent.url || '',
        status: requestEvent.status || 0,
        duration: requestEvent.duration || 0,
        time: item.createdAt,
        params: requestEvent.params || '{}',
        response: requestEvent.response || '{}'
      }
    })

    totalRequests.value = response.data.total
  } catch (error) {
    console.error('Failed to fetch requests:', error)
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

    const response = await requestService.getRequestErrorStats({
      projectId: projectStore.currentProject.id,
      startTime,
      endTime
    })

    const data = response.data

    stats.totalRequests = data.totalRequests || 0
    stats.totalRequestsTrend = data.totalRequestsTrend || 0
    stats.successRate = data.successRate || 0
    stats.successRateTrend = data.successRateTrend || 0
    stats.avgResponseTime = data.avgResponseTime || 0
    stats.avgResponseTimeTrend = data.avgResponseTimeTrend || 0
    stats.errorRequests = data.errorRequests || 0
    stats.errorRequestsTrend = data.errorRequestsTrend || 0
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
watch([searchQuery, filterStatus, filterMethod], () => {
  currentPage.value = 1
})

onMounted(() => {
  if (projectStore.currentProject) {
    refreshData()
  }
})
</script>

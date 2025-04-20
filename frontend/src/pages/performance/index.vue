<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">性能监控</h1>
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
      <!-- 性能统计卡片 -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">平均加载时间</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.avgLoadTime }}ms</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.loadTimeTrend > 0 ? '+' : '' }}{{ stats.loadTimeTrend }}% 相比上周
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">首次内容绘制</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.fcp }}ms</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.fcpTrend > 0 ? '+' : '' }}{{ stats.fcpTrend }}% 相比上周
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">最大内容绘制</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.lcp }}ms</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.lcpTrend > 0 ? '+' : '' }}{{ stats.lcpTrend }}% 相比上周
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">累积布局偏移</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.cls }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.clsTrend > 0 ? '+' : '' }}{{ stats.clsTrend }}% 相比上周
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- 性能趋势图 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>性能趋势</CardTitle>
          <CardDescription>过去 7 天的页面加载性能趋势</CardDescription>
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

      <!-- 资源加载列表 -->
      <Card>
        <CardHeader>
          <CardTitle>资源加载</CardTitle>
          <CardDescription>页面资源加载性能</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center space-x-2 mb-4">
              <Input
                v-model="searchQuery"
                placeholder="搜索资源..."
                class="max-w-sm"
              />
              <select
                v-model="filterType"
                class="h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="all">所有类型</option>
                <option value="script">脚本</option>
                <option value="stylesheet">样式表</option>
                <option value="image">图片</option>
                <option value="font">字体</option>
                <option value="fetch">Fetch</option>
                <option value="xmlhttprequest">XHR</option>
              </select>
            </div>

            <div v-if="loading" class="flex justify-center py-8">
              <div class="animate-spin">&#8635;</div>
            </div>

            <div v-else-if="filteredResources.length === 0" class="text-center py-8 text-muted-foreground">
              暂无资源数据
            </div>

            <div v-else class="border rounded-md">
              <div class="grid grid-cols-12 gap-4 p-4 font-medium border-b">
                <div class="col-span-5">资源 URL</div>
                <div class="col-span-1">类型</div>
                <div class="col-span-2">大小</div>
                <div class="col-span-2">加载时间</div>
                <div class="col-span-2">状态</div>
              </div>
              <div
                v-for="resource in filteredResources"
                :key="resource.id"
                class="grid grid-cols-12 gap-4 p-4 border-b last:border-0 hover:bg-muted/50"
              >
                <div class="col-span-5 truncate" :title="resource.url">
                  {{ resource.url }}
                </div>
                <div class="col-span-1">{{ resource.type }}</div>
                <div class="col-span-2">{{ formatSize(resource.size) }}</div>
                <div class="col-span-2">{{ resource.duration }}ms</div>
                <div class="col-span-2">
                  <span
                    :class="{
                      'text-green-500': resource.status === 'success',
                      'text-red-500': resource.status === 'error'
                    }"
                  >
                    {{ resource.status }}
                  </span>
                </div>
              </div>
            </div>

            <!-- 分页 -->
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-foreground">
                显示 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize, totalResources) }} 条，共 {{ totalResources }} 条
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
                  :disabled="currentPage * pageSize >= totalResources"
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
import { usePerformanceService } from '@/services'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const router = useRouter()
const projectStore = useProjectStore()
const performanceService = usePerformanceService()

// 统计数据
const stats = reactive({
  avgLoadTime: 0,
  loadTimeTrend: 0,
  fcp: 0,
  fcpTrend: 0,
  lcp: 0,
  lcpTrend: 0,
  cls: 0,
  clsTrend: 0
})

// 资源列表
const resources = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')
const filterType = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const totalResources = ref(0)

// 过滤后的资源列表
const filteredResources = computed(() => {
  let result = resources.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(resource =>
      resource.url.toLowerCase().includes(query)
    )
  }

  // 类型过滤
  if (filterType.value !== 'all') {
    result = result.filter(resource => resource.type === filterType.value)
  }

  return result
})

// 格式化文件大小
const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 刷新数据
const refreshData = () => {
  fetchResources()
  fetchStats()
}

// 获取资源列表
const fetchResources = async () => {
  if (!projectStore.currentProject) return

  loading.value = true

  try {
    const response = await performanceService.getPerformanceStats({
      projectId: projectStore.currentProject.id,
      startTime: format(new Date(Date.now() - 24 * 60 * 60 * 1000), 'yyyy-MM-dd'), // 最近24小时
      endTime: format(new Date(), 'yyyy-MM-dd'),
      includeResources: true
    })

    // 处理资源数据
    resources.value = (response.data.resources || []).map((resource: any, index: number) => {
      return {
        id: index + 1,
        url: resource.url || '',
        type: resource.type || 'other',
        size: resource.size || 0,
        duration: resource.duration || 0,
        status: resource.status === 200 ? 'success' : 'error'
      }
    })

    totalResources.value = resources.value.length
  } catch (error) {
    console.error('Failed to fetch resources:', error)
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

    const response = await performanceService.getPerformanceStats({
      projectId: projectStore.currentProject.id,
      startTime,
      endTime
    })

    const data = response.data

    stats.avgLoadTime = data.avgLoadTime || 0
    stats.loadTimeTrend = data.loadTimeTrend || 0
    stats.fcp = data.avgFcp || 0
    stats.fcpTrend = data.fcpTrend || 0
    stats.lcp = data.avgLcp || 0
    stats.lcpTrend = data.lcpTrend || 0
    stats.cls = data.avgCls || 0
    stats.clsTrend = data.clsTrend || 0
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

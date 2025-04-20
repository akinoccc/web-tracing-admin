<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">资源监控</h1>
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
      <!-- 资源统计卡片 -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">资源总数</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.totalResources }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.totalResourcesTrend > 0 ? '+' : '' }}{{ stats.totalResourcesTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">加载成功率</CardTitle>
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
            <CardTitle class="text-sm font-medium">平均加载时间</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.avgLoadTime }}ms</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.avgLoadTimeTrend > 0 ? '+' : '' }}{{ stats.avgLoadTimeTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">总资源大小</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ formatSize(stats.totalSize) }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.totalSizeTrend > 0 ? '+' : '' }}{{ stats.totalSizeTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- 资源类型分布图 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>资源类型分布</CardTitle>
          <CardDescription>各类型资源的数量和大小分布</CardDescription>
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

      <!-- 资源列表 -->
      <Card>
        <CardHeader>
          <CardTitle>资源列表</CardTitle>
          <CardDescription>页面加载的资源列表</CardDescription>
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
                <option value="media">媒体</option>
                <option value="other">其他</option>
              </select>
              <select
                v-model="filterStatus"
                class="h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="all">所有状态</option>
                <option value="success">成功</option>
                <option value="error">失败</option>
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
                <div class="col-span-5">URL</div>
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
import { useProjectStore } from '@/stores/project'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const router = useRouter()
const projectStore = useProjectStore()

// 统计数据
const stats = reactive({
  totalResources: 0,
  totalResourcesTrend: 0,
  successRate: 0,
  successRateTrend: 0,
  avgLoadTime: 0,
  avgLoadTimeTrend: 0,
  totalSize: 0,
  totalSizeTrend: 0
})

// 资源列表
const resources = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')
const filterType = ref('all')
const filterStatus = ref('all')
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

  // 状态过滤
  if (filterStatus.value !== 'all') {
    result = result.filter(resource => resource.status === filterStatus.value)
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
    // 模拟 API 请求
    await new Promise(resolve => setTimeout(resolve, 500))

    // 模拟数据
    resources.value = [
      {
        id: 1,
        url: 'https://example.com/main.js',
        type: 'script',
        size: 156720,
        duration: 324,
        status: 'success'
      },
      {
        id: 2,
        url: 'https://example.com/styles.css',
        type: 'stylesheet',
        size: 45280,
        duration: 128,
        status: 'success'
      },
      {
        id: 3,
        url: 'https://example.com/logo.png',
        type: 'image',
        size: 24680,
        duration: 87,
        status: 'success'
      },
      {
        id: 4,
        url: 'https://example.com/fonts/roboto.woff2',
        type: 'font',
        size: 78340,
        duration: 156,
        status: 'success'
      },
      {
        id: 5,
        url: 'https://example.com/video/intro.mp4',
        type: 'media',
        size: 1245680,
        duration: 1245,
        status: 'success'
      },
      {
        id: 6,
        url: 'https://example.com/missing-image.jpg',
        type: 'image',
        size: 0,
        duration: 45,
        status: 'error'
      }
    ]

    totalResources.value = 28 // 模拟总数
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
    // 模拟 API 请求
    await new Promise(resolve => setTimeout(resolve, 300))

    // 模拟数据
    stats.totalResources = 28
    stats.totalResourcesTrend = 5
    stats.successRate = 96.4
    stats.successRateTrend = -0.8
    stats.avgLoadTime = 215
    stats.avgLoadTimeTrend = -12
    stats.totalSize = 3567840
    stats.totalSizeTrend = 8
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
watch([searchQuery, filterType, filterStatus], () => {
  currentPage.value = 1
})

onMounted(() => {
  if (projectStore.currentProject) {
    refreshData()
  }
})
</script>

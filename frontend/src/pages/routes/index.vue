<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">路由监控</h1>
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
      <!-- 路由统计卡片 -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">页面访问总数</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.totalPageViews }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.totalPageViewsTrend > 0 ? '+' : '' }}{{ stats.totalPageViewsTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">独立访客数</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.uniqueVisitors }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.uniqueVisitorsTrend > 0 ? '+' : '' }}{{ stats.uniqueVisitorsTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">平均停留时间</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ formatDuration(stats.avgDuration) }}</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.avgDurationTrend > 0 ? '+' : '' }}{{ stats.avgDurationTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">跳出率</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ stats.bounceRate }}%</div>
            <p class="text-xs text-muted-foreground">
              {{ stats.bounceRateTrend > 0 ? '+' : '' }}{{ stats.bounceRateTrend }}% 相比昨天
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- 访问量趋势图 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>访问量趋势</CardTitle>
          <CardDescription>过去 7 天的页面访问量趋势</CardDescription>
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

      <!-- 热门页面 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>热门页面</CardTitle>
          <CardDescription>访问量最高的页面</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div v-if="loading" class="flex justify-center py-8">
              <div class="animate-spin">&#8635;</div>
            </div>

            <div v-else-if="popularPages.length === 0" class="text-center py-8 text-muted-foreground">
              暂无页面数据
            </div>

            <div v-else class="border rounded-md">
              <div class="grid grid-cols-12 gap-4 p-4 font-medium border-b">
                <div class="col-span-1">排名</div>
                <div class="col-span-5">页面 URL</div>
                <div class="col-span-2">访问量</div>
                <div class="col-span-2">平均停留时间</div>
                <div class="col-span-2">跳出率</div>
              </div>
              <div
                v-for="(page, index) in popularPages"
                :key="page.id"
                class="grid grid-cols-12 gap-4 p-4 border-b last:border-0 hover:bg-muted/50"
              >
                <div class="col-span-1">{{ index + 1 }}</div>
                <div class="col-span-5 truncate" :title="page.url">
                  {{ page.url }}
                </div>
                <div class="col-span-2">{{ page.views }}</div>
                <div class="col-span-2">{{ formatDuration(page.avgDuration) }}</div>
                <div class="col-span-2">{{ page.bounceRate }}%</div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 路由访问列表 -->
      <Card>
        <CardHeader>
          <CardTitle>路由访问记录</CardTitle>
          <CardDescription>用户页面访问记录</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center space-x-2 mb-4">
              <Input
                v-model="searchQuery"
                placeholder="搜索页面..."
                class="max-w-sm"
              />
            </div>

            <div v-if="loading" class="flex justify-center py-8">
              <div class="animate-spin">&#8635;</div>
            </div>

            <div v-else-if="filteredRoutes.length === 0" class="text-center py-8 text-muted-foreground">
              暂无路由数据
            </div>

            <div v-else class="border rounded-md">
              <div class="grid grid-cols-12 gap-4 p-4 font-medium border-b">
                <div class="col-span-5">页面 URL</div>
                <div class="col-span-3">来源页面</div>
                <div class="col-span-2">访问时间</div>
                <div class="col-span-2">停留时间</div>
              </div>
              <div
                v-for="route in filteredRoutes"
                :key="route.id"
                class="grid grid-cols-12 gap-4 p-4 border-b last:border-0 hover:bg-muted/50"
              >
                <div class="col-span-5 truncate" :title="route.url">
                  {{ route.url }}
                </div>
                <div class="col-span-3 truncate" :title="route.referrer">
                  {{ route.referrer || '直接访问' }}
                </div>
                <div class="col-span-2">{{ formatDate(route.time) }}</div>
                <div class="col-span-2">{{ formatDuration(route.duration) }}</div>
              </div>
            </div>

            <!-- 分页 -->
            <div class="flex items-center justify-between">
              <div class="text-sm text-muted-foreground">
                显示 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize, totalRoutes) }} 条，共 {{ totalRoutes }} 条
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
                  :disabled="currentPage * pageSize >= totalRoutes"
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
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const router = useRouter()
const projectStore = useProjectStore()

// 统计数据
const stats = reactive({
  totalPageViews: 0,
  totalPageViewsTrend: 0,
  uniqueVisitors: 0,
  uniqueVisitorsTrend: 0,
  avgDuration: 0,
  avgDurationTrend: 0,
  bounceRate: 0,
  bounceRateTrend: 0
})

// 热门页面
const popularPages = ref<any[]>([])

// 路由列表
const routes = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const totalRoutes = ref(0)

// 过滤后的路由列表
const filteredRoutes = computed(() => {
  let result = routes.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(route => 
      route.url.toLowerCase().includes(query) || 
      (route.referrer && route.referrer.toLowerCase().includes(query))
    )
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

// 格式化持续时间
const formatDuration = (seconds: number) => {
  if (seconds < 60) {
    return `${seconds}秒`
  } else if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60
    return `${minutes}分${remainingSeconds}秒`
  } else {
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    return `${hours}时${minutes}分`
  }
}

// 刷新数据
const refreshData = () => {
  fetchRoutes()
  fetchPopularPages()
  fetchStats()
}

// 获取路由列表
const fetchRoutes = async () => {
  if (!projectStore.currentProject) return

  loading.value = true

  try {
    // 模拟 API 请求
    await new Promise(resolve => setTimeout(resolve, 500))

    // 模拟数据
    routes.value = [
      {
        id: 1,
        url: 'https://example.com/',
        referrer: '',
        time: '2023-12-01T14:32:45Z',
        duration: 120
      },
      {
        id: 2,
        url: 'https://example.com/products',
        referrer: 'https://example.com/',
        time: '2023-12-01T14:34:45Z',
        duration: 85
      },
      {
        id: 3,
        url: 'https://example.com/products/1',
        referrer: 'https://example.com/products',
        time: '2023-12-01T14:36:10Z',
        duration: 210
      },
      {
        id: 4,
        url: 'https://example.com/cart',
        referrer: 'https://example.com/products/1',
        time: '2023-12-01T14:39:40Z',
        duration: 45
      },
      {
        id: 5,
        url: 'https://example.com/checkout',
        referrer: 'https://example.com/cart',
        time: '2023-12-01T14:40:25Z',
        duration: 180
      }
    ]

    totalRoutes.value = 128 // 模拟总数
  } catch (error) {
    console.error('Failed to fetch routes:', error)
  } finally {
    loading.value = false
  }
}

// 获取热门页面
const fetchPopularPages = async () => {
  if (!projectStore.currentProject) return

  try {
    // 模拟 API 请求
    await new Promise(resolve => setTimeout(resolve, 300))

    // 模拟数据
    popularPages.value = [
      {
        id: 1,
        url: 'https://example.com/',
        views: 1256,
        avgDuration: 95,
        bounceRate: 35.2
      },
      {
        id: 2,
        url: 'https://example.com/products',
        views: 876,
        avgDuration: 120,
        bounceRate: 28.5
      },
      {
        id: 3,
        url: 'https://example.com/about',
        views: 542,
        avgDuration: 75,
        bounceRate: 42.1
      },
      {
        id: 4,
        url: 'https://example.com/contact',
        views: 423,
        avgDuration: 60,
        bounceRate: 38.7
      },
      {
        id: 5,
        url: 'https://example.com/blog',
        views: 387,
        avgDuration: 180,
        bounceRate: 25.3
      }
    ]
  } catch (error) {
    console.error('Failed to fetch popular pages:', error)
  }
}

// 获取统计数据
const fetchStats = async () => {
  if (!projectStore.currentProject) return

  try {
    // 模拟 API 请求
    await new Promise(resolve => setTimeout(resolve, 300))

    // 模拟数据
    stats.totalPageViews = 4256
    stats.totalPageViewsTrend = 8.5
    stats.uniqueVisitors = 1842
    stats.uniqueVisitorsTrend = 6.2
    stats.avgDuration = 105
    stats.avgDurationTrend = 3.8
    stats.bounceRate = 32.5
    stats.bounceRateTrend = -2.1
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

// 监听搜索条件变化
watch(searchQuery, () => {
  currentPage.value = 1
})

onMounted(() => {
  if (projectStore.currentProject) {
    refreshData()
  }
})
</script>

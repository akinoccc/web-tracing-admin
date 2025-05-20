<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold tracking-tight">用户行为分析</h1>
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
      <!-- 行为统计卡片 -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">总PV</CardTitle>
            <CardDescription>页面浏览量</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ behaviorStore.behaviorStats?.pvStats?.totalPV || 0 }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">总UV</CardTitle>
            <CardDescription>独立访客数</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ behaviorStore.behaviorStats?.pvStats?.totalUV || 0 }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">平均停留时间</CardTitle>
            <CardDescription>页面平均访问时长</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ formatDuration(behaviorStore.behaviorStats?.pvStats?.avgStayTime) }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">跳出率</CardTitle>
            <CardDescription>访问单页就离开的比率</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ formatPercentage(behaviorStore.behaviorStats?.pvStats?.bounceRate) }}</div>
          </CardContent>
        </Card>
      </div>

      <!-- PV/UV趋势图 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>访问趋势</CardTitle>
          <CardDescription>最近7天PV/UV趋势</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="h-80">
            <!-- 这里可以使用图表库，如 Chart.js 或 ECharts -->
            <div v-if="behaviorStore.behaviorStats?.pvTrend && behaviorStore.behaviorStats.pvTrend.length > 0" class="h-full">
              <!-- 简单的柱状图示例 -->
              <div class="flex h-full items-end space-x-2">
                <div 
                  v-for="item in behaviorStore.behaviorStats.pvTrend" 
                  :key="item.date" 
                  class="flex-1 relative group"
                >
                  <!-- PV柱状图 -->
                  <div 
                    class="bg-primary/20 hover:bg-primary/30 rounded-t-md w-full absolute bottom-0"
                    :style="{ height: `${getPVBarHeight(item.pv)}%` }"
                  ></div>
                  
                  <!-- UV柱状图 -->
                  <div 
                    class="bg-primary/60 hover:bg-primary/70 rounded-t-md w-full absolute bottom-0"
                    :style="{ height: `${getUVBarHeight(item.uv)}%` }"
                  ></div>
                  
                  <div class="absolute bottom-0 left-0 right-0 text-xs text-center -mb-6 opacity-70">
                    {{ formatDate(item.date) }}
                  </div>
                  <div class="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-background border px-2 py-1 rounded text-xs hidden group-hover:block whitespace-nowrap">
                    {{ item.date }}: PV {{ item.pv }}, UV {{ item.uv }}
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

      <!-- 热门页面 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>热门页面</CardTitle>
          <CardDescription>访问量最高的页面</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="!topPages.length" class="text-center p-4 text-muted-foreground">
            暂无数据
          </div>
          <div v-else>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>页面URL</TableHead>
                  <TableHead>访问量</TableHead>
                  <TableHead>占比</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="(count, url) in topPages" :key="url">
                  <TableCell class="max-w-md truncate">{{ url }}</TableCell>
                  <TableCell>{{ count }}</TableCell>
                  <TableCell>{{ formatPercentage(count / totalPV * 100) }}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <!-- 用户行为数据 -->
      <Tabs default-value="pv" class="w-full">
        <TabsList>
          <TabsTrigger value="pv">页面访问</TabsTrigger>
          <TabsTrigger value="clicks">点击行为</TabsTrigger>
        </TabsList>
        
        <!-- 页面访问 -->
        <TabsContent value="pv">
          <Card>
            <CardHeader>
              <CardTitle>页面访问数据</CardTitle>
              <CardDescription>
                共 {{ behaviorStore.totalPageViews }} 条记录
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div v-if="behaviorStore.loading" class="flex justify-center p-8">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
              </div>
              <div v-else-if="!behaviorStore.hasPageViews" class="text-center p-8 text-muted-foreground">
                暂无页面访问数据
              </div>
              <div v-else>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>页面URL</TableHead>
                      <TableHead>标题</TableHead>
                      <TableHead>来源</TableHead>
                      <TableHead>停留时间</TableHead>
                      <TableHead>浏览器</TableHead>
                      <TableHead>时间</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-for="item in behaviorStore.pageViews?.list" :key="item.id">
                      <TableCell class="max-w-xs truncate">{{ item.pageURL }}</TableCell>
                      <TableCell class="max-w-xs truncate">{{ item.title }}</TableCell>
                      <TableCell class="max-w-xs truncate">{{ item.referrer || '-' }}</TableCell>
                      <TableCell>{{ formatDuration(item.stayTime) }}</TableCell>
                      <TableCell>{{ item.browser }}</TableCell>
                      <TableCell>{{ formatTime(item.triggerTime) }}</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>

                <!-- 分页 -->
                <div class="flex items-center justify-end space-x-2 py-4">
                  <Button
                    variant="outline"
                    size="sm"
                    :disabled="behaviorStore.pvFilters.page <= 1"
                    @click="changePVPage(behaviorStore.pvFilters.page - 1)"
                  >
                    上一页
                  </Button>
                  <span class="text-sm text-muted-foreground">
                    第 {{ behaviorStore.pvFilters.page }} 页，共 {{ behaviorStore.totalPVPages }} 页
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    :disabled="behaviorStore.pvFilters.page >= behaviorStore.totalPVPages"
                    @click="changePVPage(behaviorStore.pvFilters.page + 1)"
                  >
                    下一页
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        
        <!-- 点击行为 -->
        <TabsContent value="clicks">
          <Card>
            <CardHeader>
              <CardTitle>点击行为数据</CardTitle>
              <CardDescription>
                共 {{ behaviorStore.totalClicks }} 条记录
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div v-if="behaviorStore.loading" class="flex justify-center p-8">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
              </div>
              <div v-else-if="!behaviorStore.hasClicks" class="text-center p-8 text-muted-foreground">
                暂无点击行为数据
              </div>
              <div v-else>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>元素路径</TableHead>
                      <TableHead>元素类型</TableHead>
                      <TableHead>内容</TableHead>
                      <TableHead>页面URL</TableHead>
                      <TableHead>时间</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-for="item in behaviorStore.clicks?.list" :key="item.id">
                      <TableCell class="max-w-xs truncate">{{ item.elementPath }}</TableCell>
                      <TableCell>{{ item.elementType }}</TableCell>
                      <TableCell class="max-w-xs truncate">{{ item.innerText || '-' }}</TableCell>
                      <TableCell class="max-w-xs truncate">{{ item.pageURL }}</TableCell>
                      <TableCell>{{ formatTime(item.triggerTime) }}</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>

                <!-- 分页 -->
                <div class="flex items-center justify-end space-x-2 py-4">
                  <Button
                    variant="outline"
                    size="sm"
                    :disabled="behaviorStore.clickFilters.page <= 1"
                    @click="changeClickPage(behaviorStore.clickFilters.page - 1)"
                  >
                    上一页
                  </Button>
                  <span class="text-sm text-muted-foreground">
                    第 {{ behaviorStore.clickFilters.page }} 页，共 {{ behaviorStore.totalClickPages }} 页
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    :disabled="behaviorStore.clickFilters.page >= behaviorStore.totalClickPages"
                    @click="changeClickPage(behaviorStore.clickFilters.page + 1)"
                  >
                    下一页
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { useBehaviorStore } from '@/stores/behavior'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const router = useRouter()
const projectStore = useProjectStore()
const behaviorStore = useBehaviorStore()

// 获取用户行为数据
const fetchBehaviorData = () => {
  if (!projectStore.currentProject) return
  behaviorStore.fetchPageViews()
  behaviorStore.fetchClicks()
  behaviorStore.fetchBehaviorStats()
}

// 刷新数据
const refreshData = () => {
  fetchBehaviorData()
}

// 切换页面访问页码
const changePVPage = (page: number) => {
  behaviorStore.setPVFilters({ page })
}

// 切换点击行为页码
const changeClickPage = (page: number) => {
  behaviorStore.setClickFilters({ page })
}

// 格式化时间戳
const formatTime = (timestamp: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString()
}

// 格式化日期
const formatDate = (dateStr: string) => {
  return dateStr.split('-').slice(1).join('/')
}

// 格式化持续时间（秒）
const formatDuration = (seconds?: number) => {
  if (!seconds) return '0秒'
  if (seconds < 60) return `${seconds}秒`
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}分${remainingSeconds}秒`
}

// 格式化百分比
const formatPercentage = (value?: number) => {
  if (value === undefined || value === null) return '0%'
  return `${value.toFixed(2)}%`
}

// 热门页面
const topPages = computed(() => {
  const pages = behaviorStore.behaviorStats?.pvStats?.topPages || {}
  return Object.entries(pages)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10)
    .reduce((acc, [url, count]) => {
      acc[url] = count
      return acc
    }, {} as Record<string, number>)
})

// 总PV
const totalPV = computed(() => behaviorStore.behaviorStats?.pvStats?.totalPV || 0)

// 计算PV柱状图高度百分比
const getPVBarHeight = (count: number) => {
  const maxCount = Math.max(...(behaviorStore.behaviorStats?.pvTrend?.map(item => item.pv) || [0]))
  if (maxCount === 0) return 0
  return (count / maxCount) * 90 // 最高90%，留出空间显示数值
}

// 计算UV柱状图高度百分比
const getUVBarHeight = (count: number) => {
  const maxCount = Math.max(...(behaviorStore.behaviorStats?.pvTrend?.map(item => item.uv) || [0]))
  if (maxCount === 0) return 0
  return (count / maxCount) * 90 // 最高90%，留出空间显示数值
}

onMounted(() => {
  fetchBehaviorData()
})
</script>

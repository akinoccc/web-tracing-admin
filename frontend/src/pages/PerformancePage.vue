<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold tracking-tight">性能监控</h1>
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
      <!-- 性能指标卡片 -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">平均FCP</CardTitle>
            <CardDescription>首次内容绘制</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ formatTime(performanceStore.performanceStats?.stats?.avgFCP) }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">平均LCP</CardTitle>
            <CardDescription>最大内容绘制</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ formatTime(performanceStore.performanceStats?.stats?.avgLCP) }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">平均TTFB</CardTitle>
            <CardDescription>首字节时间</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ formatTime(performanceStore.performanceStats?.stats?.avgTTFB) }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">平均加载时间</CardTitle>
            <CardDescription>页面完全加载</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ formatTime(performanceStore.performanceStats?.stats?.avgLoad) }}</div>
          </CardContent>
        </Card>
      </div>

      <!-- 性能趋势图 -->
      <Card class="mb-6">
        <CardHeader>
          <CardTitle>性能趋势</CardTitle>
          <CardDescription>最近7天性能指标趋势</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="h-80">
            <!-- 这里可以使用图表库，如 Chart.js 或 ECharts -->
            <div v-if="performanceStore.performanceStats?.trend && performanceStore.performanceStats.trend.length > 0" class="h-full">
              <!-- 简单的折线图示例 -->
              <div class="flex flex-col h-full">
                <div class="flex justify-between mb-2">
                  <span class="text-xs text-muted-foreground">FCP</span>
                  <span class="text-xs text-muted-foreground">LCP</span>
                  <span class="text-xs text-muted-foreground">TTFB</span>
                </div>
                <div class="flex-1 relative border-b border-l">
                  <!-- 坐标轴 -->
                  <div class="absolute left-0 top-0 h-full flex flex-col justify-between">
                    <span class="text-xs text-muted-foreground -translate-x-2">0ms</span>
                    <span class="text-xs text-muted-foreground -translate-x-2">{{ maxValue }}ms</span>
                  </div>
                  
                  <!-- 数据点 -->
                  <div class="absolute inset-0 flex items-end">
                    <div v-for="(item, index) in performanceStore.performanceStats.trend" :key="index" class="flex-1 flex flex-col items-center">
                      <div class="text-xs text-muted-foreground absolute -bottom-6">{{ formatDate(item.date) }}</div>
                      
                      <!-- FCP 点 -->
                      <div 
                        class="absolute w-2 h-2 bg-blue-500 rounded-full"
                        :style="{
                          bottom: `${getPercentage(item.fcp)}%`,
                          left: `${index * (100 / (performanceStore.performanceStats.trend.length - 1))}%`
                        }"
                      ></div>
                      
                      <!-- LCP 点 -->
                      <div 
                        class="absolute w-2 h-2 bg-green-500 rounded-full"
                        :style="{
                          bottom: `${getPercentage(item.lcp)}%`,
                          left: `${index * (100 / (performanceStore.performanceStats.trend.length - 1))}%`
                        }"
                      ></div>
                      
                      <!-- TTFB 点 -->
                      <div 
                        class="absolute w-2 h-2 bg-red-500 rounded-full"
                        :style="{
                          bottom: `${getPercentage(item.ttfb)}%`,
                          left: `${index * (100 / (performanceStore.performanceStats.trend.length - 1))}%`
                        }"
                      ></div>
                    </div>
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

      <!-- 性能数据列表 -->
      <Tabs default-value="page" class="w-full">
        <TabsList>
          <TabsTrigger value="page">页面性能</TabsTrigger>
          <TabsTrigger value="resource">资源性能</TabsTrigger>
        </TabsList>
        
        <!-- 页面性能 -->
        <TabsContent value="page">
          <Card>
            <CardHeader>
              <CardTitle>页面性能数据</CardTitle>
              <CardDescription>
                共 {{ performanceStore.totalPerformanceItems }} 条记录
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div v-if="performanceStore.loading" class="flex justify-center p-8">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
              </div>
              <div v-else-if="!performanceStore.hasPerformanceData" class="text-center p-8 text-muted-foreground">
                暂无性能数据
              </div>
              <div v-else>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>页面</TableHead>
                      <TableHead>FCP</TableHead>
                      <TableHead>LCP</TableHead>
                      <TableHead>TTFB</TableHead>
                      <TableHead>加载时间</TableHead>
                      <TableHead>浏览器</TableHead>
                      <TableHead>时间</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-for="item in performanceStore.performance?.list" :key="item.id">
                      <TableCell class="max-w-xs truncate">{{ item.pageURL }}</TableCell>
                      <TableCell>{{ formatTime(item.fcp) }}</TableCell>
                      <TableCell>{{ formatTime(item.lcp) }}</TableCell>
                      <TableCell>{{ formatTime(item.ttfb) }}</TableCell>
                      <TableCell>{{ formatTime(item.load) }}</TableCell>
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
                    :disabled="performanceStore.filters.page <= 1"
                    @click="changePage(performanceStore.filters.page - 1)"
                  >
                    上一页
                  </Button>
                  <span class="text-sm text-muted-foreground">
                    第 {{ performanceStore.filters.page }} 页，共 {{ performanceStore.totalPerformancePages }} 页
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    :disabled="performanceStore.filters.page >= performanceStore.totalPerformancePages"
                    @click="changePage(performanceStore.filters.page + 1)"
                  >
                    下一页
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        
        <!-- 资源性能 -->
        <TabsContent value="resource">
          <Card>
            <CardHeader>
              <CardTitle>资源性能数据</CardTitle>
              <CardDescription>
                共 {{ performanceStore.totalResourceItems }} 条记录
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div class="mb-4">
                <div class="grid gap-4 md:grid-cols-3">
                  <div class="space-y-2">
                    <label class="text-sm font-medium">资源类型</label>
                    <Select v-model="resourceFilters.resourceType" @update:modelValue="applyResourceFilters">
                      <option value="">全部类型</option>
                      <option value="script">脚本</option>
                      <option value="stylesheet">样式表</option>
                      <option value="image">图片</option>
                      <option value="font">字体</option>
                      <option value="fetch">Fetch</option>
                      <option value="xmlhttprequest">XHR</option>
                    </Select>
                  </div>
                </div>
              </div>
              
              <div v-if="performanceStore.loading" class="flex justify-center p-8">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
              </div>
              <div v-else-if="!performanceStore.hasResourcePerformanceData" class="text-center p-8 text-muted-foreground">
                暂无资源性能数据
              </div>
              <div v-else>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>资源URL</TableHead>
                      <TableHead>类型</TableHead>
                      <TableHead>加载时间</TableHead>
                      <TableHead>大小</TableHead>
                      <TableHead>页面</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-for="item in performanceStore.resourcePerformance?.list" :key="item.id">
                      <TableCell class="max-w-xs truncate">{{ item.resourceURL }}</TableCell>
                      <TableCell>{{ formatResourceType(item.resourceType) }}</TableCell>
                      <TableCell>{{ formatTime(item.duration) }}</TableCell>
                      <TableCell>{{ formatSize(item.transferSize) }}</TableCell>
                      <TableCell class="max-w-xs truncate">{{ item.pageURL }}</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>

                <!-- 分页 -->
                <div class="flex items-center justify-end space-x-2 py-4">
                  <Button
                    variant="outline"
                    size="sm"
                    :disabled="performanceStore.resourceFilters.page <= 1"
                    @click="changeResourcePage(performanceStore.resourceFilters.page - 1)"
                  >
                    上一页
                  </Button>
                  <span class="text-sm text-muted-foreground">
                    第 {{ performanceStore.resourceFilters.page }} 页，共 {{ performanceStore.totalResourcePages }} 页
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    :disabled="performanceStore.resourceFilters.page >= performanceStore.totalResourcePages"
                    @click="changeResourcePage(performanceStore.resourceFilters.page + 1)"
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
import { usePerformanceStore } from '@/stores/performance'
import { Button } from '@/components/ui/button'
import { Select } from '@/components/ui/select'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const router = useRouter()
const projectStore = useProjectStore()
const performanceStore = usePerformanceStore()

// 本地过滤器状态
const resourceFilters = ref({
  resourceType: ''
})

// 获取性能数据
const fetchPerformanceData = () => {
  if (!projectStore.currentProject) return
  performanceStore.fetchPerformance()
  performanceStore.fetchPerformanceStats()
  performanceStore.fetchResourcePerformance()
}

// 刷新数据
const refreshData = () => {
  fetchPerformanceData()
}

// 应用资源过滤器
const applyResourceFilters = () => {
  performanceStore.setResourceFilters({
    resourceType: resourceFilters.value.resourceType || undefined
  })
}

// 切换页码
const changePage = (page: number) => {
  performanceStore.setFilters({ page })
}

// 切换资源页码
const changeResourcePage = (page: number) => {
  performanceStore.setResourceFilters({ page })
}

// 格式化时间（毫秒）
const formatTime = (ms?: number) => {
  if (!ms) return '0ms'
  if (typeof ms === 'number') {
    if (ms < 1000) return `${ms}ms`
    return `${(ms / 1000).toFixed(2)}s`
  }
  return `${ms}ms`
}

// 格式化日期
const formatDate = (dateStr: string) => {
  return dateStr.split('-').slice(1).join('/')
}

// 格式化资源类型
const formatResourceType = (type: string) => {
  const typeMap: Record<string, string> = {
    script: '脚本',
    stylesheet: '样式表',
    image: '图片',
    font: '字体',
    fetch: 'Fetch',
    xmlhttprequest: 'XHR'
  }
  return typeMap[type] || type
}

// 格式化大小
const formatSize = (bytes?: number) => {
  if (!bytes) return '0B'
  if (bytes < 1024) return `${bytes}B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)}KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)}MB`
}

// 计算图表最大值
const maxValue = computed(() => {
  if (!performanceStore.performanceStats?.trend || performanceStore.performanceStats.trend.length === 0) {
    return 1000
  }
  
  const values = performanceStore.performanceStats.trend.flatMap(item => [
    item.fcp || 0,
    item.lcp || 0,
    item.ttfb || 0
  ])
  
  return Math.max(...values, 1000)
})

// 计算百分比高度
const getPercentage = (value: number) => {
  if (!value) return 0
  return (value / maxValue.value) * 90 // 最高90%，留出空间显示坐标
}

onMounted(() => {
  fetchPerformanceData()
})
</script>

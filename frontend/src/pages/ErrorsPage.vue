<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold tracking-tight">错误监控</h1>
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
      <!-- 错误统计卡片 -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">总错误数</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ errorStore.errors?.stats?.totalErrors || 0 }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">今日错误</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ errorStore.errors?.stats?.errorsToday || 0 }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">昨日错误</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ errorStore.errors?.stats?.errorsYesterday || 0 }}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-medium">影响用户数</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ errorStore.errors?.stats?.affectedUsers || 0 }}</div>
          </CardContent>
        </Card>
      </div>

      <!-- 过滤器 -->
      <Card class="mb-6">
        <CardContent class="p-4">
          <div class="grid gap-4 md:grid-cols-4">
            <div class="space-y-2">
              <label class="text-sm font-medium">错误类型</label>
              <Select v-model="filters.errorType" @update:modelValue="applyFilters">
                <SelectTrigger class="h-9">
                  <SelectValue placeholder="全部类型" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">全部类型</SelectItem>
                  <SelectItem value="js">JavaScript错误</SelectItem>
                  <SelectItem value="promise">Promise错误</SelectItem>
                  <SelectItem value="resource">资源加载错误</SelectItem>
                  <SelectItem value="ajax">AJAX请求错误</SelectItem>
                  <SelectItem value="vue">Vue错误</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium">严重程度</label>
              <Select v-model="filters.severity" @update:modelValue="applyFilters">
                <SelectTrigger class="h-9">
                  <SelectValue placeholder="全部" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">全部</SelectItem>
                  <SelectItem value="high">高</SelectItem>
                  <SelectItem value="medium">中</SelectItem>
                  <SelectItem value="low">低</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium">日期范围</label>
              <Popover>
                <PopoverTrigger>
                  <Button
                    variant="outline"
                    class="w-full justify-start text-left font-normal h-9"
                  >
                    <CalendarIcon class="mr-2 h-4 w-4" />
                    {{ dateRangeDisplay || "选择日期范围" }}
                  </Button>
                </PopoverTrigger>
                <PopoverContent class="w-auto p-0">
                  <RangeCalendar
                    initialFocus
                    :value="dateRange"
                    @update:model-value="onDateRangeChange"
                  />
                </PopoverContent>
              </Popover>
            </div>
          </div>
          <div class="flex justify-end mt-4">
            <Button variant="outline" class="mr-2" @click="resetFilters">重置</Button>
            <Button @click="applyFilters">应用</Button>
          </div>
        </CardContent>
      </Card>

      <!-- 错误列表 -->
      <Card>
        <CardHeader>
          <CardTitle>错误列表</CardTitle>
          <CardDescription>
            共 {{ errorStore.totalErrors }} 条记录
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="errorStore.loading" class="flex justify-center p-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
          </div>
          <div v-else-if="!errorStore.hasErrors" class="text-center p-8 text-muted-foreground">
            暂无错误数据
          </div>
          <div v-else>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>错误类型</TableHead>
                  <TableHead>错误信息</TableHead>
                  <TableHead>次数</TableHead>
                  <TableHead>最后发生</TableHead>
                  <TableHead>严重程度</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="error in errorStore.errors?.list" :key="error.id">
                  <TableCell>{{ formatErrorType(error.errorType) }}</TableCell>
                  <TableCell class="max-w-md truncate">{{ error.errorMessage }}</TableCell>
                  <TableCell>{{ error.count }}</TableCell>
                  <TableCell>{{ formatTime(error.lastSeen) }}</TableCell>
                  <TableCell>
                    <Badge :variant="getSeverityVariant(error.severity)">
                      {{ formatSeverity(error.severity) }}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Button variant="ghost" size="sm" @click="viewErrorDetail(error.id!)">
                      查看
                    </Button>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>

            <!-- 分页 -->
            <div class="flex items-center justify-end space-x-2 py-4">
              <Button
                variant="outline"
                size="sm"
                :disabled="errorStore.filters.page <= 1"
                @click="changePage(errorStore.filters.page - 1)"
              >
                上一页
              </Button>
              <span class="text-sm text-muted-foreground">
                第 {{ errorStore.filters.page }} 页，共 {{ errorStore.totalPages }} 页
              </span>
              <Button
                variant="outline"
                size="sm"
                :disabled="errorStore.filters.page >= errorStore.totalPages"
                @click="changePage(errorStore.filters.page + 1)"
              >
                下一页
              </Button>
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
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectValue, SelectTrigger, SelectItem } from '@/components/ui/select'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { RangeCalendar } from '@/components/ui/range-calendar'
import { CalendarIcon } from 'lucide-vue-next'
import { type DateValue } from '@internationalized/date'

const router = useRouter()
const projectStore = useProjectStore()
const errorStore = useErrorStore()

// 本地过滤器状态
const filters = ref({
  errorType: 'all',
  severity: 'all',
  startTime: undefined as number | undefined,
  endTime: undefined as number | undefined
})

// 日期范围选择
const dateRange = ref<{ start: DateValue | undefined, end: DateValue | undefined }>({
  start: undefined,
  end: undefined
})
const dateRangeDisplay = computed(() => {
  if (dateRange.value.start && dateRange.value.end) {
    const startDate = new Date(
      dateRange.value.start.year,
      dateRange.value.start.month - 1,
      dateRange.value.start.day
    )
    const endDate = new Date(
      dateRange.value.end.year,
      dateRange.value.end.month - 1,
      dateRange.value.end.day
    )
    return `${startDate.toLocaleDateString('zh-CN')} 至 ${endDate.toLocaleDateString('zh-CN')}`
  }
  return ''
})

// 获取错误列表
const fetchErrors = () => {
  if (!projectStore.currentProject) return
  errorStore.fetchErrors()
}

// 刷新数据
const refreshData = () => {
  fetchErrors()
}

// 应用过滤器
const applyFilters = () => {
  errorStore.setFilters({
    errorType: filters.value.errorType || undefined,
    severity: filters.value.severity || undefined,
    startTime: filters.value.startTime,
    endTime: filters.value.endTime
  })
}

// 重置过滤器
const resetFilters = () => {
  filters.value = {
    errorType: 'all',
    severity: 'all',
    startTime: undefined,
    endTime: undefined
  }
  dateRange.value = {
    start: undefined,
    end: undefined
  }
  errorStore.resetFilters()
}

// 处理日期范围选择
const onDateRangeChange = (range: any) => {
  dateRange.value = range

  if (range.start) {
    // 将 DateValue 转换为 JavaScript Date 对象，设置为当天的开始时间 (00:00:00)
    const startDate = new Date(range.start.year, range.start.month - 1, range.start.day)
    startDate.setHours(0, 0, 0, 0)
    filters.value.startTime = startDate.getTime()
  } else {
    filters.value.startTime = undefined
  }

  if (range.end) {
    // 将 DateValue 转换为 JavaScript Date 对象，设置为当天的结束时间 (23:59:59)
    const endDate = new Date(range.end.year, range.end.month - 1, range.end.day)
    endDate.setHours(23, 59, 59, 999)
    filters.value.endTime = endDate.getTime()
  } else {
    filters.value.endTime = undefined
  }
}

// 切换页码
const changePage = (page: number) => {
  errorStore.setFilters({ page })
}

// 查看错误详情
const viewErrorDetail = (id: number) => {
  router.push(`/errors/${id}`)
}

// 格式化错误类型
const formatErrorType = (type?: string) => {
  if (!type) {
    return ''
  }
  const typeMap: Record<string, string> = {
    js: 'JavaScript错误',
    promise: 'Promise错误',
    resource: '资源加载错误',
    ajax: 'AJAX请求错误',
    vue: 'Vue错误'
  }
  return typeMap[type] || type
}

// 格式化严重程度
const formatSeverity = (severity?: string) => {
  if (!severity) {
    return ''
  }
  const severityMap: Record<string, string> = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return severityMap[severity] || severity
}

// 获取严重程度对应的Badge变体
const getSeverityVariant = (severity?: string) => {
  if (!severity) {
    return 'default'
  }
  const variantMap: Record<string, 'default' | 'destructive' | 'outline'> = {
    high: 'destructive',
    medium: 'default',
    low: 'outline'
  }
  return variantMap[severity] || 'default'
}

// 格式化时间戳
const formatTime = (timestamp?: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString()
}

onMounted(() => {
  fetchErrors()
})
</script>

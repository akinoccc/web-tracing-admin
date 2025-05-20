<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-2">
        <Button variant="outline" size="sm" @click="router.back()">
          返回
        </Button>
        <h1 class="text-2xl font-semibold tracking-tight">错误详情</h1>
      </div>
      <div class="flex items-center space-x-2">
        <Button variant="outline" size="sm" @click="refreshData">
          刷新
        </Button>
      </div>
    </div>

    <div v-if="errorStore.loading" class="flex justify-center p-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
    </div>

    <div v-else-if="!errorDetail" class="flex flex-col items-center justify-center p-8 border rounded-lg">
      <p class="text-lg text-center text-muted-foreground mb-4">
        未找到错误详情
      </p>
      <Button @click="router.push('/errors')">
        返回错误列表
      </Button>
    </div>

    <div v-else class="space-y-6">
      <!-- 错误概览 -->
      <Card>
        <CardHeader>
          <CardTitle>错误概览</CardTitle>
          <CardDescription>
            错误分组信息
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="grid gap-4 md:grid-cols-2">
            <div>
              <h3 class="text-sm font-medium text-muted-foreground mb-1">错误类型</h3>
              <p>{{ formatErrorType(errorDetail.group.errorType) }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-muted-foreground mb-1">错误消息</h3>
              <p class="break-words">{{ errorDetail.group.errorMessage }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-muted-foreground mb-1">发生次数</h3>
              <p>{{ errorDetail.group.count }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-muted-foreground mb-1">首次发生</h3>
              <p>{{ formatTime(errorDetail.group.firstSeen) }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-muted-foreground mb-1">最后发生</h3>
              <p>{{ formatTime(errorDetail.group.lastSeen) }}</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-muted-foreground mb-1">严重程度</h3>
              <Badge :variant="getSeverityVariant(errorDetail.group.severity)">
                {{ formatSeverity(errorDetail.group.severity) }}
              </Badge>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 错误事件列表 -->
      <Card>
        <CardHeader>
          <CardTitle>错误事件</CardTitle>
          <CardDescription>
            最近 {{ errorDetail.events.length }} 条错误事件记录
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs default-value="list" class="w-full">
            <TabsList>
              <TabsTrigger value="list">列表视图</TabsTrigger>
              <TabsTrigger value="detail">详细视图</TabsTrigger>
            </TabsList>
            <TabsContent value="list">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>时间</TableHead>
                    <TableHead>页面</TableHead>
                    <TableHead>浏览器</TableHead>
                    <TableHead>操作系统</TableHead>
                    <TableHead>操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="event in errorDetail.events" :key="event.id">
                    <TableCell>{{ formatTime(event.triggerTime) }}</TableCell>
                    <TableCell class="max-w-xs truncate">{{ event.pageURL }}</TableCell>
                    <TableCell>{{ event.browser }}</TableCell>
                    <TableCell>{{ event.os }}</TableCell>
                    <TableCell>
                      <Button variant="ghost" size="sm" @click="selectedEvent = event">
                        查看
                      </Button>
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </TabsContent>
            <TabsContent value="detail">
              <div v-if="!selectedEvent && errorDetail.events.length > 0" class="text-center p-4">
                <p class="text-muted-foreground">请从列表视图中选择一个错误事件查看详情</p>
              </div>
              <div v-else-if="selectedEvent" class="space-y-4">
                <div>
                  <h3 class="text-sm font-medium text-muted-foreground mb-1">错误消息</h3>
                  <p class="p-2 bg-muted rounded-md">{{ selectedEvent.errorMessage }}</p>
                </div>
                <div>
                  <h3 class="text-sm font-medium text-muted-foreground mb-1">错误堆栈</h3>
                  <pre class="p-2 bg-muted rounded-md overflow-auto text-xs">{{ selectedEvent.errorStack }}</pre>
                </div>
                <div class="grid gap-4 md:grid-cols-2">
                  <div>
                    <h3 class="text-sm font-medium text-muted-foreground mb-1">文件路径</h3>
                    <p>{{ selectedEvent.filePath || '-' }}</p>
                  </div>
                  <div>
                    <h3 class="text-sm font-medium text-muted-foreground mb-1">位置</h3>
                    <p>行 {{ selectedEvent.lineNumber || '-' }}，列 {{ selectedEvent.columnNumber || '-' }}</p>
                  </div>
                  <div>
                    <h3 class="text-sm font-medium text-muted-foreground mb-1">页面URL</h3>
                    <p class="break-all">{{ selectedEvent.pageURL }}</p>
                  </div>
                  <div>
                    <h3 class="text-sm font-medium text-muted-foreground mb-1">设备信息</h3>
                    <p>{{ selectedEvent.browser }} / {{ selectedEvent.os }} / {{ selectedEvent.device || '未知' }}</p>
                  </div>
                </div>
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useErrorStore } from '@/stores/error'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import type { ServiceErrorDetailResponse } from '@/types/gen/service/ErrorDetailResponse'
import type { ServiceErrorEventItem } from '@/types/gen/service/ErrorEventItem'

const route = useRoute()
const router = useRouter()
const errorStore = useErrorStore()

const errorDetail = computed<ServiceErrorDetailResponse | null>(() => errorStore.errorDetail)
const selectedEvent = ref<ServiceErrorEventItem | null>(null)

// 获取错误详情
const fetchErrorDetail = async () => {
  const id = Number(route.params.id)
  if (isNaN(id)) {
    router.push('/errors')
    return
  }
  
  await errorStore.fetchErrorDetail(id)
  
  // 默认选中第一个事件
  if (errorStore.errorDetail?.events && errorStore.errorDetail.events.length > 0) {
    selectedEvent.value = errorStore.errorDetail.events[0]
  }
}

// 刷新数据
const refreshData = () => {
  fetchErrorDetail()
}

// 格式化错误类型
const formatErrorType = (type: string) => {
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
const formatSeverity = (severity: string) => {
  const severityMap: Record<string, string> = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return severityMap[severity] || severity
}

// 获取严重程度对应的Badge变体
const getSeverityVariant = (severity: string) => {
  const variantMap: Record<string, 'default' | 'destructive' | 'outline'> = {
    high: 'destructive',
    medium: 'default',
    low: 'outline'
  }
  return variantMap[severity] || 'default'
}

// 格式化时间戳
const formatTime = (timestamp: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString()
}

onMounted(() => {
  fetchErrorDetail()
})
</script>

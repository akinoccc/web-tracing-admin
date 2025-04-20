<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <Button variant="outline" size="sm" @click="router.back()">
          返回
        </Button>
        <h1 class="text-3xl font-bold tracking-tight">错误详情</h1>
      </div>
      <div class="flex items-center space-x-2">
        <Button variant="outline" @click="refreshData">
          刷新
        </Button>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin text-2xl">&#8635;</div>
    </div>

    <div v-else-if="!error" class="flex flex-col items-center justify-center py-12">
      <h2 class="text-xl font-semibold mb-4">未找到错误信息</h2>
      <Button @click="router.push('/errors')">返回错误列表</Button>
    </div>

    <div v-else class="space-y-6">
      <!-- 错误概览 -->
      <Card>
        <CardHeader>
          <CardTitle>错误概览</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">错误类型</h3>
                <p class="text-lg font-semibold">{{ error.type }}</p>
              </div>
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">错误信息</h3>
                <p class="text-lg font-semibold break-all">{{ error.message }}</p>
              </div>
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">发生时间</h3>
                <p class="text-lg font-semibold">{{ formatDate(error.time) }}</p>
              </div>
            </div>
            <div class="space-y-4">
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">浏览器</h3>
                <p class="text-lg font-semibold">{{ error.browser }} {{ error.browserVersion }}</p>
              </div>
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">操作系统</h3>
                <p class="text-lg font-semibold">{{ error.os }} {{ error.osVersion }}</p>
              </div>
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">页面 URL</h3>
                <p class="text-lg font-semibold break-all">{{ error.url }}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 错误堆栈 -->
      <Card>
        <CardHeader>
          <CardTitle>错误堆栈</CardTitle>
        </CardHeader>
        <CardContent>
          <pre class="bg-muted p-4 rounded-md overflow-x-auto text-sm">{{ error.stack }}</pre>
        </CardContent>
      </Card>

      <!-- 用户信息 -->
      <Card>
        <CardHeader>
          <CardTitle>用户信息</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">用户 ID</h3>
                <p class="text-lg font-semibold">{{ error.userId || '匿名用户' }}</p>
              </div>
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">IP 地址</h3>
                <p class="text-lg font-semibold">{{ error.ip }}</p>
              </div>
            </div>
            <div class="space-y-4">
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">设备类型</h3>
                <p class="text-lg font-semibold">{{ error.deviceType }}</p>
              </div>
              <div>
                <h3 class="text-sm font-medium text-muted-foreground">屏幕分辨率</h3>
                <p class="text-lg font-semibold">{{ error.screenWidth }}x{{ error.screenHeight }}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 错误录屏 -->
      <Card v-if="error.recordScreen">
        <CardHeader>
          <CardTitle>错误录屏</CardTitle>
          <CardDescription>错误发生时的屏幕录制</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="bg-muted rounded-md p-4 flex items-center justify-center">
            <p class="text-muted-foreground">录屏播放器将在此处显示</p>
          </div>
        </CardContent>
      </Card>

      <!-- 相似错误 -->
      <Card>
        <CardHeader>
          <CardTitle>相似错误</CardTitle>
          <CardDescription>与此错误相似的其他错误</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="similarErrors.length === 0" class="text-center py-4 text-muted-foreground">
            暂无相似错误
          </div>
          <div
            v-else
            v-for="similarError in similarErrors"
            :key="similarError.id"
            class="flex items-center p-4 border rounded-md mb-2 last:mb-0"
          >
            <div class="flex-1 space-y-1">
              <p class="text-sm font-medium">{{ similarError.message }}</p>
              <p class="text-xs text-muted-foreground">{{ formatDate(similarError.time) }}</p>
            </div>
            <Button variant="ghost" size="sm" @click="viewErrorDetail(similarError.id)">
              查看详情
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { format } from 'date-fns'
import { useEventService } from '@/services'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

const router = useRouter()
const route = useRoute()
const eventService = useEventService()

const error = ref<any>(null)
const errorDetail = ref<any>(null)
const similarErrors = ref<any[]>([])
const loading = ref(true)

// 格式化日期
const formatDate = (dateString: string) => {
  try {
    const date = new Date(dateString)
    return format(date, 'yyyy-MM-dd HH:mm:ss')
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
  fetchErrorDetail()
}

// 获取错误详情
const fetchErrorDetail = async () => {
  const errorId = route.params.id

  if (!errorId) {
    router.push('/errors')
    return
  }

  loading.value = true

  try {
    const response = await eventService.getEventDetail(Number(errorId))
    const data = response.data
    const event = data.event || {}
    const errorEvent = data.errorEvent || {}
    const baseInfo = event.baseInfo || {}

    errorDetail.value = {
      id: data.id,
      message: errorEvent.message || '',
      type: errorEvent.type || '',
      browser: `${baseInfo.browser || ''} ${baseInfo.browserVersion || ''}`.trim(),
      os: `${baseInfo.os || ''} ${baseInfo.osVersion || ''}`.trim(),
      time: data.createdAt,
      url: errorEvent.pageUrl || '',
      stack: errorEvent.stack || '',
      userId: baseInfo.userId || '',
      ip: baseInfo.ip || '',
      deviceType: baseInfo.deviceType || '',
      screenWidth: baseInfo.screenWidth || 0,
      screenHeight: baseInfo.screenHeight || 0,
      recordScreen: !!errorEvent.recordScreen
    }

    // 将错误详情赋值给 error
    error.value = errorDetail.value

    // 获取相似错误
    await fetchSimilarErrors()
  } catch (error) {
    console.error('Failed to fetch error detail:', error)
  } finally {
    loading.value = false
  }
}

// 获取相似错误
const fetchSimilarErrors = async () => {
  if (!errorDetail.value) return

  try {
    // 获取相似类型的错误
    const response = await eventService.getEvents({
      projectId: Number(route.query.projectId) || 0,
      eventType: 'error',
      page: 1,
      pageSize: 5
    })

    // 过滤出相似的错误，但不包括当前错误
    similarErrors.value = response.data.items
      .filter((item: any) => {
        const errorEvent = item.errorEvent || {}
        return (
          item.id !== errorDetail.value.id &&
          errorEvent.type === errorDetail.value.type
        )
      })
      .slice(0, 2) // 只取前两个
      .map((item: any) => {
        const errorEvent = item.errorEvent || {}
        return {
          id: item.id,
          message: errorEvent.message || '',
          time: item.createdAt
        }
      })
  } catch (error) {
    console.error('Failed to fetch similar errors:', error)
  }
}

onMounted(() => {
  fetchErrorDetail()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">设置</h1>
      <div class="flex items-center space-x-2">
        <Button variant="outline" @click="refreshData">
          刷新
        </Button>
      </div>
    </div>

    <div v-if="!projectDetail" class="flex flex-col items-center justify-center py-12">
      <h2 class="text-xl font-semibold mb-4">请先选择一个项目</h2>
      <Button @click="router.push('/')">返回仪表盘</Button>
    </div>

    <div v-else class="space-y-6">
      <!-- 项目设置 -->
      <Card>
        <CardHeader>
          <CardTitle>项目设置</CardTitle>
          <CardDescription>管理项目基本信息</CardDescription>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="updateProject" class="space-y-4">
            <div class="space-y-2">
              <label for="project-name" class="text-sm font-medium">项目名称</label>
              <Input
                id="project-name"
                v-model="projectForm.name"
                placeholder="请输入项目名称"
                required
              />
            </div>
            <div class="space-y-2">
              <label for="project-description" class="text-sm font-medium">项目描述</label>
              <Input
                id="project-description"
                v-model="projectForm.description"
                placeholder="请输入项目描述"
              />
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium">App Key</label>
              <div class="flex items-center space-x-2">
                <Input
                  :model-value="projectDetail.appKey"
                  readonly
                  class="bg-muted"
                />
                <Button type="button" variant="outline" size="sm" @click="copyAppKey">
                  复制
                </Button>
              </div>
              <p class="text-xs text-muted-foreground">
                App Key 用于 SDK 初始化，请妥善保管
              </p>
            </div>
            <div class="flex justify-end">
              <Button type="submit" :loading="updating">保存更改</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <!-- SDK 集成指南 -->
      <Card>
        <CardHeader>
          <CardTitle>SDK 集成指南</CardTitle>
          <CardDescription>如何在您的项目中集成 web-tracing-sdk</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div>
              <h3 class="text-lg font-medium mb-2">Vue 3 项目</h3>
              <pre class="bg-muted p-4 rounded-md overflow-x-auto text-sm">
import { createApp } from 'vue'
import App from './App.vue'
import WebTracing from '@web-tracing/vue3'

const app = createApp(App)

app.use(WebTracing, {
  dsn: '{{ sdkEndpoint }}',
  appKey: '{{ projectDetail.appKey }}',
  appName: '{{ projectDetail.name }}',
  debug: true,
  pv: true,
  performance: true,
  error: true,
  event: true
})

app.mount('#app')</pre>
            </div>

            <div>
              <h3 class="text-lg font-medium mb-2">Vue 2 项目</h3>
              <pre class="bg-muted p-4 rounded-md overflow-x-auto text-sm">
import Vue from 'vue'
import App from './App.vue'
import WebTracing from '@web-tracing/vue'

Vue.use(WebTracing, {
  dsn: '{{ sdkEndpoint }}',
  appKey: '{{ projectDetail.appKey }}',
  appName: '{{ projectDetail.name }}',
  debug: true,
  pv: true,
  performance: true,
  error: true,
  event: true
})

new Vue({
  render: h => h(App)
}).$mount('#app')</pre>
            </div>

            <div>
              <h3 class="text-lg font-medium mb-2">JavaScript 项目</h3>
              <pre class="bg-muted p-4 rounded-md overflow-x-auto text-sm">
import WebTracing from '@web-tracing/js'

WebTracing.init({
  dsn: '{{ sdkEndpoint }}',
  appKey: '{{ projectDetail.appKey }}',
  appName: '{{ projectDetail.name }}',
  debug: true,
  pv: true,
  performance: true,
  error: true,
  event: true
})</pre>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 数据管理 -->
      <Card>
        <CardHeader>
          <CardTitle>数据管理</CardTitle>
          <CardDescription>管理项目数据</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center justify-between p-4 border rounded-md">
              <div>
                <h3 class="font-medium">清除所有数据</h3>
                <p class="text-sm text-muted-foreground">
                  删除项目中的所有监控数据，此操作不可恢复
                </p>
              </div>
              <Button variant="destructive" @click="confirmClearData">
                清除数据
              </Button>
            </div>

            <div class="flex items-center justify-between p-4 border rounded-md">
              <div>
                <h3 class="font-medium">导出数据</h3>
                <p class="text-sm text-muted-foreground">
                  将项目数据导出为 JSON 格式
                </p>
              </div>
              <Button variant="outline" @click="exportData">
                导出
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 删除项目 -->
      <Card>
        <CardHeader>
          <CardTitle>危险区域</CardTitle>
          <CardDescription>删除项目及其所有数据</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center justify-between p-4 border border-destructive rounded-md">
              <div>
                <h3 class="font-medium text-destructive">删除项目</h3>
                <p class="text-sm text-muted-foreground">
                  删除项目及其所有数据，此操作不可恢复
                </p>
              </div>
              <Button variant="destructive" @click="confirmDeleteProject">
                删除项目
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- 确认对话框 -->
    <div v-if="showConfirmDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <Card class="w-full max-w-md">
        <CardHeader>
          <CardTitle>{{ confirmDialogTitle }}</CardTitle>
          <CardDescription>{{ confirmDialogMessage }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="confirmDialogType === 'delete'" class="space-y-4">
            <p class="text-sm text-destructive font-medium">
              请输入项目名称以确认删除
            </p>
            <Input
              v-model="confirmInput"
              placeholder="请输入项目名称"
            />
          </div>
        </CardContent>
        <CardFooter class="flex justify-between">
          <Button variant="outline" @click="showConfirmDialog = false">取消</Button>
          <Button
            :variant="confirmDialogType === 'delete' ? 'destructive' : 'default'"
            @click="handleConfirmAction"
            :disabled="confirmDialogType === 'delete' && confirmInput !== projectDetail?.name"
            :loading="confirmLoading"
          >
            确认
          </Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useProjectService } from '@/services'
import { ModelProject } from '@/types'

const router = useRouter()
const route = useRoute()
const projectService = useProjectService()

const projectDetail = ref<ModelProject | null>(null)

// 项目表单
const projectForm = reactive({
  name: '',
  description: ''
})

// 状态
const updating = ref(false)
const confirmLoading = ref(false)
const showConfirmDialog = ref(false)
const confirmDialogType = ref('')
const confirmDialogTitle = ref('')
const confirmDialogMessage = ref('')
const confirmInput = ref('')

// SDK 端点
const sdkEndpoint = computed(() => {
  return `${window.location.origin}/trackweb`
})

// 刷新数据
const refreshData = async () => {
  const res = await projectService.getProject(Number(route.params.id))
  projectDetail.value = res.data
  projectForm.name = res.data.name || ''
  projectForm.description = res.data.description || ''
}

// 更新项目
const updateProject = async () => {
  updating.value = true

  try {
    await projectService.updateProject(
      Number(route.params.id),
      projectForm
    )
    alert('项目更新成功')
  } catch (error) {
    console.error('Failed to update project:', error)
    alert('项目更新失败')
  } finally {
    updating.value = false
  }
}

// 复制 App Key
const copyAppKey = () => {
  navigator.clipboard.writeText(projectDetail.value?.appKey || '')
    .then(() => {
      alert('App Key 已复制到剪贴板')
    })
    .catch((err) => {
      console.error('Failed to copy App Key:', err)
      alert('复制失败')
    })
}

// 确认清除数据
const confirmClearData = () => {
  confirmDialogType.value = 'clear'
  confirmDialogTitle.value = '确认清除数据'
  confirmDialogMessage.value = '您确定要清除所有监控数据吗？此操作不可恢复。'
  showConfirmDialog.value = true
}

// 确认删除项目
const confirmDeleteProject = () => {
  confirmDialogType.value = 'delete'
  confirmDialogTitle.value = '确认删除项目'
  confirmDialogMessage.value = '您确定要删除此项目及其所有数据吗？此操作不可恢复。'
  confirmInput.value = ''
  showConfirmDialog.value = true
}

// 导出数据
const exportData = async () => {
  if (!projectDetail.value) return

  try {
    // 模拟 API 请求
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 模拟数据
    const data = {
      project: projectDetail.value,
      events: [
        // 模拟事件数据
      ]
    }

    // 创建下载链接
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${projectDetail.value?.name}-data.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    alert('数据导出成功')
  } catch (error) {
    console.error('Failed to export data:', error)
    alert('数据导出失败')
  }
}

// 处理确认操作
const handleConfirmAction = async () => {
  if (!projectDetail.value) return

  confirmLoading.value = true

  try {
    if (confirmDialogType.value === 'clear') {
      // 模拟清除数据
      await new Promise(resolve => setTimeout(resolve, 1000))
      alert('数据已清除')
    } else if (confirmDialogType.value === 'delete') {
      // 删除项目
      const success = await projectService.deleteProject(Number(route.params.id))
      if (success) {
        router.push('/')
      }
    }
  } catch (error) {
    console.error('Action failed:', error)
    alert('操作失败')
  } finally {
    confirmLoading.value = false
    showConfirmDialog.value = false
  }
}

onMounted(() => {
  refreshData()
})
</script>

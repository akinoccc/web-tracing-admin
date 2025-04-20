<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">仪表盘</h1>
      <div class="flex items-center space-x-2">
        <Button v-if="projectStore.hasProjects" @click="openProjectSelector">
          {{ projectStore.currentProject?.name || '选择项目' }}
        </Button>
        <Button v-else @click="openCreateProjectDialog">创建项目</Button>
      </div>
    </div>

    <div v-if="!projectStore.currentProject" class="flex flex-col items-center justify-center py-12">
      <h2 class="text-xl font-semibold mb-4">请先创建或选择一个项目</h2>
      <Button @click="openCreateProjectDialog">创建项目</Button>
    </div>

    <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <!-- 错误统计卡片 -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">错误总数</CardTitle>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            class="h-4 w-4 text-muted-foreground"
          >
            <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
          </svg>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.errorCount }}</div>
          <p class="text-xs text-muted-foreground">
            {{ stats.errorTrend > 0 ? '+' : '' }}{{ stats.errorTrend }}% 相比上周
          </p>
        </CardContent>
      </Card>

      <!-- 性能统计卡片 -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">平均加载时间</CardTitle>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            class="h-4 w-4 text-muted-foreground"
          >
            <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
          </svg>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.avgLoadTime }}ms</div>
          <p class="text-xs text-muted-foreground">
            {{ stats.loadTimeTrend > 0 ? '+' : '' }}{{ stats.loadTimeTrend }}% 相比上周
          </p>
        </CardContent>
      </Card>

      <!-- 请求统计卡片 -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">请求成功率</CardTitle>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            class="h-4 w-4 text-muted-foreground"
          >
            <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
          </svg>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.requestSuccessRate }}%</div>
          <p class="text-xs text-muted-foreground">
            {{ stats.requestRateTrend > 0 ? '+' : '' }}{{ stats.requestRateTrend }}% 相比上周
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- 错误趋势图 -->
    <Card v-if="projectStore.currentProject">
      <CardHeader>
        <CardTitle>错误趋势</CardTitle>
        <CardDescription>过去 7 天的错误数量趋势</CardDescription>
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

    <!-- 最近错误列表 -->
    <Card v-if="projectStore.currentProject">
      <CardHeader>
        <CardTitle>最近错误</CardTitle>
        <CardDescription>最近发生的错误事件</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="space-y-4">
          <div v-if="recentErrors.length === 0" class="text-center py-4 text-muted-foreground">
            暂无错误数据
          </div>
          <div
            v-for="error in recentErrors"
            :key="error.id"
            class="flex items-center p-4 border rounded-md"
          >
            <div class="flex-1 space-y-1">
              <p class="text-sm font-medium">{{ error.message }}</p>
              <p class="text-xs text-muted-foreground">{{ error.time }}</p>
            </div>
            <Button variant="ghost" size="sm" @click="viewErrorDetail(error.id)">
              查看详情
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- 创建项目对话框 -->
    <div v-if="showCreateProjectDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <Card class="w-full max-w-md">
        <CardHeader>
          <CardTitle>创建项目</CardTitle>
          <CardDescription>创建一个新的监控项目</CardDescription>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="createProject" class="space-y-4">
            <div class="space-y-2">
              <label for="project-name" class="text-sm font-medium">项目名称</label>
              <Input
                id="project-name"
                v-model="newProject.name"
                placeholder="请输入项目名称"
                required
              />
            </div>
            <div class="space-y-2">
              <label for="project-description" class="text-sm font-medium">项目描述</label>
              <Input
                id="project-description"
                v-model="newProject.description"
                placeholder="请输入项目描述"
              />
            </div>
          </form>
        </CardContent>
        <CardFooter class="flex justify-between">
          <Button variant="outline" @click="showCreateProjectDialog = false">取消</Button>
          <Button @click="createProject" :loading="projectStore.loading">创建</Button>
        </CardFooter>
      </Card>
    </div>

    <!-- 项目选择器对话框 -->
    <div v-if="showProjectSelector" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <Card class="w-full max-w-md">
        <CardHeader>
          <CardTitle>选择项目</CardTitle>
          <CardDescription>选择要监控的项目</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-2">
            <div
              v-for="project in projectStore.projects"
              :key="project.id"
              class="flex items-center p-3 border rounded-md cursor-pointer hover:bg-accent"
              :class="{ 'bg-accent': projectStore.currentProject?.id === project.id }"
              @click="selectProject(project)"
            >
              <div class="flex-1">
                <p class="font-medium">{{ project.name }}</p>
                <p class="text-xs text-muted-foreground">{{ project.description }}</p>
              </div>
            </div>
          </div>
        </CardContent>
        <CardFooter class="flex justify-between">
          <Button variant="outline" @click="showProjectSelector = false">取消</Button>
          <Button @click="openCreateProjectDialog">创建新项目</Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const router = useRouter()
const projectStore = useProjectStore()

// 统计数据
const stats = reactive({
  errorCount: 0,
  errorTrend: 0,
  avgLoadTime: 0,
  loadTimeTrend: 0,
  requestSuccessRate: 0,
  requestRateTrend: 0
})

// 最近错误
const recentErrors = ref<any[]>([])

// 对话框状态
const showCreateProjectDialog = ref(false)
const showProjectSelector = ref(false)

// 新项目表单
const newProject = reactive({
  name: '',
  description: ''
})

// 打开创建项目对话框
const openCreateProjectDialog = () => {
  showCreateProjectDialog.value = true
  showProjectSelector.value = false
}

// 打开项目选择器
const openProjectSelector = () => {
  showProjectSelector.value = true
  showCreateProjectDialog.value = false
}

// 创建项目
const createProject = async () => {
  const project = await projectStore.createProject(newProject.name, newProject.description)
  if (project) {
    projectStore.setCurrentProject(project)
    showCreateProjectDialog.value = false
    newProject.name = ''
    newProject.description = ''
  }
}

// 选择项目
const selectProject = (project: any) => {
  projectStore.setCurrentProject(project)
  showProjectSelector.value = false
}

// 查看错误详情
const viewErrorDetail = (id: number) => {
  router.push(`/errors/${id}`)
}

// 获取统计数据
const fetchStats = async () => {
  // 模拟数据，实际应该从 API 获取
  stats.errorCount = 24
  stats.errorTrend = -5
  stats.avgLoadTime = 320
  stats.loadTimeTrend = -12
  stats.requestSuccessRate = 99.8
  stats.requestRateTrend = 0.2

  // 模拟最近错误数据
  recentErrors.value = [
    {
      id: 1,
      message: 'TypeError: Cannot read property "length" of undefined',
      time: '2023-12-01 14:32:45'
    },
    {
      id: 2,
      message: 'SyntaxError: Unexpected token in JSON at position 0',
      time: '2023-12-01 12:15:22'
    },
    {
      id: 3,
      message: 'ReferenceError: someVariable is not defined',
      time: '2023-11-30 23:45:11'
    }
  ]
}

onMounted(() => {
  fetchStats()
})
</script>

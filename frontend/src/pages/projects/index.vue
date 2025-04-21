<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useProjectStore } from '@/stores/project'
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useClipboard } from '@vueuse/core'
import type { ModelProject } from '@/types'
import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle, 
  DialogDescription,
  DialogFooter 
} from '@/components/ui/dialog'
import { showToast } from '@/utils/toast'
import { useRouter } from 'vue-router'


const router = useRouter()

// 初始化项目store
const projectStore = useProjectStore()
const loading = computed(() => projectStore.loading)
const projects = computed(() => projectStore.projects)
const error = computed(() => projectStore.error)

// 搜索相关
const searchKeyword = ref('')
const filteredProjects = computed(() => {
  if (!searchKeyword.value) {
    return projects.value
  }
  
  const keyword = searchKeyword.value.toLowerCase()
  return projects.value.filter(project => 
    project.name?.toLowerCase().includes(keyword) || 
    project.description?.toLowerCase().includes(keyword)
  )
})

// 模态框相关
const createModalVisible = ref(false)
const currentProject = ref<ModelProject | null>(null)
const form = ref({
  name: '',
  description: ''
})
const formRef = ref()

// 确认对话框相关
const confirmDialogVisible = ref(false) 
const confirmDialogTitle = ref('')
const confirmDialogMessage = ref('')
const confirmDialogAction = ref<() => Promise<void>>(() => Promise.resolve())

// 提示信息相关
const showSuccessMessage = (message: string) => {
  showToast(message, 'success')
}

const showErrorMessage = (message: string) => {
  showToast(message, 'error')
}

// 复制功能
const { copy } = useClipboard()

// 创建项目
const showCreateModal = () => {
  form.value = { name: '', description: '' }
  createModalVisible.value = true
}

const handleCreateProject = async () => {
  if (!form.value.name) {
    showErrorMessage('项目名称不能为空')
    return
  }
  
  try {
    const result = await projectStore.createProject(form.value.name, form.value.description || '')
    if (result) {
      showSuccessMessage('项目创建成功')
      createModalVisible.value = false
    } else {
      showErrorMessage('项目创建失败')
    }
  } catch (err: any) {
    showErrorMessage(err?.message || '项目创建失败')
  }
}

// 编辑项目
const showEditModal = (project: ModelProject) => {
  router.push(`/projects/${project.id}/settings`)
}

// 删除项目
const showDeleteConfirm = (id: number) => {
  confirmDialogTitle.value = '删除项目'
  confirmDialogMessage.value = '确定要删除此项目吗？此操作不可撤销，所有相关的监控数据将被永久删除。'
  confirmDialogAction.value = async () => {
    try {
      const success = await projectStore.deleteProject(id)
      if (success) {
        showSuccessMessage('项目删除成功')
      } else {
        showErrorMessage('项目删除失败')
      }
    } catch (err: any) {
      showErrorMessage(err?.message || '项目删除失败')
    } finally {
      confirmDialogVisible.value = false
    }
  }
  confirmDialogVisible.value = true
}

// 重置AppKey
const showResetAppKeyConfirm = (id: number) => {
  confirmDialogTitle.value = '重置App Key'
  confirmDialogMessage.value = '确定要重置此项目的 App Key 吗？重置后，使用旧 App Key 的应用将无法上报数据。'
  confirmDialogAction.value = async () => {
    try {
      // 待实现
      showSuccessMessage('重置AppKey功能尚未实现')
    } catch (err: any) {
      showErrorMessage(err?.message || '重置AppKey失败')
    } finally {
      confirmDialogVisible.value = false
    }
  }
  confirmDialogVisible.value = true
}

// 复制AppKey
const handleCopyAppKey = (appKey: string) => {
  copy(appKey)
  showSuccessMessage('AppKey 已复制到剪贴板')
}

// 初始化
onMounted(() => {
  projectStore.fetchProjects()
})

// AppKey 格式化显示 (显示前4位和后4位，中间用***代替)
const formatAppKey = (appKey?: string) => {
  if (!appKey) return '无'
  if (appKey.length <= 8) return appKey
  return `${appKey.slice(0, 4)}****${appKey.slice(-4)}`
}
</script>

<template>
  <div class="space-y-6 p-4">
    <!-- 顶部操作区域 -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <h1 class="text-3xl font-bold tracking-tight">项目管理</h1>
      
      <div class="flex items-center space-x-2">
        <Input
          v-model="searchKeyword"
          placeholder="搜索项目"
          class="w-full sm:w-[250px]"
        />
        
        <Button 
          @click="showCreateModal"
          :disabled="loading"
        >
          创建项目
        </Button>
      </div>
    </div>

    <!-- 项目列表 -->
    <div v-if="error" class="bg-destructive/15 text-destructive p-4 rounded-md mb-4">
      {{ error }}
    </div>
    
    <div class="projects-list min-h-[300px]">
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin">&#8635;</div>
      </div>
      
      <div v-else-if="filteredProjects.length === 0" class="flex flex-col items-center justify-center py-12">
        <p class="text-muted-foreground">暂无项目</p>
      </div>
      
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-4">
        <Card 
          v-for="project in filteredProjects" 
          :key="project.id" 
          class="shadow-sm hover:shadow transition-shadow duration-200 flex flex-col overflow-hidden"
        >
          <CardHeader class="pb-2">
            <div class="flex justify-between items-start">
              <CardTitle class="truncate max-w-[80%]">{{ project.name }}</CardTitle>
              <span class="text-xs text-muted-foreground whitespace-nowrap">{{ project.createdAt?.split('T')[0] }}</span>
            </div>
          </CardHeader>
          
          <CardContent class="flex-grow">
            <p class="text-sm text-muted-foreground line-clamp-2 mb-4 break-words">
              {{ project.description || '无项目描述' }}
            </p>
            
            <div class="text-xs space-y-1">
              <div class="flex items-center justify-between gap-2">
                <span class="text-muted-foreground whitespace-nowrap">项目ID:</span> 
                <span class="text-right">{{ project.id }}</span>
              </div>
              <div class="flex items-center justify-between gap-2">
                <span class="text-muted-foreground whitespace-nowrap">App Key:</span>
                <div class="flex items-center overflow-hidden">
                  <span class="truncate">{{ formatAppKey(project.appKey) }}</span>
                  <Button 
                    variant="ghost" 
                    class="h-6 w-6 p-0 ml-1 flex-shrink-0"
                    @click="project.appKey && handleCopyAppKey(project.appKey)"
                  >
                    <span class="sr-only">复制</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-copy"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg>
                  </Button>
                </div>
              </div>
              <div class="flex items-center justify-between gap-2">
                <span class="text-muted-foreground whitespace-nowrap">错误数量:</span>
                <span class="text-right">今日: 0 / 总计: 0</span>
              </div>
            </div>
          </CardContent>
          
          <CardFooter class="flex flex-col sm:flex-row justify-between pt-0 gap-2">
            <Button variant="default" size="sm" class="w-full sm:w-auto">查看详情</Button>
            <div class="flex gap-1 justify-end w-full sm:w-auto">
              <Button 
                variant="ghost" 
                size="icon" 
                class="h-8 w-8" 
                @click="showEditModal(project)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-pencil"><path d="M17 3a2.85 2.85 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/></svg>
                <span class="sr-only">编辑</span>
              </Button>
              
              <Button 
                variant="ghost" 
                size="icon" 
                class="h-8 w-8 text-destructive hover:text-destructive" 
                @click="project.id && showDeleteConfirm(project.id)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-trash-2"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                <span class="sr-only">删除</span>
              </Button>
              
              <Button 
                variant="ghost" 
                size="icon" 
                class="h-8 w-8 text-warning hover:text-warning" 
                @click="project.id && showResetAppKeyConfirm(project.id)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-refresh-cw"><path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/><path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/><path d="M3 21v-5h5"/></svg>
                <span class="sr-only">重置</span>
              </Button>
            </div>
          </CardFooter>
        </Card>
      </div>
    </div>
    
    <!-- 创建项目弹窗 -->
    <Dialog :open="createModalVisible" @update:open="createModalVisible = $event">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>创建新项目</DialogTitle>
        </DialogHeader>
        
        <div class="space-y-4">
          <div class="space-y-2">
            <label for="name" class="text-sm font-medium">项目名称</label>
            <Input id="name" v-model="form.name" placeholder="请输入项目名称" />
          </div>
          
          <div class="space-y-2">
            <label for="description" class="text-sm font-medium">项目描述</label>
            <textarea 
              id="description"
              v-model="form.description" 
              placeholder="请输入项目描述（选填）" 
              rows="4"
              class="w-full min-h-[80px] rounded-md border bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            ></textarea>
          </div>
        </div>
        
        <DialogFooter class="mt-6">
          <Button variant="outline" @click="createModalVisible = false">取消</Button>
          <Button @click="handleCreateProject" :disabled="loading">创建</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
    
    <!-- 确认操作弹窗 -->
    <Dialog :open="confirmDialogVisible" @update:open="confirmDialogVisible = $event">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{{ confirmDialogTitle }}</DialogTitle>
          <DialogDescription>
            {{ confirmDialogMessage }}
          </DialogDescription>
        </DialogHeader>
        
        <DialogFooter class="mt-6">
          <Button variant="outline" @click="confirmDialogVisible = false">取消</Button>
          <Button 
            :variant="confirmDialogTitle.includes('删除') ? 'destructive' : 'default'"
            @click="confirmDialogAction"
          >
            确认
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<style scoped>
.text-warning {
  color: #f59e0b;
}

.hover\:text-warning:hover {
  color: #d97706;
}
</style> 
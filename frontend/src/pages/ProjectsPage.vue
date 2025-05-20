<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold tracking-tight">项目管理</h1>
      <div class="flex items-center space-x-2">
        <Button @click="openCreateDialog">
          创建项目
        </Button>
      </div>
    </div>

    <!-- 项目列表 -->
    <div v-if="projectStore.loading" class="flex justify-center p-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
    </div>
    <div v-else-if="!projectStore.hasProjects" class="flex flex-col items-center justify-center p-8 border rounded-lg">
      <p class="text-lg text-center text-muted-foreground mb-4">
        暂无项目，请创建一个新项目
      </p>
      <Button @click="openCreateDialog">
        创建项目
      </Button>
    </div>
    <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Card 
        v-for="project in projectStore.projects" 
        :key="project.id"
        :class="{'border-primary': projectStore.currentProject?.id === project.id}"
      >
        <CardHeader>
          <CardTitle>{{ project.name }}</CardTitle>
          <CardDescription>{{ project.description || '暂无描述' }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-2">
            <div class="flex items-center">
              <span class="text-sm text-muted-foreground mr-2">项目ID:</span>
              <span>{{ project.id }}</span>
            </div>
            <div class="flex items-center">
              <span class="text-sm text-muted-foreground mr-2">创建时间:</span>
              <span>{{ formatTime(project.createdAt) }}</span>
            </div>
          </div>
        </CardContent>
        <CardFooter class="flex justify-between">
          <Button 
            variant="outline" 
            size="sm"
            @click="selectProject(project)"
          >
            {{ projectStore.currentProject?.id === project.id ? '当前项目' : '选择' }}
          </Button>
          <div class="space-x-2">
            <Button 
              variant="ghost" 
              size="sm"
              @click="openEditDialog(project)"
            >
              编辑
            </Button>
            <Button 
              variant="destructive" 
              size="sm"
              @click="openDeleteDialog(project)"
            >
              删除
            </Button>
          </div>
        </CardFooter>
      </Card>
    </div>

    <!-- 创建/编辑项目对话框 -->
    <Dialog :open="showDialog" @update:open="showDialog = $event">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ isEditing ? '编辑项目' : '创建项目' }}</DialogTitle>
          <DialogDescription>
            {{ isEditing ? '修改项目信息' : '创建一个新的项目' }}
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="name">项目名称</Label>
            <Input id="name" v-model="projectForm.name" placeholder="请输入项目名称" />
          </div>
          <div class="space-y-2">
            <Label for="description">项目描述</Label>
            <Textarea id="description" v-model="projectForm.description" placeholder="请输入项目描述" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showDialog = false">取消</Button>
          <Button :disabled="!projectForm.name" @click="saveProject">保存</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- 删除项目确认对话框 -->
    <AlertDialog :open="showDeleteDialog" @update:open="showDeleteDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除</AlertDialogTitle>
          <AlertDialogDescription>
            您确定要删除项目 "{{ projectToDelete?.name }}" 吗？此操作不可撤销，项目下的所有数据将被永久删除。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="showDeleteDialog = false">取消</AlertDialogCancel>
          <AlertDialogAction @click="deleteProject">确认删除</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useProjectStore } from '@/stores/project'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { toast } from '@/utils/toast'
import type { ServiceProject } from '@/types/gen/service/Project'

const projectStore = useProjectStore()

// 对话框状态
const showDialog = ref(false)
const showDeleteDialog = ref(false)
const isEditing = ref(false)
const projectToDelete = ref<ServiceProject | null>(null)

// 项目表单
const projectForm = ref({
  id: 0,
  name: '',
  description: ''
})

// 获取项目列表
const fetchProjects = async () => {
  await projectStore.fetchProjects()
}

// 打开创建项目对话框
const openCreateDialog = () => {
  isEditing.value = false
  projectForm.value = {
    id: 0,
    name: '',
    description: ''
  }
  showDialog.value = true
}

// 打开编辑项目对话框
const openEditDialog = (project: ServiceProject) => {
  isEditing.value = true
  projectForm.value = {
    id: project.id,
    name: project.name,
    description: project.description || ''
  }
  showDialog.value = true
}

// 打开删除项目确认对话框
const openDeleteDialog = (project: ServiceProject) => {
  projectToDelete.value = project
  showDeleteDialog.value = true
}

// 保存项目
const saveProject = async () => {
  if (!projectForm.value.name) {
    toast.error('请输入项目名称')
    return
  }

  try {
    if (isEditing.value) {
      await projectStore.updateProject(projectForm.value.id, {
        name: projectForm.value.name,
        description: projectForm.value.description
      })
      toast.success('项目更新成功')
    } else {
      await projectStore.createProject({
        name: projectForm.value.name,
        description: projectForm.value.description
      })
      toast.success('项目创建成功')
    }
    showDialog.value = false
  } catch (error) {
    console.error('保存项目失败', error)
  }
}

// 删除项目
const deleteProject = async () => {
  if (!projectToDelete.value) return

  try {
    await projectStore.deleteProject(projectToDelete.value.id)
    toast.success('项目删除成功')
    showDeleteDialog.value = false
    projectToDelete.value = null
  } catch (error) {
    console.error('删除项目失败', error)
  }
}

// 选择项目
const selectProject = (project: ServiceProject) => {
  projectStore.setCurrentProject(project)
  toast.success(`已选择项目: ${project.name}`)
}

// 格式化时间戳
const formatTime = (timestamp: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString()
}

onMounted(() => {
  fetchProjects()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold tracking-tight">设置</h1>
    </div>

    <Tabs default-value="account" class="w-full">
      <TabsList>
        <TabsTrigger value="account">账号设置</TabsTrigger>
        <TabsTrigger value="sdk">SDK配置</TabsTrigger>
      </TabsList>
      
      <!-- 账号设置 -->
      <TabsContent value="account">
        <Card>
          <CardHeader>
            <CardTitle>账号信息</CardTitle>
            <CardDescription>
              管理您的账号信息
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div class="space-y-4">
              <div class="space-y-2">
                <Label for="username">用户名</Label>
                <Input id="username" v-model="userForm.username" placeholder="用户名" />
              </div>
              <div class="space-y-2">
                <Label for="email">邮箱</Label>
                <Input id="email" v-model="userForm.email" placeholder="邮箱" type="email" />
              </div>
              <Button @click="updateUserInfo" :disabled="loading">
                {{ loading ? '保存中...' : '保存' }}
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card class="mt-6">
          <CardHeader>
            <CardTitle>修改密码</CardTitle>
            <CardDescription>
              更新您的账号密码
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div class="space-y-4">
              <div class="space-y-2">
                <Label for="currentPassword">当前密码</Label>
                <Input id="currentPassword" v-model="passwordForm.currentPassword" placeholder="当前密码" type="password" />
              </div>
              <div class="space-y-2">
                <Label for="newPassword">新密码</Label>
                <Input id="newPassword" v-model="passwordForm.newPassword" placeholder="新密码" type="password" />
              </div>
              <div class="space-y-2">
                <Label for="confirmPassword">确认新密码</Label>
                <Input id="confirmPassword" v-model="passwordForm.confirmPassword" placeholder="确认新密码" type="password" />
              </div>
              <Button @click="updatePassword" :disabled="loading">
                {{ loading ? '更新中...' : '更新密码' }}
              </Button>
            </div>
          </CardContent>
        </Card>
      </TabsContent>
      
      <!-- SDK配置 -->
      <TabsContent value="sdk">
        <Card>
          <CardHeader>
            <CardTitle>SDK配置</CardTitle>
            <CardDescription>
              配置和获取SDK接入信息
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div v-if="!projectStore.currentProject" class="flex flex-col items-center justify-center p-8 border rounded-lg">
              <p class="text-lg text-center text-muted-foreground mb-4">
                请先创建或选择一个项目
              </p>
              <Button @click="router.push('/projects')">
                前往项目管理
              </Button>
            </div>
            <div v-else class="space-y-6">
              <div>
                <h3 class="text-lg font-medium mb-2">项目信息</h3>
                <div class="grid gap-4 md:grid-cols-2">
                  <div>
                    <Label>项目ID</Label>
                    <div class="flex items-center mt-1">
                      <Input readonly :value="projectStore.currentProject.id.toString()" />
                      <Button variant="outline" class="ml-2" @click="copyToClipboard(projectStore.currentProject.id.toString())">
                        复制
                      </Button>
                    </div>
                  </div>
                  <div>
                    <Label>项目名称</Label>
                    <div class="mt-1">
                      <Input readonly :value="projectStore.currentProject.name" />
                    </div>
                  </div>
                </div>
              </div>

              <div>
                <h3 class="text-lg font-medium mb-2">SDK接入代码</h3>
                <div class="bg-muted p-4 rounded-md">
                  <pre class="text-sm overflow-auto"><code>// 安装SDK
npm install web-tracing

// 在main.js/ts中引入并初始化
import { init } from 'web-tracing'

init({
  dsn: 'http://your-domain.com/trackweb',
  apikey: '{{ projectStore.currentProject.id }}',
  userId: 'user-unique-id', // 用户唯一标识，可选
  debug: false, // 是否开启调试模式，可选
  pv: true, // 是否监控页面访问，可选
  performance: true, // 是否监控页面性能，可选
  error: true, // 是否监控页面错误，可选
  event: true, // 是否监控用户行为，可选
})</code></pre>
                </div>
                <Button variant="outline" class="mt-2" @click="copyToClipboard(sdkCode)">
                  复制代码
                </Button>
              </div>

              <div>
                <h3 class="text-lg font-medium mb-2">自定义配置</h3>
                <div class="space-y-4">
                  <div class="flex items-center space-x-2">
                    <Checkbox id="enablePV" v-model:checked="sdkConfig.pv" />
                    <Label for="enablePV">启用页面访问监控</Label>
                  </div>
                  <div class="flex items-center space-x-2">
                    <Checkbox id="enablePerformance" v-model:checked="sdkConfig.performance" />
                    <Label for="enablePerformance">启用性能监控</Label>
                  </div>
                  <div class="flex items-center space-x-2">
                    <Checkbox id="enableError" v-model:checked="sdkConfig.error" />
                    <Label for="enableError">启用错误监控</Label>
                  </div>
                  <div class="flex items-center space-x-2">
                    <Checkbox id="enableEvent" v-model:checked="sdkConfig.event" />
                    <Label for="enableEvent">启用用户行为监控</Label>
                  </div>
                  <Button @click="updateSdkConfig">
                    保存配置
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/project'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { toast } from '@/utils/toast'

const router = useRouter()
const authStore = useAuthStore()
const projectStore = useProjectStore()
const loading = ref(false)

// 用户信息表单
const userForm = ref({
  username: authStore.user?.username || '',
  email: authStore.user?.email || ''
})

// 密码表单
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// SDK配置
const sdkConfig = ref({
  pv: true,
  performance: true,
  error: true,
  event: true
})

// SDK接入代码
const sdkCode = computed(() => {
  if (!projectStore.currentProject) return ''
  
  return `// 安装SDK
npm install web-tracing

// 在main.js/ts中引入并初始化
import { init } from 'web-tracing'

init({
  dsn: 'http://your-domain.com/trackweb',
  apikey: '${projectStore.currentProject.id}',
  userId: 'user-unique-id', // 用户唯一标识，可选
  debug: false, // 是否开启调试模式，可选
  pv: ${sdkConfig.value.pv}, // 是否监控页面访问，可选
  performance: ${sdkConfig.value.performance}, // 是否监控页面性能，可选
  error: ${sdkConfig.value.error}, // 是否监控页面错误，可选
  event: ${sdkConfig.value.event}, // 是否监控用户行为，可选
})`
})

// 更新用户信息
const updateUserInfo = async () => {
  loading.value = true
  
  try {
    // 这里应该调用API更新用户信息
    // 由于后端未提供此接口，这里只做模拟
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    toast.success('用户信息更新成功')
  } catch (error) {
    toast.error('用户信息更新失败')
  } finally {
    loading.value = false
  }
}

// 更新密码
const updatePassword = async () => {
  if (!passwordForm.value.currentPassword) {
    toast.error('请输入当前密码')
    return
  }
  
  if (!passwordForm.value.newPassword) {
    toast.error('请输入新密码')
    return
  }
  
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    toast.error('两次输入的新密码不一致')
    return
  }
  
  loading.value = true
  
  try {
    // 这里应该调用API更新密码
    // 由于后端未提供此接口，这里只做模拟
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
    
    toast.success('密码更新成功')
  } catch (error) {
    toast.error('密码更新失败')
  } finally {
    loading.value = false
  }
}

// 更新SDK配置
const updateSdkConfig = async () => {
  try {
    // 这里应该调用API更新SDK配置
    // 由于后端未提供此接口，这里只做模拟
    await new Promise(resolve => setTimeout(resolve, 500))
    
    toast.success('SDK配置更新成功')
  } catch (error) {
    toast.error('SDK配置更新失败')
  }
}

// 复制到剪贴板
const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
    .then(() => {
      toast.success('复制成功')
    })
    .catch(() => {
      toast.error('复制失败')
    })
}
</script>

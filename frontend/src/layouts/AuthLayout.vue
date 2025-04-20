<template>
  <div class="min-h-screen flex items-center justify-center bg-background">
    <div class="w-full max-w-md p-6">
      <div class="flex flex-col space-y-2 text-center mb-8">
        <h1 class="text-2xl font-semibold tracking-tight">Web Tracing Admin</h1>
        <p class="text-sm text-muted-foreground">
          错误监控后台管理系统
        </p>
      </div>
      <router-view />
    </div>
  </div>
</template>

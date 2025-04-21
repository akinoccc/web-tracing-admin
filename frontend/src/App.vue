<template>
  <router-view />
  <Toaster position="top-right" richColors closeButton />
</template>

<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'
</script>

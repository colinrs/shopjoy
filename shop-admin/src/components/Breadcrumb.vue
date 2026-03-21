<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index" :to="item.path">
      {{ item.title }}
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

interface BreadcrumbItem {
  title: string
  path?: string
}

const breadcrumbs = computed<BreadcrumbItem[]>(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  const result: BreadcrumbItem[] = []

  // Add home
  result.push({ title: '首页', path: '/dashboard' })

  // Add matched routes
  matched.forEach(item => {
    if (item.meta?.title) {
      result.push({
        title: item.meta.title as string,
        path: item.path !== route.path ? item.path : undefined
      })
    }
  })

  return result
})
</script>

<style scoped>
.el-breadcrumb {
  font-size: 14px;
}
</style>
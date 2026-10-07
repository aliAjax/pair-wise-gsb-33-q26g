<template>
  <el-form inline class="search-filter">
    <slot name="filters" />
    <el-form-item>
      <el-input v-model="keyword" placeholder="输入关键词搜索" clearable @keyup.enter="emit('search', keyword)" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="emit('search', keyword)">搜索</el-button>
      <el-button @click="emit('reset')">重置</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const keyword = ref('')
const emit = defineEmits<{ (e: 'search', keyword: string): void; (e: 'reset'): void }>()
</script>

<style scoped>
.search-filter { padding: 8px 0; }
</style>

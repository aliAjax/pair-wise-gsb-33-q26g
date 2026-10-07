<template>
  <el-card class="disease-card" shadow="hover">
    <div class="name">{{ pest.name }}</div>
    <div class="section"><b>症状：</b>{{ pest.symptoms }}</div>
    <div class="section"><b>病因：</b>{{ pest.cause }}</div>
    <div class="section"><b>防治：</b>{{ pest.treatment }}</div>
    <div class="section medicine"><b>推荐用药：</b>{{ pest.recommended_medicine || '-' }}</div>
  </el-card>
</template>

<script setup lang="ts">
import type { DiseasePest } from '@/types/api'
defineProps<{ pest: DiseasePest }>()
</script>

<style scoped>
.disease-card { margin-bottom: 12px; }
.name { font-weight: 700; font-size: 16px; margin-bottom: 8px; color: #b3541e; }
.section { margin: 6px 0; line-height: 1.5; }
.medicine { color: #3c8d5c; }
</style>

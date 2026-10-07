<template>
  <el-card class="article-card" shadow="hover" @click="$router.push(`/articles/${article.id}`)">
    <el-image v-if="article.cover" :src="article.cover" fit="cover" class="cover" lazy>
      <template #error><div class="img-placeholder">📖</div></template>
    </el-image>
    <div class="body">
      <el-tag size="small" type="success">{{ CareTopicTagMap[article.topic_tag] }}</el-tag>
      <div class="title">{{ article.title }}</div>
      <div class="meta">{{ formatDate(article.created_at) }} · 阅读 {{ article.view_count }}</div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { CareTopicTagMap, type CareArticle } from '@/constants/article'
import { formatDate } from '@/utils/dateFormat'

defineProps<{ article: CareArticle }>()
</script>

<style scoped>
.article-card { cursor: pointer; }
.cover { width: 100%; height: 140px; border-radius: 6px; }
.img-placeholder { height: 140px; display: flex; align-items: center; justify-content: center; font-size: 40px; background: #f5f0e8; }
.body { padding-top: 10px; }
.title { font-weight: 700; margin-top: 6px; line-height: 1.4; }
.meta { color: #888; font-size: 12px; margin-top: 6px; }
</style>

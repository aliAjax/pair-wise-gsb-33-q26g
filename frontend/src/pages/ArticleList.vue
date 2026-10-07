<template>
  <div class="page">
    <h1>养护知识文章</h1>
    <div class="topic-tabs">
      <el-radio-group v-model="topicTag" @change="onTopicChange">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button v-for="(label, value) in CareTopicTagMap" :key="value" :value="value">{{ label }}</el-radio-button>
      </el-radio-group>
    </div>
    <SearchFilter @search="onSearch" @reset="onReset" />
    <el-row :gutter="16">
      <el-col v-for="a in articles" :key="a.id" :xs="12" :sm="8" :md="6">
        <CareArticleCard :article="a" />
      </el-col>
    </el-row>
    <el-empty v-if="!articles.length" description="暂无文章" />
    <el-pagination v-if="total > 0" layout="prev, pager, next" :total="total" :page-size="pageSize" :current-page="page" @current-change="onPage" class="pager" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import CareArticleCard from '@/components/common/CareArticleCard.vue'
import SearchFilter from '@/components/common/SearchFilter.vue'
import { useArticleStore } from '@/stores/articleStore'
import { CareTopicTagMap } from '@/constants/article'

const store = useArticleStore()
const topicTag = ref('')
const keyword = ref('')
const page = ref(1)
const pageSize = 8
const articles = computed(() => store.articles)
const total = computed(() => store.total)

onMounted(() => load())

async function load() {
  await store.load({ page: page.value, page_size: pageSize, topic_tag: topicTag.value, keyword: keyword.value })
}
function onTopicChange() {
  page.value = 1
  load()
}
function onSearch(kw: string) {
  keyword.value = kw
  page.value = 1
  load()
}
function onReset() {
  topicTag.value = ''
  keyword.value = ''
  page.value = 1
  load()
}
function onPage(p: number) {
  page.value = p
  load()
}
</script>

<style scoped>
.page { max-width: 1200px; margin: 0 auto; }
.topic-tabs { margin: 12px 0; }
.pager { margin-top: 20px; justify-content: center; }
</style>

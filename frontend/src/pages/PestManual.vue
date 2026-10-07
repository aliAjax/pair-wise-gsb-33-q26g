<template>
  <div class="page">
    <h1>病虫害诊断手册</h1>
    <SearchFilter @search="onSearch" @reset="onReset" />
    <el-row :gutter="16">
      <el-col v-for="p in pests" :key="p.id" :xs="24" :sm="12" :md="8">
        <DiseaseCard :pest="p" />
      </el-col>
    </el-row>
    <el-empty v-if="!pests.length" description="未找到相关病虫害" />
    <el-pagination v-if="total > 0" layout="prev, pager, next" :total="total" :page-size="pageSize" :current-page="page" @current-change="onPage" class="pager" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import SearchFilter from '@/components/common/SearchFilter.vue'
import DiseaseCard from '@/components/common/DiseaseCard.vue'
import { listPests } from '@/api/pest'
import type { DiseasePest } from '@/types/api'

const pests = ref<DiseasePest[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 9
const keyword = ref('')

onMounted(() => load())

async function load() {
  const res = await listPests({ page: page.value, page_size: pageSize, keyword: keyword.value })
  pests.value = res.list
  total.value = res.total
}
function onSearch(kw: string) {
  keyword.value = kw
  page.value = 1
  load()
}
function onReset() {
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
.pager { margin-top: 20px; justify-content: center; }
</style>

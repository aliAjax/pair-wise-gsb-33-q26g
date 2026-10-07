<template>
  <div class="home">
    <section class="hero">
      <h1>植物养护知识百科平台</h1>
      <p class="sub">品种库 · 养护文章 · 病虫害防治 · 季节日历 · 问答社区</p>
      <el-input v-model="keyword" size="large" placeholder="搜索植物品种 / 养护文章 / 病虫害" class="search" @keyup.enter="doSearch">
        <template #append><el-button @click="doSearch">搜索</el-button></template>
      </el-input>
    </section>

    <section class="season-banner">
      <el-alert :title="`当季养护重点：${seasonTask}`" type="success" :closable="false" show-icon />
    </section>

    <section>
      <h2>🔥 热门品种</h2>
      <el-row :gutter="16">
        <el-col v-for="p in hotPlants" :key="p.id" :xs="12" :sm="8" :md="6">
          <PlantCard :plant="p" />
        </el-col>
      </el-row>
    </section>

    <section>
      <h2>📖 最新养护文章</h2>
      <el-row :gutter="16">
        <el-col v-for="a in latestArticles" :key="a.id" :xs="12" :sm="8" :md="6">
          <CareArticleCard :article="a" />
        </el-col>
      </el-row>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import PlantCard from '@/components/common/PlantCard.vue'
import CareArticleCard from '@/components/common/CareArticleCard.vue'
import { currentSeasonTask } from '@/utils/season'
import type { PlantSpecies } from '@/constants/plant'
import type { CareArticle } from '@/constants/article'

const router = useRouter()
const keyword = ref('')
const seasonTask = ref(currentSeasonTask())
const hotPlants = ref<PlantSpecies[]>([])
const latestArticles = ref<CareArticle[]>([])

onMounted(async () => {
  const res = await axios.get('/api/v1/home/overview')
  hotPlants.value = res.data.data.hot_plants || []
  latestArticles.value = res.data.data.latest_articles || []
  seasonTask.value = res.data.data.season_task || seasonTask.value
})

function doSearch() {
  if (keyword.value) {
    router.push({ path: '/plants', query: { keyword: keyword.value } })
  }
}
</script>

<style scoped>
.home { max-width: 1200px; margin: 0 auto; }
.hero { text-align: center; padding: 40px 0 20px; }
.hero h1 { color: #2c6e49; font-size: 32px; }
.sub { color: #888; }
.search { max-width: 560px; margin: 16px auto; }
.season-banner { margin: 8px 0 24px; }
h2 { color: #333; margin: 24px 0 16px; }
</style>

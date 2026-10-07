<template>
  <div class="page">
    <h1>植物品种库</h1>
    <SearchFilter @search="onSearch" @reset="onReset">
      <template #filters>
        <el-form-item label="类型">
          <el-select v-model="type" placeholder="全部类型" clearable style="width: 140px" @change="load">
            <el-option v-for="(label, value) in PlantTypeMap" :key="value" :label="label" :value="value" />
          </el-select>
        </el-form-item>
        <el-form-item label="科">
          <el-input v-model="family" placeholder="如：蔷薇科" clearable style="width: 160px" @change="load" />
        </el-form-item>
      </template>
    </SearchFilter>
    <el-row :gutter="16">
      <el-col v-for="p in plants" :key="p.id" :xs="12" :sm="8" :md="6">
        <PlantCard :plant="p" />
      </el-col>
    </el-row>
    <el-empty v-if="!plants.length && !loading" description="暂无品种数据" />
    <el-pagination v-if="total > 0" layout="prev, pager, next" :total="total" :page-size="pageSize" :current-page="page" @current-change="onPage" class="pager" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import PlantCard from '@/components/common/PlantCard.vue'
import SearchFilter from '@/components/common/SearchFilter.vue'
import { usePlantStore } from '@/stores/plantStore'
import { PlantTypeMap } from '@/constants/plant'

const store = usePlantStore()
const route = useRoute()
const type = ref('')
const family = ref('')
const keyword = ref('')
const page = ref(1)
const pageSize = 12
const loading = ref(false)

const plants = computed(() => store.plants)
const total = computed(() => store.total)

onMounted(async () => {
  keyword.value = (route.query.keyword as string) || ''
  await load()
})

async function load() {
  loading.value = true
  try {
    await store.load({ page: page.value, page_size: pageSize, type: type.value, family: family.value, keyword: keyword.value })
  } finally {
    loading.value = false
  }
}

function onSearch(kw: string) {
  keyword.value = kw
  page.value = 1
  load()
}
function onReset() {
  type.value = ''
  family.value = ''
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

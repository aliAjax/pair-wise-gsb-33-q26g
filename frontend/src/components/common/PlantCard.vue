<template>
  <el-card class="plant-card" shadow="hover" @click="$router.push(`/plants/${plant.id}`)">
    <el-image :src="cover" fit="cover" class="cover" lazy>
      <template #error><div class="img-placeholder">🌿</div></template>
    </el-image>
    <div class="info">
      <div class="name">{{ plant.name }} <el-tag size="small" :type="tagType">{{ PlantTypeMap[plant.type] }}</el-tag></div>
      <div class="meta">{{ plant.family }} · {{ plant.genus }}</div>
      <div class="temp">{{ formatTemp(plant.temp_min, plant.temp_max) }}</div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { PlantTypeMap, parseImages, type PlantSpecies } from '@/constants/plant'

const props = defineProps<{ plant: PlantSpecies }>()
const cover = computed(() => parseImages(props.plant.image_urls)[0] || '')
const tagType = computed(() => {
  const map: Record<string, 'success' | 'warning' | 'info' | 'primary'> = {
    flower: 'primary',
    foliage: 'success',
    succulent: 'warning',
    aquatic: 'info',
  }
  return map[props.plant.type] || 'info'
})
function formatTemp(min: number, max: number): string {
  if (!min && !max) return '温度不限'
  return `${min}°C ~ ${max}°C`
}
</script>

<style scoped>
.plant-card { cursor: pointer; }
.cover { width: 100%; height: 160px; border-radius: 6px; }
.img-placeholder { height: 160px; display: flex; align-items: center; justify-content: center; font-size: 40px; background: #eef5ec; }
.info { padding-top: 10px; }
.name { font-weight: 700; display: flex; align-items: center; gap: 6px; }
.meta { color: #888; font-size: 12px; margin-top: 4px; }
.temp { color: #3c8d5c; font-size: 13px; margin-top: 4px; }
</style>

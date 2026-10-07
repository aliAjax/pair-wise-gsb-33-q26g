import { defineStore } from 'pinia'
import { ref } from 'vue'
import { listPlants } from '@/api/plant'
import type { PlantSpecies } from '@/constants/plant'

export const usePlantStore = defineStore('plant', () => {
  const plants = ref<PlantSpecies[]>([])
  const total = ref(0)

  async function load(params: { page?: number; page_size?: number; type?: string; family?: string; keyword?: string } = {}) {
    const res = await listPlants(params)
    plants.value = res.list
    total.value = res.total
  }

  return { plants, total, load }
})

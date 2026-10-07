export type PlantType = 'flower' | 'foliage' | 'succulent' | 'aquatic'

export const PlantTypeMap: Record<PlantType, string> = {
  flower: '观花',
  foliage: '观叶',
  succulent: '多肉',
  aquatic: '水生',
}

export const PLANT_TYPES = Object.keys(PlantTypeMap) as PlantType[]

export interface PlantSpecies {
  id: number
  family: string
  genus: string
  name: string
  alias: string
  type: PlantType
  origin: string
  temp_min: number
  temp_max: number
  light_requirement: string
  water_frequency: string
  description: string
  image_urls: string
  created_at: string
}

export function parseImages(raw: string): string[] {
  try {
    const arr = JSON.parse(raw || '[]')
    return Array.isArray(arr) ? arr : []
  } catch {
    return []
  }
}

import { ref } from 'vue'
import { addFavorite, removeFavorite } from '@/api/favorite'
import type { FavoriteTargetType } from '@/constants/favorite'

export function useFavorite(targetType: FavoriteTargetType) {
  const favorited = ref(false)
  const loading = ref(false)

  async function toggle(targetId: number): Promise<boolean> {
    loading.value = true
    try {
      if (favorited.value) {
        await removeFavorite(targetType, targetId)
        favorited.value = false
        return false
      }
      await addFavorite(targetType, targetId)
      favorited.value = true
      return true
    } finally {
      loading.value = false
    }
  }

  return { favorited, loading, toggle }
}

import request from '@/utils/request'
import type { Favorite, FavoriteTargetType } from '@/constants/favorite'

export function listFavorites(targetType?: FavoriteTargetType | '') {
  return request.get<never, Favorite[]>('/favorites', { params: { target_type: targetType } })
}

export function addFavorite(targetType: FavoriteTargetType, targetId: number) {
  return request.post<never, Favorite>('/favorites', { target_type: targetType, target_id: targetId })
}

export function removeFavorite(targetType: FavoriteTargetType, targetId: number) {
  return request.delete<never, { removed: boolean }>(`/favorites/${targetType}/${targetId}`)
}

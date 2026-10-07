export type FavoriteTargetType = 'plant' | 'article'

export const FavoriteTargetTypeMap: Record<FavoriteTargetType, string> = {
  plant: '植物品种',
  article: '养护文章',
}

export interface Favorite {
  id: number
  user_id: number
  target_type: FavoriteTargetType
  target_id: number
  created_at: string
}

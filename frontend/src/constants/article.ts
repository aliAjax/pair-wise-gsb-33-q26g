export type CareTopicTag = 'fertilizing' | 'pruning' | 'repotting' | 'pest_control' | 'propagation'

export const CareTopicTagMap: Record<CareTopicTag, string> = {
  fertilizing: '施肥',
  pruning: '修剪',
  repotting: '换盆',
  pest_control: '病虫害',
  propagation: '繁殖',
}

export const TOPIC_TAGS = Object.keys(CareTopicTagMap) as CareTopicTag[]

export interface CareArticle {
  id: number
  user_id: number
  title: string
  content: string
  cover: string
  topic_tag: CareTopicTag
  status: string
  view_count: number
  created_at: string
  updated_at: string
}

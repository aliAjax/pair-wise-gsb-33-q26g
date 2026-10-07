import request from '@/utils/request'
import type { CareArticle, CareTopicTag } from '@/constants/article'
import type { PageData } from '@/types/api'

export function listArticles(params: { page?: number; page_size?: number; topic_tag?: string; keyword?: string }) {
  return request.get<never, PageData<CareArticle>>('/articles', { params })
}

export function getArticle(id: number | string) {
  return request.get<never, CareArticle>(`/articles/${id}`)
}

export function createArticle(payload: { title: string; content: string; cover?: string; topic_tag: CareTopicTag }) {
  return request.post<never, CareArticle>('/articles', payload)
}

export function updateArticle(id: number, payload: { title: string; content: string; cover?: string; topic_tag: CareTopicTag }) {
  return request.put<never, CareArticle>(`/articles/${id}`, payload)
}

export function deleteArticle(id: number) {
  return request.delete<never, { deleted: boolean }>(`/articles/${id}`)
}

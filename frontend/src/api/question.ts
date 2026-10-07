import request from '@/utils/request'
import type { Question, Answer } from '@/types/api'
import type { PageData } from '@/types/api'

export function listQuestions(params: { page?: number; page_size?: number }) {
  return request.get<never, PageData<Question>>('/questions', { params })
}

export function getQuestion(id: number | string) {
  return request.get<never, Question>(`/questions/${id}`)
}

export function createQuestion(payload: { title: string; content: string; images?: string }) {
  return request.post<never, Question>('/questions', payload)
}

export function listAnswers(questionId: number) {
  return request.get<never, Answer[]>(`/questions/${questionId}/answers`)
}

export function createAnswer(questionId: number, content: string) {
  return request.post<never, Answer>(`/questions/${questionId}/answers`, { content })
}

export function adoptAnswer(questionId: number, answerId: number) {
  return request.put<never, Answer>(`/questions/${questionId}/adopt`, { answer_id: answerId })
}

export function likeAnswer(answerId: number) {
  return request.put<never, Answer>(`/answers/${answerId}/like`)
}

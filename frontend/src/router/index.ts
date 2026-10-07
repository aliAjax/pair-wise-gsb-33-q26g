import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'home', component: () => import('@/pages/Home.vue'), meta: { title: '首页' } },
  { path: '/plants', name: 'plants', component: () => import('@/pages/PlantLibrary.vue'), meta: { title: '品种库' } },
  { path: '/plants/:id', name: 'plantDetail', component: () => import('@/pages/PlantDetail.vue'), meta: { title: '品种详情' } },
  { path: '/articles', name: 'articles', component: () => import('@/pages/ArticleList.vue'), meta: { title: '养护文章' } },
  { path: '/articles/:id', name: 'articleDetail', component: () => import('@/pages/ArticleDetail.vue'), meta: { title: '文章详情' } },
  { path: '/pests', name: 'pests', component: () => import('@/pages/PestManual.vue'), meta: { title: '病虫害手册' } },
  { path: '/calendar', name: 'calendar', component: () => import('@/pages/SeasonCalendar.vue'), meta: { title: '季节养护日历', requiresAuth: true } },
  { path: '/garden', name: 'garden', component: () => import('@/pages/Garden.vue'), meta: { title: '我的花园', requiresAuth: true } },
  { path: '/questions', name: 'questions', component: () => import('@/pages/QuestionCommunity.vue'), meta: { title: '问答社区' } },
  { path: '/questions/:id', name: 'questionDetail', component: () => import('@/pages/QuestionDetail.vue'), meta: { title: '问题详情' } },
  { path: '/quiz', name: 'quiz', component: () => import('@/pages/Quiz.vue'), meta: { title: '养护测验' } },
  { path: '/profile', name: 'profile', component: () => import('@/pages/Profile.vue'), meta: { title: '个人中心', requiresAuth: true } },
  { path: '/login', name: 'login', component: () => import('@/pages/Login.vue'), meta: { title: '登录' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  document.title = `${to.meta.title || ''} - 植物养护知识百科平台`
  if (to.meta.requiresAuth) {
    const auth = useAuthStore()
    if (!auth.token) {
      return { path: '/login', query: { redirect: to.fullPath } }
    }
  }
  return true
})

export default router

<template>
  <div class="page" v-if="question">
    <el-page-header @back="$router.back()" content="问题详情" />
    <el-card class="question-card">
      <h2>{{ question.title }}</h2>
      <p class="content">{{ question.content }}</p>
      <p class="meta">{{ formatDateTime(question.created_at) }}</p>
    </el-card>
    <el-card class="answers-card">
      <template #header>回答（{{ answers.length }}）</template>
      <div v-for="a in answers" :key="a.id" class="answer">
        <div class="answer-head">
          <span>用户 #{{ a.user_id }}</span>
          <el-tag v-if="a.is_best" type="success" size="small">最佳回答</el-tag>
        </div>
        <p class="answer-content">{{ a.content }}</p>
        <div class="answer-actions">
          <el-button size="small" @click="like(a.id)">👍 {{ a.like_count }}</el-button>
          <el-button v-if="isOwner && !a.is_best" size="small" type="warning" @click="adopt(a.id)">采纳为最佳</el-button>
        </div>
      </div>
      <el-empty v-if="!answers.length" description="暂无回答" />
      <div class="reply">
        <el-input v-model="reply" type="textarea" :rows="3" placeholder="写下你的回答…" />
        <el-button type="primary" :loading="replying" @click="submitReply">提交回答</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getQuestion, listAnswers, createAnswer, adoptAnswer, likeAnswer } from '@/api/question'
import { useAuth } from '@/hooks/useAuth'
import { formatDateTime } from '@/utils/dateFormat'
import type { Answer, Question } from '@/types/api'

const route = useRoute()
const router = useRouter()
const { isLoggedIn, user } = useAuth()
const question = ref<Question | null>(null)
const answers = ref<Answer[]>([])
const reply = ref('')
const replying = ref(false)

const isOwner = computed(() => !!user.value && question.value?.user_id === user.value.id)

onMounted(async () => {
  const id = Number(route.params.id)
  question.value = await getQuestion(id)
  answers.value = await listAnswers(id)
})

async function submitReply() {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  if (!reply.value.trim()) return
  replying.value = true
  try {
    await createAnswer(question.value!.id, reply.value)
    answers.value = await listAnswers(question.value!.id)
    reply.value = ''
  } finally {
    replying.value = false
  }
}
async function adopt(answerId: number) {
  await adoptAnswer(question.value!.id, answerId)
  answers.value = await listAnswers(question.value!.id)
  ElMessage.success('已采纳该回答')
}
async function like(answerId: number) {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  await likeAnswer(answerId)
  answers.value = await listAnswers(question.value!.id)
}
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.question-card { margin-bottom: 16px; }
.content { line-height: 1.7; }
.meta { color: #999; font-size: 12px; }
.answers-card .answer { border-bottom: 1px solid #f0f0f0; padding: 12px 0; }
.answer-head { display: flex; gap: 8px; align-items: center; color: #888; font-size: 13px; }
.answer-content { line-height: 1.6; }
.answer-actions { display: flex; gap: 8px; }
.reply { margin-top: 16px; display: flex; flex-direction: column; gap: 12px; }
</style>

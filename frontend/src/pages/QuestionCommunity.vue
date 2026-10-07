<template>
  <div class="page">
    <h1>问答互动社区</h1>
    <el-card class="ask-card">
      <el-input v-model="form.title" placeholder="问题标题" />
      <el-input v-model="form.content" type="textarea" :rows="3" placeholder="详细描述你的养护问题" class="ask-body" />
      <el-button type="primary" :loading="submitting" @click="submit">发布问题</el-button>
    </el-card>
    <QuestionItem v-for="q in questions" :key="q.id" :question="q" @open="(id) => $router.push(`/questions/${id}`)" />
    <el-pagination v-if="total > 0" layout="prev, pager, next" :total="total" :page-size="pageSize" :current-page="page" @current-change="onPage" class="pager" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import QuestionItem from '@/components/common/QuestionItem.vue'
import { listQuestions, createQuestion } from '@/api/question'
import { useAuth } from '@/hooks/useAuth'
import type { Question } from '@/types/api'

const router = useRouter()
const { isLoggedIn } = useAuth()
const questions = ref<Question[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const submitting = ref(false)
const form = reactive({ title: '', content: '' })

onMounted(() => load())

async function load() {
  const res = await listQuestions({ page: page.value, page_size: pageSize })
  questions.value = res.list
  total.value = res.total
}
async function submit() {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  if (!form.title || !form.content) {
    ElMessage.warning('请填写标题和内容')
    return
  }
  submitting.value = true
  try {
    await createQuestion({ title: form.title, content: form.content })
    ElMessage.success('问题发布成功')
    form.title = ''
    form.content = ''
    page.value = 1
    await load()
  } finally {
    submitting.value = false
  }
}
function onPage(p: number) {
  page.value = p
  load()
}
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.ask-card { margin-bottom: 16px; }
.ask-body { margin: 12px 0; }
.pager { margin-top: 16px; justify-content: center; }
</style>

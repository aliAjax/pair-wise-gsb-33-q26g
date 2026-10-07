<template>
  <div class="page">
    <h1>养护知识小测验</h1>
    <el-alert title="共 6 题，交卷后展示成绩与错题解析" type="info" :closable="false" show-icon />
    <QuizCard
      v-for="(q, i) in QUIZ_BANK"
      :key="q.id"
      :question="q"
      :index="i"
      :submitted="submitted"
      :model-value="getAnswer(q.id)"
      @update:model-value="setAnswer(q.id, $event)"
    />
    <div class="actions">
      <el-button v-if="!submitted" type="primary" size="large" @click="submitQuiz">交卷</el-button>
      <template v-else>
        <el-result :icon="score >= 4 ? 'success' : 'warning'" :title="`得分：${score} / ${QUIZ_BANK.length}`">
          <template #sub-title>正确 {{ correct }} 题，错题解析见每道题下方</template>
          <template #extra>
            <el-button type="primary" @click="restart">重新作答</el-button>
          </template>
        </el-result>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import QuizCard from '@/components/common/QuizCard.vue'
import { QUIZ_BANK, useQuiz } from '@/hooks/useQuiz'

const { answers, submitted, submit, correctCount } = useQuiz()
const score = ref(0)
const correct = ref(0)

function getAnswer(qid: number): number | undefined {
  return answers.value[qid]
}
function setAnswer(qid: number, v: number) {
  answers.value[qid] = v
}

function submitQuiz() {
  score.value = submit()
  correct.value = correctCount()
}
function restart() {
  Object.keys(answers.value).forEach((k) => delete answers.value[Number(k)])
  submitted.value = false
  score.value = 0
  correct.value = 0
}
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.actions { margin: 20px 0; text-align: center; }
</style>

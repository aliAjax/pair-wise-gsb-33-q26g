import { ref } from 'vue'

export interface QuizQuestion {
  id: number
  question: string
  options: string[]
  answer: number
  explanation: string
}

export const QUIZ_BANK: QuizQuestion[] = [
  { id: 1, question: '以下哪种植物属于多肉植物？', options: ['月季', '吉娃娃', '碗莲', '龟背竹'], answer: 1, explanation: '吉娃娃为景天科拟石莲属多肉植物。' },
  { id: 2, question: '多肉植物夏季施肥的原则是？', options: ['薄肥勤施', '大量施肥', '停止施肥', '只施氮肥'], answer: 2, explanation: '夏季高温多数多肉休眠，应停止施肥避免肥害。' },
  { id: 3, question: '月季黑斑病的典型症状是？', options: ['叶片白粉', '黑色圆形斑点', '叶背蛛网', '叶片卷曲'], answer: 1, explanation: '黑斑病叶片出现黑色圆形斑点，边缘放射状。' },
  { id: 4, question: '龟背竹适合的光照条件是？', options: ['全日照', '散射光', '完全黑暗', '强直射光'], answer: 1, explanation: '龟背竹耐阴，适合明亮散射光环境。' },
  { id: 5, question: '换盆的最佳季节通常是？', options: ['夏季', '深冬', '春季', '雨季'], answer: 2, explanation: '春季气温回升、根系活跃，是换盆最佳时机。' },
  { id: 6, question: '“见干见湿”的浇水原则适用于？', options: ['所有植物', '多肉植物', '绝大多数盆栽植物', '水生植物'], answer: 2, explanation: '绝大多数盆栽植物遵循见干见湿原则。' },
]

export function useQuiz() {
  const answers = ref<Record<number, number>>({})
  const submitted = ref(false)

  function submit(): number {
    submitted.value = true
    let score = 0
    for (const q of QUIZ_BANK) {
      if (answers.value[q.id] === q.answer) score++
    }
    return score
  }

  const correctCount = () => Object.entries(answers.value).filter(([id, a]) => {
    const q = QUIZ_BANK.find((x) => x.id === Number(id))
    return q && q.answer === a
  }).length

  return { answers, submitted, submit, correctCount }
}

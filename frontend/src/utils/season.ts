import { PlantType } from '@/constants/plant'

export interface SeasonTask {
  month: number
  label: string
  task: string
}

const SEASON_TASKS: SeasonTask[] = [
  { month: 1, label: '一月', task: '冬季防寒：注意保温控水，避免冻伤' },
  { month: 2, label: '二月', task: '冬季收尾：清理枯叶，准备春季换盆' },
  { month: 3, label: '三月', task: '春季焕新：换盆施肥，病虫害早防' },
  { month: 4, label: '四月', task: '春季生长期：增加光照与浇水频率' },
  { month: 5, label: '五月', task: '初夏：追施磷钾肥促进开花' },
  { month: 6, label: '六月', task: '夏季遮阴：避免暴晒，加强通风' },
  { month: 7, label: '七月', task: '盛夏：早晚浇水，暂停施肥' },
  { month: 8, label: '八月', task: '夏末：修剪残花，防治红蜘蛛' },
  { month: 9, label: '九月', task: '初秋：恢复施肥，补充养分' },
  { month: 10, label: '十月', task: '秋季养护：减少施肥，做好越冬准备' },
  { month: 11, label: '十一月', task: '晚秋：控水控肥，移入室内' },
  { month: 12, label: '十二月', task: '冬季防寒：注意保温控水，避免冻伤' },
]

export function currentSeasonTask(): string {
  const month = new Date().getMonth() + 1
  return SEASON_TASKS.find((t) => t.month === month)?.task ?? '季节养护'
}

export function getSeasonTasks(): SeasonTask[] {
  return SEASON_TASKS
}

export function plantTypeOptions(): { label: string; value: PlantType }[] {
  return [
    { label: '观花', value: 'flower' },
    { label: '观叶', value: 'foliage' },
    { label: '多肉', value: 'succulent' },
    { label: '水生', value: 'aquatic' },
  ]
}

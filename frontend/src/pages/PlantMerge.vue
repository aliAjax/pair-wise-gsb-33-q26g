<template>
  <div class="page">
    <h1>品种合并</h1>
    <el-alert type="info" :closable="false" class="tip"
      title="因别名重复建卡的品种在此合并：收藏、我的花园、养护提醒、病虫害四类关联改挂到保留卡，待合并卡退出。同一用户重复收藏/养护只留一条，养护提醒按最近日期归并。" />

    <el-card class="block">
      <template #header>1. 选择保留卡与待合并卡</template>
      <el-form label-width="100px">
        <el-form-item label="保留卡">
          <el-select v-model="keepId" filterable remote clearable :remote-method="searchKeep"
            placeholder="搜索名称/别名，选择要保留的品种卡" style="width: 420px" @change="resetPrecheck">
            <el-option v-for="p in keepOptions" :key="p.id" :value="p.id" :label="plantLabel(p)" />
          </el-select>
        </el-form-item>
        <el-form-item label="待合并卡">
          <el-select v-model="sourceIds" multiple filterable remote :remote-method="searchSources"
            placeholder="搜索名称/别名，选择要合并退出的品种卡" style="width: 420px" @change="resetPrecheck">
            <el-option v-for="p in sourceOptions" :key="p.id" :value="p.id" :label="plantLabel(p)" :disabled="p.id === keepId" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="prechecking" :disabled="!keepId || !sourceIds.length" @click="runPrecheck">预检</el-button>
          <el-button type="danger" :loading="merging" :disabled="!canMerge" @click="runMerge">执行合并</el-button>
          <span v-if="precheck && !precheck.issues.length" class="hint">预检通过，可执行合并</span>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card v-if="precheck" class="block">
      <template #header>2. 预检结果</template>
      <template v-if="precheck.issues.length">
        <el-alert v-for="issue in precheck.issues" :key="issue" type="error" :title="issue" :closable="false" class="issue" />
      </template>
      <template v-else>
        <el-descriptions :column="2" border class="precheck-meta">
          <el-descriptions-item label="保留卡">{{ precheck.keep_name }}（#{{ precheck.keep_id }}）</el-descriptions-item>
          <el-descriptions-item label="待合并卡">{{ precheck.source_names.join('、') }}</el-descriptions-item>
        </el-descriptions>
        <el-table :data="precheckRows" border>
          <el-table-column prop="label" label="关联类型" width="160" />
          <el-table-column prop="total" label="待迁移" width="120" />
          <el-table-column prop="dedupe" label="去重/归并" width="120" />
          <el-table-column prop="note" label="说明" />
        </el-table>
      </template>
    </el-card>

    <el-card v-if="mergeTask" class="block">
      <template #header>3. 迁移结果</template>
      <el-alert :type="MergeTaskStatusAlertType[mergeTask.status]" :closable="false" class="issue"
        :title="`任务 #${mergeTask.id}：${MergeTaskStatusMap[mergeTask.status]}${mergeMessage ? ' — ' + mergeMessage : ''}`" />
      <el-alert v-if="mergeTask.last_error" type="error" :closable="false" class="issue" :title="mergeTask.last_error" />
      <el-table v-if="mergeTask.result" :data="resultRows(mergeTask.result)" border>
        <el-table-column prop="label" label="项目" width="200" />
        <el-table-column prop="value" label="数量" />
      </el-table>
    </el-card>

    <el-card class="block">
      <template #header>
        <div class="tasks-header">
          <span>合并任务记录</span>
          <el-button size="small" @click="loadTasks">刷新</el-button>
        </div>
      </template>
      <el-table :data="tasks" empty-text="暂无合并任务">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="保留卡" width="110">
          <template #default="{ row }">#{{ row.keep_id }}</template>
        </el-table-column>
        <el-table-column label="待合并卡" width="140">
          <template #default="{ row }">{{ row.source_ids.map((i: number) => '#' + i).join('、') }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="MergeTaskStatusTagType[row.status as MergeTaskStatus]">{{ MergeTaskStatusMap[row.status as MergeTaskStatus] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="迁移结果">
          <template #default="{ row }">
            <span v-if="row.result">{{ resultSummary(row.result) }}</span>
            <span v-else-if="row.last_error" class="error-text">{{ row.last_error }}</span>
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column label="时间" width="170">
          <template #default="{ row }">{{ formatDate(row.updated_at) }}</template>
        </el-table-column>
      </el-table>
      <el-pagination v-if="taskTotal > 0" layout="prev, pager, next" :total="taskTotal" :page-size="taskPageSize"
        :current-page="taskPage" class="pager" @current-change="onTaskPage" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listPlants } from '@/api/plant'
import { precheckMerge, submitMerge, getMergeTask, listMergeTasks } from '@/api/plantMerge'
import type { MergePrecheck, MergeResult, MergeTask } from '@/api/plantMerge'
import { MergeTaskStatusMap, MergeTaskStatusTagType, MergeTaskStatusAlertType, type MergeTaskStatus } from '@/constants/merge'
import type { PlantSpecies } from '@/constants/plant'
import { formatDate } from '@/utils/dateFormat'

const keepId = ref<number>()
const sourceIds = ref<number[]>([])
const keepOptions = ref<PlantSpecies[]>([])
const sourceOptions = ref<PlantSpecies[]>([])
const precheck = ref<MergePrecheck | null>(null)
const prechecking = ref(false)
const merging = ref(false)
const mergeTask = ref<MergeTask | null>(null)
const mergeMessage = ref('')
const tasks = ref<MergeTask[]>([])
const taskTotal = ref(0)
const taskPage = ref(1)
const taskPageSize = 10
let pollTimer: number | undefined

const canMerge = computed(() =>
  !!keepId.value && sourceIds.value.length > 0 && !!precheck.value && precheck.value.issues.length === 0 && !merging.value,
)

const precheckRows = computed(() => {
  if (!precheck.value) return []
  const s = precheck.value.stats
  return [
    { label: '收藏', total: s.favorites_total, dedupe: s.favorites_dedupe, note: '同一用户重复收藏只保留一条' },
    { label: '我的花园', total: s.gardens_total, dedupe: s.gardens_dedupe, note: '同一用户重复养护只保留一条' },
    { label: '养护提醒', total: s.reminders_total, dedupe: s.reminders_merge, note: '同用户同任务按最近日期归并' },
    { label: '病虫害关联', total: s.pests_total, dedupe: 0, note: '全部改挂到保留卡' },
  ]
})

function plantLabel(p: PlantSpecies) {
  return `${p.name}（${p.alias || '无别名'} / #${p.id}）`
}

async function searchPlants(keyword: string, target: typeof keepOptions) {
  const data = await listPlants({ page: 1, page_size: 20, keyword: keyword || '' })
  target.value = data.list
}

function searchKeep(kw: string) { searchPlants(kw, keepOptions) }
function searchSources(kw: string) { searchPlants(kw, sourceOptions) }

function resetPrecheck() {
  precheck.value = null
}

async function runPrecheck() {
  if (!keepId.value || !sourceIds.value.length) return
  prechecking.value = true
  try {
    precheck.value = await precheckMerge({ keep_id: keepId.value, source_ids: sourceIds.value })
    if (precheck.value.issues.length) {
      ElMessage.warning('预检未通过，请根据提示调整选择')
    }
  } finally {
    prechecking.value = false
  }
}

async function runMerge() {
  if (!keepId.value || !sourceIds.value.length) return
  try {
    await ElMessageBox.confirm(
      `将把 ${sourceIds.value.length} 张待合并卡的收藏、花园、提醒、病虫害关联改挂到保留卡 #${keepId.value}，并退出待合并卡。确认执行？`,
      '确认合并',
      { type: 'warning', confirmButtonText: '执行合并', cancelButtonText: '取消' },
    )
  } catch {
    return // user cancelled
  }
  merging.value = true
  try {
    const res = await submitMerge({ keep_id: keepId.value, source_ids: sourceIds.value })
    mergeTask.value = res.task
    mergeMessage.value = res.message
    if (res.task.status === 'succeeded') {
      ElMessage.success(res.message || '品种合并完成')
    }
    await loadTasks()
    if (res.task.status === 'processing' || res.task.status === 'pending') {
      pollTask(res.task.id)
    }
  } finally {
    merging.value = false
  }
}

function pollTask(id: number) {
  stopPoll()
  let attempts = 0
  pollTimer = window.setInterval(async () => {
    attempts += 1
    try {
      const task = await getMergeTask(id)
      mergeTask.value = task
      if (task.status === 'succeeded' || task.status === 'failed' || attempts >= 40) {
        stopPoll()
        if (task.status === 'succeeded') ElMessage.success('品种合并完成')
        if (task.status === 'failed') ElMessage.error('品种合并失败，可重试从检查点恢复')
        loadTasks()
      }
    } catch {
      stopPoll()
    }
  }, 1500)
}

function stopPoll() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = undefined
  }
}

async function loadTasks() {
  const data = await listMergeTasks({ page: taskPage.value, page_size: taskPageSize })
  tasks.value = data.list
  taskTotal.value = data.total
}

function onTaskPage(p: number) {
  taskPage.value = p
  loadTasks()
}

function resultRows(r: MergeResult) {
  return [
    { label: '收藏改挂 / 去重', value: `${r.favorites_moved} / ${r.favorites_deduped}` },
    { label: '花园改挂 / 去重', value: `${r.gardens_moved} / ${r.gardens_deduped}` },
    { label: '提醒改挂 / 归并', value: `${r.reminders_moved} / ${r.reminders_merged}` },
    { label: '病虫害改挂', value: `${r.pests_moved}` },
    { label: '退出品种卡', value: `${r.retired_plants}` },
  ]
}

function resultSummary(r: MergeResult) {
  return `收藏+${r.favorites_moved}/-${r.favorites_deduped} 花园+${r.gardens_moved}/-${r.gardens_deduped} 提醒+${r.reminders_moved}/-${r.reminders_merged} 病虫害+${r.pests_moved} 退卡${r.retired_plants}`
}

onMounted(() => {
  searchPlants('', keepOptions)
  searchPlants('', sourceOptions)
  loadTasks()
})

onUnmounted(stopPoll)
</script>

<style scoped>
.page { max-width: 1000px; margin: 0 auto; }
.tip { margin-bottom: 16px; }
.block { margin-bottom: 16px; }
.issue { margin-bottom: 8px; }
.hint { margin-left: 12px; color: #67c23a; }
.precheck-meta { margin-bottom: 12px; }
.tasks-header { display: flex; justify-content: space-between; align-items: center; }
.error-text { color: #f56c6c; }
.pager { margin-top: 12px; justify-content: center; }
</style>

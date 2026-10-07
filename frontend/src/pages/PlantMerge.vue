<template>
  <div class="page">
    <h1>品种合并</h1>
    <el-alert type="info" :closable="false" class="tip">
      用于清理因别名重复建卡的品种：收藏、我的花园、养护提醒、病虫害四类关联改挂到保留卡，旧卡退出品种库。
      同一用户重复收藏或重复养护只保留一条；养护提醒按最近日期归并。同批合并重复提交不会重复执行。
    </el-alert>

    <el-card class="block">
      <template #header>1. 选择品种卡</template>
      <el-form inline>
        <el-form-item label="保留卡">
          <el-select
            v-model="keepId"
            filterable
            remote
            clearable
            placeholder="搜索要保留的品种"
            :remote-method="searchKeep"
            :loading="searching"
            style="width: 260px"
            @change="onSelectionChange"
          >
            <el-option v-for="p in keepOptions" :key="p.id" :label="`${p.name}（${p.alias || '无别名'} #${p.id}）`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="待合并卡">
          <el-select
            v-model="sourceIds"
            multiple
            filterable
            remote
            placeholder="搜索要合并退出的品种"
            :remote-method="searchSources"
            :loading="searching"
            style="width: 360px"
            @change="onSelectionChange"
          >
            <el-option
              v-for="p in sourceOptions"
              :key="p.id"
              :label="`${p.name}（${p.alias || '无别名'} #${p.id}）`"
              :value="p.id"
              :disabled="p.id === keepId"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="!canSubmit" :loading="prechecking" @click="runPrecheck">预检</el-button>
          <el-button type="primary" :disabled="!canSubmit" :loading="merging" @click="runMerge">执行合并</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card v-if="precheck" class="block">
      <template #header>2. 预检结果（保留卡：{{ precheck.keep_name }}）</template>
      <el-table :data="precheck.items">
        <el-table-column label="待合并卡" min-width="140">
          <template #default="{ row }">#{{ row.source_id }} {{ row.source_name }}</template>
        </el-table-column>
        <el-table-column label="收藏" min-width="130">
          <template #default="{ row }">
            {{ row.favorites_total }} 条
            <el-tag v-if="row.favorites_dedup" size="small" type="warning">{{ row.favorites_dedup }} 人去重</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="我的花园" min-width="130">
          <template #default="{ row }">
            {{ row.gardens_total }} 条
            <el-tag v-if="row.gardens_dedup" size="small" type="warning">{{ row.gardens_dedup }} 人去重</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="养护提醒" min-width="130">
          <template #default="{ row }">
            {{ row.reminders_total }} 条
            <el-tag v-if="row.reminders_merge" size="small" type="warning">{{ row.reminders_merge }} 条归并</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="病虫害" min-width="90">
          <template #default="{ row }">{{ row.pests_total }} 条</template>
        </el-table-column>
        <el-table-column label="备注" min-width="120">
          <template #default="{ row }">
            <el-tag v-if="row.already_merged" size="small" type="success">已完成过合并</el-tag>
            <span v-else>将改挂到保留卡</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card v-if="results.length" class="block">
      <template #header>3. 迁移结果</template>
      <el-table :data="results">
        <el-table-column label="合并" min-width="180">
          <template #default="{ row }">#{{ row.source_id }} {{ row.source_name }} → {{ row.keep_name }}</template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ MergeStatusMap[row.status as MergeTaskStatus] }}</el-tag>
            <div v-if="row.status === 'failed' && row.stage" class="sub">止于：{{ MergeStageMap[row.stage] || '开始前' }}</div>
          </template>
        </el-table-column>
        <el-table-column label="收藏（移/去重）" width="130">
          <template #default="{ row }">{{ row.stats.favorites_moved }} / {{ row.stats.favorites_deduped }}</template>
        </el-table-column>
        <el-table-column label="花园（移/去重）" width="130">
          <template #default="{ row }">{{ row.stats.gardens_moved }} / {{ row.stats.gardens_deduped }}</template>
        </el-table-column>
        <el-table-column label="提醒（移/归并）" width="130">
          <template #default="{ row }">{{ row.stats.reminders_moved }} / {{ row.stats.reminders_merged }}</template>
        </el-table-column>
        <el-table-column label="病虫害" width="80">
          <template #default="{ row }">{{ row.stats.pests_moved }}</template>
        </el-table-column>
        <el-table-column label="失败原因" min-width="160">
          <template #default="{ row }">{{ row.last_error || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="90">
          <template #default="{ row }">
            <el-button v-if="row.status === 'failed'" size="small" type="warning" @click="retry(row)">重试</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="block">
      <template #header>合并记录</template>
      <el-table :data="tasks" empty-text="暂无合并记录">
        <el-table-column prop="id" label="任务" width="70" />
        <el-table-column label="合并" min-width="180">
          <template #default="{ row }">#{{ row.source_id }} {{ row.source_name }} → {{ row.keep_name }}</template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ MergeStatusMap[row.status as MergeTaskStatus] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="完成时间" width="150">
          <template #default="{ row }">{{ row.finished_at ? formatDateTime(row.finished_at) : '-' }}</template>
        </el-table-column>
        <el-table-column label="发起时间" width="150">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-if="taskTotal > 0"
        layout="prev, pager, next"
        :total="taskTotal"
        :page-size="taskPageSize"
        :current-page="taskPage"
        class="pager"
        @current-change="loadTasks"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listPlants } from '@/api/plant'
import {
  getMergeTask,
  listMergeTasks,
  mergePlants,
  precheckPlantMerge,
  MergeStageMap,
  MergeStatusMap,
  type MergePrecheckResult,
  type MergeTaskStatus,
  type PlantMergeTask,
} from '@/api/plantMerge'
import type { PlantSpecies } from '@/constants/plant'
import { formatDateTime } from '@/utils/dateFormat'

const keepId = ref<number>()
const sourceIds = ref<number[]>([])
const keepOptions = ref<PlantSpecies[]>([])
const sourceOptions = ref<PlantSpecies[]>([])
const searching = ref(false)
const prechecking = ref(false)
const merging = ref(false)
const precheck = ref<MergePrecheckResult | null>(null)
const results = ref<PlantMergeTask[]>([])
const tasks = ref<PlantMergeTask[]>([])
const taskTotal = ref(0)
const taskPage = ref(1)
const taskPageSize = 10

const canSubmit = computed(() => !!keepId.value && sourceIds.value.length > 0)

onMounted(async () => {
  await Promise.all([searchKeep(''), searchSources(''), loadTasks(1)])
})

async function searchPlants(keyword: string): Promise<PlantSpecies[]> {
  searching.value = true
  try {
    const res = await listPlants({ keyword, page: 1, page_size: 20 })
    return res.list
  } finally {
    searching.value = false
  }
}

async function searchKeep(keyword: string) {
  keepOptions.value = await searchPlants(keyword)
}

async function searchSources(keyword: string) {
  sourceOptions.value = await searchPlants(keyword)
}

function onSelectionChange() {
  if (keepId.value && sourceIds.value.includes(keepId.value)) {
    sourceIds.value = sourceIds.value.filter((id) => id !== keepId.value)
  }
  precheck.value = null
  results.value = []
}

async function runPrecheck() {
  if (!canSubmit.value) return
  prechecking.value = true
  try {
    precheck.value = await precheckPlantMerge({ keep_id: keepId.value!, source_ids: [...sourceIds.value] })
  } finally {
    prechecking.value = false
  }
}

async function runMerge() {
  if (!canSubmit.value) return
  try {
    await ElMessageBox.confirm(
      `将把 ${sourceIds.value.length} 张待合并卡的四类关联改挂到保留卡，旧卡退出品种库。确认执行？`,
      '执行品种合并',
      { confirmButtonText: '执行', cancelButtonText: '取消', type: 'warning' },
    )
  } catch {
    return
  }
  merging.value = true
  try {
    results.value = await mergePlants({ keep_id: keepId.value!, source_ids: [...sourceIds.value] })
    await pollRunning()
    const failed = results.value.filter((r) => r.status === 'failed').length
    if (failed) {
      ElMessage.warning(`${failed} 张卡合并失败，可从检查点重试`)
    } else {
      ElMessage.success('合并完成')
    }
    await loadTasks(taskPage.value)
  } finally {
    merging.value = false
  }
}

// pollRunning waits for tasks still owned by a concurrent executor.
async function pollRunning() {
  for (let i = 0; i < 15; i++) {
    const pending = results.value.filter((r) => r.status === 'pending' || r.status === 'running')
    if (!pending.length) return
    await new Promise((resolve) => setTimeout(resolve, 1000))
    for (const p of pending) {
      const fresh = await getMergeTask(p.id)
      const idx = results.value.findIndex((r) => r.id === p.id)
      if (idx >= 0) results.value[idx] = fresh
    }
  }
}

async function retry(row: PlantMergeTask) {
  const [updated] = await mergePlants({ keep_id: row.keep_id, source_ids: [row.source_id] })
  const idx = results.value.findIndex((r) => r.id === row.id)
  if (idx >= 0 && updated) results.value[idx] = updated
  await pollRunning()
  await loadTasks(taskPage.value)
}

async function loadTasks(page: number) {
  taskPage.value = page
  const res = await listMergeTasks({ page, page_size: taskPageSize })
  tasks.value = res.list
  taskTotal.value = res.total
}

function statusTagType(status: MergeTaskStatus): 'success' | 'warning' | 'danger' | 'info' {
  if (status === 'done') return 'success'
  if (status === 'failed') return 'danger'
  if (status === 'running') return 'warning'
  return 'info'
}
</script>

<style scoped>
.page { max-width: 1200px; margin: 0 auto; }
.tip { margin-bottom: 16px; }
.block { margin-bottom: 16px; }
.sub { font-size: 12px; color: #909399; margin-top: 4px; }
.pager { margin-top: 16px; justify-content: center; }
</style>

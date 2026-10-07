<template>
  <div class="page">
    <h1>我的花园</h1>
    <el-row :gutter="16">
      <el-col :xs="24" :md="16">
        <el-card>
          <template #header>花园清单</template>
          <el-table :data="gardenItems" empty-text="花园还是空的，去品种库添加吧">
            <el-table-column label="植物">
              <template #default="{ row }">
                <span class="garden-name">{{ row.nickname || row.plant_species_id }}</span>
              </template>
            </el-table-column>
            <el-table-column label="位置" prop="location" />
            <el-table-column label="拥有时间">
              <template #default="{ row }">{{ formatDate(row.owned_since) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="remove(row.id)">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
        <el-card class="block">
          <template #header>我的养护提醒</template>
          <ReminderList :reminders="reminders" @done="markDone" @remove="removeReminder" />
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card>
          <template #header>我的收藏</template>
          <el-table :data="favorites" empty-text="暂无收藏">
            <el-table-column prop="target_type" label="类型" width="80">
              <template #default="{ row }">{{ FavoriteTargetTypeMap[row.target_type as FavoriteTargetType] }}</template>
            </el-table-column>
            <el-table-column prop="target_id" label="目标 ID" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import ReminderList from '@/components/common/ReminderList.vue'
import { listGardens, removeGarden } from '@/api/garden'
import { listFavorites } from '@/api/favorite'
import { listReminders, deleteReminder, updateReminderStatus } from '@/api/reminder'
import { FavoriteTargetTypeMap, type Favorite, type FavoriteTargetType } from '@/constants/favorite'
import { formatDate } from '@/utils/dateFormat'
import type { CareReminder, UserGarden } from '@/types/api'

const gardenItems = ref<UserGarden[]>([])
const favorites = ref<Favorite[]>([])
const reminders = ref<CareReminder[]>([])

onMounted(async () => {
  gardenItems.value = await listGardens()
  favorites.value = await listFavorites()
  reminders.value = await listReminders()
})

async function remove(id: number) {
  await removeGarden(id)
  gardenItems.value = await listGardens()
  ElMessage.success('已移除')
}
async function markDone(id: number) {
  await updateReminderStatus(id, 'done')
  reminders.value = await listReminders()
}
async function removeReminder(id: number) {
  await deleteReminder(id)
  reminders.value = await listReminders()
}
</script>

<style scoped>
.page { max-width: 1200px; margin: 0 auto; }
.block { margin-top: 16px; }
.garden-name { font-weight: 600; }
</style>

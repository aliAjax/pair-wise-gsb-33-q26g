<template>
  <el-button :type="favorited ? 'warning' : 'default'" :loading="loading" @click="onToggle">
    <el-icon><component :is="favorited ? StarFilled : Star" /></el-icon>&nbsp;
    {{ favorited ? '已收藏' : '收藏' }}
  </el-button>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { Star, StarFilled } from '@element-plus/icons-vue'
import { useAuth } from '@/hooks/useAuth'
import { useFavorite } from '@/hooks/useFavorite'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import type { FavoriteTargetType } from '@/constants/favorite'

const props = defineProps<{ targetType: FavoriteTargetType; targetId: number; initial?: boolean }>()
const router = useRouter()
const { isLoggedIn } = useAuth()
const { favorited, loading, toggle } = useFavorite(props.targetType)

watch(() => props.initial, (v) => { favorited.value = !!v }, { immediate: true })

async function onToggle() {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  const added = await toggle(props.targetId)
  ElMessage.success(added ? '收藏成功' : '已取消收藏')
}
</script>

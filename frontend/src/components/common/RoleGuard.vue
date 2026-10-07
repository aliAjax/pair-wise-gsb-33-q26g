<template>
  <slot v-if="allowed" />
  <el-result v-else icon="warning" title="无权限访问" sub-title="当前角色无权查看该内容">
    <template #extra><el-button type="primary" @click="$router.push('/')">返回首页</el-button></template>
  </el-result>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { hasRole } from '@/router/guards'

const props = defineProps<{ roles: string[] }>()
const allowed = computed(() => hasRole(props.roles))
</script>

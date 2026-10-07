<template>
  <div class="uploader">
    <el-upload
      :show-file-list="false"
      :http-request="doUpload"
      accept="image/*"
      :before-upload="beforeUpload"
    >
      <el-button :loading="loading" type="primary" plain>
        <el-icon><Plus /></el-icon>&nbsp;上传图片
      </el-button>
    </el-upload>
    <el-image v-if="modelValue" :src="modelValue" fit="cover" class="preview" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { useAuthStore } from '@/stores/authStore'

defineProps<{ modelValue: string }>()
const emit = defineEmits<{ (e: 'update:modelValue', url: string): void }>()
const loading = ref(false)

function beforeUpload(file: File): boolean {
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('图片不能超过 5MB')
    return false
  }
  return true
}

async function doUpload(option: { file: File }) {
  const auth = useAuthStore()
  const form = new FormData()
  form.append('file', option.file)
  loading.value = true
  try {
    const res = await axios.post('/api/v1/uploads', form, {
      headers: { Authorization: `Bearer ${auth.token}`, 'Content-Type': 'multipart/form-data' },
    })
    emit('update:modelValue', res.data.data.url)
    ElMessage.success('上传成功')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.uploader { display: flex; gap: 12px; align-items: center; }
.preview { width: 120px; height: 90px; border-radius: 6px; }
</style>

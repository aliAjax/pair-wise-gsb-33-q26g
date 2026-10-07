<template>
  <el-carousel v-if="images.length" height="320px" class="carousel">
    <el-carousel-item v-for="(img, i) in images" :key="i">
      <el-image :src="img" fit="cover" class="slide-img" />
    </el-carousel-item>
  </el-carousel>
  <div v-else class="empty-carousel">暂无图片</div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { parseImages } from '@/constants/plant'

const props = defineProps<{ imageUrls: string }>()
const images = computed(() => parseImages(props.imageUrls))
</script>

<style scoped>
.carousel { border-radius: 8px; overflow: hidden; }
.slide-img { width: 100%; height: 100%; }
.empty-carousel { height: 200px; display: flex; align-items: center; justify-content: center; background: #f0f0f0; color: #999; border-radius: 8px; }
</style>

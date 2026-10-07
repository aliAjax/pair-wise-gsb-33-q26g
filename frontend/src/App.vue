<template>
  <el-container class="app-shell">
    <el-header class="app-header">
      <div class="brand" @click="$router.push('/')">🌿 植物养护百科</div>
      <el-menu mode="horizontal" :ellipsis="false" router :default-active="$route.path">
        <el-menu-item index="/">首页</el-menu-item>
        <el-menu-item index="/plants">品种库</el-menu-item>
        <el-menu-item index="/articles">养护文章</el-menu-item>
        <el-menu-item index="/pests">病虫害手册</el-menu-item>
        <el-menu-item index="/calendar">季节日历</el-menu-item>
        <el-menu-item index="/garden">我的花园</el-menu-item>
        <el-menu-item index="/questions">问答社区</el-menu-item>
        <el-menu-item index="/quiz">养护测验</el-menu-item>
      </el-menu>
      <div class="user-area">
        <template v-if="auth.token">
          <el-dropdown @command="onCommand">
            <span class="user-name">{{ auth.user?.nickname || auth.user?.username }}</span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <el-button v-else type="primary" size="small" @click="$router.push('/login')">登录/注册</el-button>
      </div>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

const auth = useAuthStore()
const router = useRouter()

function onCommand(cmd: string) {
  if (cmd === 'profile') {
    router.push('/profile')
  } else if (cmd === 'logout') {
    auth.logout()
    router.push('/')
  }
}
</script>

<style>
html, body, #app { margin: 0; height: 100%; background: #f6f8f4; font-family: "PingFang SC", "Microsoft YaHei", sans-serif; }
.app-shell { min-height: 100vh; }
.app-header { display: flex; align-items: center; gap: 16px; background: #fff; box-shadow: 0 2px 8px rgba(0,0,0,.06); }
.brand { font-size: 20px; font-weight: 700; color: #3c8d5c; cursor: pointer; white-space: nowrap; }
.user-area { margin-left: auto; }
.user-name { cursor: pointer; color: #3c8d5c; font-weight: 600; }
</style>

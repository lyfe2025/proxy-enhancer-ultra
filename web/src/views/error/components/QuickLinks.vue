<template>
  <div class="quick-links">
    <h3>快速链接</h3>
    <div class="links-grid">
      <router-link
        v-for="link in links"
        :key="link.path"
        :to="link.path"
        class="quick-link"
      >
        <el-icon>
          <component :is="link.icon" />
        </el-icon>
        <span>{{ link.title }}</span>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  House,
  User,
  Setting,
  Document,
  DataAnalysis,
  Connection
} from '@element-plus/icons-vue'

interface QuickLink {
  path: string
  title: string
  icon: any
}

const props = withDefaults(defineProps<{
  links?: QuickLink[]
}>(), {
  links: () => [
    { path: '/', title: '首页', icon: House },
    { path: '/dashboard', title: '仪表盘', icon: DataAnalysis },
    { path: '/proxy', title: '代理管理', icon: Connection },
    { path: '/rules', title: '规则配置', icon: Setting },
    { path: '/data', title: '数据收集', icon: Document },
    { path: '/profile', title: '个人中心', icon: User }
  ]
})
</script>

<style scoped>
.quick-links {
  color: #ffffff;
  margin-bottom: 40px;
}

.quick-links h3 {
  margin: 0 0 20px 0;
  color: #00ff88;
  font-size: 18px;
  text-align: center;
}

.links-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  max-width: 600px;
  margin: 0 auto;
}

.quick-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(0, 255, 136, 0.2);
  border-radius: 12px;
  text-decoration: none;
  color: #cccccc;
  transition: all 0.3s ease;
}

.quick-link:hover {
  background: rgba(0, 255, 136, 0.1);
  border-color: rgba(0, 255, 136, 0.4);
  color: #ffffff;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 255, 136, 0.2);
}

.quick-link .el-icon {
  font-size: 24px;
  color: #00ff88;
}

.quick-link span {
  font-size: 14px;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .links-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .links-grid {
    grid-template-columns: 1fr;
  }
}

/* 深色主题适配 */
.dark .quick-link {
  background: rgba(255, 255, 255, 0.03);
  border-color: rgba(0, 255, 136, 0.15);
}

.dark .quick-link:hover {
  background: rgba(0, 255, 136, 0.08);
}
</style>

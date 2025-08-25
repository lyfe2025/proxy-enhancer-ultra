<template>
  <div class="error-content">
    <h1 class="error-title">{{ title }}</h1>
    <p class="error-description">{{ description }}</p>
    
    <!-- 可能的原因 -->
    <div class="error-reasons">
      <h3>可能的原因：</h3>
      <ul>
        <li v-for="reason in reasons" :key="reason">{{ reason }}</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string
  description?: string
  reasons?: string[]
}

withDefaults(defineProps<Props>(), {
  title: '页面未找到',
  description: '抱歉，您访问的页面不存在或已被移除。请检查URL是否正确，或返回首页继续浏览。',
  reasons: () => [
    '页面链接已过期或被删除',
    'URL地址输入错误',
    '您没有访问该页面的权限',
    '服务器配置问题'
  ]
})
</script>

<style scoped>
.error-content {
  color: #ffffff;
}

.error-title {
  font-size: 48px;
  font-weight: 700;
  margin: 0 0 20px 0;
  background: linear-gradient(45deg, #ffffff, #cccccc);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.error-description {
  font-size: 18px;
  line-height: 1.6;
  margin: 0 0 30px 0;
  color: #cccccc;
}

.error-reasons {
  text-align: left;
  max-width: 500px;
  margin: 0 auto 40px auto;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(0, 255, 136, 0.2);
}

.error-reasons h3 {
  margin: 0 0 15px 0;
  color: #00ff88;
  font-size: 16px;
}

.error-reasons ul {
  margin: 0;
  padding-left: 20px;
  color: #cccccc;
}

.error-reasons li {
  margin-bottom: 8px;
  line-height: 1.5;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .error-title {
    font-size: 32px;
  }
  
  .error-description {
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .error-title {
    font-size: 24px;
  }
}

/* 深色主题适配 */
.dark .error-reasons {
  background: rgba(255, 255, 255, 0.03);
  border-color: rgba(0, 255, 136, 0.15);
}
</style>

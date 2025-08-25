<template>
  <div class="social-login">
    <div class="social-login-divider">
      <span class="divider-text">或使用社交账号登录</span>
    </div>
    
    <div class="social-buttons">
      <el-tooltip content="微信登录" placement="top">
        <el-button circle @click="handleSocialLogin('wechat')">
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="20" height="20">
              <path fill="currentColor" d="M679.2 465.088c-18.816 0-37.632-12.032-37.632-31.168 0-18.816 18.816-31.168 37.632-31.168 25.28 0 43.776 12.352 43.776 31.168 0 19.136-18.496 31.168-43.776 31.168zM565.248 402.752c25.28 0 43.776 12.352 43.776 31.168 0 19.136-18.496 31.168-43.776 31.168-18.816 0-37.632-12.032-37.632-31.168 0-18.816 18.816-31.168 37.632-31.168z"/>
              <path fill="currentColor" d="M1024 576.32c0 100.032-93.888 181.248-209.6 181.248-31.168 0-62.336-6.336-87.296-18.816L621.888 789.12l24.96-75.2c-62.336-43.776-99.968-100.032-99.968-162.688 0-119.168 112.32-216.32 249.6-216.32S1024 457.152 1024 576.32z"/>
              <path fill="currentColor" d="M725.12 695.68c18.816 0 31.168 12.032 31.168 25.28 0 18.816-12.352 31.168-31.168 31.168-12.352 0-31.168-12.352-31.168-31.168 0-13.248 18.816-25.28 31.168-25.28zM631.168 695.68c18.816 0 31.168 12.032 31.168 25.28 0 18.816-12.352 31.168-31.168 31.168-12.352 0-31.168-12.352-31.168-31.168 0-13.248 18.816-25.28 31.168-25.28z"/>
            </svg>
          </el-icon>
        </el-button>
      </el-tooltip>

      <el-tooltip content="QQ登录" placement="top">
        <el-button circle @click="handleSocialLogin('qq')">
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="20" height="20">
              <path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm263.4 667.4c-18.7-14.4-37.8-31.2-37.8-31.2s-14.1 25.1-46.6 25.1-46.6-25.1-46.6-25.1-7 12.6-30.7 25.1c-24.3 12.8-48.8 12.8-48.8 12.8s33.2-31.9 33.2-75.7c0-47.1-33.2-75.7-33.2-75.7s24.5 0 48.8 12.8c23.7 12.5 30.7 25.1 30.7 25.1s14.1-25.1 46.6-25.1 46.6 25.1 46.6 25.1 19.1 16.8 37.8 31.2c38.4 29.5 38.4 75.7 0 75.7z"/>
            </svg>
          </el-icon>
        </el-button>
      </el-tooltip>

      <el-tooltip content="GitHub登录" placement="top">
        <el-button circle @click="handleSocialLogin('github')">
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="20" height="20">
              <path fill="currentColor" d="M512 42.666667A464.64 464.64 0 0 0 42.666667 502.186667 460.373333 460.373333 0 0 0 363.52 938.666667c23.466667 4.266667 32-9.813333 32-22.186667v-78.08c-130.56 27.733333-158.293333-61.44-158.293333-61.44a122.026667 122.026667 0 0 0-52.053334-67.413333c-42.666667-28.16 3.413333-27.733333 3.413334-27.733334a98.56 98.56 0 0 1 71.68 47.36 101.12 101.12 0 0 0 136.533333 37.973334 99.413333 99.413333 0 0 1 29.866667-61.44c-104.106667-11.52-213.333333-50.773333-213.333334-226.986667a177.066667 177.066667 0 0 1 47.36-124.16 161.28 161.28 0 0 1 4.693334-121.173333s39.68-12.8 128 46.933333a455.68 455.68 0 0 1 234.666666 0c89.6-59.733333 128-46.933333 128-46.933333a161.28 161.28 0 0 1 4.693334 121.173333A177.066667 177.066667 0 0 1 810.666667 477.866667c0 176.64-110.08 215.466667-213.333334 226.986666a106.666667 106.666667 0 0 1 32 85.333334v125.866666c0 14.933333 8.533333 27.733333 32 22.186667A460.8 460.8 0 0 0 981.333333 502.186667 464.64 464.64 0 0 0 512 42.666667"/>
            </svg>
          </el-icon>
        </el-button>
      </el-tooltip>

      <el-tooltip content="Google登录" placement="top">
        <el-button circle @click="handleSocialLogin('google')">
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="20" height="20">
              <path fill="currentColor" d="M512 224c70.4 0 134.4 25.6 184.32 67.84l137.6-137.6C744.96 71.04 635.52 32 512 32 310.08 32 142.08 161.28 70.4 342.72l160 124.16C284.16 344.32 388.48 224 512 224z"/>
              <path fill="currentColor" d="M512 800c-124.16 0-228.48-120.32-281.6-243.12l-160 124.16C142.08 862.72 310.08 992 512 992c118.4 0 224-43.52 302.72-114.56L675.84 755.84C625.92 780.8 572.16 800 512 800z"/>
              <path fill="currentColor" d="M825.6 441.6H512v176h181.12c-24.96 66.56-85.76 116.48-181.12 116.48v155.52c188.16 0 340.48-152.32 340.48-340.48 0-22.72-2.56-44.8-6.88-65.92z"/>
              <path fill="currentColor" d="M512 441.6V265.6H312.32c-24.96 66.56-24.96 140.8 0 207.36H512z"/>
            </svg>
          </el-icon>
        </el-button>
      </el-tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'

type SocialProvider = 'wechat' | 'qq' | 'github' | 'google'

const emit = defineEmits<{
  socialLogin: [provider: SocialProvider]
}>()

const handleSocialLogin = (provider: SocialProvider) => {
  // 这里可以根据不同的提供商显示不同的消息
  const providerNames = {
    wechat: '微信',
    qq: 'QQ',
    github: 'GitHub',
    google: 'Google'
  }
  
  ElMessage.info(`${providerNames[provider]}登录功能开发中...`)
  emit('socialLogin', provider)
}
</script>

<style scoped>
.social-login {
  margin-top: 24px;
}

.social-login-divider {
  position: relative;
  text-align: center;
  margin-bottom: 20px;
}

.social-login-divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--el-border-color-light);
}

.divider-text {
  background: var(--el-bg-color);
  padding: 0 16px;
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  position: relative;
  z-index: 1;
}

.social-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.social-buttons .el-button {
  width: 44px;
  height: 44px;
  transition: all 0.3s ease;
}

.social-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.social-buttons .el-button .el-icon {
  font-size: 20px;
}

.social-buttons .el-button .el-icon svg {
  width: 20px;
  height: 20px;
  fill: currentColor;
}

/* 深色主题适配 */
.dark .divider-text {
  background: var(--el-bg-color-overlay);
}
</style>

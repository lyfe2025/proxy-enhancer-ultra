<template>
  <div class="error-icon">
    <div class="number-404">
      <span class="four">4</span>
      <span class="zero">0</span>
      <span class="four">4</span>
    </div>
    <div class="error-illustration">
      <div class="planet">
        <div class="ring"></div>
        <div class="cover-ring"></div>
        <div class="spots">
          <span></span>
          <span></span>
          <span></span>
          <span></span>
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>
      <div class="glitch">
        <div class="glitch__line"></div>
        <div class="glitch__line"></div>
        <div class="glitch__line"></div>
        <div class="glitch__line"></div>
        <div class="glitch__line"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 纯动画组件，无需逻辑
</script>

<style scoped>
.error-icon {
  margin-bottom: 40px;
  position: relative;
}

/* 404数字 */
.number-404 {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.number-404 span {
  font-size: 120px;
  font-weight: 900;
  text-shadow: 0 0 20px rgba(0, 255, 136, 0.5);
  animation: glow 2s ease-in-out infinite alternate;
}

.four {
  color: #00ff88;
}

.zero {
  color: #0088ff;
  position: relative;
}

@keyframes glow {
  from {
    text-shadow: 0 0 20px rgba(0, 255, 136, 0.5);
  }
  to {
    text-shadow: 0 0 30px rgba(0, 255, 136, 0.8), 0 0 40px rgba(0, 255, 136, 0.6);
  }
}

/* 错误插图 */
.error-illustration {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

/* 行星 */
.planet {
  position: relative;
  width: 100px;
  height: 100px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  box-shadow: 
    inset -10px -10px 20px rgba(0, 0, 0, 0.3),
    0 0 20px rgba(102, 126, 234, 0.4);
}

.ring {
  position: absolute;
  top: -30px;
  left: -30px;
  width: 160px;
  height: 160px;
  border: 3px solid rgba(0, 255, 136, 0.6);
  border-radius: 50%;
  animation: rotate 10s linear infinite;
  transform-style: preserve-3d;
}

.cover-ring {
  position: absolute;
  top: -20px;
  left: -20px;
  width: 140px;
  height: 140px;
  border: 2px solid rgba(0, 136, 255, 0.4);
  border-radius: 50%;
  animation: rotate 15s linear infinite reverse;
}

.spots {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.spots span {
  position: absolute;
  width: 4px;
  height: 4px;
  background: #ffffff;
  border-radius: 50%;
  animation: twinkle 2s ease-in-out infinite alternate;
}

.spots span:nth-child(1) { top: 20%; left: 30%; animation-delay: 0s; }
.spots span:nth-child(2) { top: 40%; left: 70%; animation-delay: 0.3s; }
.spots span:nth-child(3) { top: 60%; left: 20%; animation-delay: 0.6s; }
.spots span:nth-child(4) { top: 80%; left: 60%; animation-delay: 0.9s; }
.spots span:nth-child(5) { top: 30%; left: 80%; animation-delay: 1.2s; }
.spots span:nth-child(6) { top: 70%; left: 40%; animation-delay: 1.5s; }
.spots span:nth-child(7) { top: 50%; left: 50%; animation-delay: 1.8s; }

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes twinkle {
  from { opacity: 0.3; transform: scale(0.8); }
  to { opacity: 1; transform: scale(1.2); }
}

/* 故障效果 */
.glitch {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 200px;
  height: 200px;
  pointer-events: none;
}

.glitch__line {
  position: absolute;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #00ff88, transparent);
  animation: glitch-line 3s linear infinite;
}

.glitch__line:nth-child(1) { top: 20%; animation-delay: 0s; }
.glitch__line:nth-child(2) { top: 40%; animation-delay: 0.5s; }
.glitch__line:nth-child(3) { top: 60%; animation-delay: 1s; }
.glitch__line:nth-child(4) { top: 80%; animation-delay: 1.5s; }
.glitch__line:nth-child(5) { top: 100%; animation-delay: 2s; }

@keyframes glitch-line {
  0%, 90%, 100% { opacity: 0; transform: scaleX(0); }
  5%, 85% { opacity: 1; transform: scaleX(1); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .number-404 span {
    font-size: 80px;
  }
  
  .planet {
    width: 80px;
    height: 80px;
  }
  
  .ring {
    width: 120px;
    height: 120px;
    top: -20px;
    left: -20px;
  }
  
  .cover-ring {
    width: 100px;
    height: 100px;
    top: -10px;
    left: -10px;
  }
}

@media (max-width: 480px) {
  .number-404 {
    gap: 10px;
  }
  
  .number-404 span {
    font-size: 60px;
  }
}
</style>

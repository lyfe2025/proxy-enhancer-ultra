<template>
  <div class="login-background">
    <div class="bg-pattern"></div>
    <div class="bg-gradient"></div>
  </div>
</template>

<script setup lang="ts">
// 无需逻辑，纯装饰组件
</script>

<style scoped>
.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 1;
}

.bg-pattern {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    radial-gradient(circle at 20% 30%, rgba(0, 255, 136, 0.15) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(0, 255, 136, 0.08) 0%, transparent 40%),
    radial-gradient(circle at 70% 80%, rgba(74, 222, 128, 0.12) 0%, transparent 45%),
    radial-gradient(circle at 30% 70%, rgba(0, 255, 136, 0.06) 0%, transparent 35%),
    radial-gradient(circle at 90% 60%, rgba(74, 222, 128, 0.05) 0%, transparent 30%);
  animation: float 25s ease-in-out infinite;
}

.bg-gradient {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(
    ellipse at center,
    rgba(0, 255, 136, 0.03) 0%,
    rgba(26, 26, 26, 0.8) 50%,
    var(--bg-primary) 100%
  );
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg) scale(1);
  }
  25% {
    transform: translateY(-15px) rotate(0.5deg) scale(1.02);
  }
  50% {
    transform: translateY(-8px) rotate(-0.3deg) scale(0.98);
  }
  75% {
    transform: translateY(-20px) rotate(0.8deg) scale(1.01);
  }
}
</style>

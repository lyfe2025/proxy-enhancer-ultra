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
    radial-gradient(circle at 25% 25%, rgba(120, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 75% 75%, rgba(255, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 75% 25%, rgba(119, 198, 255, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 25% 75%, rgba(198, 255, 119, 0.3) 0%, transparent 50%);
  animation: float 20s ease-in-out infinite;
}

.bg-gradient {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    135deg,
    rgba(67, 56, 202, 0.1) 0%,
    rgba(139, 92, 246, 0.1) 25%,
    rgba(59, 130, 246, 0.1) 50%,
    rgba(16, 185, 129, 0.1) 75%,
    rgba(245, 158, 11, 0.1) 100%
  );
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  33% {
    transform: translateY(-20px) rotate(1deg);
  }
  66% {
    transform: translateY(-10px) rotate(-1deg);
  }
}
</style>

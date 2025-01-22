<template>
  <main class="hero bg-cover font-customFont">
    <div class="flex h-screen w-screen">
      <slot />
    </div>
  </main>
</template>

<script setup></script>

<style scoped>
.hero {
  background: url("../public/images/bg.png") center center/cover no-repeat;
}
</style>

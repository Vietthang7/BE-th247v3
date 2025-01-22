<template>
  <div class="flex h-[60px] w-full bg-gray-200">
    <!-- <Sidebar /> -->

    <!-- Main Content -->
    <div class="flex-1">
      <!-- Top Navigation -->
      <header class="flex items-center justify-between bg-white p-4 shadow">
        <h1 class="text-lg font-semibold text-gray-700">Danh sách nhân sự</h1>
        <div class="flex items-center space-x-4">
          <div class="relative">
            <input
              type="text"
              placeholder="Tìm kiếm"
              class="flex w-full rounded-lg border px-3 py-2 text-sm text-black focus:outline-none focus:ring focus:ring-blue-200"
            />
            <div
              class="pointer-events-none absolute inset-y-0 right-2 flex items-center pl-3"
            >
              <i class="fas fa-magnifying-glass text-gray-400"></i>
            </div>
          </div>

          <NuxtLink to="user">
            <img
              src="../public/images/barking-up-the-wrong-tree.webp"
              alt="Profile"
              class="h-10 w-10 rounded-full"
              href="../Pages/content/user.vue"
            />
          </NuxtLink>
        </div>
      </header>
    </div>
  </div>
</template>

<script setup lang="ts">
import Sidebar from "./sidebar.vue";

/* No special logic required for this static layout */
</script>

<style scoped>
/* Add custom styles if needed */
</style>

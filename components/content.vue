<template>
  <div class="flex h-full w-full overflow-auto rounded-2xl bg-gray-50">
    <!-- Main Content -->
    <div class="flex-1">
      <!-- Content Area -->
      <div class="text-black">
        <main class="flex items-center justify-center p-8">
          <div class="text-center">
            <div class="mb-4 text-blue-500">
              <i class="fas fa-box-open text-6xl"></i>
            </div>
            <p class="text-lg text-gray-700">Page content home</p>
            <button
              class="mt-4 rounded-lg bg-blue-500 px-4 py-2 font-medium text-black hover:bg-blue-600"
              v-on:click="alert('')"
            >
              Thêm mới
            </button>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    alert(message) {
      window.alert(message);
    },
  },
};
</script>

<script setup>
const formRef = ref(null);
const formValue = reactive({
  email: null,
  password: null,
});
</script>

<template>
  <div class="relative mx-auto flex min-h-screen items-center justify-center">
    <!-- Background Decorations -->

    <!-- Login Card -->
    <div class="relative w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
      <div class="mb-6 text-center">
        <h1 class="text-2xl font-bold text-blue-600">
          <img src="../public/images/logo.png" class="mx-auto size-1/2" />
        </h1>
        <p class="mt-5 text-4xl font-bold text-blue-900">Đăng nhập</p>
      </div>
      <form :model="formValue" :rules="rules">
        <div class="relative z-0 mb-4">
          <input
            type="email"
            v-model="formValue.email"
            class="peer block w-full appearance-none border-0 border-b-2 border-gray-300 bg-transparent px-0 py-2.5 text-sm text-gray-900 focus:border-blue-600 focus:outline-none focus:ring-0 dark:border-gray-600 dark:focus:border-blue-500"
            placeholder=" "
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm text-gray-500 duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4 dark:text-gray-400 peer-focus:dark:text-blue-500"
            >Email</label
          >
        </div>
        <div class="relative z-0 mb-4">
          <input
            type="password"
            v-model="formValue.password"
            class="peer block w-full appearance-none border-0 border-b-2 border-gray-300 bg-transparent px-0 py-2.5 text-sm text-gray-900 focus:border-blue-600 focus:outline-none focus:ring-0 dark:border-gray-600 dark:focus:border-blue-500"
            placeholder=" "
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm text-gray-500 duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4 dark:text-gray-400 peer-focus:dark:text-blue-500"
            >Mật khẩu</label
          >
        </div>
        <div class="mb-4 flex items-center justify-between">
          <label class="flex items-center text-gray-600">
            <input
              type="checkbox"
              v-model="rememberMe"
              class="rounded border-gray-300 text-blue-500 accent-blue-500 focus:ring-blue-300"
            />
            <span class="ml-2">Ghi nhớ đăng nhập</span>
          </label>
          <NuxtLink to="/forgot" class="text-sm text-blue-500"
            >Quên mật khẩu?</NuxtLink
          >
        </div>
        <button
          type="submit"
          class="w-full rounded-lg bg-blue-500 px-4 py-2 font-semibold text-white duration-75 hover:bg-blue-600 focus:outline-none focus:ring focus:ring-blue-300"
        >
          Đăng nhập
        </button>
      </form>
      <div class="mt-4 text-center text-gray-600">
        Bạn chưa có tài khoản?
        <NuxtLink to="/auth/register" class="text-blue-500"
          >Đăng ký ngay!</NuxtLink
        >
      </div>
    </div>
  </div>
</template>

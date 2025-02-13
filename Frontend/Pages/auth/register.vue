<script setup>
import { NuxtLink } from "#components";

definePageMeta({
  layout: "auth",
});
</script>

<template>
  <div
    class="relative mx-auto flex min-h-screen w-2/6 items-center justify-center"
  >
    <div class="mx-auto rounded-3xl bg-white bg-opacity-60 shadow">
      <n-form
        :label-width="200"
        class="p-6"
        label-align="left"
        require-mark-placement="right"
      >
        <div class="flex w-full items-center justify-center">
          <img src="@/public/images/log0.png" class="mx-auto size-1/2" />
        </div>
        <div class="mt-8 flex justify-center">
          <h1 class="text-3xl font-bold text-[#133D85]">Đăng ký</h1>
        </div>
        <n-form-item path="name" class="relative z-0 -mb-8">
          <input
            type="text"
            class="peer block w-full appearance-none border-0 border-b-2 border-[#4D6FA8] bg-transparent px-0 py-2.5 text-sm text-[#4D6FA8] focus:border-blue-600 focus:outline-none focus:ring-0"
            placeholder=" "
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm font-medium text-[#4D6FA8] duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4"
            >Họ và tên</label
          >
        </n-form-item>
        <n-form-item path="email" class="relative z-0 -mb-8">
          <input
            type="email"
            class="peer block w-full appearance-none border-0 border-b-2 border-[#4D6FA8] bg-transparent px-0 py-2.5 text-sm text-[#4D6FA8] focus:border-blue-600 focus:outline-none focus:ring-0"
            placeholder=" "
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm font-medium text-[#4D6FA8] duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4"
            >Email</label
          >
        </n-form-item>
        <n-form-item path="phone" class="relative z-0 -mb-8">
          <input
            type="tel"
            class="peer block w-full appearance-none border-0 border-b-2 border-[#4D6FA8] bg-transparent px-0 py-2.5 text-sm text-[#4D6FA8] focus:border-blue-600 focus:outline-none focus:ring-0"
            placeholder=" "
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm font-medium text-[#4D6FA8] duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4"
            >Số điện thoại</label
          >
        </n-form-item>
        <n-form-item path="password" class="relative z-0 -mb-8">
          <input
            type="password"
            class="peer block w-full appearance-none border-0 border-b-2 border-[#4D6FA8] bg-transparent px-0 py-2.5 text-sm text-[#4D6FA8] focus:border-blue-600 focus:outline-none focus:ring-0"
            placeholder=" "
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm font-medium text-[#4D6FA8] duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4"
            >Nhập mật khẩu</label
          >
        </n-form-item>
        <n-form-item path="repassword" class="relative z-0">
          <input
            type="password"
            class="peer block w-full appearance-none border-0 border-b-2 border-[#4D6FA8] bg-transparent px-0 py-2.5 text-sm text-[#4D6FA8] focus:border-blue-600 focus:outline-none focus:ring-0"
            placeholder=" "
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm font-medium text-[#4D6FA8] duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4"
            >Nhập lại mật khẩu</label
          >
        </n-form-item>
        <div class="flex items-center justify-between pb-4"></div>
        <div class="flex flex-row justify-center">
          <n-button
            type="info"
            class="h-[40px] w-full rounded-xl text-base font-medium"
          >
            Đăng ký
          </n-button>
        </div>
        <div class="relative mt-5 text-center text-[#133D85]">
          Bạn đã có tài khoản?
          <NuxtLink
            to="login"
            class="text-primary-600 font-small text-[#00A2EB] hover:underline"
          >
            Đăng nhập ngay
          </NuxtLink>
        </div>
      </n-form>
    </div>
  </div>
  <!-- </div> -->
</template>

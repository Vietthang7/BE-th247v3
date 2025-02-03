<script setup>
definePageMeta({
  layout: "auth",
});
</script>

<template>
  <div
    class="relative mx-auto flex min-h-screen w-1/4 items-center justify-center"
  >
    <div class="mx-auto rounded-3xl bg-white bg-opacity-60 shadow">
      <n-form
        :label-width="200"
        class="p-6"
        label-align="left"
        require-mark-placement="right"
      >
        <div class="flex w-full">
          <img src="@/public/images/log0.png" class="mx-auto size-3/4" />
        </div>
        <div class="relative mt-8 text-center">
          <h1 class="text-3xl font-bold text-[#133D85]">Nhập mã xác thực</h1>
          <p class="mt-5 text-lg text-blue-900">
            Bạn vui lòng nhập mã xác thực được gửi đến email đăng ký
          </p>
        </div>
        <n-form-item path="code" class="relative z-0">
          <input
            type="text"
            class="peer block w-full appearance-none border-0 border-b-2 border-[#4D6FA8] bg-transparent px-0 py-2.5 text-sm text-[#4D6FA8] focus:border-blue-600 focus:outline-none focus:ring-0"
            placeholder=" "
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm font-medium text-[#4D6FA8] duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4"
            >Mã xác thực</label
          >
        </n-form-item>
        <div class="flex flex-row justify-center">
          <n-button
            type="info"
            class="h-[40px] w-full rounded-xl text-base font-medium"
          >
            Xác nhận
          </n-button>
        </div>
        <div class="mt-4 text-center">
          Không nhận được mã xác thực?
          <Nuxtlink
            to="register"
            class="text-primary-600 font-small cursor-pointer text-[#00A2EB] hover:underline"
          >
            Gửi lại
          </Nuxtlink>
        </div>
      </n-form>
    </div>
  </div>
</template>

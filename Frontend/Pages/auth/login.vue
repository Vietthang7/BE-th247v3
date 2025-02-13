<script setup>
definePageMeta({
  layout: "auth",
});

const formValue = ref({
  email: "",
  password: "",
});
import axios from "axios";
import { useRouter } from "vue-router";

const isLoading = ref(false);
const router = useRouter();

const handleSubmit = async () => {
  isLoading.value = true;
  try {
    const response = await axios.post("/api/auth/login", {
      email: formValue.value.email,
      password: formValue.value.password,
    });
    if (response.data.success) {
      router.push("/");
    } else {
      alert("Invalid email or password");
    }
  } catch (error) {
    console.error("Error during login:", error);
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="mx-auto my-auto h-2/3 w-1/4 rounded-2xl">
    <div class="mx-auto rounded-3xl bg-white bg-opacity-60 shadow">
      <n-form
        :label-width="200"
        class="p-6"
        label-align="left"
        require-mark-placement="right"
        :model="formValue"
      >
        <div class="flex w-full items-center justify-center">
          <img src="@/public/images/log0.png" class="mx-auto size-3/4" />
        </div>
        <div class="mt-8 flex justify-center">
          <h1 class="text-3xl font-bold text-[#133D85]">Đăng nhập</h1>
        </div>
        <n-form-item path="email" class="relative z-0 -mb-8">
          <input
            type="email"
            class="peer block w-full appearance-none border-0 border-b-2 border-[#4D6FA8] bg-transparent px-0 py-2.5 text-sm text-[#4D6FA8] focus:border-blue-600 focus:outline-none focus:ring-0"
            placeholder=" "
            v-model="formValue.email"
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm font-medium text-[#4D6FA8] duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4"
            >Email</label
          >
        </n-form-item>
        <n-form-item path="password" type="password" class="relative z-0 mb-5">
          <input
            type="password"
            class="peer block w-full appearance-none border-0 border-b-2 border-[#4D6FA8] bg-transparent px-0 py-2.5 text-sm text-[#4D6FA8] focus:border-blue-600 focus:outline-none focus:ring-0"
            placeholder=" "
            :value="formValue.password"
          />
          <label
            class="absolute top-3 -z-10 origin-[0] -translate-y-6 scale-75 transform text-sm font-medium text-[#4D6FA8] duration-300 peer-placeholder-shown:translate-y-0 peer-placeholder-shown:scale-100 peer-focus:start-0 peer-focus:-translate-y-6 peer-focus:scale-75 peer-focus:font-medium peer-focus:text-blue-600 rtl:peer-focus:translate-x-1/4"
            >Nhập mật khẩu</label
          >
        </n-form-item>
        <div class="flex items-center justify-between pb-4">
          <div class="flex items-start">
            <div class="flex-initial">
              <n-checkbox></n-checkbox>
              <span
                class="text-primary-600 font-small ml-2 text-sm text-[#133D85]"
                >Ghi nhớ đăng nhập</span
              >
            </div>
          </div>
          <NuxtLink
            to="forgot"
            class="text-primary-600 font-small text-sm text-[#00A2EB] hover:underline"
          >
            Quên mật khẩu?
          </NuxtLink>
        </div>
        <div class="flex flex-row justify-center">
          <n-button
            type="info"
            class="h-[40px] w-full rounded-xl text-base font-medium"
            :loading="isLoading"
            @click="handleSubmit"
          >
            Đăng nhập
          </n-button>
        </div>
        <div class="relative mt-5 text-center text-[#133D85]">
          Bạn chưa có tài khoản?
          <NuxtLink
            to="register"
            class="text-primary-600 font-small text-[#00A2EB] hover:underline"
          >
            Đăng ký ngay
          </NuxtLink>
        </div>
      </n-form>
    </div>
  </div>
</template>

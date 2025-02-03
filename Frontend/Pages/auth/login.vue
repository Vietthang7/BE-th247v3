<script setup>
definePageMeta({
  layout: "auth",
});

const { restAPI } = useApi();
const formRef = ref(null);
const formValue = reactive({
  email: null,
  password: null,
});

const autoLogin = ref(true);

function handleKeyup(e) {
  switch (e.key) {
    case "Enter":
      handleSubmit(e);
      break;
  }
}

const isLoading = ref(false);
const handleSubmit = async (e) => {
  if (isLoading.value) return;
  e.preventDefault();
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      const { email, password } = formValue;
      isLoading.value = true;
      const body = { email, password };
      const { data: resVerify, error } = await restAPI.cms.adminLogin({
        body,
      });
      if (resVerify.value?.status) {
        if (resVerify.value?.data?.domain) {
          await navigateTo(
            `${resVerify.value?.data?.domain}/callback?accessToken=${resVerify.value?.data?.token}`,
            {
              external: true,
            },
          );
        }
      }
      isLoading.value = false;
    }
  });
};

onMounted(() => {
  window.addEventListener("keyup", handleKeyup);
});
onBeforeUnmount(() => {
  window.removeEventListener("keyup", handleKeyup);
});
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
        :rules="rules"
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
            :value="formValue.email"
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
              <n-checkbox :checked="autoLogin"></n-checkbox>
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

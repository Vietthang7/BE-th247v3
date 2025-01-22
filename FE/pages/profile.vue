<script setup>
const displaypage = ref(true);
definePageMeta({
  layout: "page",
});
</script>

<template>
  <div class="flex min-h-screen w-full flex-row">
    <div
      class="mx-5 my-5 h-full w-full flex-grow overflow-auto rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
    >
      <!-- Tabs -->
      <div class="border-b border-gray-200">
        <ul class="-mb-px flex flex-wrap font-bold text-gray-500">
          <li class="me-2 cursor-pointer duration-500 hover:text-[#00A2EB]">
            <button
              :class="{
                'border-b-2 border-[#00A2EB] text-[#00A2EB]': displaypage,
              }"
              class="inline-block rounded-t-lg p-4"
              @click="
                if (displaypage == true) {
                } else {
                  displaypage = !displaypage;
                }
              "
            >
              Thông tin tài khoản
            </button>
          </li>
          <li class="me-2 cursor-pointer duration-500 hover:text-[#00A2EB]">
            <button
              :class="{
                'border-b-2 border-[#00A2EB] text-[#00A2EB]': !displaypage,
              }"
              class="inline-block rounded-t-lg p-4"
              @click="
                if (displaypage == false) {
                } else {
                  displaypage = !displaypage;
                }
              "
            >
              Đổi mật khẩu
            </button>
          </li>
        </ul>
      </div>
      <div>
        <AccsettingProfileInfo
          class="w-3/4 md:w-5/6 2xl:w-3/5"
          v-if="displaypage"
        ></AccsettingProfileInfo>
        <AccsettingChangePass
          class="w-3/4 md:w-5/6 2xl:w-3/5"
          v-if="!displaypage"
        ></AccsettingChangePass>
      </div>
    </div>
  </div>
</template>

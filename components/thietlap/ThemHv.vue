<script lang="ts">
import { defineComponent, ref } from "vue";

import type { GlobalThemeOverrides } from "naive-ui";

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: "#FF8C00",
    textColor1: "#4D6FA8",
  },
};
export default defineComponent({
  setup() {
    return {
      value: ref(null),
      options: [
        {
          label: "Trạng thái",
          value: "song0",
        },
        {
          label: "Trạng thái1",
          value: "song1",
        },
        {
          label: "Trạng thái2",
          value: "song2",
        },
      ],
      themeOverrides,
    };
  },
});
</script>
<template>
  <div
    class="mx-5 my-5 w-full flex-grow rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
  >
    <div
      class="mt-5 grid h-4/6 w-5/6 grid-flow-col auto-rows-min grid-rows-6 gap-x-5 text-lg font-bold"
    >
      <div>
        <label class="block">Họ tên nhân sự *</label>
        <n-input
          size="large"
          required
          class="rounded-2xl focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
          placeholder="Nguyễn Tú Nam"
        />
      </div>
      <div>
        <label class="block">Email *</label>
        <input
          type="text"
          class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
          placeholder="thuchisue@gmail.com"
          required
        />
      </div>
      <div>
        <label class="block">Tên tài khoản *</label>
        <input
          type="text"
          class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
          placeholder="thuchisue"
          required
        />
      </div>
      <div class="rows col-span-2 row-span-2">
        <label class="block">Ghi chú</label>
        <textarea
          class="h-2/3 w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
          placeholder="Giới thiệu"
          required
        />
      </div>
      <div class="col-start-2 row-start-6 flex-row gap-10">
        <n-space class="gap-10">
          <n-button round class="h-14 w-32 text-2xl text-blue-400"
            >Hủy</n-button
          >
          <n-button round type="info" class="h-14 w-32 text-2xl">
            Lưu
          </n-button>
        </n-space>
      </div>
      <div class="col-start-2 row-start-1">
        <label class="block">Trạng thái học viên *</label>
        <select
          name=""
          id=""
          class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
        >
          <option value="">1</option>
          <option value="">2</option>
          <option value="">3</option>
        </select>
      </div>
      <div class="col-start-2 row-start-2">
        <label class="block">Số điện thoại *</label>
        <input
          type="text"
          class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
          placeholder="0123456789"
          required
        />
      </div>
      <div class="col-start-2 row-start-3">
        <label class="block">Mật khẩu *</label>
        <input
          type="password"
          class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
          placeholder="abc123456"
          required
        />
      </div>
      <div class="row-span-5 grid grid-rows-subgrid gap-4">
        <div class="row-start-3 my-9 flex">
          <label class="mr-3">Trạng thái hoạt động:</label>
          <label class="cursor-pointer items-center">
            <input type="checkbox" value="" class="peer sr-only" />
            <div
              class="peer relative h-6 w-11 rounded-full bg-gray-200 after:absolute after:start-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white"
            ></div>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

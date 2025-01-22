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
    class="mx-5 my-5 h-full w-full flex-grow rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
  >
    <!-- <div
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
    </div> -->

    <n-grid class="my-5 w-full" :x-gap="70" cols="1" responsive="screen">
      <n-grid-item span="1">
        <n-grid :x-gap="30" cols="1 m:3" responsive="screen">
          <n-gi span="1 m:1">
            <n-form-item label="Họ tên nhân sự *" path="name">
              <n-input placeholder="Họ tên" />
            </n-form-item>
          </n-gi>
          <n-gi span="1 m:1">
            <n-form-item label="Trạng thái học viên *" path="status">
              <n-select
                v-model="value"
                :options="options"
                placeholder="Chọn trạng thái học viên"
              />
            </n-form-item>
          </n-gi>
          <n-gi span="1"> </n-gi>
          <n-gi span="1 m:1">
            <n-form-item label="Email *" path="email">
              <n-input placeholder="Email" />
            </n-form-item>
          </n-gi>
          <n-gi span="1 m:1">
            <n-form-item label="Số điện thoại *" path="phone">
              <n-input placeholder="0123456789" />
            </n-form-item>
          </n-gi>
          <n-gi span="1"> </n-gi>
          <n-gi span="1">
            <n-form-item label="Tên tài khoản *" path="username">
              <n-input placeholder="thuchisue" />
            </n-form-item>
          </n-gi>
          <n-gi span="1">
            <n-form-item label="Mật khẩu *" path="password">
              <n-input placeholder="******" />
            </n-form-item>
          </n-gi>
          <n-gi class="mt-7 text-2xl">
            <n-form-item
              label-placement="left"
              label="Trạng thái hoạt động: "
              path="switchValue"
            >
              <n-switch />
            </n-form-item>
          </n-gi>
          <n-gi span="1 m:2">
            <n-form-item label="Ghi chú" path="note">
              <n-input
                placeholder="Giới thiệu"
                type="textarea"
                :autosize="{
                  minRows: 3,
                  maxRows: 5,
                }"
              />
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-grid-item>
    </n-grid>
    <div class="gap-25px flex w-full items-center justify-center">
      <n-space>
        <n-button
          class="flex h-10 w-20 items-center justify-center text-center text-xl text-blue-400"
          >Hủy</n-button
        >
        <n-button
          type="info"
          class="flex h-10 w-20 items-center justify-center text-center text-xl"
        >
          Lưu
        </n-button>
      </n-space>
    </div>
  </div>
</template>

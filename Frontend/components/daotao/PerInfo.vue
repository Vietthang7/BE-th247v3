<script lang="ts">
import { defineComponent, ref } from "vue";
import type { CSSProperties } from "vue";
export default defineComponent({
  setup() {
    return {
      railStyle: ({
        focused,
        checked,
      }: {
        focused: boolean;
        checked: boolean;
      }) => {
        const style: CSSProperties = {};
        if (checked) {
          style.background = "#2080f0";
          if (focused) {
            style.boxShadow = "0 0 0 2px #2080f0";
          }
        } else {
        }
        return style;
      },

      timestamp: ref(""),
      value: ref(null),
      options: [
        {
          label: "Nam",
          value: "Nam",
        },
        {
          label: "Nữ",
          value: "Nữ",
        },
        {
          label: "Khác",
          value: "song2",
        },
      ],
    };
  },
});
</script>

<template>
  <div class="px-5">
    <n-grid :x-gap="30" cols="1 m:2" responsive="screen" class="my-5">
      <n-gi>
        <n-form-item label="Họ và tên học viên *" path="name">
          <n-input placeholder="Điền tên nhóm học viên" />
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-form-item
          class="mt-8"
          label="Trạng thái tài khoản"
          label-placement="left"
          path="switchValue"
        >
          <n-switch :rail-style="railStyle" />
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-form-item label="Giới tính">
          <n-select
            v-model="value"
            :options="options"
            placeholder="Chọn giới tính"
          />
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-form-item label="Ngày sinh">
          <n-date-picker
            class="w-full"
            v-model="timestamp"
            type="date"
            placeholder="Chọn ngày"
          />
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-form-item label="Email *" path="email">
          <n-input placeholder="Nhập email" />
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-form-item label="Số điện thoại" path="phone">
          <n-input placeholder="Nhập số điện thoại" type="text" />
        </n-form-item>
      </n-gi>
      <n-gi span="1 m:2">
        <n-form-item label="Địa chỉ" path="address">
          <n-grid :x-gap="30" cols="1 m:2" responsive="screen">
            <n-gi>
              <n-select placeholder="Tỉnh/Thành phố" />
            </n-gi>
            <n-gi>
              <n-select placeholder="Quận/Huyện" />
            </n-gi>
            <n-gi span="1 m:2" class="mt-2">
              <n-input placeholder="Nhập địa chỉ" />
            </n-gi>
          </n-grid>
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-form-item label="Nguồn khách hàng *">
          <n-select placeholder="Chọn nguồn khách hàng" />
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-form-item label="Kênh liên hệ">
          <n-select placeholder="Chọn kênh liên hệ" />
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-button round type="info" class="h-12 w-52 rounded-2xl text-lg">
          Lưu
        </n-button>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup>
definePageMeta({
  layout: "page",
});
</script>
<template>
  <div
    class="mx-5 my-5 w-full flex-grow rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
  >
    <n-grid class="my-10 w-full" :x-gap="70" cols="1 m:5" responsive="screen">
      <n-grid-item span="1 m:3">
        <n-gi span="1">
          <n-form
            ref="formRef"
            :model="model"
            :rules="rules"
            label-placement="left"
            label-align="left"
            require-mark-placement="right-hanging"
            :size="size"
            label-width="auto"
            :style="{
              maxWidth: '640px',
            }"
          >
            <n-form-item label="Tên ca làm *" path="inputValue">
              <n-input placeholder="Input" />
            </n-form-item>
            <n-form-item label="Giờ làm *" path="inputValue">
              <n-input placeholder="Input" />
            </n-form-item>
            <n-form-item label="Trạng thái hoạt động" path="switchValue">
              <n-switch />
            </n-form-item>
          </n-form>
        </n-gi>
      </n-grid-item>
    </n-grid>
    <div class="gap-25px flex w-2/3 justify-center">
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

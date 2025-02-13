<script setup>
const disabledInput = ref(false);
const mobile = ref(false);
</script>
<template>
  <div
    class="mx-5 my-5 w-full flex-grow rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
  >
    <div class="text-md mx-5 my-10">
      <h1 class="text-4xl font-bold">Hoc Vien</h1>
      <n-grid class="my-5 w-full" :x-gap="70" cols="1 m:5" responsive="screen">
        <n-grid-item span="1 m:3">
          <n-grid :x-gap="30" cols="1 m:2" responsive="screen">
            <n-gi span="1">
              <n-form-item label="Email" path="email">
                <n-input
                  :class="{ inputbg: !disabledInput }"
                  :readonly="true"
                  placeholder="Nhập email"
                />
              </n-form-item>
            </n-gi>

            <n-gi span="1">
              <n-form-item label="Số điện thoại" path="phone">
                <n-input
                  :class="{ inputbg: !disabledInput }"
                  :readonly="true"
                  placeholder="Nhập số điện thoại"
                />
              </n-form-item>
            </n-gi>

            <n-gi span="1 m:2">
              <n-form-item label="Địa chỉ" path="address">
                <n-input
                  :class="{ inputbg: !disabledInput }"
                  :readonly="true"
                  placeholder="Nhập địa chỉ"
                />
              </n-form-item>
            </n-gi>

            <n-gi span="1 m:2">
              <n-form-item label="Điểm nổi bật" path="highlight">
                <div
                  v-if="disabledInput"
                  v-html="'Nhập điểm nổi bật'"
                  class="rounded-15px w-full border border-[#E7EDFF] py-4 pl-5 pr-8 text-[#8BA3CB]"
                ></div>
                <n-input
                  placeholder="Textarea"
                  type="textarea"
                  :autosize="{
                    minRows: 3,
                    maxRows: 5,
                  }"
                />
              </n-form-item>
            </n-gi>

            <n-gi span="1 m:2">
              <n-form-item label="Chương trình dạy" path="highlight">
                <div
                  v-if="disabledInput"
                  v-html="'Nhập điểm nổi bật'"
                  class="rounded-15px w-full border border-[#E7EDFF] py-4 pl-5 pr-8 text-[#8BA3CB]"
                ></div>
                <n-input
                  placeholder="Textarea"
                  type="textarea"
                  :autosize="{
                    minRows: 3,
                    maxRows: 5,
                  }"
                />
              </n-form-item>
            </n-gi>

            <n-gi span="1 m:2">
              <n-form-item label="Mô tả chung" path="highlight">
                <div
                  v-if="disabledInput"
                  v-html="'Nhập điểm nổi bật'"
                  class="rounded-15px w-full border border-[#E7EDFF] py-4 pl-5 pr-8 text-[#8BA3CB]"
                ></div>
                <n-input
                  placeholder="Textarea"
                  type="textarea"
                  :autosize="{
                    minRows: 3,
                    maxRows: 5,
                  }"
                />
              </n-form-item>
            </n-gi>

            <n-gi span="1 m:2">
              <n-form-item label="Độ tuổi" path="highlight">
                <div
                  v-if="disabledInput"
                  v-html="'Nhập điểm nổi bật'"
                  class="rounded-15px w-full border border-[#E7EDFF] py-4 pl-5 pr-8 text-[#8BA3CB]"
                ></div>
                <n-input
                  placeholder="Textarea"
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

        <!-- right -->
        <n-grid-item span="1 m:2">
          <!-- Logo, Favicon, Cover on desktop -->
          <template v-if="!mobile">
            <n-form-item label="Logo" path="logo">
              <NuxtImg
                :size="3.1"
                :disabled="disabledInput"
                :hideSize="isStudent"
              />
            </n-form-item>
            <n-form-item label="Favicon" path="favicon">
              <NuxtImg
                :size="1"
                :disabled="disabledInput"
                :hideSize="isStudent"
              />
            </n-form-item>
            <n-form-item label="Ảnh bìa" path="cover">
              <NuxtImg
                :square="false"
                :width="400"
                :height="150"
                :rounded="15"
                :disabled="disabledInput"
                :hideSize="isStudent"
              />
            </n-form-item>
          </template>
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
  </div>
</template>

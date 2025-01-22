<script lang="ts">
import { defineComponent, ref } from "vue";

export default defineComponent({
  setup() {
    return {
      components: {},
      value: ref(null),
      options: [
        {
          label: "Everybody's Got Something to Hide Except Me and My Monkey",
          value: "song0",
          disabled: true,
        },
        {
          label: "Drive My Car",
          value: "song1",
        },
        {
          label: "Norwegian Wood",
          value: "song2",
        },
        {
          label: "You Won't See",
          value: "song3",
          disabled: true,
        },
        {
          label: "Nowhere Man",
          value: "song4",
        },
        {
          label: "Think For Yourself",
          value: "song5",
        },
        {
          label: "The Word",
          value: "song6",
        },
        {
          label: "Michelle",
          value: "song7",
          disabled: true,
        },
      ],
    };
  },
});
</script>
<template>
  <div
    class="mx-5 my-5 w-full flex-grow rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
  >
    <n-grid class="my-10 w-full" :x-gap="50" cols="1 m:2" responsive="screen">
      <n-grid-item span="1 m:1">
        <n-grid :x-gap="30" cols="1 m:4" responsive="screen">
          <n-gi span="2">
            <n-form>
              <n-input type="text" placeholder="Tìm kiếm cơ bản" />
            </n-form>
          </n-gi>
          <n-gi span="1">
            <n-form>
              <n-select v-model="value" :options="options" />
            </n-form>
          </n-gi>
          <n-gi span="1">
            <n-form>
              <n-select v-model="value" :options="options" />
            </n-form>
          </n-gi>
        </n-grid>
      </n-grid-item>
      <n-grid-item span="1 m:1">
        <n-grid :x-gap="5" cols="1 m:2" responsive="screen">
          <n-gi span="1">
            <n-float-button :left="0" :bottom="0" shape="square">
              <n-icon>
                <CashIcon />
              </n-icon>
            </n-float-button>
          </n-gi>
          <n-gi span="1">
            <n-button :right="0" icon-placement="right">
              <ion-icon name="add-outline"></ion-icon>
              Xuất file
            </n-button>
          </n-gi>
        </n-grid>
      </n-grid-item>
    </n-grid>
  </div>
</template>

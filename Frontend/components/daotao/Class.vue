<script lang="ts">
import type { DataTableColumns, DataTableRowKey } from "naive-ui";
import { defineComponent, ref, reactive, h } from "vue";
import { NDataTable } from "naive-ui";
interface RowData {
  chtrinh: string;
  mhoc: string;
  lhoc: string;
  tghoc: string;
  sbhoc: number;
  shvien: number;
}
export default defineComponent({
  setup() {
    const paginationReactive = reactive({
      page: 2,
      pageSize: 5,
      showSizePicker: true,
      pageSizes: [3, 5, 7],
      onChange: (page: number) => {
        paginationReactive.page = page;
      },
      onUpdatePageSize: (pageSize: number) => {
        paginationReactive.pageSize = pageSize;
        paginationReactive.page = 1;
      },
    });
    return {
      data,
      columns: [
        {
          title: "Chương trình học",
          key: "chtrinh",
        },
        {
          title: "Môn học",
          key: "mhoc",
        },
        {
          title: "Lớp học",
          key: "lhoc",
        },
        {
          title: "Thời gian học",
          key: "tghoc",
        },
        {
          title: "Số buổi học",
          key: "sbhoc",
        },
        {
          title: "Số học viên trong lớp",
          key: "shvien",
        },
      ],
      pagination: paginationReactive,
    };
  },
});

const data = Array.from({ length: 46 }).map((_, index) => ({
  chtrinh: "Thiết kế",
  mhoc: "Design Thinking",
  lhoc: "Lớp Design Thinking K01(21/09/2024)",
  tghoc: "(17:00-19:00)",
  sbhoc: `${index}`,
  shvien: `${index}`,
}));
</script>
<template>
  <div class="mt-5 px-5">
    <n-grid
      :x-gap="30"
      :y-gap="8"
      cols="1 m:5"
      responsive="screen"
      class="my-5"
    >
      <n-gi span="1 m:2">
        <label>Số môn học chưa xếp lớp: </label>
      </n-gi>
      <n-gi span="1" offset="2" class="justify-self-end">
        <n-button round type="info" class="h-12 w-40 rounded-2xl text-lg">
          Xếp lớp
          <i class="fa-regular fa-flag pl-3"></i>
        </n-button>
      </n-gi>
      <n-gi span="1 m:5">
        <n-data-table
          ref="dataTableInst"
          :bordered="false"
          :single-line="false"
          :columns="columns"
          :data="data"
          :pagination="pagination"
          :max-height="350"
        />
      </n-gi>
      <n-gi span="1 m:1" offset="2" class="justify-center">
        <n-button round type="info" class="h-12 w-52 rounded-2xl text-lg">
          Lưu
        </n-button>
      </n-gi>
    </n-grid>
  </div>
</template>

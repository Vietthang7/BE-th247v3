<script lang="ts">
import type { DataTableColumns } from "naive-ui";
import { defineComponent, ref, h } from "vue";
import { NButton, NTag, useMessage } from "naive-ui";
interface RowData {
  key: number;
  ca: string;
  coso: string;
  time: string;
  status: string[];
}

function createColumns({
  sendMail,
}: {
  sendMail: (rowData: RowData) => void;
}): DataTableColumns<RowData> {
  return [
    {
      title: "STT",
      key: "key",
    },
    {
      title: "Ca làm",
      key: "ca",
    },
    {
      title: "Cơ sở áp dụng",
      key: "coso",
    },
    {
      title: "Thời gian",
      key: "time",
    },
    {
      title: "Trạng thái",
      key: "status",
      render(row) {
        const tags = row.status.map((tagKey) => {
          return h(
            NTag,
            {
              style: {
                marginRight: "6px",
              },
              type: "info",
              bordered: false,
            },
            {
              default: () => tagKey,
            },
          );
        });
        return tags;
      },
    },
    {
      title: "Hành động",
      key: "Edit",
      render(row) {
        return h(
          NButton,
          {
            size: "small",
            onClick: () => sendMail(row),
          },
          {
            default: () => "Edit",
          },
        );
      },
    },
    {
      title: "Hành động",
      key: "Del",
      render(row) {
        return h(
          NButton,
          {
            size: "small",
            onClick: () => sendMail(row),
          },
          {
            default: () => "del",
          },
        );
      },
    },
  ];
}

function createData(): RowData[] {
  return [
    {
      key: 1,
      ca: "Sáng",
      coso: "Cơ sở 1",
      time: "08:00 - 12:00",
      status: ["Hoạt động"],
    },
    {
      key: 2,
      ca: "Chiều",
      coso: "Cơ sở 2",
      time: "01:00 - 05:00",
      status: ["Dừng hoạt động"],
    },
    {
      key: 3,
      ca: "Sáng",
      coso: "Cơ sở 3",
      time: "08:00 - 12:00",
      status: ["Hoạt động"],
    },
  ];
}

export default defineComponent({
  setup() {
    return {
      data: createData(),
      columns: createColumns({
        sendMail(rowData) {
          message.info(`send mail to ${rowData.ca}`);
        },
      }),
      pagination: {
        pageSize: 10,
      },
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
        },
      ],
    };
  },
});
</script>
<template>
  <div
    class="mx-5 my-5 h-full w-full flex-grow rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
  >
    <n-grid class="my-10 w-full" :x-gap="50" cols="1 m:4" responsive="screen">
      <n-grid-item span="1 m:3">
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
        <n-grid
          :x-gap="5"
          class="justify-items-end"
          cols="1 m:3"
          responsive="screen"
        >
          <n-gi span="1 m:2">
            <n-button icon-placement="right" :right="0">
              <ion-icon name="add-outline"></ion-icon>
              Thêm ca làm
            </n-button>
          </n-gi>
          <n-gi span="1">
            <n-button :right="0" icon-placement="right">
              <ion-icon name="add-outline"></ion-icon>
              Xuất file
            </n-button>
          </n-gi>
        </n-grid>
      </n-grid-item>
      <n-grid-item span="4" class="my-5">
        <n-data-table
          :bordered="false"
          :single-line="false"
          :columns="columns"
          :data="data"
          :pagination="pagination"
        />
      </n-grid-item>
    </n-grid>
  </div>
</template>

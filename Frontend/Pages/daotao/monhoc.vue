<!-- <template>
    <div class="flex h-full w-full overflow-auto rounded-2xl bg-gray-50">
     
      <div class="flex-1">
   
        <div class="text-black">
          <main class="flex items-center justify-center p-8">
            <div class="text-center">
              <div class="mb-4 text-blue-500">
                <i class="fas fa-box-open text-6xl"></i>
              </div>
              <p class="text-lg text-gray-700">Page nhansu home</p>
              <button
                class="mt-4 rounded-lg bg-blue-500 px-4 py-2 font-medium text-black hover:bg-blue-600"
              >
                Thêm mới
              </button>
            </div>
          </main>
        </div>
      </div>
    </div>
  </template> -->
<!-- 
<script  lang="ts">
import type { DataTableColumns } from "naive-ui";
import { defineComponent, ref, h } from "vue";
import { NButton, NTag, useMessage } from "naive-ui";
interface RowData {
  key: number;
  name: string;
  age: number;
  address: string;
  tags: string[];
}

function createColumns({
  sendMail,
}: {
  sendMail: (rowData: RowData) => void;
}): DataTableColumns<RowData> {
  return [
    {
      title: "Name",
      key: "name",
    },
    {
      title: "Age",
      key: "age",
    },
    {
      title: "Address",
      key: "address",
    },
    {
      title: "Tags",
      key: "tags",
      render(row) {
        const tags = row.tags.map((tagKey) => {
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
      title: "Action",
      key: "actions",
      render(row) {
        return h(
          NButton,
          {
            size: "small",
            onClick: () => sendMail(row),
          },
          { default: () => "Send Email" },
        );
      },
    },
  ];
}

function createData(): RowData[] {
  return [
    {
      key: 0,
      name: "John Brown",
      age: 32,
      address: "New York No. 1 Lake Park",
      tags: ["nice", "developer"],
    },
    {
      key: 1,
      name: "Jim Green",
      age: 42,
      address: "London No. 1 Lake Park",
      tags: ["wow"],
    },
    {
      key: 2,
      name: "Joe Black",
      age: 32,
      address: "Sidney No. 1 Lake Park",
      tags: ["cool", "teacher"],
    },
  ];
}

export default defineComponent({
  setup() {
    const message = useMessage();
    return {
      data: createData(),
      columns: createColumns({
        sendMail(rowData) {
          message.info(`send mail to ${rowData.name}`);
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
</script> -->
<template>
  <div
    class="mx-5 my-5 w-full flex-grow rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
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
              <!-- <n-select v-model="value" :options="options" /> -->
            </n-form>
          </n-gi>
          <n-gi span="1">
            <n-form>
              <!-- <n-select v-model="value" :options="options" /> -->
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
      <n-grid-item span="4" class="">
        <!-- <n-data-table
          :bordered="false"
          :single-line="false"
          :columns="columns"
          :data="data"
          :pagination="pagination"
        /> -->
      </n-grid-item>
    </n-grid>
  </div>
</template>

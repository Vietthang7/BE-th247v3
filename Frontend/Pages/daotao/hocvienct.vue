<script lang="ts">
import type { DataTableColumns, DataTableRowKey, SelectOption } from "naive-ui";
import { defineComponent, ref, h } from "vue";
import { c, NButton, NDataTable, NDropdown } from "naive-ui";
import { onMounted } from "vue";

export default defineComponent({
  setup() {
    const delModal = ref(false);
    const checkedRowKeysRef = ref<DataTableRowKey[]>([]);
    const dataTableInstRef = ref<InstanceType<typeof NDataTable> | null>(null);
    const activeItem = ref("Tất cả trạng thái");
    const accountStatus = ref("");
    function filterStatus() {
      if (dataTableInstRef.value) {
        return dataTableInstRef.value.filter(null);
      }
    }
    onMounted(() => {
      filterStatus();
    });

    return {
      delModal,
      activeItem,
      accountStatus,
      data,
      columns: createColumns(),
      dataTableInst: dataTableInstRef,
      checkedRowKeys: checkedRowKeysRef,
      pagination: {
        pageSize: 5,
      },
      filterStatus() {
        if (dataTableInstRef.value) {
          if (activeItem.value === "Tất cả trạng thái") {
            if (accountStatus.value == "All") {
              dataTableInstRef.value.filter(null);
            } else {
              dataTableInstRef.value.filter({
                accstatus: [accountStatus.value || ""],
              });
            }
          } else {
            if (accountStatus.value == "All") {
              dataTableInstRef.value.filter({
                status: [activeItem.value || ""],
              });
            } else {
              dataTableInstRef.value.filter({
                status: [activeItem.value || ""],
                accstatus: [accountStatus.value || ""],
              });
            }
          }
        }
      },
      rowKey: (row: RowData) => row.status,
      handleCheck(rowKeys: DataTableRowKey[]) {
        checkedRowKeysRef.value = rowKeys;
      },
      time: ref(null),
      value: ref(null),
      options: [
        {
          label: "Select",
          value: "song0",
          disabled: true,
        },
        {
          label: "A",
          value: "A",
        },
        {
          label: "B",
          value: "B",
        },
        {
          label: "C",
          value: "C",
        },
      ],
      statusoptions: [
        {
          label: "All",
          value: "All",
        },
        {
          label: "Hoạt động",
          value: "Hoạt động",
        },
        {
          label: "Không hoạt động",
          value: "Không hoạt động",
        },
      ],
    };
  },
});
interface RowData {
  stt: number;
  tthv: string;
  khts: string;
  mhht: string;
  mdxl: string;
  ncn: string;
  ncs: string;
  status: string;
  accstatus: string;
}
const actionMenu = [
  {
    title: "Xếp lớp",
    key: "Xep",
  },
  {
    title: "Dừng hoạt động",
    key: "Stop",
  },
];
function createColumns(): DataTableColumns<RowData> {
  return [
    {
      type: "selection",
    },
    {
      title: "STT",
      key: "stt",
      defaultSortOrder: "ascend",
      sorter: "default",
      titleAlign: "center",
    },
    {
      title: "Thông tin học viên",
      key: "tthv",
      defaultSortOrder: "ascend",
      sorter: "default",
    },
    {
      title: "Kế hoạch tuyển sinh",
      key: "khts",
    },
    {
      title: "Số môn học hoàn thành",
      key: "mhht",
      defaultSortOrder: "ascend",
      sorter: "default",
    },
    {
      title: "Môn đã xếp lớp",
      key: "mdxl",
      defaultSortOrder: "ascend",
      sorter: "default",
    },
    {
      title: "Ngày cập nhật",
      key: "ncn",
      defaultSortOrder: "ascend",
      sorter: "default",
    },
    {
      title: "Người chăm sóc",
      key: "ncs",
      defaultSortOrder: "ascend",
      sorter: "default",
    },
    {
      title: "Tình trạng",
      key: "status",
      render(row) {
        let color = "";
        let background = "";
        switch (row.status) {
          case "Đang học":
            color = "#00974F";
            background = "#F0FFF8";
            break;
          case "Chưa xếp lớp":
            color = "#FF7A00";
            background = "#FFF6E1";
            break;
          case "Bảo lưu":
            color = "#4D6FA8";
            background = "#ECF1F9";
            break;
          default:
            color = "gray";
        }
        return h(
          "span",
          {
            style: {
              padding: "5px 10px",
              borderRadius: "10px",
              color,
              background,
            },
          },
          row.status,
        );
      },
      defaultFilterOptionValues: ["Đang học", "Chưa xếp lớp", "Bảo lưu"],
      filter(value, row) {
        return row.status.includes(value as string);
      },
    },
    {
      title: "Trạng thái tài khoản",
      key: "accstatus",
      render(row) {
        let color = "";
        let background = "";
        switch (row.accstatus) {
          case "Hoạt động":
            color = "#00974F";
            background = "#F0FFF8";
            break;
          case "Không hoạt động":
            color = "#4D6FA8";
            background = "#ECF1F9";
            break;
          default:
            color = "gray";
        }
        return h(
          "span",
          {
            style: {
              padding: "5px 10px",
              borderRadius: "10px",
              color,
              background,
            },
          },
          row.accstatus,
        );
      },
      defaultFilterOptionValues: ["Hoạt động", "Không hoạt động"],
      filter(value, row) {
        return row.accstatus.includes(value as string);
      },
    },
    {
      title: "Action",
      key: "actions",
      titleAlign: "center",
      render(row) {
        return h("div", [
          h(
            NButton,
            {
              size: "small",
              quaternary: true,
              style: { backgroundColor: "transparent", color: "green" },
            },
            {
              default: () =>
                h("i", {
                  class: "fa-regular fa-pen-to-square",
                }),
            },
          ),
          h(
            NButton,
            {
              size: "small",
              quaternary: true,
              style: { backgroundColor: "transparent", color: "red" },
            },
            { default: () => h("i", { class: "fa-solid fa-trash" }) },
          ),
          h(
            NDropdown,
            {
              trigger: "click",
              options: actionMenu,
              quaternary: true,
              style: { color: "gray" },
              onSelect(key) {
                if (key === "Xếp lớp") {
                  editRow(row);
                } else if (key === "Dừng hoạt động") {
                  deleteRow(row);
                }
              },
            },
            { default: () => h("i", { class: "fa-solid fa-ellipsis-v" }) },
          ),
        ]);
      },
    },
  ];
}

const data = [
  {
    stt: 1,
    tthv: "Nguyen Van A",
    khts: "Thiết kế đồ họa",
    mhht: "4/5",
    mdxl: "4/5",
    ncn: "21/07/2023",
    ncs: "Lê Anh Nuôi",
    status: "Đang học",
    accstatus: "Hoạt động",
  },
  {
    stt: 2,
    tthv: "Nguyen Van B",
    khts: "Thiết kế đồ họa",
    mhht: "3/5",
    mdxl: "1/5",
    ncn: "21/07/2022",
    ncs: "Lê Anh Nuôi",
    status: "Chưa xếp lớp",
    accstatus: "Hoạt động",
  },
  {
    stt: 3,
    tthv: "Nguyen Van C",
    khts: "Thiết kế đồ họa",
    mhht: "4/5",
    mdxl: "1/5",
    ncn: "21/07/2022",
    ncs: "Lê Anh Nuôi",
    status: "Bảo lưu",
    accstatus: "Không hoạt động",
  },
];
</script>

<template>
  <div class="flex h-full w-full overflow-auto rounded-2xl bg-gray-50">
    <!-- Main Content -->
    <div class="flex-1">
      <!-- Content Area -->
      <div class="h-full text-black">
        <n-card class="h-full bg-gray-50">
          <div>
            <nav>
              <ul class="mt-5 flex flex-row gap-5 text-xl text-gray-400">
                <li class="cursor-pointer duration-75 hover:text-blue-500">
                  <NuxtLink
                    to="#"
                    @click="
                      (() => {
                        activeItem = 'Tất cả trạng thái';
                        filterStatus();
                      })()
                    "
                    :class="{
                      'text-blue-500 underline underline-offset-8 duration-150':
                        activeItem === 'Tất cả trạng thái',
                    }"
                  >
                    Tất cả trạng thái
                  </NuxtLink>
                </li>
                <li class="cursor-pointer duration-75 hover:text-blue-500">
                  <NuxtLink
                    to="#"
                    @click="
                      (() => {
                        activeItem = 'Đang học';
                        filterStatus();
                      })()
                    "
                    :class="{
                      'text-blue-500 underline underline-offset-8 duration-150':
                        activeItem === 'Đang học',
                    }"
                  >
                    Đang học
                  </NuxtLink>
                </li>
                <li class="cursor-pointer duration-75 hover:text-blue-500">
                  <NuxtLink
                    to="#"
                    @click="
                      (() => {
                        activeItem = 'Chưa xếp lớp';
                        filterStatus();
                      })()
                    "
                    :class="{
                      'text-blue-500 underline underline-offset-8 duration-150':
                        activeItem === 'Chưa xếp lớp',
                    }"
                  >
                    Chưa xếp lớp
                  </NuxtLink>
                </li>
                <li class="cursor-pointer duration-75 hover:text-blue-500">
                  <NuxtLink
                    to="#"
                    @click="
                      (() => {
                        activeItem = 'Bảo lưu';
                        filterStatus();
                      })()
                    "
                    :class="{
                      'text-blue-500 underline underline-offset-8 duration-150':
                        activeItem === 'Bảo lưu',
                    }"
                  >
                    Bảo lưu
                  </NuxtLink>
                </li>
              </ul>
            </nav>
          </div>

          <n-grid
            class="min-h-fit w-full"
            :x-gap="70"
            cols="1 m:3"
            responsive="screen"
          >
            <n-gi span="1">
              <n-form-item>
                <n-input
                  type="text"
                  placeholder="Tìm kiếm tên, số điện thoại, email học viên"
                />
              </n-form-item>
            </n-gi>
            <n-gi span="1">
              <n-form-item>
                <n-select
                  v-model="time"
                  :options="options"
                  placeholder="Tuần này"
                />
              </n-form-item>
            </n-gi>
            <n-gi span="1">
              <n-form-item>
                <n-select
                  v-model="accountStatus"
                  :options="statusoptions"
                  @change="
                    (() => {
                      accountStatus = $event;
                      filterStatus();
                    })()
                  "
                  placeholder="Trạng thái"
                />
              </n-form-item>
            </n-gi>
            <n-gi span="1 m:3">
              <n-grid
                class="-mt-8 w-full"
                :x-gap="30"
                cols="1 m:3"
                responsive="screen"
              >
                <n-gi span="1 m:2" class="flex-cols-3 flex gap-x-10">
                  <n-form-item>
                    <n-button type="info"> Dừng hoạt động </n-button>
                  </n-form-item>
                  <n-form-item>
                    <n-button type="info"> Hoạt động </n-button>
                  </n-form-item>
                  <n-form-item>
                    <n-button type="error" ghost>
                      <i class="fa-solid fa-trash mr-1"></i>
                      Xóa lựa chọn
                    </n-button>
                  </n-form-item>
                </n-gi>
                <n-gi
                  span="1"
                  class="flex-cols-3 flex gap-x-10 justify-self-end"
                >
                  <n-form-item>
                    <n-button type="info">
                      Import học viên
                      <i class="fa-solid fa-cloud-arrow-up ml-1"></i>
                    </n-button>
                  </n-form-item>
                  <n-form-item>
                    <n-button type="info">
                      Xuất file
                      <i class="fa-solid fa-file-export ml-1"></i>
                    </n-button>
                  </n-form-item>
                </n-gi>
              </n-grid>
            </n-gi>
            <n-gi span="1 m:3">
              <n-data-table
                ref="dataTableInst"
                :bordered="false"
                :single-line="false"
                :columns="columns"
                :data="data"
                :scroll-x="1800"
                :pagination="pagination"
                :row-key="rowKey"
                @update:checked-row-keys="handleCheck"
              />
            </n-gi>
          </n-grid>
        </n-card>
      </div>
    </div>
  </div>
</template>

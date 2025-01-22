<script setup>
const items = ref([
  {
    hocvien: "Nguyen Van A",
    phone: "0123456789",
    email: "A@gmail.com",
    trangthaihv: "Hoc vien tiem nang",
    ngaytao: "08/03/2024",
    trangthai: "hoạt động",
  },
  {
    hocvien: "Nguyen Van B",
    phone: "0123456789",
    email: "B@gmail.com",
    trangthaihv: "Hoc vien dang cham soc",
    ngaytao: "08/03/2024",
    trangthai: "khóa",
  },
  {
    hocvien: "Nguyen Van B",
    phone: "0123456789",
    email: "B@gmail.com",
    trangthaihv: "Hoc vien dang cham soc",
    ngaytao: "08/03/2024",
    trangthai: "khóa",
  },
  {
    hocvien: "Nguyen Van B",
    phone: "0123456789",
    email: "B@gmail.com",
    trangthaihv: "Hoc vien dang cham soc",
    ngaytao: "08/03/2024",
    trangthai: "khóa",
  },
  {
    hocvien: "Nguyen Van B",
    phone: "0123456789",
    email: "B@gmail.com",
    trangthaihv: "Hoc vien dang cham soc",
    ngaytao: "08/03/2024",
    trangthai: "khóa",
  },
  {
    hocvien: "Nguyen Van B",
    phone: "0123456789",
    email: "B@gmail.com",
    trangthaihv: "Hoc vien dang cham soc",
    ngaytao: "08/03/2024",
    trangthai: "khóa",
  },
]);
</script>
<template>
  <div
    class="mx-5 my-5 w-full h-full flex-grow rounded-2xl bg-white px-5 shadow-md 2xl:pl-10"
  >
    <div class="text-md my-10">
      <div class="mx-auto my-5 grid grid-cols-6 gap-10">
        <Icon
          class="absolute ml-4 mt-4"
          name="weui:search-filled"
          width="24"
          height="24"
        />
        <input
          type="text"
          class="col-span-2 rounded-xl border bg-gray-200 px-10 py-3 font-bold text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-600"
          placeholder="Tìm kiếm học viên"
        />
        <select
          class="rounded-xl bg-gray-200 px-5 font-bold text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-600"
        >
          <option selected disabled value="">Trạng thái học viên</option>
          <option value=""></option>
          <option value=""></option>
        </select>
        <select
          class="rounded-xl bg-gray-200 px-5 font-bold text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-600"
        >
          <option selected disabled value="">Trạng thái tài khoản</option>
          <option value=""></option>
          <option value=""></option>
        </select>
        <select
          class="rounded-xl bg-gray-200 px-5 font-bold text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-600"
        >
          <option selected disabled value="">Tất cả thời gian</option>
          <option value=""></option>
          <option value=""></option>
        </select>
        <button
          class="right-10 mx-auto gap-10 rounded-xl bg-blue-500 px-10 py-3 text-white"
        >
          Xuất file
          <Icon name="fe:export" width="24" height="24" />
        </button>
      </div>
      <div class="w-full">
        <table class="w-5/6 table-auto border-2 text-left">
          <thead class="border-2 px-5">
            <tr>
              <th class="px-5">
                <input type="checkbox" />
              </th>
              <th><label>STT</label></th>
              <th><label>Học viên</label></th>
              <th><label>Email</label></th>
              <th><label>Trạng thái học viên</label></th>
              <th><label>Ngày tạo</label></th>
              <th><label>Trạng thái</label></th>
              <th><label>Hoạt động</label></th>
            </tr>
          </thead>
          <tbody class="cursor-default">
            <tr
              v-for="(item, index) in items"
              :key="item.id"
              class="my-5 h-14 hover:bg-gray-200"
              :class="{ 'bg-gray-100': index % 2 == 1 }"
            >
              <td class="px-5"><input type="checkbox" /></td>
              <td>
                <P>{{ index + 1 }}</P>
              </td>
              <td class="flex-col-2 flex">
                <img
                  :src="epreview || defaultImag"
                  class="row-span-2 mt-1 h-12 w-12 rounded-full border border-gray-300 object-cover"
                />
                <div class="px-3">
                  <P>{{ item.hocvien }}</P>
                  <p>{{ item.phone }}</p>
                </div>
              </td>
              <td>
                <P>{{ item.email }}</P>
              </td>
              <td>
                <P>{{ item.trangthaihv }}</P>
              </td>
              <td>
                <P>{{ item.ngaytao }}</P>
              </td>
              <td>
                <span
                  :class="{
                    'rounded-2xl bg-green-400 px-2 py-1 text-green-700':
                      item.trangthai == 'hoạt động',
                    'rounded-2xl bg-gray-400 px-2 py-1 text-gray-700':
                      item.trangthai == 'khóa',
                  }"
                  >{{ item.trangthai }}</span
                >
              </td>
              <td class="items-center">
                <button class="bg-transparent">
                  <Icon
                    class="mx-5 text-red-500"
                    name="weui:delete-filled"
                    width="24"
                    height="24"
                  />
                </button>
                <button>
                  <Icon
                    class="text-green-600"
                    name="weui:pencil-filled"
                    width="24"
                    height="24"
                  />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div></div>
    </div>
  </div>
</template>

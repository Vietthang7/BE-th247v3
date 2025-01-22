<script setup lang="ts">
import { reactive } from "vue";

const dropdowns = reactive({
  daoTao: false,
  thietLap: false,
});
const changesize = (payload: MouseEvent) => {};

const toggleDropdown = (menu: "daoTao" | "thietLap") => {
  dropdowns[menu] = !dropdowns[menu];
};
</script>

<template>
  <aside class="w-64 bg-white shadow-md">
    <!-- User Info -->
    <div class="flex items-center space-x-3 border-b p-4">
      <img
        src="https://via.placeholder.com/40"
        alt="User Avatar"
        class="h-10 w-10 rounded-full"
      />
      <div>
        <h3 class="text-sm font-semibold">Đào Thị Hồng Thư</h3>
        <p class="text-xs text-gray-500">Giám đốc trung tâm</p>
      </div>
    </div>

    <!-- Menu -->
    <nav class="mt-4">
      <ul class="space-y-2">
        <!-- Menu Item -->
        <li>
          <NuxtLink
            to="/dashboard"
            class="flex items-center px-4 py-2.5 text-gray-600 hover:bg-blue-50 hover:text-blue-500"
            active-class="text-blue-500"
          >
            <i class="fas fa-home mr-3"></i>
            Không gian chung
          </NuxtLink>
        </li>

        <!-- Dropdown Menu -->
        <li>
          <div
            @click="toggleDropdown('daoTao')"
            class="flex cursor-pointer items-center justify-between px-4 py-2.5 text-gray-600 hover:bg-blue-50 hover:text-blue-500"
            active-class="text-blue-500"
          >
            <span class="flex items-center">
              <i class="fas fa-graduation-cap mr-3"></i>
              Đào tạo
            </span>
            <i
              :class="
                dropdowns.daoTao
                  ? 'fas fa-chevron-down'
                  : 'fas fa-chevron-right'
              "
            ></i>
          </div>
          <ul
            v-if="dropdowns.daoTao"
            class="ml-6 mt-2 space-y-1 text-sm text-gray-500"
          >
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i
                  class="fas fa-circle-dot fa-2xs mr-5"
                  v-on:click="changesize"
                >
                </i>
                Môn học
              </a>
            </li>
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Lớp học
              </a>
            </li>
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Học viên
              </a>
            </li>
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Giảng viên
              </a>
            </li>
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Lịch Học</a
              >
            </li>
          </ul>
        </li>

        <!-- Another Dropdown Menu -->
        <li>
          <div
            @click="toggleDropdown('thietLap')"
            class="flex cursor-pointer items-center justify-between px-4 py-2.5 text-gray-600 hover:bg-blue-50 hover:text-blue-500"
          >
            <span class="flex items-center">
              <i class="fas fa-cogs mr-3"></i>
              Thiết lập
            </span>
            <i
              :class="
                dropdowns.thietLap
                  ? 'fas fa-chevron-down'
                  : 'fas fa-chevron-right'
              "
            ></i>
          </div>
          <ul
            v-if="dropdowns.thietLap"
            class="ml-6 mt-2 space-y-1 text-sm text-gray-500"
          >
            <li>
              <NuxtLink
                to="/thietlap/themmoi"
                active-class="text-blue-500"
                class="block px-2 py-1 hover:text-blue-500"
              >
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Tài khoản học viên
              </NuxtLink>
            </li>
            <li>
              <NuxtLink
                to="/thietlap/tlhv"
                active-class="text-blue-500"
                class="block px-2 py-1 hover:text-blue-500"
              >
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Thông tin đơn vị
              </NuxtLink>
            </li>
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Tài khoản
              </a>
            </li>
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Phân quyền
              </a>
            </li>
          </ul>
        </li>

        <!-- Single Items -->
        <li>
          <a
            href="#"
            class="flex items-center px-4 py-2.5 text-gray-600 hover:bg-blue-50 hover:text-blue-500"
          >
            <i class="fas fa-bell mr-3"></i>
            Thông báo
          </a>
        </li>
        <li>
          <a
            href="#"
            class="flex items-center px-4 py-2.5 text-gray-600 hover:bg-blue-50 hover:text-blue-500"
          >
            <i class="fas fa-chart-bar mr-3"></i>
            Báo cáo đào tạo
          </a>
        </li>
      </ul>
    </nav>
  </aside>
</template>

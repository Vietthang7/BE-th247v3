<script setup lang="ts">
import { reactive, ref, onMounted, onUnmounted } from "vue";

const dropdowns = reactive<{ [key: string]: boolean }>({
  tuyensinh: false,
  daoTao: false,
  taichinh: false,
  nhansu: false,
  baocao: false,
  thietlap: false,
});

const activeItem = ref();

const setActive = (menu: string) => {
  activeItem.value = menu;
};
const toggleDropdown = (menu: string) => {
  dropdowns[menu] = !dropdowns[menu];
};

const isCollapsed = ref(false);

const handleResize = () => {
  isCollapsed.value = window.innerWidth < 1500; // Collapse sidebar on small screens
};

onMounted(() => {
  window.addEventListener("resize", handleResize);
  handleResize(); // Trigger on load
});

onUnmounted(() => {
  window.removeEventListener("resize", handleResize);
});
</script>

<template>
  <aside class="h-full w-1/6 overflow-auto bg-white shadow-md">
    <!-- User Info -->
    <div class="flex items-center space-x-3 border-b p-4">
      <img src="" alt="User Avatar" class="h-10 w-10 rounded-full" />
      <div v-if="!isCollapsed">
        <h3 class="text-sm font-semibold text-black">Đào Thị Hồng Thư</h3>
        <p class="text-xs text-gray-500">Giám đốc trung tâm</p>
      </div>
    </div>

    <!-- Menu -->
    <nav class="mt-4">
      <ul class="space-y-2">
        <!-- Menu Item -->
        <li>
          <NuxtLink
            to="/"
            class="flex items-center px-4 py-2.5 text-gray-600 hover:bg-blue-50 hover:text-blue-500"
            active-class="text-blue-500 bg-blue-50 "
          >
            <i class="fas fa-home mr-3"></i>
            <div v-if="!isCollapsed">Không gian chung</div>
          </NuxtLink>
        </li>

        <!-- Dropdown Menu -->
        <li>
          <div
            @click="toggleDropdown('tuyensinh')"
            :class="[
              dropdowns.tuyensinh
                ? 'bg-blue-50 font-semibold text-blue-500'
                : 'text-gray-600',
              'flex cursor-pointer items-center justify-between px-4 py-2.5 hover:bg-blue-50 hover:text-blue-500',
            ]"
          >
            <span class="flex items-center">
              <i class="fas fa-tower-broadcast mr-3"></i>
              <div v-if="!isCollapsed">Tuyển Sinh</div>
            </span>
            <i
              :class="
                dropdowns.tuyensinh
                  ? 'fas fa-chevron-down'
                  : 'fas fa-chevron-right'
              "
            ></i>
          </div>
          <ul
            v-if="dropdowns.tuyensinh"
            class="ml-6 mt-2 space-y-1 text-sm text-gray-500"
          >
            <li>
              <NuxtLink
                to="/blog"
                @click="setActive('lopHoc')"
                :class="[
                  activeItem === 'lopHoc' ? 'font-semibold text-blue-500' : '',
                  'block px-2 py-1 hover:text-blue-500',
                ]"
              >
                <!-- 
                <NuxtLink
                class="cursor-pointer duration-500 hover:text-yellow-400"
                to="/lophoc"
                active-class="text-blue-500"
                >Radar</NuxtLink
                 > -->
                <i class="fas fa-circle-dot fa-2xs mr-5" v-if="!isCollapsed">
                </i>

                tuyen sinh page
              </NuxtLink>
            </li>
          </ul>
        </li>

        <li>
          <div
            @click="toggleDropdown('daoTao')"
            :class="[
              dropdowns.daoTao
                ? 'bg-blue-50 font-semibold text-blue-500'
                : 'text-gray-600',
              'flex cursor-pointer items-center justify-between px-4 py-2.5 hover:bg-blue-50 hover:text-blue-500',
            ]"
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
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>
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

        <li>
          <div
            @click="toggleDropdown('taichinh')"
            :class="[
              dropdowns.taichinh
                ? 'bg-blue-50 font-semibold text-blue-500'
                : 'text-gray-600',
              'flex cursor-pointer items-center justify-between px-4 py-2.5 hover:bg-blue-50 hover:text-blue-500',
            ]"
          >
            <span class="flex items-center">
              <i class="fas fa-arrow-trend-up mr-3"></i>
              Tài Chính
            </span>
            <i
              :class="
                dropdowns.taichinh
                  ? 'fas fa-chevron-down'
                  : 'fas fa-chevron-right'
              "
            ></i>
          </div>
          <ul
            v-if="dropdowns.taichinh"
            class="ml-6 mt-2 space-y-1 text-sm text-gray-500"
          >
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>
                something
              </a>
            </li>
          </ul>
        </li>

        <li>
          <div
            @click="toggleDropdown('nhansu')"
            :class="[
              dropdowns.nhansu
                ? 'bg-blue-50 font-semibold text-blue-500'
                : 'text-gray-600',
              'flex cursor-pointer items-center justify-between px-4 py-2.5 hover:bg-blue-50 hover:text-blue-500',
            ]"
          >
            <span class="flex items-center">
              <i class="fas fa-user mr-3"></i>
              Nhân Sự
            </span>
            <i
              :class="
                dropdowns.nhansu
                  ? 'fas fa-chevron-down'
                  : 'fas fa-chevron-right'
              "
            ></i>
          </div>
          <ul
            v-if="dropdowns.nhansu"
            class="ml-6 mt-2 space-y-1 text-sm text-gray-500"
          >
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>
                Danh sách nhân sự
              </a>
            </li>
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>
                Bảng chấm công
              </a>
            </li>
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>
                Cài đặt nhiệm vụ
              </a>
            </li>
          </ul>
        </li>

        <!-- Another Dropdown Menu -->
        <li>
          <div
            @click="toggleDropdown('baocao')"
            :class="[
              dropdowns.baocao
                ? 'bg-blue-50 font-semibold text-blue-500'
                : 'text-gray-600',
              'flex cursor-pointer items-center justify-between px-4 py-2.5 hover:bg-blue-50 hover:text-blue-500',
            ]"
          >
            <span class="flex items-center">
              <i class="fas fa-chart-simple mr-3"></i>
              Báo cáo
            </span>
            <i
              :class="
                dropdowns.baocao
                  ? 'fas fa-chevron-down'
                  : 'fas fa-chevron-right'
              "
            ></i>
          </div>
          <ul
            v-if="dropdowns.baocao"
            class="ml-6 mt-2 space-y-1 text-sm text-gray-500"
          >
            <li>
              <a href="#" class="block px-2 py-1 hover:text-blue-500">
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>
                Something
              </a>
            </li>
          </ul>
        </li>
        <li>
          <div
            @click="toggleDropdown('thietlap')"
            :class="[
              dropdowns.thietlap
                ? 'bg-blue-50 font-semibold text-blue-500'
                : 'text-gray-600',
              'flex cursor-pointer items-center justify-between px-4 py-2.5 hover:bg-blue-50 hover:text-blue-500',
            ]"
          >
            <span class="flex items-center">
              <i class="fas fa-cogs mr-3"></i>
              Thiết lập
            </span>
            <i
              :class="
                dropdowns.thietlap
                  ? 'fas fa-chevron-down'
                  : 'fas fa-chevron-right'
              "
            ></i>
          </div>
          <ul
            v-if="dropdowns.thietlap"
            class="ml-6 mt-2 space-y-1 text-sm text-gray-500"
          >
            <li>
              <NuxtLink
                to="/thietlap/hocvien"
                active-class="text-blue-500"
                class="block px-2 py-1 hover:text-blue-500"
              >
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Tài khoản học viên
              </NuxtLink>
            </li>
            <li>
              <NuxtLink
                to="/thietlap/addhocvien"
                active-class="text-blue-500"
                class="block px-2 py-1 hover:text-blue-500"
              >
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Thêm học viên
              </NuxtLink>
            </li>
            <li></li>
            <li>
              <NuxtLink
                to="/thietlap/trungtam"
                active-class="text-blue-500"
                class="block px-2 py-1 hover:text-blue-500"
              >
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Thông tin đơn vị
              </NuxtLink>
            </li>
            <li>
              <NuxtLink
                to="/thietlap/calam"
                active-class="text-blue-500"
                class="block px-2 py-1 hover:text-blue-500"
              >
                <i class="fas fa-circle-dot fa-2xs mr-5"> </i>

                Cài đặt ca làm
              </NuxtLink>
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
        <li>
          <a
            href="#"
            class="flex items-center px-4 py-2.5 text-gray-600 hover:bg-blue-50 hover:text-blue-500"
          >
            <i class="fas fa-circle-question mr-3"></i>
            Trung tâm trợ giúp
          </a>
        </li>
      </ul>
    </nav>
  </aside>
</template>

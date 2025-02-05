<template>
  <aside class="h-screen w-64 bg-white shadow-md">
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
        <!-- Single Menu Item -->
        <li>
          <a
            href="#"
            :class="[
              'flex items-center px-4 py-2.5 hover:bg-blue-50',
              activeItem === 'khongGianChung'
                ? 'bg-blue-100 text-blue-500'
                : 'text-gray-600',
            ]"
            @click="setActive('khongGianChung')"
          >
            <i class="fas fa-home mr-3"></i>
            Không gian chung
          </a>
        </li>

        <!-- Dropdown Menu -->
        <li>
          <div
            @click="toggleDropdown('daoTao')"
            :class="[
              'flex cursor-pointer items-center justify-between px-4 py-2.5 hover:bg-blue-50',
              activeDropdown === 'daoTao'
                ? 'bg-blue-100 text-blue-500'
                : 'text-gray-600',
            ]"
          >
            <span class="flex items-center">
              <i class="fas fa-graduation-cap mr-3"></i>
              Đào tạo
            </span>
            <i
              :class="
                dropdowns.daoTao ? 'fas fa-chevron-up' : 'fas fa-chevron-down'
              "
            ></i>
          </div>
          <ul
            v-if="dropdowns.daoTao"
            class="ml-6 mt-2 space-y-1 text-sm text-gray-500"
          >
            <li>
              <a
                href="#"
                :class="[
                  'block px-2 py-1 hover:text-blue-500',
                  activeItem === 'monHoc' ? 'font-semibold text-blue-500' : '',
                ]"
                @click="setActive('monHoc', 'daoTao')"
              >
                Môn học
              </a>
            </li>
            <li>
              <a
                href="#"
                :class="[
                  'block px-2 py-1 hover:text-blue-500',
                  activeItem === 'lopHoc' ? 'font-semibold text-blue-500' : '',
                ]"
                @click="setActive('lopHoc', 'daoTao')"
              >
                Lớp học
              </a>
            </li>
          </ul>
        </li>

        <!-- Single Menu Item -->
        <li>
          <a
            href="#"
            :class="[
              'flex items-center px-4 py-2.5 hover:bg-blue-50',
              activeItem === 'thongBao'
                ? 'bg-blue-100 text-blue-500'
                : 'text-gray-600',
            ]"
            @click="setActive('thongBao')"
          >
            <i class="fas fa-bell mr-3"></i>
            Thông báo
          </a>
        </li>
      </ul>
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";

// State for dropdowns
const dropdowns = reactive<{ [key: string]: boolean }>({
  daoTao: false,
});

// Track the active dropdown and active item
const activeDropdown = ref(""); // Tracks the open dropdown
const activeItem = ref(""); // Tracks the active menu item

// Toggle dropdowns
const toggleDropdown = (menu: string) => {
  if (dropdowns[menu]) {
    // Close the dropdown if already open
    dropdowns[menu] = false;
    activeDropdown.value = ""; // Unset the active dropdown
  } else {
    // Close all dropdowns
    Object.keys(dropdowns).forEach((key) => {
      dropdowns[key] = false;
    });
    // Open the clicked dropdown
    dropdowns[menu] = true;
    activeDropdown.value = menu; // Set the active dropdown
  }
};

// Set active item and optionally the active dropdown
const setActive = (menu: string, dropdown?: string) => {
  activeItem.value = menu;

  // If the item is part of a dropdown, ensure the dropdown stays open
  if (dropdown) {
    dropdowns[dropdown] = true;
    activeDropdown.value = dropdown;
  } else {
    // If a standalone menu item is clicked, close all dropdowns
    Object.keys(dropdowns).forEach((key) => {
      dropdowns[key] = false;
    });
    activeDropdown.value = ""; // Reset active dropdown
  }
};
</script>

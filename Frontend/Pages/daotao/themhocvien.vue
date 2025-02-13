<script setup lang="ts">
import { reactive, ref, onMounted, onUnmounted } from "vue";

const dropdowns = reactive<{ [key: string]: boolean }>({
  canhan: false,
  nhucau: false,
  danhsach: false,
  phuhuynh: false,
  chamsoc: false,
});
const isCollapsed = ref(false);
const showForm = ref(false);
const activeDropdown = ref(""); // Tracks the open dropdown

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
</script>
<template>
  <div class="flex h-full w-full">
    <div class="my-5 flex h-full w-1/6">
      <nav>
        <ul class="flex flex-col gap-y-3">
          <li
            @click="toggleDropdown('canhan')"
            :class="[
              'cursor-pointer py-2 pl-2 pr-16',
              activeDropdown === 'canhan'
                ? 'bg-[#E6F7FF] text-[#133D85]'
                : 'text-[#4D6FA8]',
            ]"
          >
            Thông tin cá nhân
          </li>
          <li
            @click="toggleDropdown('nhucau')"
            :class="[
              'cursor-pointer py-2 pl-2 pr-16',
              activeDropdown === 'nhucau'
                ? 'bg-[#E6F7FF] text-[#133D85]'
                : 'text-[#4D6FA8]',
            ]"
          >
            Nhu cầu học tập
          </li>
          <li
            @click="toggleDropdown('danhsach')"
            :class="[
              'cursor-pointer py-2 pl-2 pr-16',
              activeDropdown === 'danhsach'
                ? 'bg-[#E6F7FF] text-[#133D85]'
                : 'text-[#4D6FA8]',
            ]"
          >
            Danh sách lớp học
          </li>
          <li
            @click="toggleDropdown('phuhuynh')"
            :class="[
              'cursor-pointer py-2 pl-2 pr-16',
              activeDropdown === 'phuhuynh'
                ? 'bg-[#E6F7FF] text-[#133D85]'
                : 'text-[#4D6FA8]',
            ]"
          >
            Thông tin phụ huynh
          </li>
          <li
            @click="toggleDropdown('chamsoc')"
            :class="[
              'cursor-pointer py-2 pl-2 pr-16',
              activeDropdown === 'chamsoc'
                ? 'bg-[#E6F7FF] text-[#133D85]'
                : 'text-[#4D6FA8]',
            ]"
          >
            Thông tin chăm sóc
          </li>
        </ul>
      </nav>
    </div>
    <div class="flex h-full w-5/6 overflow-auto rounded-2xl bg-gray-50">
      <!-- Main Content -->
      <div class="flex-1">
        <nav>
          <ul class="mx-5 my-5 flex flex-col gap-y-5 text-xl">
            <li class="rounded-2xl text-[#133D85]">
              <div
                @click="toggleDropdown('canhan')"
                :class="[
                  'flex cursor-pointer items-center justify-between px-4 py-2.5',

                  activeDropdown === 'canhan'
                    ? 'rounded-2xl bg-[#E6F7FF] text-[#133D85]'
                    : 'text-gray-600',
                ]"
              >
                <span class="flex items-center">
                  <div v-if="!isCollapsed">Thông tin cá nhân</div>
                </span>
                <i
                  :class="
                    dropdowns.canhan
                      ? 'fas fa-chevron-up'
                      : 'fas fa-chevron-down'
                  "
                ></i>
              </div>
              <ul v-if="dropdowns.canhan" class="w-ful h-full">
                <li>
                  <div class="w-ful h-full" v-if="!isCollapsed">
                    <DaotaoPerInfo />
                  </div>
                </li>
              </ul>
            </li>
            <li class="rounded-2xl text-[#133D85]">
              <div
                @click="toggleDropdown('nhucau')"
                :class="[
                  'flex cursor-pointer items-center justify-between px-4 py-2.5',

                  activeDropdown === 'nhucau'
                    ? 'rounded-2xl bg-[#E6F7FF] text-[#133D85]'
                    : 'text-gray-600',
                ]"
              >
                <span class="flex items-center">
                  <div v-if="!isCollapsed">Nhu cầu học tập</div>
                </span>
                <i
                  :class="
                    dropdowns.nhucau
                      ? 'fas fa-chevron-up'
                      : 'fas fa-chevron-down'
                  "
                ></i>
              </div>
              <ul v-if="dropdowns.nhucau" class="w-ful h-full">
                <li>
                  <div class="w-ful h-full" v-if="!isCollapsed">
                    <div class="px5 mt-5">
                      <n-button
                        round
                        type="info"
                        class="h-12 w-48 rounded-2xl text-xl"
                        @click="showForm = !showForm"
                      >
                        Thêm mới
                        <i class="fa-solid fa-plus ml-3"></i>
                      </n-button>
                      <div v-if="showForm" class="mt-5">
                        <DaotaoNeed />
                      </div>
                    </div>
                  </div>
                </li>
              </ul>
            </li>
            <li class="rounded-2xl text-[#133D85]">
              <div
                @click="toggleDropdown('danhsach')"
                :class="[
                  'flex cursor-pointer items-center justify-between px-4 py-2.5',

                  activeDropdown === 'danhsach'
                    ? 'rounded-2xl bg-[#E6F7FF] text-[#133D85]'
                    : 'text-gray-600',
                ]"
              >
                <span class="flex items-center">
                  <div v-if="!isCollapsed">Danh sách lớp học</div>
                </span>
                <i
                  :class="
                    dropdowns.danhsach
                      ? 'fas fa-chevron-up'
                      : 'fas fa-chevron-down'
                  "
                ></i>
              </div>
              <ul v-if="dropdowns.danhsach" class="w-ful h-full">
                <li>
                  <div class="w-ful h-full" v-if="!isCollapsed">
                    <DaotaoClass />
                  </div>
                </li>
              </ul>
            </li>
            <li class="rounded-2xl text-[#133D85]">
              <div
                @click="toggleDropdown('phuhuynh')"
                :class="[
                  'flex cursor-pointer items-center justify-between px-4 py-2.5',

                  activeDropdown === 'phuhuynh'
                    ? 'rounded-2xl bg-[#E6F7FF] text-[#133D85]'
                    : 'text-gray-600',
                ]"
              >
                <span class="flex items-center">
                  <div v-if="!isCollapsed">Thông tin phụ huynh</div>
                </span>
                <i
                  :class="
                    dropdowns.phuhuynh
                      ? 'fas fa-chevron-up'
                      : 'fas fa-chevron-down'
                  "
                ></i>
              </div>
            </li>
            <li class="rounded-2xl text-[#133D85]">
              <div
                @click="toggleDropdown('chamsoc')"
                :class="[
                  'flex cursor-pointer items-center justify-between px-4 py-2.5',

                  activeDropdown === 'chamsoc'
                    ? 'rounded-2xl bg-[#E6F7FF] text-[#133D85]'
                    : 'text-gray-600',
                ]"
              >
                <span class="flex items-center">
                  <div v-if="!isCollapsed">Thông tin chăm sóc</div>
                </span>
                <i
                  :class="
                    dropdowns.chamsoc
                      ? 'fas fa-chevron-up'
                      : 'fas fa-chevron-down'
                  "
                ></i>
              </div>
            </li>
          </ul>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";

const activeTab = ref("draw"); // Options: "draw", "upload"
const canvas = ref<HTMLCanvasElement | null>(null);

const switchToDraw = () => {
  activeTab.value = "draw";
  clearCanvas();
};

const handleUpload = (event: Event) => {
  activeTab.value = "upload";
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      if (canvas.value) {
        const ctx = canvas.value.getContext("2d");
        if (ctx) {
          const img = new Image();
          img.onload = () => {
            ctx.clearRect(0, 0, canvas.value!.width, canvas.value!.height);
            ctx.drawImage(img, 0, 0, canvas.value!.width, canvas.value!.height);
          };
          img.src = e.target?.result as string;
        }
      }
    };
    reader.readAsDataURL(file);
  }
};

const clearSignature = () => {
  if (activeTab.value === "draw") {
    clearCanvas();
  } else {
    activeTab.value = "draw";
  }
};

const clearCanvas = () => {
  if (canvas.value) {
    const ctx = canvas.value.getContext("2d");
    if (ctx) {
      ctx.clearRect(0, 0, canvas.value.width, canvas.value.height);
    }
  }
};

// Drawing functionality
onMounted(() => {
  if (canvas.value) {
    const ctx = canvas.value.getContext("2d");
    if (ctx) {
      let drawing = false;

      const startDrawing = (event: MouseEvent) => {
        drawing = true;
        ctx.beginPath();
        ctx.moveTo(event.offsetX, event.offsetY);
      };

      const draw = (event: MouseEvent) => {
        if (!drawing) return;
        ctx.lineTo(event.offsetX, event.offsetY);
        ctx.stroke();
      };

      const stopDrawing = () => {
        drawing = false;
        ctx.closePath();
      };

      canvas.value.addEventListener("mousedown", startDrawing);
      canvas.value.addEventListener("mousemove", draw);
      canvas.value.addEventListener("mouseup", stopDrawing);
      canvas.value.addEventListener("mouseleave", stopDrawing);
    }
  }
});

// Active tab
</script>

<template>
  <div class="h-min-fit flex w-full overflow-auto rounded-2xl bg-gray-50">
    <!-- Main Content -->
    <div class="flex-1">
      <!-- Content Area -->
      <div class="text-black">
        <main class="flex items-center justify-center p-8">
          <div class="min-h-fit bg-gray-50 p-8">
            <!-- Tabs -->
            <div class="mb-6 flex border-b">
              <button
                @click="activeTab = 'accountInfo'"
                :class="[
                  'px-6 py-2 font-medium',
                  activeTab === 'accountInfo'
                    ? 'border-b-2 border-blue-500 text-blue-500'
                    : 'text-gray-600',
                ]"
              >
                Thông tin tài khoản
              </button>
              <button
                @click="activeTab = 'changePassword'"
                :class="[
                  'px-6 py-2 font-medium',
                  activeTab === 'changePassword'
                    ? 'border-b-2 border-blue-500 text-blue-500'
                    : 'text-gray-600',
                ]"
              >
                Đổi mật khẩu
              </button>
            </div>
            <!-- Account Information -->
            <div v-if="activeTab === 'accountInfo'">
              <div class="grid grid-cols-3 gap-6">
                <!-- Profile Picture -->
                <div>
                  <div class="mb-4">
                    <img
                      src="https://via.placeholder.com/120"
                      alt="Profile Picture"
                      class="mx-auto h-32 w-32 rounded-full"
                    />
                    <div class="mt-2 text-center">
                      <label
                        for="profile-picture-upload"
                        class="cursor-pointer text-sm text-blue-500 hover:underline"
                      >
                        Chỉnh sửa ảnh
                      </label>
                      <input
                        type="file"
                        id="profile-picture-upload"
                        class="hidden"
                      />
                    </div>
                  </div>
                  <p class="text-center text-xs text-gray-500">
                    Cho phép ảnh jpeg, jpg, png. Size ảnh tối đa 3.1 MB
                  </p>
                </div>

                <!-- siganture -->
                <div class="mt-6">
                  <label class="block text-sm font-medium text-gray-700"
                    >Chữ ký cá nhân</label
                  >
                  <div class="mt-2 flex items-start space-x-4">
                    <!-- Canvas or Placeholder -->
                    <div class="relative">
                      <canvas
                        v-show="activeTab.value === 'draw'"
                        ref="canvas"
                        class="h-32 w-64 rounded border border-gray-300 bg-gray-50"
                      ></canvas>
                      <div
                        v-show="activeTab.value === 'upload'"
                        class="flex h-32 w-64 items-center justify-center rounded border border-gray-300 bg-gray-50"
                      >
                        <p class="text-sm text-gray-500">No image uploaded</p>
                      </div>
                    </div>

                    <!-- Actions -->
                    <div class="flex flex-col space-y-2">
                      <div class="flex space-x-2">
                        <button
                          @click="switchToDraw"
                          :class="[
                            'rounded border px-4 py-2 text-sm',
                            activeTab.value === 'draw'
                              ? 'bg-blue-500 text-white'
                              : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-100',
                          ]"
                        >
                          Vẽ
                        </button>
                        <label
                          for="upload-input"
                          :class="[
                            'cursor-pointer rounded border px-4 py-2 text-sm',
                            activeTab.value === 'upload'
                              ? 'bg-blue-500 text-white'
                              : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-100',
                          ]"
                        >
                          Tải ảnh lên
                          <input
                            type="file"
                            id="upload-input"
                            @change="handleUpload"
                            class="hidden"
                            accept="image/*"
                          />
                        </label>
                      </div>
                      <button
                        @click="clearSignature"
                        class="rounded bg-red-500 px-4 py-2 text-sm text-white hover:bg-red-600"
                      >
                        Xóa
                      </button>
                    </div>
                  </div>
                </div>

                <!-- Personal Information -->
                <div class="col-span-2">
                  <div class="grid grid-cols-2 gap-4">
                    <div>
                      <label
                        for="name"
                        class="block text-sm font-medium text-gray-700"
                      >
                        Họ và tên *
                      </label>
                      <input
                        type="text"
                        id="name"
                        class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                        placeholder="Nguyễn Tú Nam"
                      />
                    </div>
                    <div>
                      <label
                        for="phone"
                        class="block text-sm font-medium text-gray-700"
                      >
                        Số điện thoại *
                      </label>
                      <input
                        type="text"
                        id="phone"
                        class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                        placeholder="0123456789"
                      />
                    </div>
                    <div>
                      <label
                        for="email"
                        class="block text-sm font-medium text-gray-700"
                      >
                        Email *
                      </label>
                      <input
                        type="email"
                        id="email"
                        class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                        placeholder="tuname@gmail.com"
                      />
                    </div>
                    <div>
                      <label
                        for="birthday"
                        class="block text-sm font-medium text-gray-700"
                      >
                        Ngày sinh
                      </label>
                      <input
                        type="date"
                        id="birthday"
                        class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                      />
                    </div>
                    <div>
                      <label
                        for="gender"
                        class="block text-sm font-medium text-gray-700"
                      >
                        Giới tính
                      </label>
                      <select
                        id="gender"
                        class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                      >
                        <option value="male">Nam</option>
                        <option value="female">Nữ</option>
                      </select>
                    </div>
                    <div>
                      <label
                        for="city"
                        class="block text-sm font-medium text-gray-700"
                      >
                        Địa chỉ
                      </label>
                      <select
                        id="city"
                        class="mt-1 block w-full rounded border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                      >
                        <option value="hanoi">Hà Nội</option>
                        <option value="hochiminh">TP. Hồ Chí Minh</option>
                      </select>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Signature -->

              <!-- Save Button -->
              <div class="mt-8 text-center">
                <button
                  type="button"
                  class="rounded bg-blue-500 px-6 py-2 text-white shadow hover:bg-blue-600"
                >
                  Lưu
                </button>
              </div>
            </div>

            <!-- Change Password -->
            <div v-if="activeTab === 'changePassword'">
              <p class="text-center text-gray-600">
                Mục đổi mật khẩu sẽ được xây dựng sau.
              </p>
            </div>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

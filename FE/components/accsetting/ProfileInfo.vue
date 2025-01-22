<script>
export default {
  data() {
    return {
      formData: {
        name: "",
        phone: "",
        email: "",
        birthDate: "",
        gender: "Giới tính",
        city: "Tỉnh/Thành phố",
        district: "Quận/Huyện",
        address: "",
        signature: "",
        preview: null, // To display selected image
        defaultImage: "", // Provide a default image URL
      },
    };
  },
};
</script>
<template>
  <div>
    <!-- Form -->
    <form @submit.prevent="onSubmit" class="placeholder:text-gray-400">
      <!-- Profile Picture -->
      <div class="mt-3 grid grid-cols-1 gap-6 md:grid-cols-2">
        <div>
          <div class="relative">
            <label class="mb-1 block text-lg font-bold">Ảnh đại diện</label>
            <img
              :src="preview || defaultImage"
              class="h-24 w-24 rounded-2xl border border-gray-300 object-cover"
            />
            <label
              for="profilePic"
              class="absolute bottom-0 left-0 ml-20 cursor-pointer rounded-full bg-blue-500 p-1 text-white"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="2"
                stroke="currentColor"
                class="h-4 w-4"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M15.232 5.232l3.536 3.536M9 13.5l6.364-6.364a2 2 0 112.828 2.828L11.828 16.5H9v-3z"
                />
              </svg>
            </label>
          </div>

          <!-- File input -->
          <input
            id="profilePic"
            type="file"
            accept="image/jpeg, image/png"
            class="hidden"
          />
          <p class="text-gray-400">
            Cho phép ảnh jpeg, jpg, png.<br />
            Size ảnh tối đa 3.1 MB
          </p>
        </div>
        <!-- Signature -->
        <div>
          <label class="mb-2 block text-lg font-bold">Chữ ký cá nhân</label>
          <textarea
            v-model="formData.signature"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            rows="3"
            placeholder="Ký tại đây hoặc tải ảnh lên"
          ></textarea>
        </div>
      </div>
      <!-- Input Fields -->
      <div class="mt-3 grid grid-cols-1 gap-5 md:grid-cols-2">
        <div>
          <label class="mb-1 block text-lg font-bold">Họ và tên *</label>
          <input
            type="text"
            v-model="formData.name"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            placeholder="Nguyễn Tú Nam"
            required
          />
        </div>
        <div>
          <label class="mb-1 block text-lg font-bold">Số điện thoại *</label>
          <input
            type="text"
            v-model="formData.phone"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            placeholder="0123456789"
            required
          />
        </div>
        <div class="md:col-span-2">
          <label class="mb-1 text-lg font-bold">Email *</label>
          <input
            type="email"
            v-model="formData.email"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            placeholder="tunam@gmail.com"
            required
          />
        </div>
        <div>
          <label class="mb-1 block text-lg font-bold">Giới tính</label>
          <select
            v-model="formData.gender"
            class="w-full rounded-2xl border border-gray-300 px-3 py-3 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
          >
            <option selected disabled>Giới tính</option>
            <option value="male">Nam</option>
            <option value="female">Nữ</option>
            <option value="other">Khác</option>
          </select>
        </div>
        <div>
          <label class="mb-1 block text-lg font-bold">Ngày sinh</label>
          <input
            type="date"
            v-model="formData.birthDate"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
          />
        </div>
        <div class="md:col-span-2">
          <label class="mb-1 block text-lg font-bold">Địa chỉ</label>
          <div class="mb-5 grid grid-cols-1 gap-5 md:grid-cols-2">
            <select
              v-model="formData.city"
              class="w-full rounded-2xl border border-gray-300 px-3 py-3 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
              placeholder="Tỉnh/Thành phố"
            >
              <option selected disabled>Tỉnh/Thành phố</option>
              <option value="male">Hà Nội</option>
              <option value="female">Hải Phòng</option>
              <option value="other">TP Hồ Chí Minh</option>
            </select>
            <select
              v-model="formData.district"
              class="w-full rounded-2xl border border-gray-300 px-3 py-3 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
              placeholder="Quận/Huyện"
            >
              <option selected disabled>Quận/Huyện</option>
              <option value="male">Hà Nội</option>
              <option value="female">Hải Phòng</option>
              <option value="other">TP Hồ Chí Minh</option>
            </select>
          </div>
          <input
            type="text"
            v-model="formData.address"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            placeholder="Địa chỉ cụ thể"
          />
        </div>
      </div>

      <!-- Submit Button -->
      <div class="mt-1 flex items-center">
        <button
          type="submit"
          class="mx-auto my-5 w-1/3 rounded-2xl bg-[#00A2EB] px-6 py-3 text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2"
        >
          Lưu
        </button>
      </div>
    </form>
  </div>
</template>
<style scoped>
label {
  color: #133d85;
}
</style>

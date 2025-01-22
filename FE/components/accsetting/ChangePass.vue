<script>
export default {
  data() {
    return {
      formData: {
        email: "",
        oldpass: "",
        newpass: "",
        compass: "",
      },
    };
  },
};
</script>
<template>
  <div>
    <form @submit.prevent="onSubmit" class="placeholder:text-gray-400">
      <div class="mt-5 grid grid-cols-1 gap-6 md:grid-cols-2">
        <div>
          <label class="mb-1 block text-lg font-bold"
            >Tài khoản đăng nhập</label
          >
          <input
            type="email"
            v-model="formData.email"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            placeholder="tunam@gmail.com"
            required
          />
        </div>
        <div>
          <label class="mb-1 block text-lg font-bold"
            >Mật khẩu hiện tại *</label
          >
          <input
            type="password"
            v-model="formData.oldpass"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            placeholder="Nhập mật khẩu hiện tại"
            required
          />
        </div>
        <div>
          <label class="mb-1 block text-lg font-bold">Mật khẩu mới *</label>
          <input
            type="password"
            v-model="formData.newpass"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            placeholder="Nhập mật khẩu mới"
            required
          />
        </div>
        <div>
          <label class="mb-1 block text-lg font-bold"
            >Xác nhận mật khẩu *</label
          >
          <input
            type="password"
            v-model="formData.compass"
            class="w-full rounded-2xl border border-gray-300 px-3 py-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600"
            placeholder="Xác nhận mật khẩu mới"
            required
          />
        </div>
      </div>
      <div class="mt-5 flex items-center">
        <button
          type="submit"
          class="mx-auto w-1/4 rounded-2xl bg-[#00A2EB] px-6 py-3 text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2"
        >
          Lưu
        </button>
      </div>
    </form>
  </div>
</template>

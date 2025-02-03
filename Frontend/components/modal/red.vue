<template>
  <div v-if="visible" class="popup-overlay">
    <div class="popup-container">
      <div class="popup-header">
        <div class="popup-icon">
          <!-- <img src="delete-icon.png" alt="Delete Icon" /> -->
        </div>
        <h3>Xác nhận xóa tài khoản?</h3>
      </div>
      <div class="popup-actions">
        <button class="cancel-button" @click="cancel">Hủy</button>
        <button class="confirm-button" @click="confirm">Xóa</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["confirm", "cancel"],
  methods: {
    cancel() {
      this.$emit("cancel");
    },
    confirm() {
      this.$emit("confirm");
    },
  },
};
</script>

<style scoped>
.popup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.popup-container {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 300px;
  text-align: center;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.popup-header {
  margin-bottom: 20px;
}

.popup-icon img {
  width: 40px;
  margin-bottom: 10px;
}

.popup-actions {
  display: flex;
  justify-content: space-between;
}

.cancel-button {
  background: #f1f1f1;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}

.cancel-button:hover {
  background: #e2e2e2;
}

.confirm-button {
  background: #ff5252;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}

.confirm-button:hover {
  background: #ff4040;
}
</style>

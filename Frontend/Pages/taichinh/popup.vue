<script>
definePageMeta({
  layout: "default",
});

export default {
  components: { ModalRed },
  data() {
    return {
      showModal: false,
    };
  },
};
</script>
<template>
  <div>
    <ModalRed v-show="showModal" />
    <div class="save-btn">
      <button @click="showModal = true">Save</button>
    </div>
  </div>
</template>

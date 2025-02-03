<script lang="ts">
import { defineComponent, ref } from "vue";

export default defineComponent({
  setup() {
    return {
      timestamp: ref(""),
      Mtime: ref(1183135260000),
      Mtime2: ref(1183135260000),
      Ntime: ref(1183135260000),
      Ntime2: ref(1183135260000),
      Etime: ref(1183135260000),
      Etime2: ref(1183135260000),
    };
  },
  methods: {
    createArray(i: number, j: number) {
      return Array.from({ length: j - i + 1 }, (_, k) => i + k);
    },
  },
});
</script>
<template>
  <div class="mt-5 px-5">
    <div class="h-full w-full rounded-2xl bg-gray-200 px-1 py-1">
      <p class="px-5 py-2">Nhu cầu 1</p>
      <div class="w-full rounded-b-2xl bg-white px-5 py-5">
        <n-grid :x-gap="30" cols="1 m:2" responsive="screen">
          <n-gi span="1 m:2">
            <n-form-item
              label="Chương trình học, môn học *"
              label-placement="left"
            >
              <n-select placeholder="Chọn môn học" />
            </n-form-item>
          </n-gi>
          <n-gi span="1 m:2">
            <n-form-item label="Mục tiêu học tập" path="goal">
              <n-input placeholder="Mục tiêu học tập" />
            </n-form-item>
          </n-gi>
          <n-gi span="1 m:2">
            <n-form-item label="Yêu cầu về giáo viên" path="demand">
              <n-input placeholder="Yêu cầu về giáo viên" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="Hình thức học mong muốn" label-placement="left">
              <n-checkbox-group>
                <n-space item-style="display: flex;">
                  <n-checkbox value="Học online" label="Học online" />
                  <n-checkbox value="Học offline" label="Học offline" />
                </n-space>
              </n-checkbox-group>
            </n-form-item>
          </n-gi>
          <n-gi></n-gi>
          <n-gi>
            <n-form-item
              label="Ngày mong muốn bắt đầu học:"
              label-placement="left"
            >
              <n-date-picker
                v-model="timestamp"
                type="date"
                placeholder="Chọn ngày"
              />
            </n-form-item>
          </n-gi>

          <n-gi span="1 m:2">
            <n-grid :x-gap="5" class="mb-5" cols="1 m:11" responsive="screen">
              <n-gi span="1 m:4">
                <label>Lịch trống:</label>
              </n-gi>
              <n-gi span="1 m:7">
                <n-checkbox-group v-model="Bcheck">
                  <n-grid cols="1 m:7" responsive="screen">
                    <n-gi
                      v-for="day in [
                        'Thứ 2',
                        'Thứ 3',
                        'Thứ 4',
                        'Thứ 5',
                        'Thứ 6',
                        'Thứ 7',
                        'Chủ nhật',
                      ]"
                      :key="day"
                    >
                      <n-checkbox
                        :value="day"
                        :label="day"
                        @update:checked="toggleCheckbox(day, $event)"
                      />
                    </n-gi>
                  </n-grid>
                </n-checkbox-group>
              </n-gi>
            </n-grid>
            <n-grid :x-gap="5" class="mb-5" cols="1 m:11" responsive="screen">
              <n-gi span="1 m:4">
                <n-form-item label="Sáng: " label-placement="left">
                  <n-time-picker
                    v-model="Mtime"
                    :hours="createArray(5, 11)"
                    :seconds="[0]"
                  />
                  <i class="fa-solid fa-arrow-right"></i>
                  <n-time-picker
                    v-model="Mtime2"
                    :hours="[12]"
                    :minutes="0"
                    :seconds="[0]"
                  />
                </n-form-item>
              </n-gi>
              <n-gi span="1 m:7">
                <n-checkbox-group v-model="morning">
                  <n-grid cols="1 m:7" responsive="screen">
                    <n-gi
                      v-for="day in [
                        'Thứ 2',
                        'Thứ 3',
                        'Thứ 4',
                        'Thứ 5',
                        'Thứ 6',
                        'Thứ 7',
                        'Chủ nhật',
                      ]"
                      :key="day"
                    >
                      <n-checkbox :value="day" />
                    </n-gi>
                  </n-grid>
                </n-checkbox-group>
              </n-gi>
            </n-grid>
            <n-grid :x-gap="5" class="mb-5" cols="1 m:11" responsive="screen">
              <n-gi span="1 m:4">
                <n-form-item label="Chiều: " label-placement="left">
                  <n-time-picker
                    v-model="Ntime"
                    :hours="createArray(13, 16)"
                    :seconds="[0]"
                  />
                  <i class="fa-solid fa-arrow-right"></i>
                  <n-time-picker
                    v-model="Ntime2"
                    :hours="[17]"
                    :minutes="createArray(0, 30)"
                    :seconds="[0]"
                  />
                </n-form-item>
              </n-gi>
              <n-gi span="1 m:7">
                <n-checkbox-group v-model="noon">
                  <n-grid cols="1 m:7" responsive="screen">
                    <n-gi
                      v-for="day in [
                        'Thứ 2',
                        'Thứ 3',
                        'Thứ 4',
                        'Thứ 5',
                        'Thứ 6',
                        'Thứ 7',
                        'Chủ nhật',
                      ]"
                      :key="day"
                    >
                      <n-checkbox :value="day" />
                    </n-gi>
                  </n-grid>
                </n-checkbox-group>
              </n-gi>
            </n-grid>
            <n-grid :x-gap="5" class="mb-5" cols="1 m:11" responsive="screen">
              <n-gi span="1 m:4">
                <n-form-item label="Tối: " label-placement="left">
                  <n-time-picker
                    v-model="Etime"
                    :hours="createArray(19, 21)"
                    :seconds="[0]"
                  />
                  <i class="fa-solid fa-arrow-right"></i>
                  <n-time-picker
                    v-model="Etime2"
                    :hours="[22]"
                    :minutes="createArray(0, 30)"
                    :seconds="[0]"
                  />
                </n-form-item>
              </n-gi>
              <n-gi span="1 m:7">
                <n-checkbox-group v-model="evening">
                  <n-grid cols="1 m:7" responsive="screen">
                    <n-gi
                      v-for="day in [
                        'Thứ 2',
                        'Thứ 3',
                        'Thứ 4',
                        'Thứ 5',
                        'Thứ 6',
                        'Thứ 7',
                        'Chủ nhật',
                      ]"
                      :key="day"
                    >
                      <n-checkbox :value="day" />
                    </n-gi>
                  </n-grid>
                </n-checkbox-group>
              </n-gi>
            </n-grid>
          </n-gi>

          <n-gi span="1 m:2">
            <n-form-item label="Ghi chú" path="note">
              <n-input placeholder="Ghi chú" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-button round type="info" class="h-12 w-52 rounded-2xl text-lg">
              Lưu
            </n-button>
          </n-gi>
        </n-grid>
      </div>
    </div>
  </div>
</template>

<template>
  <div class="flex h-full w-full overflow-auto">
    <SidebarBar></SidebarBar>
    <n-config-provider
      :theme-overrides="themeOverrides"
      class="flex h-full w-full overflow-x-hidden"
    >
      <main class="font-customFont w-full bg-gray-100 bg-cover">
        <div class="flex h-full w-full overflow-x-hidden">
          <slot />
        </div>
      </main>
    </n-config-provider>
  </div>
</template>

<script setup>
import { NConfigProvider } from "naive-ui";

const themeOverrides = {
  Form: {
    labelFontWeight: "500",
    labelTextColor: "#133D85",
    labelFontSizeTopLarge: "16px",
    labelFontSizeTopMedium: "16px",
    labelFontSizeTopSmall: "16px",
    asteriskColor: "#133D85FF",
  },
  Input: {
    heightLarge: "44px",
    heightMedium: "44px",
    paddingLarge: "0 20px",
    paddingMedium: "0 20px",
    borderRadius: "15px",
    textColor: "#8BA3CBFF",
    textColorDisabled: "#8BA3CBFF",
    colorDisabled: "rgba(255, 255, 255, 1)",
    border: "1px solid #E7EDFFFF",
    borderDisabled: "1px solid #E7EDFFFF",
    placeholderColor: "#8ba3cb",
    placeholderColorDisabled: "#8ba3cb",
  },
  DataTable: {
    borderRadius: "0px",
    tdTextColor: "#133D85FF",
    fontSizeLarge: "15px",
    fontSizeMedium: "15px",
    thTextColor: "#133D85FF",
    thIconColor: "#133D85FF",
    tdColorStriped: "#FBFBFBFF",
    thColor: "#FFFFFFFF",
  },
  Dropdown: {
    optionTextColor: "#8BA3CBFF",
    optionTextColorHover: "#133D85FF",
    optionTextColorActive: "#133D85FF",
    optionTextColorChildActive: "#133D85FF",
    suffixColor: "#133D85FF",
    optionColorHover: "#E6F7FFFF",
    borderRadius: "4px",
  },
};
</script>

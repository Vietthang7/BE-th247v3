<script setup>
const getThemeOverrides = computed(() => {
  // const lightenStr = lighten(designStore.appTheme, 6)
  return {
    common: {
      fontFamily: "Be Vietnam Pro",
      heightLarge: "44px",
      heightMedium: "44px",
      paddingLarge: "0 20px",
      paddingMedium: "0 20px",
      borderRadius: "15px",
    },
    Menu: {
      itemTextColor: "#8BA3CBFF",
      itemTextColorHover: "#8BA3CBFF",
      fontSize: "15px",
      itemHeight: "54px",
      itemTextColorActive: "#133D85FF",
      itemTextColorChildActive: "#00A2EBFF",
      itemIconColor: "#8BA3CBFF",
      itemIconColorHover: "#8BA3CBFF",
      itemIconColorChildActive: "#00A2EBFF",
      itemIconColorActive: "#00A2EBFF",
      itemIconColorActiveHover: "#00A2EBFF",
      itemIconColorChildActiveHover: "#00A2EBFF",
      itemTextColorActiveHover: "#133D85FF",
      itemTextColorChildActiveHover: "#00A2EBFF",
      itemColorActive: "#fff",
      itemIconColorCollapsed: "#8BA3CBFF",
      itemTextColorInverted: "#8BA3CBFF",
      arrowColor: "#8BA3CBFF",
      arrowColorHover: "#8BA3CBFF",
      arrowColorActive: "#00A2EBFF",
      arrowColorChildActive: "#00A2EBFF",
    },
    Form: {
      labelFontWeight: "500",
      labelTextColor: "#133D85FF",
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
      colorDisabled: "rgba(255, 255, 255, 1)",
      border: "1px solid #E7EDFFFF",
      borderDisabled: "1px solid #E7EDFFFF",
      placeholderColor: "#8ba3cb",
      placeholderColorDisabled: "#8ba3cb",
    },
    Checkbox: {
      colorChecked: "#00A2EBFF",
      borderColorChecked: "#00A2EBFF",
      checkMarkColor: "#FFFFFF",
    },
    Tag: {
      textColor: "#00A2EBFF",
      color: "#E6F7FFFF",
      colorBordered: "#E6F7FFFF",
      closeIconColor: "#00A2EBFF",
      closeIconColorHover: "#00A2EBFF",
      border: "1px solid #E6F7FFFF",
      colorCheckable: "#F5F5F5FF",
      textColorCheckable: "#8BA3CBFF",
      textColorHoverCheckable: "#8BA3CBFF",
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
    Tabs: {
      tabTextColorBar: "#8BA3CBFF",
    },
    Dropdown: {
      optionTextColor: "#00A2EBFF",
      optionTextColorHover: "#00A2EBFF",
      optionTextColorActive: "#00A2EBFF",
      optionTextColorChildActive: "#00A2EBFF",
      suffixColor: "#00A2EBFF",
      optionColorHover: "#E6F7FFFF",
      borderRadius: "15px",
    },
    Upload: {
      itemTextColorSuccess: "#00A2EBFF",
    },
  };
});
</script>
<template>
  <n-config-provider :theme-overrides="getThemeOverrides">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </n-config-provider>
</template>

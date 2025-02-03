<template>
  <div class="h-full w-full bg-gray-200">
    <header class="hero w-100px xl-bg-bottom flex bg-cover bg-top lg:bg-center">
      <TopNav />
    </header>
    <div class="mt-3 flex w-[calc(100%-20px)] overflow-auto">
      <Sidebar />
      <!-- <Test /> -->
      <main class="mb-5 ml-5 mt-5 w-full rounded-2xl">
        <slot />
        <!-- <content /> -->
      </main>
    </div>
  </div>
  <!-- <h2 class="pt-24 text-9xl font-bold tracking-wide text-yellow-100">
      Hello
    </h2> -->

  <!-- <MFooter />
  </div> -->
</template>

<template>
  <div class="min-h-screen bg-background dark">
    <header class="border-b border-border bg-card px-4 py-3 sm:px-6 sm:py-4">
      <div class="max-w-7xl mx-auto">
        <h1 class="text-xl sm:text-2xl font-bold text-foreground">Products Management</h1>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 py-6 sm:px-6 sm:py-8">
      <div class="w-full">
        <div class="flex flex-col sm:grid sm:grid-cols-2 mb-6 sm:mb-8 bg-muted p-1 rounded-lg gap-1 sm:gap-0">
          <button
            @click="activeTab = 'discounted'"
            :class="[
              'px-3 py-2 sm:px-4 sm:py-2 rounded-md text-sm font-medium transition-all duration-200 relative',
              activeTab === 'discounted'
                ? 'bg-background text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground hover:bg-background/50'
            ]"
          >
            Products & Discounts
            <span class="absolute top-1 right-1 text-xs opacity-50">1</span>
          </button>
          <button
            @click="activeTab = 'manage'"
            :class="[
              'px-3 py-2 sm:px-4 sm:py-2 rounded-md text-sm font-medium transition-all duration-200 relative',
              activeTab === 'manage'
                ? 'bg-background text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground hover:bg-background/50'
            ]"
          >
            Manage Products
            <span class="absolute top-1 right-1 text-xs opacity-50">2</span>
          </button>
        </div>

        <div v-if="activeTab === 'discounted'">
          <DiscountedProducts ref="discountedRef" />
        </div>

        <div v-if="activeTab === 'manage'">
          <ProductManager ref="managerRef" />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
  import { ref, onMounted, onUnmounted } from 'vue'
  import DiscountedProducts from './components/DiscountedProducts.vue'
  import ProductManager from './components/ProductManager.vue'

  const activeTab = ref('discounted')
</script>
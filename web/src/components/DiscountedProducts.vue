<template>
  <div class="space-y-4 sm:space-y-6">
    <!-- Toast Notifications -->
    <ToastNotifications :toasts="toasts" @remove="removeToast"/>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-0">
      <div class="flex items-center gap-4">
        <h2 class="text-lg sm:text-xl font-semibold text-foreground">Products & Discounts</h2>
        <div class="flex gap-2">
          <div class="text-sm text-muted-foreground bg-muted border border-border px-3 py-1 rounded-full">
            {{ filteredBonusProducts.length }} bonus
          </div>
          <div class="text-sm text-muted-foreground bg-muted border border-border px-3 py-1 rounded-full">
            {{ filteredRegularProducts.length }} regular
          </div>
        </div>
      </div>
      <div class="flex gap-2">
        <button @click="refreshProducts" :disabled="loading"
                class="inline-flex items-center justify-center px-4 py-2 border border-border bg-card text-card-foreground hover:bg-accent hover:text-accent-foreground hover:border-accent-foreground/20 rounded-md text-sm font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed min-h-[44px] sm:min-h-auto shadow-sm hover:shadow-md">
          <span class="mr-2" :class="{ 'animate-spin': loading }">↻</span>
          Refresh
        </button>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="relative">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <svg class="h-5 w-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
        </svg>
      </div>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search products by name, friendly name, or ID..."
        class="w-full pl-10 pr-4 py-3 border border-border bg-card text-card-foreground rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-all duration-200 shadow-sm hover:shadow-md hover:border-primary/50"
      />
      <div v-if="searchQuery" class="absolute inset-y-0 right-0 pr-3 flex items-center">
        <button @click="searchQuery = ''"
                class="text-muted-foreground hover:text-foreground transition-colors duration-200">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
    </div>

    <!-- Search Results Summary -->
    <div v-if="searchQuery && !loading && products.length > 0"
         class="flex items-center justify-between text-sm text-muted-foreground bg-muted/30 px-4 py-2 rounded-lg">
      <span>
        {{ filteredBonusProducts.length + filteredRegularProducts.length }} of {{
          products.length
        }} product{{ (filteredBonusProducts.length + filteredRegularProducts.length) !== 1 ? 's' : '' }} found
        <span v-if="filteredBonusProducts.length > 0" class="ml-2">
          ({{ filteredBonusProducts.length }} bonus)
        </span>
      </span>
      <button @click="searchQuery = ''" class="text-primary hover:text-primary/80 transition-colors duration-200">
        Clear search
      </button>
    </div>

    <div v-if="loading" class="flex flex-col items-center justify-center py-12">
      <div class="animate-spin rounded-full h-10 w-10 border-2 border-primary border-t-transparent mb-4"></div>
      <p class="text-muted-foreground text-sm font-medium">Loading products...</p>
    </div>

    <div v-else-if="error" class="text-center py-8">
      <p class="text-destructive mb-4 font-medium">{{ error }}</p>
      <button @click="fetchProducts"
              class="inline-flex items-center justify-center px-4 py-2 border border-border bg-card text-card-foreground hover:bg-accent hover:text-accent-foreground hover:border-accent-foreground/20 rounded-md text-sm font-medium transition-all duration-200 min-h-[44px] sm:min-h-auto shadow-sm hover:shadow-md">
        Try Again
      </button>
    </div>

    <div v-else-if="products.length === 0" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
        <span class="text-2xl">🛒</span>
      </div>
      <p class="text-muted-foreground font-medium mb-2">No products found</p>
      <p class="text-muted-foreground text-sm">Products will appear here when available</p>
    </div>

    <div v-else-if="searchQuery && filteredBonusProducts.length === 0 && filteredRegularProducts.length === 0"
         class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
        <span class="text-2xl">🔍</span>
      </div>
      <p class="text-muted-foreground font-medium mb-2">No products found</p>
      <p class="text-muted-foreground text-sm">Try adjusting your search terms</p>
    </div>

    <div v-else class="space-y-6">
      <!-- Bonus Products Section -->
      <BonusProductsSection
        v-if="filteredBonusProducts.length > 0"
        :products="filteredBonusProducts"
      />

      <!-- Regular Products Accordion -->
      <RegularProductsSection
        v-if="filteredRegularProducts.length > 0"
        :products="filteredRegularProducts"
        :show="showRegularProducts"
        @toggle="showRegularProducts = !showRegularProducts"
      />

      <!-- Show message if only one type exists or if search filtered everything -->
      <div v-if="filteredBonusProducts.length === 0 && filteredRegularProducts.length > 0" class="text-center py-6">
        <p class="text-muted-foreground text-sm">
          {{ searchQuery ? 'No bonus products match your search' : 'No bonus products available at the moment' }}
        </p>
      </div>

      <div v-if="filteredBonusProducts.length > 0 && filteredRegularProducts.length === 0 && !searchQuery"
           class="text-center py-6">
        <p class="text-muted-foreground text-sm">All available products are currently in the bonus program!</p>
      </div>
    </div>
  </div>
</template>

<script setup>
  import {
    ref,
    computed,
  } from 'vue';
  import ToastNotifications from './ToastNotifications.vue';
  import BonusProductsSection from './BonusProductsSection.vue';
  import RegularProductsSection from './RegularProductsSection.vue';

  const products = ref([]);
  const loading = ref(false);
  const error = ref('');
  const showRegularProducts = ref(false);
  const searchQuery = ref('');

  // Toast notifications
  const toasts = ref([]);
  let toastId = 0;

  // Toast notification system
  const showToast = (message, type = 'info') => {
    const toast = {
      id: ++toastId,
      message,
      type,
      show: false,
    };
    toasts.value.push(toast);

    // Trigger animation
    setTimeout(() => {
      toast.show = true;
    }, 100);

    // Auto remove after 4 seconds
    setTimeout(() => {
      removeToast(toast.id);
    }, 4000);
  };

  const removeToast = (id) => {
    const index = toasts.value.findIndex(t => t.id === id);
    if (index > -1) {
      toasts.value[index].show = false;
      setTimeout(() => {
        toasts.value.splice(index, 1);
      }, 300);
    }
  };

  // Computed properties to filter products based on search
  const filteredProducts = computed(() => {
    if (!searchQuery.value) {
      return products.value;
    }

    const query = searchQuery.value.toLowerCase();
    return products.value.filter(product =>
      product.api_name.toLowerCase().includes(query) ||
      product.friendly_name.toLowerCase().includes(query) ||
      product.appie_id.toString().includes(query) ||
      product.referer_url.toLowerCase().includes(query),
    );
  });

  // Computed properties to separate bonus and regular products from filtered results
  const filteredBonusProducts = computed(() => {
    return filteredProducts.value.filter(product => product.discount.in_bonus);
  });

  const filteredRegularProducts = computed(() => {
    return filteredProducts.value.filter(product => !product.discount.in_bonus);
  });

  const fetchProducts = async () => {
    loading.value = true;
    error.value = '';

    try {
      const response = await fetch('/api/products');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${ response.status }`);
      }
      const data = await response.json();
      let productsArray = Array.isArray(data) ? data : [data];

      // Sort products so bonus items appear first (though they'll be separated anyway)
      productsArray.sort((a, b) => {
        if (a.discount.in_bonus && !b.discount.in_bonus) return -1;
        if (!a.discount.in_bonus && b.discount.in_bonus) return 1;
        return 0;
      });

      products.value = productsArray;

      const bonusCount = productsArray.filter(p => p.discount.in_bonus).length;
      showToast(`Loaded ${ productsArray.length } products (${ bonusCount } bonus)`, 'success');

    } catch (err) {
      error.value = `Failed to fetch products: ${ err.message }`;
      showToast('Failed to load products', 'error');
      console.error('Error fetching products:', err);
    } finally {
      loading.value = false;
    }
  };

  const refreshProducts = async () => {
    loading.value = true;
    error.value = '';

    try {
      const response = await fetch('/api/products/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${ response.status }`);
      }
      const data = await response.json();
      let productsArray = Array.isArray(data) ? data : [data];

      // Sort products so bonus items appear first (though they'll be separated anyway)
      productsArray.sort((a, b) => {
        if (a.discount.in_bonus && !b.discount.in_bonus) return -1;
        if (!a.discount.in_bonus && b.discount.in_bonus) return 1;
        return 0;
      });

      products.value = productsArray;

      const bonusCount = productsArray.filter(p => p.discount.in_bonus).length;
      showToast(`Loaded ${ productsArray.length } products (${ bonusCount } bonus)`, 'success');

    } catch (err) {
      error.value = `Failed to fetch products: ${ err.message }`;
      showToast('Failed to load products', 'error');
      console.error('Error fetching products:', err);
    } finally {
      loading.value = false;
    }
  };

  // Expose fetchProducts for parent component
  defineExpose({
    fetchProducts,
  });

  // Auto-fetch on mount
  fetchProducts();
</script>
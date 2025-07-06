<template>
  <div class="space-y-4 sm:space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-0">
      <h2 class="text-lg sm:text-xl font-semibold text-foreground">Products & Discounts</h2>
      <button @click="fetchProducts" :disabled="loading" class="inline-flex items-center justify-center px-4 py-2 border border-border bg-card text-card-foreground hover:bg-accent hover:text-accent-foreground hover:border-accent-foreground/20 rounded-md text-sm font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed min-h-[44px] sm:min-h-auto shadow-sm hover:shadow-md">
        <span class="mr-2" :class="{ 'animate-spin': loading }">↻</span>
        Refresh
      </button>
    </div>

    <!-- Search Bar -->
    <div class="relative">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <svg class="h-5 w-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
        </svg>
      </div>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search products by name, friendly name, or ID..."
        class="w-full pl-10 pr-4 py-3 border border-border bg-card text-card-foreground rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-all duration-200 shadow-sm hover:shadow-md hover:border-primary/50"
      />
      <div v-if="searchQuery" class="absolute inset-y-0 right-0 pr-3 flex items-center">
        <button @click="searchQuery = ''" class="text-muted-foreground hover:text-foreground transition-colors duration-200">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
    </div>

    <!-- Search Results Summary -->
    <div v-if="searchQuery && !loading && products.length > 0" class="flex items-center justify-between text-sm text-muted-foreground bg-muted/30 px-4 py-2 rounded-lg">
      <span>
        {{ filteredBonusProducts.length + filteredRegularProducts.length }} of {{ products.length }} product{{ (filteredBonusProducts.length + filteredRegularProducts.length) !== 1 ? 's' : '' }} found
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
      <button @click="fetchProducts" class="inline-flex items-center justify-center px-4 py-2 border border-border bg-card text-card-foreground hover:bg-accent hover:text-accent-foreground hover:border-accent-foreground/20 rounded-md text-sm font-medium transition-all duration-200 min-h-[44px] sm:min-h-auto shadow-sm hover:shadow-md">Try Again</button>
    </div>

    <div v-else-if="products.length === 0" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
        <span class="text-2xl">🛒</span>
      </div>
      <p class="text-muted-foreground font-medium mb-2">No products found</p>
      <p class="text-muted-foreground text-sm">Products will appear here when available</p>
    </div>

    <div v-else-if="searchQuery && filteredBonusProducts.length === 0 && filteredRegularProducts.length === 0" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
        <span class="text-2xl">🔍</span>
      </div>
      <p class="text-muted-foreground font-medium mb-2">No products found</p>
      <p class="text-muted-foreground text-sm">Try adjusting your search terms</p>
    </div>

    <div v-else class="space-y-6">
      <!-- Bonus Products Section -->
      <div v-if="filteredBonusProducts.length > 0">
        <div class="flex items-center gap-3 mb-4">
          <div class="flex items-center gap-2">
            <span class="text-2xl">⭐</span>
            <h3 class="text-lg font-semibold text-foreground">Bonus Products</h3>
          </div>
          <div class="px-3 py-1 bg-gradient-to-r from-yellow-400/20 to-orange-400/20 border border-yellow-400/30 rounded-full">
            <span class="text-sm font-medium text-foreground">{{ filteredBonusProducts.length }} item{{ filteredBonusProducts.length !== 1 ? 's' : '' }}</span>
          </div>
        </div>

        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="product in filteredBonusProducts" :key="product.appie_id" class="border border-yellow-400/30 bg-gradient-to-br from-yellow-50/10 to-orange-50/10 rounded-lg p-4 sm:p-6 shadow-sm hover:shadow-lg hover:border-yellow-400/50 transition-all duration-200 transform hover:scale-[1.02]">
            <div class="space-y-3 sm:space-y-4">
              <div class="flex items-start sm:items-center justify-between gap-2">
                <h3 class="font-semibold text-card-foreground text-sm sm:text-base leading-tight">{{ product.api_name }}</h3>
                <span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-gradient-to-r from-yellow-400 to-orange-400 text-white font-medium whitespace-nowrap shadow-sm">
                  ⭐ Bonus
                </span>
              </div>

              <p class="text-sm text-muted-foreground">
                {{ product.friendly_name }}
              </p>

              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 pt-2 border-t border-yellow-400/20">
                <span class="text-xs text-muted-foreground font-medium">
                  In Bonus Program
                </span>
                <span class="text-xs text-muted-foreground bg-muted px-2 py-1 rounded-full">
                  ID: {{ product.appie_id }}
                </span>
              </div>

              <div class="pt-2">
                <a :href="product.referer_url" target="_blank" class="inline-flex items-center text-xs text-primary hover:text-primary/80 hover:underline transition-colors duration-150 min-h-[44px] sm:min-h-auto">
                  View Product →
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Regular Products Accordion -->
      <div v-if="filteredRegularProducts.length > 0">
        <div class="border border-border rounded-lg bg-card shadow-sm">
          <button
            @click="showRegularProducts = !showRegularProducts"
            class="w-full flex items-center justify-between p-4 text-left hover:bg-muted/30 transition-colors duration-200 rounded-t-lg"
          >
            <div class="flex items-center gap-3">
              <span class="text-xl">📦</span>
              <h3 class="text-lg font-semibold text-foreground">Regular Products</h3>
              <div class="px-3 py-1 bg-muted rounded-full">
                <span class="text-sm font-medium text-muted-foreground">{{ filteredRegularProducts.length }} item{{ filteredRegularProducts.length !== 1 ? 's' : '' }}</span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-sm text-muted-foreground">
                {{ showRegularProducts ? 'Hide' : 'Show' }}
              </span>
              <div :class="['transition-transform duration-200', showRegularProducts ? 'rotate-180' : '']">
                <svg class="w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                </svg>
              </div>
            </div>
          </button>

          <div v-if="showRegularProducts" class="border-t border-border">
            <div class="p-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              <div v-for="product in filteredRegularProducts" :key="product.appie_id" class="border border-border rounded-lg p-4 bg-card shadow-sm hover:shadow-md hover:border-accent-foreground/20 transition-all duration-200 transform hover:scale-[1.02]">
                <div class="space-y-3 sm:space-y-4">
                  <div class="flex items-start sm:items-center justify-between gap-2">
                    <h3 class="font-semibold text-card-foreground text-sm sm:text-base leading-tight">{{ product.api_name }}</h3>
                    <span class="text-xs text-muted-foreground bg-muted px-2 py-1 rounded-full whitespace-nowrap">
                      Regular
                    </span>
                  </div>

                  <p class="text-sm text-muted-foreground">
                    {{ product.friendly_name }}
                  </p>

                  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 pt-2 border-t border-border">
                    <span class="text-xs text-muted-foreground">
                      Regular Product
                    </span>
                    <span class="text-xs text-muted-foreground bg-muted px-2 py-1 rounded-full">
                      ID: {{ product.appie_id }}
                    </span>
                  </div>

                  <div class="pt-2">
                    <a :href="product.referer_url" target="_blank" class="inline-flex items-center text-xs text-primary hover:text-primary/80 hover:underline transition-colors duration-150 min-h-[44px] sm:min-h-auto">
                      View Product →
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Show message if only one type exists or if search filtered everything -->
      <div v-if="filteredBonusProducts.length === 0 && filteredRegularProducts.length > 0" class="text-center py-6">
        <p class="text-muted-foreground text-sm">
          {{ searchQuery ? 'No bonus products match your search' : 'No bonus products available at the moment' }}
        </p>
      </div>

      <div v-if="filteredBonusProducts.length > 0 && filteredRegularProducts.length === 0 && !searchQuery" class="text-center py-6">
        <p class="text-muted-foreground text-sm">All available products are currently in the bonus program!</p>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted, computed } from 'vue'

  const products = ref([])
  const loading = ref(false)
  const error = ref('')
  const showRegularProducts = ref(false)
  const searchQuery = ref('')

  // Computed properties to filter products based on search
  const filteredProducts = computed(() => {
    if (!searchQuery.value) {
      return products.value
    }

    const query = searchQuery.value.toLowerCase()
    return products.value.filter(product =>
      product.api_name.toLowerCase().includes(query) ||
      product.friendly_name.toLowerCase().includes(query) ||
      product.appie_id.toString().includes(query) ||
      product.referer_url.toLowerCase().includes(query)
    )
  })

  // Computed properties to separate bonus and regular products from filtered results
  const filteredBonusProducts = computed(() => {
    return filteredProducts.value.filter(product => product.discount.in_bonus)
  })

  const filteredRegularProducts = computed(() => {
    return filteredProducts.value.filter(product => !product.discount.in_bonus)
  })

  const fetchProducts = async () => {
    loading.value = true
    error.value = ''

    try {
      const response = await fetch('/api/products')
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      let productsArray = Array.isArray(data) ? data : [data]

      // Sort products so bonus items appear first (though they'll be separated anyway)
      productsArray.sort((a, b) => {
        if (a.discount.in_bonus && !b.discount.in_bonus) return -1
        if (!a.discount.in_bonus && b.discount.in_bonus) return 1
        return 0
      })

      products.value = productsArray
    } catch (err) {
      error.value = `Failed to fetch products: ${err.message}`
      console.error('Error fetching products:', err)
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchProducts()
  })
</script>
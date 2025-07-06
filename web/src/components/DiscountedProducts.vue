<template>
  <div class="space-y-4 sm:space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-0">
      <h2 class="text-lg sm:text-xl font-semibold text-foreground">Products & Discounts</h2>
      <button @click="fetchProducts" :disabled="loading" class="inline-flex items-center justify-center px-4 py-2 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-md text-sm font-medium transition-colors disabled:opacity-50 min-h-[44px] sm:min-h-auto">
        <span class="mr-2" :class="{ 'animate-spin': loading }">↻</span>
        Refresh
      </button>
    </div>

    <div v-if="loading" class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
    </div>

    <div v-else-if="error" class="text-center py-8">
      <p class="text-destructive mb-4">{{ error }}</p>
      <button @click="fetchProducts" class="inline-flex items-center justify-center px-4 py-2 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-md text-sm font-medium transition-colors min-h-[44px] sm:min-h-auto">Try Again</button>
    </div>

    <div v-else-if="products.length === 0" class="text-center py-8">
      <p class="text-muted-foreground">No products found.</p>
    </div>

    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="product in products" :key="product.appie_id" class="border border-border rounded-lg p-4 sm:p-6 bg-card">
        <div class="space-y-3 sm:space-y-4">
          <div class="flex items-start sm:items-center justify-between gap-2">
            <h3 class="font-semibold text-card-foreground text-sm sm:text-base leading-tight">{{ product.api_name }}</h3>
            <span v-if="product.discount.in_bonus" class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-secondary text-secondary-foreground whitespace-nowrap">
              ⭐ Bonus
            </span>
          </div>

          <p class="text-sm text-muted-foreground">
            {{ product.friendly_name }}
          </p>

          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 pt-2">
            <span class="text-xs text-muted-foreground">
              {{ product.discount.in_bonus ? 'In Bonus Program' : 'Regular Product' }}
            </span>
            <span class="text-xs text-muted-foreground">
              ID: {{ product.appie_id }}
            </span>
          </div>

          <div class="pt-2">
            <a :href="product.referer_url" target="_blank" class="inline-flex items-center text-xs text-primary hover:underline min-h-[44px] sm:min-h-auto">
              View Product →
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'

  const products = ref([])
  const loading = ref(false)
  const error = ref('')

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

      // Sort products so bonus items appear first
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
<template>
  <div class="space-y-4 sm:space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-0">
      <h2 class="text-lg sm:text-xl font-semibold text-foreground">Product Database</h2>
      <div class="flex flex-col sm:flex-row gap-2">
        <button @click="fetchProducts" :disabled="loading" class="inline-flex items-center justify-center px-4 py-2 border border-border bg-card text-card-foreground hover:bg-accent hover:text-accent-foreground hover:border-accent-foreground/20 rounded-md text-sm font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed min-h-[44px] sm:min-h-auto shadow-sm hover:shadow-md">
          <span class="mr-2" :class="{ 'animate-spin': loading }">↻</span>
          Refresh
        </button>
        <button @click="openAddDialog" class="inline-flex items-center justify-center px-4 py-2 bg-secondary text-secondary-foreground hover:bg-secondary/80 hover:shadow-md rounded-md text-sm font-medium transition-all duration-200 min-h-[44px] sm:min-h-auto shadow-sm">
          <span class="mr-2">+</span>
          Add Product
        </button>
      </div>
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
    <div v-if="searchQuery && !loading" class="flex items-center justify-between text-sm text-muted-foreground bg-muted/30 px-4 py-2 rounded-lg">
      <span>
        {{ filteredProducts.length }} of {{ products.length }} product{{ filteredProducts.length !== 1 ? 's' : '' }} found
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
        <span class="text-2xl">📦</span>
      </div>
      <p class="text-muted-foreground font-medium mb-2">No products in database</p>
      <p class="text-muted-foreground text-sm">Click "Add Product" to get started</p>
    </div>

    <div v-else-if="searchQuery && filteredProducts.length === 0" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
        <span class="text-2xl">🔍</span>
      </div>
      <p class="text-muted-foreground font-medium mb-2">No products found</p>
      <p class="text-muted-foreground text-sm">Try adjusting your search terms</p>
    </div>

    <!-- Single Products Display - CSS controls responsive layout -->
    <div v-else>
      <!-- Desktop Table View -->
      <div class="hidden sm:block border border-border rounded-lg p-6 bg-card shadow-sm hover:shadow-md transition-shadow duration-200">
        <div class="overflow-x-auto">
          <table class="w-full border-collapse">
            <thead>
            <tr class="border-b border-border">
              <th class="text-left p-3 font-medium text-card-foreground bg-muted/50 first:rounded-tl-md last:rounded-tr-md">API Name</th>
              <th class="text-left p-3 font-medium text-card-foreground bg-muted/50">Friendly Name</th>
              <th class="text-left p-3 font-medium text-card-foreground bg-muted/50">Referer URL</th>
              <th class="text-left p-3 font-medium text-card-foreground bg-muted/50">Appie ID</th>
              <th class="text-left p-3 font-medium text-card-foreground bg-muted/50 first:rounded-tl-md last:rounded-tr-md">Actions</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="product in filteredProducts" :key="product.appie_id" class="border-b border-border last:border-b-0 hover:bg-muted/30 transition-colors duration-150">
              <td class="p-3 font-mono text-sm text-card-foreground">{{ product.api_name }}</td>
              <td class="p-3 text-card-foreground">{{ product.friendly_name }}</td>
              <td class="p-3">
                <a :href="product.referer_url" target="_blank" class="text-primary hover:text-primary/80 hover:underline text-sm transition-colors duration-150 break-all">
                  {{ product.referer_url }}
                </a>
              </td>
              <td class="p-3 text-card-foreground font-medium">{{ product.appie_id }}</td>
              <td class="p-3">
                <div class="flex gap-2">
                  <button @click="editProduct(product)" class="inline-flex items-center px-3 py-1.5 border border-border bg-card text-card-foreground hover:bg-accent hover:text-accent-foreground hover:border-accent-foreground/20 rounded-md text-xs font-medium transition-all duration-200 shadow-sm hover:shadow-md">
                    <span class="mr-1">✏️</span>
                    Edit
                  </button>
                  <button @click="deleteProduct(product)" class="inline-flex items-center px-3 py-1.5 bg-red-600 text-white hover:bg-red-700 hover:shadow-lg rounded-md text-xs font-medium transition-all duration-200 shadow-sm transform hover:scale-105">
                    <span class="mr-1">🗑️</span>
                    Delete
                  </button>
                </div>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Mobile Card View -->
      <div class="sm:hidden space-y-4">
        <div v-for="product in filteredProducts" :key="product.appie_id" class="border border-border rounded-lg p-4 bg-card shadow-sm hover:shadow-md hover:border-accent-foreground/20 transition-all duration-200 transform hover:scale-[1.02]">
          <div class="space-y-4">
            <div class="flex items-start justify-between gap-3">
              <div class="flex-1 min-w-0">
                <h3 class="font-semibold text-card-foreground text-sm mb-1">{{ product.api_name }}</h3>
                <p class="text-sm text-muted-foreground">{{ product.friendly_name }}</p>
              </div>
              <span class="text-xs text-muted-foreground whitespace-nowrap bg-muted px-2 py-1 rounded-full font-medium">ID: {{ product.appie_id }}</span>
            </div>

            <div class="border-t border-border pt-3">
              <a :href="product.referer_url" target="_blank" class="text-xs text-primary hover:text-primary/80 hover:underline break-all transition-colors duration-150 block">
                {{ product.referer_url }}
              </a>
            </div>

            <div class="flex gap-2 pt-2">
              <button @click="editProduct(product)" class="flex-1 inline-flex items-center justify-center px-3 py-2.5 border border-border bg-card text-card-foreground hover:bg-accent hover:text-accent-foreground hover:border-accent-foreground/20 rounded-md text-sm font-medium transition-all duration-200 min-h-[44px] shadow-sm hover:shadow-md">
                <span class="mr-1">✏️</span>
                Edit
              </button>
              <button @click="deleteProduct(product)" class="flex-1 inline-flex items-center justify-center px-3 py-2.5 bg-red-600 text-white hover:bg-red-700 hover:shadow-lg rounded-md text-sm font-medium transition-all duration-200 min-h-[44px] shadow-sm transform hover:scale-105">
                <span class="mr-1">🗑️</span>
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Product Modal -->
    <ProductModal
      v-if="showDialog"
      :show="showDialog"
      :editing-product="editingProduct"
      :form-data="formData"
      :submitting="submitting"
      @close="showDialog = false"
      @save="saveProduct"
    />
  </div>
</template>

<script setup>
  import { ref, reactive, computed } from 'vue'
  import ProductModal from './ProductModal.vue'

  const products = ref([])
  const loading = ref(false)
  const error = ref('')
  const showDialog = ref(false)
  const editingProduct = ref(null)
  const submitting = ref(false)
  const searchQuery = ref('')

  const formData = reactive({
    api_name: '',
    friendly_name: '',
    referer_url: '',
    appie_id: 0
  })

  // Computed property for filtered products
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

  const resetForm = () => {
    formData.api_name = ''
    formData.friendly_name = ''
    formData.referer_url = ''
    formData.appie_id = 0
  }

  const fetchProducts = async () => {
    loading.value = true
    error.value = ''

    try {
      const response = await fetch('/api/db/products')
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      products.value = Array.isArray(data) ? data : []
    } catch (err) {
      error.value = `Failed to fetch products: ${err.message}`
      console.error('Error fetching products:', err)
    } finally {
      loading.value = false
    }
  }

  const openAddDialog = () => {
    editingProduct.value = null
    resetForm()
    showDialog.value = true
  }

  const editProduct = (product) => {
    editingProduct.value = product
    formData.api_name = product.api_name
    formData.friendly_name = product.friendly_name
    formData.referer_url = product.referer_url
    formData.appie_id = product.appie_id
    showDialog.value = true
  }

  const deleteProduct = async (product) => {
    if (!confirm(`Are you sure you want to delete "${product.friendly_name}"?`)) {
      return
    }

    try {
      // Remove the product from the array
      const updatedProducts = products.value.filter(p => p.appie_id !== product.appie_id)

      // Send the updated array to the API
      const response = await fetch('/api/db/products', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedProducts)
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      // Update local state
      products.value = updatedProducts
    } catch (err) {
      error.value = `Failed to delete product: ${err.message}`
      console.error('Error deleting product:', err)
    }
  }

  const saveProduct = async () => {
    submitting.value = true

    try {
      // Create a copy of the current products array
      let updatedProducts = [...products.value]

      if (editingProduct.value) {
        // Edit existing product - preserve all existing fields and update only form fields
        const index = updatedProducts.findIndex(p => p.appie_id === editingProduct.value.appie_id)
        if (index !== -1) {
          updatedProducts[index] = {
            ...editingProduct.value, // Keep all existing fields
            api_name: formData.api_name,
            friendly_name: formData.friendly_name,
            referer_url: formData.referer_url,
            appie_id: formData.appie_id
          }
        }
      } else {
        // Add new product - include all required fields with defaults
        const newProduct = {
          id: 0, // Default ID
          api_name: formData.api_name,
          friendly_name: formData.friendly_name,
          referer_url: formData.referer_url,
          appie_id: formData.appie_id,
          discount: {
            in_bonus: false,
            description: "",
            label: ""
          }
        }
        updatedProducts.push(newProduct)
      }

      // Send the entire products array
      const response = await fetch('/api/db/products', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedProducts)
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      // Update local state with the new products array
      products.value = updatedProducts

      showDialog.value = false
      resetForm()
    } catch (err) {
      error.value = `Failed to save product: ${err.message}`
      console.error('Error saving product:', err)
    } finally {
      submitting.value = false
    }
  }

  // Expose fetchProducts for parent component
  defineExpose({
    fetchProducts
  })

  // Auto-fetch on mount
  fetchProducts()
</script>
<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-foreground">Product Database</h2>
      <div class="flex gap-2">
        <button @click="fetchProducts" :disabled="loading" class="inline-flex items-center px-4 py-2 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-md text-sm font-medium transition-colors disabled:opacity-50">
          <span class="mr-2" :class="{ 'animate-spin': loading }">↻</span>
          Refresh
        </button>
        <button @click="openAddDialog" class="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground hover:bg-primary/90 rounded-md text-sm font-medium transition-colors">
          <span class="mr-2">+</span>
          Add Product
        </button>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
    </div>

    <div v-else-if="error" class="text-center py-8">
      <p class="text-destructive mb-4">{{ error }}</p>
      <button @click="fetchProducts" class="inline-flex items-center px-4 py-2 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-md text-sm font-medium transition-colors">Try Again</button>
    </div>

    <div v-else-if="products.length === 0" class="text-center py-8">
      <p class="text-muted-foreground">No products in database.</p>
    </div>

    <div v-else class="border border-border rounded-lg p-6 bg-card">
      <div class="overflow-x-auto">
        <table class="w-full border-collapse">
          <thead>
          <tr class="border-b border-border">
            <th class="text-left p-2 font-medium text-card-foreground">API Name</th>
            <th class="text-left p-2 font-medium text-card-foreground">Friendly Name</th>
            <th class="text-left p-2 font-medium text-card-foreground">Referer URL</th>
            <th class="text-left p-2 font-medium text-card-foreground">Appie ID</th>
            <th class="text-left p-2 font-medium text-card-foreground">Actions</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="product in products" :key="product.appie_id" class="border-b border-border last:border-b-0">
            <td class="p-2 font-mono text-sm text-card-foreground">{{ product.api_name }}</td>
            <td class="p-2 text-card-foreground">{{ product.friendly_name }}</td>
            <td class="p-2">
              <a :href="product.referer_url" target="_blank" class="text-primary hover:underline text-sm">
                {{ product.referer_url }}
              </a>
            </td>
            <td class="p-2 text-card-foreground">{{ product.appie_id }}</td>
            <td class="p-2">
              <div class="flex gap-1">
                <button @click="editProduct(product)" class="inline-flex items-center px-2 py-1 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded text-xs font-medium transition-colors">
                  <span class="mr-1">✏️</span>
                  Edit
                </button>
                <button @click="deleteProduct(product)" class="inline-flex items-center px-2 py-1 border border-destructive text-destructive hover:bg-destructive hover:text-destructive-foreground rounded text-xs font-medium transition-colors">
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

    <!-- Add/Edit Product Modal -->
    <div v-if="showDialog" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="fixed inset-0 bg-black/50" @click="showDialog = false"></div>
      <div class="relative bg-card border border-border rounded-lg p-6 max-w-md w-full mx-4 shadow-lg">
        <div class="mb-4">
          <h3 class="text-lg font-semibold text-card-foreground">
            {{ editingProduct ? 'Edit Product' : 'Add New Product' }}
          </h3>
        </div>

        <form @submit.prevent="saveProduct" class="space-y-4">
          <div>
            <label for="api_name" class="block text-sm font-medium text-card-foreground mb-1">API Name</label>
            <input
              id="api_name"
              v-model="formData.api_name"
              placeholder="Enter API name"
              required
              class="w-full px-3 py-2 border border-input bg-background text-foreground rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            />
          </div>

          <div>
            <label for="friendly_name" class="block text-sm font-medium text-card-foreground mb-1">Friendly Name</label>
            <input
              id="friendly_name"
              v-model="formData.friendly_name"
              placeholder="Enter friendly name"
              required
              class="w-full px-3 py-2 border border-input bg-background text-foreground rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            />
          </div>

          <div>
            <label for="referer_url" class="block text-sm font-medium text-card-foreground mb-1">Referer URL</label>
            <input
              id="referer_url"
              v-model="formData.referer_url"
              placeholder="https://example.com"
              type="url"
              required
              class="w-full px-3 py-2 border border-input bg-background text-foreground rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            />
          </div>

          <div>
            <label for="appie_id" class="block text-sm font-medium text-card-foreground mb-1">Appie ID</label>
            <input
              id="appie_id"
              v-model.number="formData.appie_id"
              placeholder="Enter ID"
              type="number"
              required
              class="w-full px-3 py-2 border border-input bg-background text-foreground rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            />
          </div>

          <div class="flex justify-end gap-2 pt-4">
            <button type="button" @click="showDialog = false" class="inline-flex items-center px-4 py-2 border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-md text-sm font-medium transition-colors">
              Cancel
            </button>
            <button type="submit" :disabled="submitting" class="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground hover:bg-primary/90 rounded-md text-sm font-medium transition-colors disabled:opacity-50">
              {{ submitting ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted, reactive } from 'vue'

  const products = ref([])
  const loading = ref(false)
  const error = ref('')
  const showDialog = ref(false)
  const editingProduct = ref(null)
  const submitting = ref(false)

  const formData = reactive({
    api_name: '',
    friendly_name: '',
    referer_url: '',
    appie_id: 0
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

  onMounted(() => {
    fetchProducts()
  })
</script>
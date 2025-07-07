<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm" @click="$emit('close')"></div>
    <div class="relative bg-card border border-border rounded-xl p-6 max-w-md w-full shadow-2xl transform transition-all duration-300 scale-100">
      <div class="mb-6">
        <h3 class="text-xl font-semibold text-card-foreground">
          {{ editingProduct ? 'Edit Product' : 'Add New Product' }}
        </h3>
        <p class="text-sm text-muted-foreground mt-1">
          {{ editingProduct ? 'Update the product information below.' : 'Fill in the details for the new product.' }}
        </p>
      </div>

      <form @submit.prevent="$emit('save')" class="space-y-5">
        <div>
          <label for="api_name" class="block text-sm font-medium text-card-foreground mb-2">API Name</label>
          <input
            id="api_name"
            v-model="formData.api_name"
            placeholder="Enter API name"
            required
            class="w-full px-4 py-3 border border-input bg-background text-foreground rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-all duration-200 min-h-[44px] sm:min-h-auto hover:border-primary/50"
          />
        </div>

        <div>
          <label for="friendly_name" class="block text-sm font-medium text-card-foreground mb-2">Friendly Name</label>
          <input
            id="friendly_name"
            v-model="formData.friendly_name"
            placeholder="Enter friendly name"
            required
            class="w-full px-4 py-3 border border-input bg-background text-foreground rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-all duration-200 min-h-[44px] sm:min-h-auto hover:border-primary/50"
          />
        </div>

        <div>
          <label for="referer_url" class="block text-sm font-medium text-card-foreground mb-2">Referer URL</label>
          <input
            id="referer_url"
            v-model="formData.referer_url"
            placeholder="https://example.com"
            type="url"
            required
            class="w-full px-4 py-3 border border-input bg-background text-foreground rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-all duration-200 min-h-[44px] sm:min-h-auto hover:border-primary/50"
          />
        </div>

        <div>
          <label for="appie_id" class="block text-sm font-medium text-card-foreground mb-2">Appie ID</label>
          <input
            id="appie_id"
            v-model.number="formData.appie_id"
            placeholder="Enter ID"
            type="number"
            required
            class="w-full px-4 py-3 border border-input bg-background text-foreground rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-all duration-200 min-h-[44px] sm:min-h-auto hover:border-primary/50"
          />
        </div>

        <div class="flex flex-col sm:flex-row justify-end gap-3 pt-6 border-t border-border">
          <button type="button" @click="$emit('close')" class="inline-flex items-center justify-center px-6 py-3 border border-border bg-card text-card-foreground hover:bg-accent hover:text-accent-foreground hover:border-accent-foreground/20 rounded-lg text-sm font-medium transition-all duration-200 min-h-[44px] sm:min-h-auto shadow-sm hover:shadow-md">
            Cancel
          </button>
          <button type="submit" :disabled="submitting" class="inline-flex items-center justify-center px-6 py-3 bg-secondary text-secondary-foreground hover:bg-secondary/80 hover:shadow-md rounded-lg text-sm font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed min-h-[44px] sm:min-h-auto shadow-sm">
            {{ submitting ? 'Saving...' : 'Save' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
  defineProps({
    show: {
      type: Boolean,
      required: true
    },
    editingProduct: {
      type: Object,
      default: null
    },
    formData: {
      type: Object,
      required: true
    },
    submitting: {
      type: Boolean,
      default: false
    }
  })

  defineEmits(['close', 'save'])
</script>
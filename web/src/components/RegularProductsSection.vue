<template>
  <div class="border border-border rounded-lg bg-card shadow-sm">
    <button
      @click="$emit('toggle')"
      class="w-full flex items-center justify-between p-4 text-left hover:bg-muted/30 transition-colors duration-200 rounded-t-lg"
    >
      <div class="flex items-center gap-3">
        <span class="text-xl">📦</span>
        <h3 class="text-lg font-semibold text-foreground">Regular Products</h3>
        <div class="px-3 py-1 bg-muted rounded-full">
          <span class="text-sm font-medium text-muted-foreground">{{ products.length }} item{{ products.length !== 1 ? 's' : '' }}</span>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-sm text-muted-foreground">
          {{ show ? 'Hide' : 'Show' }}
        </span>
        <div :class="['transition-transform duration-200', show ? 'rotate-180' : '']">
          <svg class="w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
          </svg>
        </div>
      </div>
    </button>

    <div v-if="show" class="border-t border-border">
      <div class="p-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="product in products" :key="product.appie_id" class="border border-border rounded-lg p-4 bg-card shadow-sm hover:shadow-md hover:border-accent-foreground/20 transition-all duration-200 transform hover:scale-[1.02]">
          <div class="space-y-3 sm:space-y-4">
            <img v-if="product.image" :src="product.image" :alt="product.api_name" class="w-full h-32 object-contain rounded-md bg-muted" />
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
</template>

<script setup>
  defineProps({
    products: {
      type: Array,
      required: true
    },
    show: {
      type: Boolean,
      required: true
    }
  })

  defineEmits(['toggle'])
</script>
<template>
  <div>
    <div class="flex items-center gap-3 mb-4">
      <div class="flex items-center gap-2">
        <span class="text-2xl">⭐</span>
        <h3 class="text-lg font-semibold text-foreground">Bonus Products</h3>
      </div>
      <div class="px-3 py-1 bg-muted border border-border rounded-full">
        <span class="text-sm font-medium text-muted-foreground">{{ products.length }} item{{ products.length !== 1 ? 's' : '' }}</span>
      </div>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="product in products" :key="product.appie_id" class="bonus-card rounded-lg p-4 sm:p-6 transition-all duration-200">
        <div class="space-y-3 sm:space-y-4">
          <img v-if="product.image" :src="product.image" :alt="product.api_name" class="w-full h-32 object-contain rounded-md bg-muted/50" />
          <div class="flex items-start sm:items-center justify-between gap-2">
            <h3 class="font-semibold text-card-foreground text-sm sm:text-base leading-tight">{{ product.api_name }}</h3>
          </div>

          <p class="text-sm text-muted-foreground">
            {{ product.friendly_name }}
          </p>

          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 pt-2 border-t border-border/50">
            <span class="text-xs text-muted-foreground font-medium">
              {{ product.discount.description }}
            </span>
          </div>

          <div class="pt-2">
            <a :href="product.referer_url" target="_blank" class="inline-flex items-center text-xs text-muted-foreground hover:text-muted-foreground/80 hover:underline transition-colors duration-150 min-h-[44px] sm:min-h-auto">
              View Product →
            </a>
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
    }
  })
</script>
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
      <div v-for="product in products" :key="product.appie_id" class="border border-accent-foreground/10 bg-accent/30 rounded-lg p-4 sm:p-6 shadow-sm hover:shadow-md hover:border-accent-foreground/20 transition-all duration-200">
        <div class="space-y-3 sm:space-y-4">
          <div class="flex items-start sm:items-center justify-between gap-2">
            <h3 class="font-semibold text-card-foreground text-sm sm:text-base leading-tight">{{ product.api_name }}</h3>
            <span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-muted text-muted-foreground font-medium whitespace-nowrap">
              ⭐ Bonus
            </span>
          </div>

          <p class="text-sm text-muted-foreground">
            {{ product.friendly_name }}
          </p>

          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 pt-2 border-t border-border">
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
</template>

<script setup>
  defineProps({
    products: {
      type: Array,
      required: true
    }
  })
</script>
<template>
  <div class="fixed top-4 right-4 z-50 space-y-2">
    <div v-for="toast in toasts" :key="toast.id"
         :class="[
           'max-w-sm w-full px-4 py-3 rounded-lg shadow-lg border transform transition-all duration-300',
           toast.type === 'success' ? 'bg-green-50 border-green-200 text-green-800' :
           toast.type === 'error' ? 'bg-red-50 border-red-200 text-red-800' :
           'bg-blue-50 border-blue-200 text-blue-800',
           toast.show ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'
         ]">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span v-if="toast.type === 'success'">✅</span>
          <span v-else-if="toast.type === 'error'">❌</span>
          <span v-else>ℹ️</span>
          <span class="text-sm font-medium">{{ toast.message }}</span>
        </div>
        <button @click="$emit('remove', toast.id)" class="ml-2 text-current opacity-60 hover:opacity-100">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"></path>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
  defineProps({
    toasts: {
      type: Array,
      required: true
    }
  })

  defineEmits(['remove'])
</script>
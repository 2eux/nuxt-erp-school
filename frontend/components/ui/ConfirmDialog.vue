<template>
  <UModal v-model:open="isOpen">
    <template #header>
      <div class="flex items-center gap-3">
        <div
          class="w-10 h-10 rounded-full flex items-center justify-center"
          :class="{
            'bg-red-100 dark:bg-red-900/30': variant === 'danger',
            'bg-amber-100 dark:bg-amber-900/30': variant === 'warning',
            'bg-brand-100 dark:bg-brand-900/30': variant === 'info',
          }"
        >
          <UIcon
            :name="iconMap[variant]"
            class="w-5 h-5"
            :class="{
              'text-red-600': variant === 'danger',
              'text-amber-600': variant === 'warning',
              'text-brand-600': variant === 'info',
            }"
          />
        </div>
        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ title }}
          </h3>
          <p v-if="description" class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
            {{ description }}
          </p>
        </div>
      </div>
    </template>

    <template #body>
      <slot />
    </template>

    <template #footer>
      <div class="flex items-center justify-end gap-3">
        <UButton
          color="gray"
          variant="ghost"
          :disabled="loading"
          @click="handleCancel"
        >
          {{ cancelLabel || $t('common.cancel') }}
        </UButton>
        <UButton
          :color="variant === 'danger' ? 'red' : 'primary'"
          :loading="loading"
          :disabled="loading"
          @click="handleConfirm"
        >
          {{ confirmLabel || $t('common.confirm') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    title: string
    description?: string
    confirmLabel?: string
    cancelLabel?: string
    variant?: 'danger' | 'warning' | 'info'
    loading?: boolean
    modelValue?: boolean
  }>(),
  {
    modelValue: false,
    variant: 'danger',
    loading: false,
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  confirm: []
  cancel: []
}>()

const iconMap = {
  danger: 'i-heroicons-exclamation-triangle',
  warning: 'i-heroicons-exclamation-circle',
  info: 'i-heroicons-information-circle',
}

const isOpen = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  isOpen.value = val
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const handleConfirm = () => {
  emit('confirm')
}

const handleCancel = () => {
  emit('cancel')
  isOpen.value = false
}

const open = () => { isOpen.value = true }
const close = () => { isOpen.value = false }

defineExpose({ open, close })
</script>

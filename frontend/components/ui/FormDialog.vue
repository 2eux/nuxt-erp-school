<template>
  <UModal v-model:open="isOpen" :title="title">
    <template #body>
      <div v-if="description" class="text-sm text-gray-500 dark:text-gray-400 mb-4">
        {{ description }}
      </div>
      <slot />
    </template>

    <template #footer>
      <div class="flex items-center justify-end gap-3">
        <UButton
          color="gray"
          variant="ghost"
          @click="handleCancel"
        >
          {{ cancelLabel || $t('common.cancel') }}
        </UButton>
        <UButton
          color="primary"
          :loading="loading"
          :disabled="loading || disabled"
          @click="handleSubmit"
        >
          {{ submitLabel || $t('common.save') }}
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
    submitLabel?: string
    cancelLabel?: string
    loading?: boolean
    disabled?: boolean
    modelValue?: boolean
  }>(),
  {
    modelValue: false,
    loading: false,
    disabled: false,
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: []
  cancel: []
}>()

const isOpen = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  isOpen.value = val
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = () => {
  emit('submit')
}

const handleCancel = () => {
  emit('cancel')
  isOpen.value = false
}

const open = () => {
  isOpen.value = true
}

const close = () => {
  isOpen.value = false
}

defineExpose({ open, close })
</script>

<template>
  <div class="space-y-2">
    <div
      class="relative border-2 border-dashed rounded-xl p-8 transition-colors cursor-pointer"
      :class="[
        isDragging
          ? 'border-brand-400 bg-brand-50 dark:bg-brand-900/10'
          : 'border-gray-300 dark:border-gray-600 hover:border-brand-400',
        error ? 'border-red-400 bg-red-50 dark:bg-red-900/10' : '',
      ]"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="handleDrop"
      @click="triggerFileInput"
    >
      <input
        ref="fileInput"
        type="file"
        class="hidden"
        :accept="accept"
        :multiple="multiple"
        @change="handleFileSelect"
      />

      <div class="flex flex-col items-center justify-center gap-3">
        <div
          class="w-14 h-14 rounded-full flex items-center justify-center"
          :class="error ? 'bg-red-100' : 'bg-gray-100 dark:bg-gray-800'"
        >
          <UIcon
            :name="error ? 'i-heroicons-exclamation-triangle' : 'i-heroicons-cloud-arrow-up'"
            class="w-7 h-7"
            :class="error ? 'text-red-500' : 'text-gray-400 dark:text-gray-500'"
          />
        </div>
        <div class="text-center">
          <p class="text-sm font-medium text-gray-900 dark:text-white">
            {{ $t('upload.drop_files') }}
          </p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
            {{ $t('upload.or') }}
            <span class="text-brand-600 font-medium cursor-pointer hover:underline">
              {{ $t('upload.browse') }}
            </span>
          </p>
          <p v-if="acceptHint" class="text-xs text-gray-400 mt-1">
            {{ acceptHint }}
          </p>
        </div>
      </div>
    </div>

    <p v-if="error" class="text-xs text-red-500">{{ error }}</p>

    <div v-if="files.length > 0" class="space-y-2">
      <div
        v-for="(file, index) in files"
        :key="index"
        class="flex items-center gap-3 p-3 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700"
      >
        <div class="w-10 h-10 rounded-lg bg-gray-200 dark:bg-gray-700 flex items-center justify-center shrink-0 overflow-hidden">
          <img
            v-if="isImage(file)"
            :src="previews[index]"
            class="w-full h-full object-cover"
            alt=""
          />
          <UIcon
            v-else
            name="i-heroicons-document"
            class="w-5 h-5 text-gray-400"
          />
        </div>

        <div class="flex-1 min-w-0">
          <p class="text-sm text-gray-900 dark:text-white truncate">
            {{ file.name }}
          </p>
          <p class="text-xs text-gray-500">
            {{ formatSize(file.size) }}
          </p>

          <div v-if="uploading" class="mt-1">
            <div class="w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
              <div
                class="h-full bg-brand-500 rounded-full transition-all duration-300"
                :style="{ width: `${progress}%` }"
              />
            </div>
            <p class="text-xs text-brand-600 mt-0.5">{{ progress }}%</p>
          </div>
        </div>

        <UButton
          v-if="!uploading"
          color="gray"
          variant="ghost"
          size="xs"
          icon="i-heroicons-x-mark"
          @click="removeFile(index)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    accept?: string
    multiple?: boolean
    maxSize?: number
    acceptHint?: string
  }>(),
  {
    multiple: false,
    maxSize: 10,
  }
)

const emit = defineEmits<{
  'files-selected': [files: File[]]
  'upload': [files: File[]]
}>()

const fileInput = ref<HTMLInputElement>()
const files = ref<File[]>([])
const previews = ref<string[]>([])
const isDragging = ref(false)
const uploading = ref(false)
const progress = ref(0)
const error = ref('')

const triggerFileInput = () => {
  fileInput.value?.click()
}

const isImage = (file: File): boolean => {
  return file.type.startsWith('image/')
}

const formatSize = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

const validateFile = (file: File): boolean => {
  if (props.maxSize && file.size > props.maxSize * 1024 * 1024) {
    error.value = `File size exceeds ${props.maxSize}MB limit`
    return false
  }
  error.value = ''
  return true
}

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files) {
    const selectedFiles = Array.from(input.files)
    processFiles(selectedFiles)
  }
}

const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  if (event.dataTransfer?.files) {
    const droppedFiles = Array.from(event.dataTransfer.files)
    processFiles(droppedFiles)
  }
}

const processFiles = (newFiles: File[]) => {
  const validFiles = newFiles.filter(validateFile)
  if (validFiles.length === 0) return

  files.value = props.multiple ? [...files.value, ...validFiles] : [validFiles[0]]

  previews.value = files.value.map(file => {
    if (isImage(file)) {
      return URL.createObjectURL(file)
    }
    return ''
  })

  emit('files-selected', files.value)
}

const removeFile = (index: number) => {
  if (previews.value[index]) {
    URL.revokeObjectURL(previews.value[index])
  }
  files.value.splice(index, 1)
  previews.value.splice(index, 1)
}

const uploadFiles = async () => {
  if (files.value.length === 0) return
  uploading.value = true
  progress.value = 0

  const interval = setInterval(() => {
    progress.value = Math.min(progress.value + 10, 90)
  }, 200)

  try {
    emit('upload', files.value)
    progress.value = 100
    setTimeout(() => {
      files.value = []
      previews.value.forEach(p => URL.revokeObjectURL(p))
      previews.value = []
      progress.value = 0
      uploading.value = false
    }, 500)
  } catch {
    uploading.value = false
  } finally {
    clearInterval(interval)
  }
}

const clearFiles = () => {
  files.value = []
  previews.value.forEach(p => URL.revokeObjectURL(p))
  previews.value = []
  error.value = ''
}

defineExpose({ uploadFiles, clearFiles })

onUnmounted(() => {
  previews.value.forEach(p => URL.revokeObjectURL(p))
})
</script>

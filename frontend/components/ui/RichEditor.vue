<template>
  <div class="rich-editor border border-gray-300 dark:border-gray-600 rounded-lg overflow-hidden focus-within:border-brand-500 focus-within:ring-1 focus-within:ring-brand-500">
    <div v-if="!minimal" class="flex items-center gap-0.5 px-2 py-1.5 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 flex-wrap">
      <UButton
        v-for="action in toolbarActions"
        :key="action.key"
        color="gray"
        variant="ghost"
        size="2xs"
        :icon="action.icon"
        :class="{ 'bg-gray-200 dark:bg-gray-700': activeFormats[action.key] }"
        @click="action.action"
      />
      <div class="w-px h-5 bg-gray-300 dark:bg-gray-600 mx-1" />
      <UButton color="gray" variant="ghost" size="2xs" icon="i-heroicons-list-bullet" @click="insertList('unordered')" />
      <UButton color="gray" variant="ghost" size="2xs" icon="i-heroicons-numbered-list" @click="insertList('ordered')" />
      <div class="w-px h-5 bg-gray-300 dark:bg-gray-600 mx-1" />
      <UButton color="gray" variant="ghost" size="2xs" icon="i-heroicons-link" @click="insertLink" />
      <UButton color="gray" variant="ghost" size="2xs" icon="i-heroicons-photo" @click="insertImage" />
    </div>

    <div
      ref="editorRef"
      class="prose prose-sm dark:prose-invert max-w-none min-h-[150px] p-3 outline-none"
      :contenteditable="!disabled"
      :class="{ 'cursor-not-allowed opacity-60': disabled }"
      @input="handleInput"
      @paste="handlePaste"
      @keydown="handleKeydown"
    />

    <input
      ref="linkInput"
      v-model="linkUrl"
      type="url"
      class="hidden"
      @keydown.enter="applyLink"
    />

    <input
      ref="imageInput"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleImageSelect"
    />
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    disabled?: boolean
    minimal?: boolean
  }>(),
  {
    modelValue: '',
    placeholder: '',
    disabled: false,
    minimal: false,
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const editorRef = ref<HTMLDivElement>()
const linkInput = ref<HTMLInputElement>()
const imageInput = ref<HTMLInputElement>()
const linkUrl = ref('')

const activeFormats = ref<Record<string, boolean>>({})

const execFormat = (command: string, value?: string) => {
  document.execCommand(command, false, value)
  editorRef.value?.focus()
  updateActiveFormats()
}

const toolbarActions = computed(() => [
  { key: 'bold', icon: 'i-heroicons-bold', action: () => execFormat('bold') },
  { key: 'italic', icon: 'i-heroicons-italic', action: () => execFormat('italic') },
  { key: 'underline', icon: 'i-heroicons-underline', action: () => execFormat('underline') },
  { key: 'strikeThrough', icon: 'i-heroicons-strikethrough', action: () => execFormat('strikeThrough') },
])

const insertList = (type: 'ordered' | 'unordered') => {
  execFormat(type === 'ordered' ? 'insertOrderedList' : 'insertUnorderedList')
}

const insertLink = () => {
  const selection = window.getSelection()
  if (selection && selection.toString().length > 0) {
    linkUrl.value = ''
    linkInput.value?.click()
  }
}

const applyLink = () => {
  if (linkUrl.value) {
    execFormat('createLink', linkUrl.value)
  }
  linkInput.value?.blur()
}

const insertImage = () => {
  imageInput.value?.click()
}

const handleImageSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files && input.files[0]) {
    const reader = new FileReader()
    reader.onload = (e) => {
      execFormat('insertImage', e.target?.result as string)
    }
    reader.readAsDataURL(input.files[0])
  }
}

const handleInput = () => {
  if (editorRef.value) {
    emit('update:modelValue', editorRef.value.innerHTML)
    updateActiveFormats()
  }
}

const handlePaste = (e: ClipboardEvent) => {
  e.preventDefault()
  const text = e.clipboardData?.getData('text/plain') || ''
  document.execCommand('insertText', false, text)
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Tab') {
    e.preventDefault()
    document.execCommand('insertText', false, '\t')
  }
}

const updateActiveFormats = () => {
  const formats = ['bold', 'italic', 'underline', 'strikeThrough']
  for (const fmt of formats) {
    activeFormats.value[fmt] = document.queryCommandState(fmt)
  }
}

const setContent = (html: string) => {
  if (editorRef.value) {
    editorRef.value.innerHTML = html
  }
}

const getContent = (): string => {
  return editorRef.value?.innerHTML || ''
}

watch(() => props.modelValue, (val) => {
  if (editorRef.value && editorRef.value.innerHTML !== val) {
    editorRef.value.innerHTML = val
  }
})

onMounted(() => {
  if (editorRef.value) {
    editorRef.value.innerHTML = props.modelValue
    if (props.placeholder) {
      editorRef.value.setAttribute('data-placeholder', props.placeholder)
    }
  }
})

defineExpose({ setContent, getContent })
</script>

<style scoped>
[contenteditable]:empty:before {
  content: attr(data-placeholder);
  color: #9ca3af;
}

.prose :deep(a) {
  color: #059669;
  text-decoration: underline;
}

.prose :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin: 0.5rem 0;
}

.prose :deep(ul) {
  list-style-type: disc;
  padding-left: 1.5rem;
}

.prose :deep(ol) {
  list-style-type: decimal;
  padding-left: 1.5rem;
}
</style>

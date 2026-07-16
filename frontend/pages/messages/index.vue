<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('messages.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('messages.subtitle') }}</p></div>
      <UButton color="primary" size="sm" icon="i-heroicons-plus" @click="showCompose=true">{{ $t('messages.compose') }}</UButton>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="card !p-0 overflow-hidden">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700">
          <UInput v-model="searchQuery" :placeholder="$t('common.search')" icon="i-heroicons-magnifying-glass" color="gray" size="sm" />
        </div>
        <div class="divide-y divide-gray-100 dark:divide-gray-750 max-h-[600px] overflow-y-auto">
          <div v-for="msg in filteredMessages" :key="msg.id" class="p-3 hover:bg-gray-50 dark:hover:bg-gray-800/50 cursor-pointer" :class="{ 'bg-brand-50 dark:bg-brand-900/10': selectedMessage?.id === msg.id }" @click="selectMessage(msg)">
            <div class="flex items-center justify-between">
              <span class="text-sm font-semibold text-gray-900 dark:text-white" :class="{ 'text-gray-500': msg.isRead }">{{ msg.senderName }}</span>
              <span class="text-xs text-gray-400">{{ $dayjs(msg.createdAt).format('DD/MM') }}</span>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 truncate mt-0.5" :class="{ 'font-normal': msg.isRead, 'font-semibold': !msg.isRead }">{{ msg.subject }}</p>
          </div>
        </div>
        <EmptyState v-if="filteredMessages.length === 0" :title="$t('messages.no_messages')" icon="i-heroicons-inbox" />
      </div>

      <div class="lg:col-span-2">
        <div v-if="selectedMessage" class="card">
          <div class="card-header">
            <h3 class="card-title">{{ selectedMessage.subject }}</h3>
            <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-arrow-uturn-left" @click="replyToMessage">{{ $t('messages.reply') }}</UButton>
          </div>
          <div class="text-sm text-gray-500 mb-4">
            <span>{{ $t('messages.from') }}: {{ selectedMessage.senderName }}</span> |
            <span>{{ $dayjs(selectedMessage.createdAt).format('DD MMM YYYY HH:mm') }}</span>
          </div>
          <div class="prose prose-sm dark:prose-invert max-w-none" v-html="selectedMessage.content"></div>
        </div>
        <EmptyState v-else :title="$t('messages.select_conversation')" icon="i-heroicons-chat-bubble-left-right" />
      </div>
    </div>

    <FormDialog v-model="showCompose" :title="$t('messages.compose')" :loading="sending" @submit="sendMessage" @cancel="showCompose=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('messages.recipient_type')"><USelect v-model="composeForm.recipientType" :options="[{label:t('messages.individual'),value:'individual'},{label:t('messages.class'),value:'class'},{label:t('messages.role'),value:'role'},{label:t('messages.all'),value:'all'}]" color="gray" /></UFormGroup>
        <UFormGroup v-if="composeForm.recipientType === 'individual'" :label="$t('messages.recipients')"><USelect v-model="composeForm.recipientIds" :options="recipientOptions" color="gray" :multiple="true" /></UFormGroup>
        <UFormGroup v-if="composeForm.recipientType === 'class'" :label="$t('academic.class')"><USelect v-model="composeForm.recipientIds" :options="classOptions" color="gray" :multiple="true" /></UFormGroup>
        <UFormGroup :label="$t('messages.subject')" required><UInput v-model="composeForm.subject" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('messages.content')" required><RichEditor v-model="composeForm.content" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const sending = ref(false); const showCompose = ref(false)
const messages = ref<Record<string, unknown>[]>([]); const selectedMessage = ref<Record<string, unknown> | null>(null)
const searchQuery = ref('')
const recipientOptions = ref<{ label: string; value: string }[]>([]); const classOptions = ref<{ label: string; value: string }[]>([])

const filteredMessages = computed(() => {
  if (!searchQuery.value) return messages.value
  const q = searchQuery.value.toLowerCase()
  return messages.value.filter(m => (m.subject as string)?.toLowerCase().includes(q) || (m.senderName as string)?.toLowerCase().includes(q))
})

const composeForm = reactive({ recipientType: 'individual', recipientIds: [] as string[], subject: '', content: '' })

const fetchMessages = async () => { loading.value = true; try { messages.value = await api.paginate('/messages').then(r => r.data) } catch {} finally { loading.value = false } }
const selectMessage = (msg: Record<string, unknown>) => { selectedMessage.value = msg; if (!msg.isRead) { api.patch(`/messages/${msg.id}/read`); msg.isRead = true } }
const replyToMessage = () => { composeForm.subject = `Re: ${selectedMessage.value?.subject || ''}`; composeForm.recipientType = 'individual'; composeForm.recipientIds = [selectedMessage.value?.senderId as string]; showCompose.value = true }
const sendMessage = async () => { sending.value = true; try { await api.post('/messages', composeForm); toast.add({ title: t('messages.sent'), color: 'success' }); showCompose.value = false; fetchMessages() } catch {} finally { sending.value = false } }
const fetchOptions = async () => { try { recipientOptions.value = (await api.get<{id:string;fullName:string}[]>('/users')).map(u => ({ label: u.fullName, value: u.id })) } catch {}; try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })) } catch {} }
onMounted(() => { fetchMessages(); fetchOptions() })
</script>

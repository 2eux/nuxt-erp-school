<template>
  <div class="space-y-6">
    <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('ai_assistant.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('ai_assistant.subtitle') }}</p></div>

    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <div class="lg:col-span-1 space-y-4">
        <div class="card">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-3">{{ $t('ai_assistant.quick_actions') }}</h3>
          <div class="space-y-2">
            <UButton v-for="action in quickActions" :key="action.key" color="gray" variant="outline" block size="sm" icon="i-heroicons-sparkles" @click="sendPrompt(action.prompt)">{{ action.label }}</UButton>
          </div>
        </div>
        <div class="card">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-3">{{ $t('ai_assistant.history') }}</h3>
          <div class="space-y-1 max-h-64 overflow-y-auto">
            <button v-for="h in conversationHistory" :key="h.id" class="w-full text-left p-2 rounded text-xs text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800/50 truncate" @click="loadConversation(h)">{{ h.title }}</button>
          </div>
          <EmptyState v-if="conversationHistory.length === 0" :title="$t('ai_assistant.no_history')" icon="i-heroicons-chat-bubble-left" />
        </div>
      </div>

      <div class="lg:col-span-3 card flex flex-col">
        <div class="flex-1 space-y-4 mb-4 max-h-[500px] overflow-y-auto p-4" ref="chatContainer">
          <div v-if="messages.length === 0" class="flex items-center justify-center h-64">
            <div class="text-center">
              <UIcon name="i-heroicons-sparkles" class="w-12 h-12 text-brand-400 mx-auto mb-3" />
              <p class="text-gray-500">{{ $t('ai_assistant.greeting') }}</p>
            </div>
          </div>
          <div v-for="(msg, i) in messages" :key="i" class="flex gap-3" :class="msg.role === 'user' ? 'justify-end' : ''">
            <div class="max-w-[80%] rounded-xl p-3" :class="msg.role === 'user' ? 'bg-brand-500 text-white' : 'bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white'">
              <div class="text-sm prose prose-sm dark:prose-invert max-w-none" v-html="renderMarkdown(msg.content)"></div>
              <div v-if="msg.attachments?.length" class="flex flex-wrap gap-1 mt-2">
                <span v-for="a in msg.attachments" :key="a" class="text-xs bg-white/20 dark:bg-black/20 px-2 py-0.5 rounded">{{ a }}</span>
              </div>
            </div>
          </div>
          <div v-if="thinking" class="flex gap-3">
            <div class="bg-gray-100 dark:bg-gray-800 rounded-xl p-3"><div class="flex items-center gap-1"><span class="w-2 h-2 bg-gray-400 rounded-full animate-pulse" /><span class="w-2 h-2 bg-gray-400 rounded-full animate-pulse" style="animation-delay:0.15s" /><span class="w-2 h-2 bg-gray-400 rounded-full animate-pulse" style="animation-delay:0.3s" /></div></div>
          </div>
        </div>
        <div class="border-t border-gray-200 dark:border-gray-700 pt-4">
          <div class="flex items-center gap-2">
            <UButton color="gray" variant="ghost" size="sm" icon="i-heroicons-paper-clip" @click="uploadDocument" />
            <UInput v-model="inputText" :placeholder="$t('ai_assistant.placeholder')" color="gray" class="flex-1" @keydown.enter="sendMessage" />
            <UButton color="primary" icon="i-heroicons-paper-airplane" :loading="thinking" :disabled="!inputText.trim()" @click="sendMessage">{{ $t('common.send') }}</UButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const toast = useToast()
const chatContainer = ref<HTMLElement>()
const inputText = ref(''); const thinking = ref(false)
const messages = ref<{ role: string; content: string; attachments?: string[] }[]>([])
const conversationHistory = ref<{ id: string; title: string; messages: { role: string; content: string }[] }[]>([])
const currentConversationId = ref<string | null>(null)

const quickActions = [
  { key: 'lesson_plan', label: t('ai_assistant.generate_lesson_plan'), prompt: 'Buatkan RPP untuk mata pelajaran...' },
  { key: 'quiz', label: t('ai_assistant.create_quiz'), prompt: 'Buatkan soal ujian pilihan ganda...' },
  { key: 'summary', label: t('ai_assistant.summarize'), prompt: 'Ringkas data kehadiran siswa...' },
  { key: 'analysis', label: t('ai_assistant.data_analysis'), prompt: 'Analisis perkembangan hafalan...' },
]

const sendPrompt = (prompt: string) => { inputText.value = prompt; sendMessage() }
const sendMessage = async () => {
  const text = inputText.value.trim(); if (!text) return
  messages.value.push({ role: 'user', content: text }); inputText.value = ''; thinking.value = true
  scrollToBottom()
  try {
    const res = await api.post<{ response: string }>('/ai/chat', { message: text, conversationId: currentConversationId.value })
    messages.value.push({ role: 'assistant', content: res.response })
    if (!currentConversationId.value) { currentConversationId.value = res.currentConversationId; fetchHistory() }
  } catch (e) { messages.value.push({ role: 'assistant', content: 'Maaf, terjadi kesalahan. Silakan coba lagi.' }) }
  finally { thinking.value = false; scrollToBottom() }
}
const loadConversation = (conv: { id: string; title: string; messages: { role: string; content: string }[] }) => { currentConversationId.value = conv.id; messages.value = conv.messages }
const uploadDocument = () => { const input = document.createElement('input'); input.type = 'file'; input.accept = '.pdf,.doc,.docx,.txt'; input.onchange = () => { if (input.files?.[0]) { toast.add({ title: 'Dokumen siap dianalisis', color: 'info' }) } }; input.click() }
const fetchHistory = async () => { try { conversationHistory.value = await api.get('/ai/conversations') } catch {} }
const renderMarkdown = (text: string) => text.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>').replace(/\*(.*?)\*/g, '<em>$1</em>').replace(/`(.*?)`/g, '<code>$1</code>').replace(/\n/g, '<br/>')
const scrollToBottom = () => { nextTick(() => { if (chatContainer.value) chatContainer.value.scrollTop = chatContainer.value.scrollHeight }) }
onMounted(() => { if (messages.value.length === 0) messages.value.push({ role: 'assistant', content: 'Assalamu\'alaikum! Saya asisten AI untuk membantu Anda. Silakan ajukan pertanyaan atau gunakan tombol aksi cepat di samping.' }) })
</script>

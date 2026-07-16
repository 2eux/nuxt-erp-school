<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('letters.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('letters.subtitle') }}</p></div>
      <UButton v-if="permissions.can('letters.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openCreate">{{ $t('letters.create') }}</UButton>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('letters.templates') }}</h3></div>
        <div class="grid grid-cols-2 gap-3">
          <button v-for="tmpl in templates" :key="tmpl.key" class="flex flex-col items-center gap-2 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 hover:border-brand-300 transition-all text-center" @click="useTemplate(tmpl)">
            <UIcon :name="tmpl.icon" class="w-6 h-6 text-gray-500" />
            <span class="text-xs font-medium text-gray-900 dark:text-white">{{ tmpl.label }}</span>
          </button>
        </div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('letters.recent') }}</h3></div>
        <div class="space-y-3">
          <div v-for="letter in recentLetters" :key="letter.id" class="flex items-center justify-between p-2 rounded hover:bg-gray-50 dark:hover:bg-gray-800/50">
            <div><p class="text-sm font-medium text-gray-900 dark:text-white">{{ letter.title }}</p><p class="text-xs text-gray-500">{{ $dayjs(letter.createdAt).format('DD MMM YYYY') }}</p></div>
            <div class="flex items-center gap-1">
              <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-printer" @click="printLetter(letter)" />
              <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-arrow-down-tray" @click="downloadLetter(letter)" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <DataTable :columns="columns" :rows="letters" :loading="loading" :empty-title="$t('letters.no_letters')" :show-export="false">
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewLetter(row as Record<string, unknown>)" />
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-printer" @click="printLetter(row)" />
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="$t('letters.create')" :loading="saving" @submit="generateLetter" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('letters.template')"><USelect v-model="form.template" :options="templateOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('letters.title')" required><UInput v-model="form.title" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('letters.number')"><UInput v-model="form.number" color="gray" /></UFormGroup>
        <div v-for="(val, key) in form.variables" :key="key">
          <UFormGroup :label="key"><UInput v-model="form.variables[key]" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('letters.content')"><RichEditor v-model="form.content" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const saving = ref(false); const showForm = ref(false)
const letters = ref<Record<string, unknown>[]>([])

const templates = [
  { key: 'certificate', label: t('letters.certificate'), icon: 'i-heroicons-document-check', variables: { studentName: '', program: '', date: '' } },
  { key: 'warning', label: t('letters.warning_letter'), icon: 'i-heroicons-exclamation-triangle', variables: { studentName: '', reason: '', parentName: '' } },
  { key: 'transfer', label: t('letters.transfer_letter'), icon: 'i-heroicons-arrow-right-on-rectangle', variables: { studentName: '', destinationSchool: '' } },
  { key: 'permit', label: t('letters.permit_letter'), icon: 'i-heroicons-document-text', variables: { studentName: '', reason: '', date: '' } },
  { key: 'invitation', label: t('letters.invitation'), icon: 'i-heroicons-envelope', variables: { recipient: '', event: '', date: '', time: '', location: '' } },
]
const templateOpts = templates.map(t => ({ label: t.label, value: t.key }))
const columns: TableColumn[] = [
  { key: 'number', label: 'letters.number' }, { key: 'title', label: 'letters.title', sortable: true },
  { key: 'type', label: 'letters.template' }, { key: 'createdAt', label: 'common.created_at', type: 'date' },
]
const recentLetters = computed(() => letters.value.slice(0, 5))

const form = reactive({ template: '', title: '', number: '', content: '', variables: {} as Record<string, string> })

const fetchLetters = async () => { loading.value = true; try { letters.value = await api.paginate('/letters').then(r => r.data) } catch {} finally { loading.value = false } }
const openCreate = () => { Object.assign(form, { template: '', title: '', number: '', content: '', variables: {} }); showForm.value = true }
const useTemplate = (tmpl: typeof templates[0]) => { form.template = tmpl.key; form.variables = { ...tmpl.variables }; form.title = tmpl.label; showForm.value = true }
const generateLetter = async () => { saving.value = true; try { await api.post('/letters', form); toast.add({ title: t('letters.created'), color: 'success' }); showForm.value = false; fetchLetters() } catch {} finally { saving.value = false } }
const viewLetter = (row: Record<string, unknown>) => { window.open(`/api/v1/letters/${row.id}/view`, '_blank') }
const printLetter = (row: Record<string, unknown>) => { window.open(`/api/v1/letters/${row.id}/pdf`, '_blank') }
const downloadLetter = (row: Record<string, unknown>) => { window.open(`/api/v1/letters/${row.id}/pdf?download=1`, '_blank') }
onMounted(() => fetchLetters())
</script>

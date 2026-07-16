<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.report_cards') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.report_cards_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedClass" :options="classOptions" :placeholder="$t('academic.select_class')" color="gray" size="sm" class="w-36" @change="fetchReportCards" />
        <UButton v-if="permissions.can('academic.report_cards.generate')" color="primary" size="sm" icon="i-heroicons-plus" @click="openGenerate">{{ $t('academic.generate_report_card') }}</UButton>
      </div>
    </div>

    <DataTable :columns="columns" :rows="reportCards" :loading="loading" :empty-title="$t('academic.no_report_cards')" :show-export="true" @export="handleExport">
      <template #cell-gpa="{ row }">
        <span class="text-sm font-semibold font-mono text-gray-900 dark:text-white">{{ row.gpa }}</span>
      </template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewReportCard(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-printer" @click="printReportCard(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showGenerate" :title="$t('academic.generate_report_card')" :loading="generating" @submit="generateReportCards" @cancel="showGenerate=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('academic.class')" required><USelect v-model="genForm.classId" :options="classOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.term')" required><USelect v-model="genForm.termId" :options="termOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.teacher_note')"><UTextarea v-model="genForm.teacherNote" color="gray" :rows="2" /></UFormGroup>
        <UFormGroup :label="$t('academic.principal_note')"><UTextarea v-model="genForm.principalNote" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('academic.delete_report_card')" :loading="deleting" @confirm="deleteReportCard" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const loading = ref(false)
const generating = ref(false)
const deleting = ref(false)
const showGenerate = ref(false)
const showDelete = ref(false)
const selectedClass = ref('')
const reportCards = ref<Record<string, unknown>[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const termOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)

const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name', sortable: true },
  { key: 'className', label: 'academic.class' },
  { key: 'term', label: 'academic.term' },
  { key: 'gpa', label: 'academic.gpa', sortable: true },
  { key: 'rank', label: 'academic.rank' },
  { key: 'createdAt', label: 'common.created_at', type: 'date' },
]
const genForm = reactive({ classId: '', termId: '', teacherNote: '', principalNote: '' })

const fetchReportCards = async () => {
  if (!selectedClass.value) return; loading.value = true
  try { reportCards.value = await api.paginate('/report-cards', { classId: selectedClass.value }).then(r => r.data) }
  catch {} finally { loading.value = false }
}

const openGenerate = () => { genForm.classId = selectedClass.value; genForm.termId = ''; genForm.teacherNote = ''; genForm.principalNote = ''; showGenerate.value = true }

const generateReportCards = async () => {
  generating.value = true
  try { await api.post('/report-cards/generate', genForm); toast.add({ title: t('academic.report_cards_generated'), color: 'success' }); showGenerate.value = false; fetchReportCards() }
  catch {} finally { generating.value = false }
}

const viewReportCard = (row: Record<string, unknown>) => { window.open(`/api/v1/report-cards/${row.id}/view`, '_blank') }
const printReportCard = (row: Record<string, unknown>) => { window.open(`/api/v1/report-cards/${row.id}/pdf`, '_blank') }

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteReportCard = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/report-cards/${deleteTarget.value.id}`); toast.add({ title: t('academic.report_card_deleted'), color: 'success' }); showDelete.value = false; fetchReportCards() } catch {} finally { deleting.value = false } }

const handleExport = (format: string) => { window.open(`/api/v1/report-cards/export?classId=${selectedClass.value}&format=${format}`, '_blank') }

const fetchOptions = async () => {
  try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })) } catch {}
  try { termOptions.value = (await api.get<{id:string;name:string}[]>('/terms')).map(t => ({ label: t.name, value: t.id })) } catch {}
}

onMounted(() => fetchOptions())
</script>

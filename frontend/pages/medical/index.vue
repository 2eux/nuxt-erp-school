<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('medical.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('medical.subtitle') }}</p></div>
      <UButton v-if="permissions.can('medical.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('medical.add_record') }}</UButton>
    </div>
    <DataFilter :filter-fields="filterFields" :searchable="true" @apply="handleFilter" />
    <DataTable :columns="columns" :rows="records" :loading="loading" :empty-title="$t('medical.no_records')" :show-export="true" @export="handleExport">
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewRecord(row as Record<string, unknown>)" />
      </template>
    </DataTable>
    <FormDialog v-model="showForm" :title="$t('medical.add_record')" :loading="saving" @submit="saveRecord" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('tahfidz.student')" required><USelect v-model="form.studentId" :options="studentOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.date')"><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('medical.complaint')" required><UTextarea v-model="form.complaint" color="gray" :rows="2" /></UFormGroup>
        <UFormGroup :label="$t('medical.diagnosis')"><UTextarea v-model="form.diagnosis" color="gray" :rows="2" /></UFormGroup>
        <UFormGroup :label="$t('medical.treatment')"><UTextarea v-model="form.treatment" color="gray" :rows="2" /></UFormGroup>
        <UFormGroup :label="$t('medical.medication')"><UInput v-model="form.medication" color="gray" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const saving = ref(false); const showForm = ref(false)
const records = ref<Record<string, unknown>[]>([]); const studentOptions = ref<{ label: string; value: string }[]>([])
const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name', sortable: true },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'complaint', label: 'medical.complaint' },
  { key: 'diagnosis', label: 'medical.diagnosis' },
  { key: 'recordedByName', label: 'medical.recorded_by' },
]
const filterFields = [{ key: 'studentId', label: 'tahfidz.student', type: 'select' as const, options: [] }]
const form = reactive({ studentId: '', date: $dayjs().format('YYYY-MM-DD'), complaint: '', diagnosis: '', treatment: '', medication: '' })

const fetchRecords = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { records.value = await api.paginate('/medical-records', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchRecords(filters)
const openAdd = () => { Object.assign(form, { studentId: '', date: $dayjs().format('YYYY-MM-DD'), complaint: '', diagnosis: '', treatment: '', medication: '' }); showForm.value = true }
const saveRecord = async () => { saving.value = true; try { await api.post('/medical-records', form); toast.add({ title: t('medical.record_created'), color: 'success' }); showForm.value = false; fetchRecords() } catch {} finally { saving.value = false } }
const viewRecord = (row: Record<string, unknown>) => { navigateTo(`/medical/${row.id}`) }
const handleExport = (format: string) => { window.open(`/api/v1/medical-records/export?format=${format}`, '_blank') }
const fetchOptions = async () => { try { studentOptions.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })); filterFields[0].options = studentOptions.value } catch {} }
onMounted(() => { fetchRecords(); fetchOptions() })
</script>

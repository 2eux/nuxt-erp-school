<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('islamic.tasmi') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('islamic.tasmi_subtitle') }}</p></div>
      <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('islamic.record_tasmi') }}</UButton>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="records" :loading="loading" :empty-title="$t('islamic.no_tasmi')" :show-export="false">
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'passed'" /></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewDetail(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editRecord(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('islamic.edit_tasmi') : $t('islamic.record_tasmi')" :loading="saving" @submit="saveRecord" @cancel="showForm=false">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('tahfidz.student')" required><USelect v-model="form.studentId" :options="studentOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('tahfidz.date')" required><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('quran.surah')" required><UInput v-model="form.surah" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('quran.ayah_range')" required><UInput v-model="form.ayahRange" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('islamic.examiner')"><UInput v-model="form.examiner" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('islamic.score')"><UInput v-model="form.score" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('islamic.status')"><USelect v-model="form.status" :options="[{label:t('islamic.passed'),value:'passed'},{label:t('islamic.not_passed'),value:'not_passed'}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.note')" class="sm:col-span-2"><UTextarea v-model="form.note" color="gray" :rows="3" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('islamic.delete_tasmi')" :loading="deleting" @confirm="deleteRecord" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const { $dayjs } = useNuxtApp()
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const showForm = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const records = ref<Record<string, unknown>[]>([])
const studentOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)

const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name', sortable: true },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'surah', label: 'quran.surah' },
  { key: 'ayahRange', label: 'quran.ayah' },
  { key: 'examiner', label: 'islamic.examiner' },
  { key: 'score', label: 'islamic.score' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [
  { key: 'studentId', label: 'tahfidz.student', type: 'select' as const, options: [] },
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('islamic.passed'), value: 'passed' }, { label: t('islamic.not_passed'), value: 'not_passed' }] },
]
const form = reactive({ studentId: '', date: $dayjs().format('YYYY-MM-DD'), surah: '', ayahRange: '', examiner: '', score: '', status: 'passed', note: '' })

const fetchRecords = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { records.value = await api.paginate('/tasmi', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchRecords(filters)

const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { studentId: '', date: $dayjs().format('YYYY-MM-DD'), surah: '', ayahRange: '', examiner: '', score: '', status: 'passed', note: '' }); showForm.value = true }
const editRecord = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const viewDetail = (row: Record<string, unknown>) => { /* detail view */ }

const saveRecord = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/tasmi/${editId.value}`, form); toast.add({ title: t('islamic.tasmi_updated'), color: 'success' }) }
    else { await api.post('/tasmi', form); toast.add({ title: t('islamic.tasmi_recorded'), color: 'success' }) }
    showForm.value = false; fetchRecords()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteRecord = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/tasmi/${deleteTarget.value.id}`); toast.add({ title: t('islamic.tasmi_deleted'), color: 'success' }); showDelete.value = false; fetchRecords() } catch {} finally { deleting.value = false } }

const fetchStudents = async () => { try { studentOptions.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })); filterFields[0].options = studentOptions.value } catch {} }

onMounted(() => { fetchRecords(); fetchStudents() })
</script>

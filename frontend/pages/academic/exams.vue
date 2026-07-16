<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.exams') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.exams_subtitle') }}</p></div>
      <UButton v-if="permissions.can('academic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('academic.add_exam') }}</UButton>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="exams" :loading="loading" :empty-title="$t('academic.no_exams')" :show-export="false">
      <template #cell-type="{ row }">
        <StatusBadge :status="(row.type as string) || 'midterm'" />
      </template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editExam(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-table-cells" @click="openScoreEntry(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('academic.edit_exam') : $t('academic.add_exam')" :loading="saving" @submit="saveExam" @cancel="closeDialog">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('academic.exam_name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.exam_type')" required><USelect v-model="form.type" :options="typeOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('subjects.title')" required><USelect v-model="form.subjectId" :options="subjectOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.class')" required><USelect v-model="form.classId" :options="classOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.exam_date')" required><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.term')"><USelect v-model="form.termId" :options="termOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.start_time')"><UInput v-model="form.startTime" type="time" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.end_time')"><UInput v-model="form.endTime" type="time" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.max_score')"><UInput v-model.number="form.maxScore" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.passing_score')"><UInput v-model.number="form.passingScore" type="number" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('academic.delete_exam')" :loading="deleting" @confirm="deleteExam" @cancel="showDelete=false" />
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
const saving = ref(false)
const deleting = ref(false)
const exams = ref<Record<string, unknown>[]>([])
const showDialog = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const deleteTarget = ref<Record<string, unknown> | null>(null)
const classOptions = ref<{ label: string; value: string }[]>([])
const subjectOptions = ref<{ label: string; value: string }[]>([])
const termOptions = ref<{ label: string; value: string }[]>([])

const columns: TableColumn[] = [
  { key: 'name', label: 'academic.exam_name', sortable: true },
  { key: 'type', label: 'academic.exam_type', sortable: true },
  { key: 'subjectName', label: 'subjects.title' },
  { key: 'className', label: 'academic.class' },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'startTime', label: 'academic.start_time' },
  { key: 'endTime', label: 'academic.end_time' },
]
const filterFields = [
  { key: 'classId', label: 'academic.class', type: 'select' as const, options: [] },
  { key: 'type', label: 'academic.exam_type', type: 'select' as const, options: [
    { label: t('academic.midterm'), value: 'midterm' }, { label: t('academic.final'), value: 'final' },
    { label: t('academic.daily'), value: 'daily' }, { label: t('academic.semester'), value: 'semester' },
  ]},
]
const typeOptions = [
  { label: t('academic.midterm'), value: 'midterm' }, { label: t('academic.final'), value: 'final' },
  { label: t('academic.daily'), value: 'daily' }, { label: t('academic.semester'), value: 'semester' },
]

const form = reactive({ name: '', type: 'midterm', subjectId: '', classId: '', date: '', termId: '', startTime: '07:30', endTime: '09:00', maxScore: 100, passingScore: 75 })

const fetchExams = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { exams.value = await api.paginate('/exams', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchExams(filters)

const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { name: '', type: 'midterm', subjectId: '', classId: '', date: '', termId: '', startTime: '07:30', endTime: '09:00', maxScore: 100, passingScore: 75 }); showDialog.value = true }
const editExam = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false }
const openScoreEntry = (row: Record<string, unknown>) => { navigateTo(`/academic/exams/${row.id}/scores`) }

const saveExam = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/exams/${editId.value}`, form); toast.add({ title: t('academic.exam_updated'), color: 'success' }) }
    else { await api.post('/exams', form); toast.add({ title: t('academic.exam_created'), color: 'success' }) }
    closeDialog(); fetchExams()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteExam = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/exams/${deleteTarget.value.id}`); toast.add({ title: t('academic.exam_deleted'), color: 'success' }); showDelete.value = false; fetchExams() } catch {} finally { deleting.value = false } }

const fetchOptions = async () => {
  try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); filterFields[0].options = classOptions.value } catch {}
  try { subjectOptions.value = (await api.get<{id:string;name:string}[]>('/subjects')).map(s => ({ label: s.name, value: s.id })) } catch {}
  try { termOptions.value = (await api.get<{id:string;name:string}[]>('/terms')).map(t => ({ label: t.name, value: t.id })) } catch {}
}

onMounted(() => { fetchExams(); fetchOptions() })
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.assignments') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.assignments_subtitle') }}</p></div>
      <UButton v-if="permissions.can('academic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('academic.add_assignment') }}</UButton>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="assignments" :loading="loading" :empty-title="$t('academic.no_assignments')" :show-export="false">
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editAssignment(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-clipboard-document-check" @click="viewSubmissions(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('academic.edit_assignment') : $t('academic.add_assignment')" :loading="saving" @submit="saveAssignment" @cancel="closeDialog">
      <div class="space-y-4">
        <UFormGroup :label="$t('academic.title')" required><UInput v-model="form.title" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.description')"><RichEditor v-model="form.description" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('subjects.title')" required><USelect v-model="form.subjectId" :options="subjectOptions" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('academic.class')" required><USelect v-model="form.classId" :options="classOptions" color="gray" /></UFormGroup>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('academic.due_date')" required><UInput v-model="form.dueDate" type="datetime-local" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('academic.max_score')"><UInput v-model.number="form.maxScore" type="number" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('academic.attachment')"><FileUpload accept=".pdf,.doc,.docx,.jpg,.png" :multiple="true" @files-selected="handleFilesSelected" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('academic.delete_assignment')" :loading="deleting" @confirm="deleteAssignment" @cancel="showDelete=false" />
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
const showDialog = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const assignments = ref<Record<string, unknown>[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const subjectOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)
const selectedFiles = ref<File[]>([])

const columns: TableColumn[] = [
  { key: 'title', label: 'academic.title', sortable: true },
  { key: 'subjectName', label: 'subjects.title' },
  { key: 'className', label: 'academic.class' },
  { key: 'dueDate', label: 'academic.due_date', type: 'date' },
  { key: 'maxScore', label: 'academic.max_score', type: 'number' },
  { key: 'submissionCount', label: 'academic.submissions', type: 'number' },
  { key: 'createdAt', label: 'common.created_at', type: 'date' },
]
const filterFields = [
  { key: 'classId', label: 'academic.class', type: 'select' as const, options: [] },
  { key: 'subjectId', label: 'subjects.title', type: 'select' as const, options: [] },
]

const form = reactive({ title: '', description: '', subjectId: '', classId: '', dueDate: '', maxScore: 100 })

const handleFilesSelected = (files: File[]) => { selectedFiles.value = files }

const fetchAssignments = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { assignments.value = await api.paginate('/assignments', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchAssignments(filters)

const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { title: '', description: '', subjectId: '', classId: '', dueDate: '', maxScore: 100 }); showDialog.value = true }
const editAssignment = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false }
const viewSubmissions = (row: Record<string, unknown>) => { navigateTo(`/academic/assignments/${row.id}/submissions`) }

const saveAssignment = async () => {
  saving.value = true
  try {
    const formData = new FormData()
    formData.append('data', JSON.stringify(form))
    selectedFiles.value.forEach(f => formData.append('files', f))
    if (editing.value && editId.value) { await api.put(`/assignments/${editId.value}`, form); toast.add({ title: t('academic.assignment_updated'), color: 'success' }) }
    else { await api.upload('/assignments', formData); toast.add({ title: t('academic.assignment_created'), color: 'success' }) }
    closeDialog(); fetchAssignments()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteAssignment = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/assignments/${deleteTarget.value.id}`); toast.add({ title: t('academic.assignment_deleted'), color: 'success' }); showDelete.value = false; fetchAssignments() } catch {} finally { deleting.value = false } }

const fetchOptions = async () => {
  try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); filterFields[0].options = classOptions.value } catch {}
  try { subjectOptions.value = (await api.get<{id:string;name:string}[]>('/subjects')).map(s => ({ label: s.name, value: s.id })); filterFields[1].options = subjectOptions.value } catch {}
}

onMounted(() => { fetchAssignments(); fetchOptions() })
</script>

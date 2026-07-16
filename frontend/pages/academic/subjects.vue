<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.subjects') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.subjects_subtitle') }}</p></div>
      <UButton v-if="permissions.can('academic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('academic.add_subject') }}</UButton>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="subjects" :loading="loading" :empty-title="$t('academic.no_subjects')" :show-export="false">
      <template #cell-category="{ row }">
        <span class="inline-flex px-2 py-0.5 rounded-full text-xs font-medium"
          :class="categoryClass(row.category as string)">
          {{ $t(`subjects.category_${row.category}`) || row.category }}
        </span>
      </template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editSubject(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('academic.edit_subject') : $t('academic.add_subject')" :loading="saving" @submit="saveSubject" @cancel="closeDialog">
      <div class="space-y-4">
        <UFormGroup :label="$t('subjects.code')" required><UInput v-model="form.code" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('subjects.name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('subjects.description')" class="sm:col-span-2"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('subjects.category')" required>
            <USelect v-model="form.category" :options="categoryOptions" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('subjects.credit_hours')"><UInput v-model.number="form.creditHours" type="number" color="gray" /></UFormGroup>
        </div>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('academic.delete_subject')" :loading="deleting" @confirm="deleteSubject" @cancel="showDelete=false" />
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
const subjects = ref<Record<string, unknown>[]>([])
const showDialog = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const deleteTarget = ref<Record<string, unknown> | null>(null)

const columns: TableColumn[] = [
  { key: 'code', label: 'subjects.code', sortable: true },
  { key: 'name', label: 'subjects.name', sortable: true },
  { key: 'category', label: 'subjects.category' },
  { key: 'creditHours', label: 'subjects.credit_hours', type: 'number' },
]
const filterFields = [
  { key: 'category', label: 'subjects.category', type: 'select' as const, options: [
    { label: t('subjects.category_academic'), value: 'academic' }, { label: t('subjects.category_islamic'), value: 'islamic' },
    { label: t('subjects.category_language'), value: 'language' }, { label: t('subjects.category_extracurricular'), value: 'extracurricular' },
  ]},
]
const categoryOptions = [
  { label: t('subjects.category_academic'), value: 'academic' }, { label: t('subjects.category_islamic'), value: 'islamic' },
  { label: t('subjects.category_language'), value: 'language' }, { label: t('subjects.category_extracurricular'), value: 'extracurricular' }, { label: t('subjects.category_other'), value: 'other' },
]

const form = reactive({ code: '', name: '', description: '', category: 'academic', creditHours: 2 })

const categoryClass = (cat: string) => {
  const map: Record<string, string> = {
    academic: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400',
    islamic: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400',
    language: 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400',
    extracurricular: 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400',
  }
  return map[cat] || 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400'
}

const fetchSubjects = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { subjects.value = await api.paginate('/subjects', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchSubjects(filters)

const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { code: '', name: '', description: '', category: 'academic', creditHours: 2 }); showDialog.value = true }
const editSubject = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false }

const saveSubject = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/subjects/${editId.value}`, form); toast.add({ title: t('academic.subject_updated'), color: 'success' }) }
    else { await api.post('/subjects', form); toast.add({ title: t('academic.subject_created'), color: 'success' }) }
    closeDialog(); fetchSubjects()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteSubject = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/subjects/${deleteTarget.value.id}`); toast.add({ title: t('academic.subject_deleted'), color: 'success' }); showDelete.value = false; fetchSubjects() } catch {} finally { deleting.value = false } }

onMounted(() => fetchSubjects())
</script>

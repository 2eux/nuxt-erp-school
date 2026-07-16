<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.curriculum') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.curriculum_subtitle') }}</p></div>
      <UButton v-if="permissions.can('academic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('academic.add_curriculum') }}</UButton>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="curriculumItems" :loading="loading" :empty-title="$t('academic.no_curriculum')" :show-export="false">
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editItem(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('academic.edit_curriculum') : $t('academic.add_curriculum')" :loading="saving" @submit="saveCurriculum" @cancel="closeDialog">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('academic.class')" required><USelect v-model="form.classId" :options="classOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('subjects.title')" required><USelect v-model="form.subjectId" :options="subjectOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.title')"><USelect v-model="form.teacherId" :options="teacherOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.day')"><USelect v-model="form.scheduleDay" :options="dayOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.start_time')"><UInput v-model="form.scheduleStart" type="time" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.end_time')"><UInput v-model="form.scheduleEnd" type="time" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.room')"><UInput v-model="form.room" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('academic.delete_curriculum')" :loading="deleting" @confirm="deleteItem" @cancel="showDelete=false" />
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
const curriculumItems = ref<Record<string, unknown>[]>([])
const showDialog = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const deleteTarget = ref<Record<string, unknown> | null>(null)
const classOptions = ref<{ label: string; value: string }[]>([])
const subjectOptions = ref<{ label: string; value: string }[]>([])
const teacherOptions = ref<{ label: string; value: string }[]>([])

const columns: TableColumn[] = [
  { key: 'className', label: 'academic.class', sortable: true },
  { key: 'subjectName', label: 'subjects.title' },
  { key: 'teacherName', label: 'teachers.title' },
  { key: 'scheduleDay', label: 'academic.day' },
  { key: 'scheduleStart', label: 'academic.start_time' },
  { key: 'scheduleEnd', label: 'academic.end_time' },
  { key: 'room', label: 'academic.room' },
]
const filterFields = [
  { key: 'classId', label: 'academic.class', type: 'select' as const, options: [] },
  { key: 'subjectId', label: 'subjects.title', type: 'select' as const, options: [] },
]
const dayOptions = [
  { label: t('day.monday'), value: 'monday' }, { label: t('day.tuesday'), value: 'tuesday' }, { label: t('day.wednesday'), value: 'wednesday' },
  { label: t('day.thursday'), value: 'thursday' }, { label: t('day.friday'), value: 'friday' }, { label: t('day.saturday'), value: 'saturday' },
]

const form = reactive({ classId: '', subjectId: '', teacherId: '', scheduleDay: 'monday', scheduleStart: '', scheduleEnd: '', room: '' })

const fetchCurriculum = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { curriculumItems.value = await api.paginate('/curriculum', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchCurriculum(filters)

const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { classId: '', subjectId: '', teacherId: '', scheduleDay: 'monday', scheduleStart: '', scheduleEnd: '', room: '' }); showDialog.value = true }
const editItem = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false }

const saveCurriculum = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/curriculum/${editId.value}`, form); toast.add({ title: t('academic.curriculum_updated'), color: 'success' }) }
    else { await api.post('/curriculum', form); toast.add({ title: t('academic.curriculum_created'), color: 'success' }) }
    closeDialog(); fetchCurriculum()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteItem = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/curriculum/${deleteTarget.value.id}`); toast.add({ title: t('academic.curriculum_deleted'), color: 'success' }); showDelete.value = false; fetchCurriculum() } catch {} finally { deleting.value = false } }

const fetchOptions = async () => {
  try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); filterFields[0].options = classOptions.value } catch {}
  try { subjectOptions.value = (await api.get<{id:string;name:string}[]>('/subjects')).map(s => ({ label: s.name, value: s.id })); filterFields[1].options = subjectOptions.value } catch {}
  try { teacherOptions.value = (await api.get<{id:string;fullName:string}[]>('/teachers', { limit: 100 })).map(t => ({ label: t.fullName, value: t.id })) } catch {}
}

onMounted(() => { fetchCurriculum(); fetchOptions() })
</script>

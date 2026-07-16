<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.classes') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.classes_subtitle') }}</p></div>
      <UButton v-if="permissions.can('academic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('academic.add_class') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('academic.total_classes')" :value="stats.total" icon="i-heroicons-building-office" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('academic.total_students')" :value="stats.totalStudents" icon="i-heroicons-users" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('academic.avg_class_size')" :value="stats.avgClassSize" icon="i-heroicons-chart-bar" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('academic.active')" :value="stats.active" icon="i-heroicons-check-circle" color="purple" :loading="statsLoading" />
    </div>

    <DataTable :columns="columns" :rows="classes" :loading="loading" :empty-title="$t('academic.no_classes')" :show-export="false">
      <template #cell-name="{ row }">
        <NuxtLink :to="`/academic/classes/${row.id}`" class="text-sm font-medium text-brand-600 dark:text-brand-400 hover:underline">{{ row.name }}</NuxtLink>
      </template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" :to="`/academic/classes/${row.id}`" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editClass(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('academic.edit_class') : $t('academic.add_class')" :loading="saving" @submit="saveClass" @cancel="closeDialog">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('academic.class_name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.grade')" required><USelect v-model="form.grade" :options="gradeOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.program')"><UInput v-model="form.program" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.capacity')"><UInput v-model.number="form.capacity" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.homeroom_teacher')"><USelect v-model="form.homeRoomTeacherId" :options="teacherOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.room')"><UInput v-model="form.room" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.academic_year')"><USelect v-model="form.academicYearId" :options="yearOptions" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('academic.delete_class')" :loading="deleting" @confirm="deleteClass" @cancel="showDelete=false" />
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
const statsLoading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const classes = ref<Record<string, unknown>[]>([])
const showDialog = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const deleteTarget = ref<Record<string, unknown> | null>(null)
const teacherOptions = ref<{ label: string; value: string }[]>([])
const yearOptions = ref<{ label: string; value: string }[]>([])

const stats = reactive({ total: 0, totalStudents: 0, avgClassSize: 0, active: 0 })
const columns: TableColumn[] = [
  { key: 'name', label: 'academic.class_name', sortable: true },
  { key: 'grade', label: 'academic.grade', sortable: true },
  { key: 'program', label: 'academic.program' },
  { key: 'homeRoomTeacherName', label: 'academic.homeroom_teacher' },
  { key: 'studentCount', label: 'academic.students', type: 'number' },
  { key: 'capacity', label: 'academic.capacity', type: 'number' },
  { key: 'room', label: 'academic.room' },
]
const gradeOptions = Array.from({ length: 12 }, (_, i) => ({ label: `${t('academic.grade')} ${i + 1}`, value: String(i + 1) }))
const form = reactive({ name: '', grade: '1', program: '', capacity: 30, homeRoomTeacherId: '', room: '', academicYearId: '' })

const fetchClasses = async () => { loading.value = true; try { classes.value = await api.paginate('/classes').then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get<{total:number;totalStudents:number;avgClassSize:number;active:number}>('/classes/stats')) } catch {} finally { statsLoading.value = false } }
const fetchTeachers = async () => { try { teacherOptions.value = (await api.get<{id:string;fullName:string}[]>('/teachers', { limit: 100 })).map(t => ({ label: t.fullName, value: t.id })) } catch {} }
const fetchYears = async () => { try { yearOptions.value = (await api.get<{id:string;name:string}[]>('/academic-years')).map(y => ({ label: y.name, value: y.id })) } catch {} }

const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { name: '', grade: '1', program: '', capacity: 30, homeRoomTeacherId: '', room: '', academicYearId: '' }); showDialog.value = true }
const editClass = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false; editId.value = null }

const saveClass = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/classes/${editId.value}`, form); toast.add({ title: t('academic.class_updated'), color: 'success' }) }
    else { await api.post('/classes', form); toast.add({ title: t('academic.class_created'), color: 'success' }) }
    closeDialog(); fetchClasses(); fetchStats()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteClass = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/classes/${deleteTarget.value.id}`); toast.add({ title: t('academic.class_deleted'), color: 'success' }); showDelete.value = false; fetchClasses(); fetchStats() } catch {} finally { deleting.value = false } }

onMounted(() => { fetchClasses(); fetchStats(); fetchTeachers(); fetchYears() })
</script>

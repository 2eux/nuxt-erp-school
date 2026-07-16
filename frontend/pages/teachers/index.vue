<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('teachers.title') }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('teachers.subtitle') }}</p>
      </div>
      <UButton v-if="permissions.can('teachers.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAddDialog">
        {{ $t('teachers.add_teacher') }}
      </UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <StatCard :label="$t('teachers.total_teachers')" :value="stats.total" icon="i-heroicons-user-circle" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('teachers.active')" :value="stats.active" icon="i-heroicons-check-circle" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('teachers.certified')" :value="stats.certified" icon="i-heroicons-academic-cap" color="purple" :loading="statsLoading" />
    </div>

    <DataFilter :filter-fields="filterFields" :searchable="true" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="teachers" :loading="loading" :empty-title="$t('teachers.no_teachers')" :show-export="false">
      <template #cell-fullName="{ row }">
        <div class="flex items-center gap-3">
          <UAvatar :src="(row.photo as string) || undefined" size="sm" />
          <div>
            <span class="text-sm font-medium text-gray-900 dark:text-white">{{ row.fullName }}</span>
            <p class="text-xs text-gray-500">{{ row.nip }}</p>
          </div>
        </div>
      </template>
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'active'" /></template>
      <template #cell-subjects="{ row }">
        <div class="flex flex-wrap gap-1">
          <span v-for="s in ((row.subjects as string[]) || [])" :key="s" class="inline-flex px-2 py-0.5 rounded-full bg-brand-50 dark:bg-brand-900/20 text-xs text-brand-700 dark:text-brand-400">{{ s }}</span>
        </div>
      </template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewTeacher(row)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editTeacher(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('teachers.edit_teacher') : $t('teachers.add_teacher')" :loading="saving" @submit="saveTeacher" @cancel="closeDialog">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('teachers.nip')" required><UInput v-model="form.nip" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.full_name')" required class="sm:col-span-2"><UInput v-model="form.fullName" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.gender')"><USelect v-model="form.gender" :options="[{label:t('common.male'),value:'male'},{label:t('common.female'),value:'female'}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.birth_place')"><UInput v-model="form.birthPlace" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.birth_date')"><UInput v-model="form.birthDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.education')"><UInput v-model="form.education" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.major')"><UInput v-model="form.major" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.phone')"><UInput v-model="form.phone" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.email')"><UInput v-model="form.email" type="email" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('teachers.address')" class="sm:col-span-2"><UTextarea v-model="form.address" color="gray" :rows="2" /></UFormGroup>
        <UFormGroup :label="$t('teachers.join_date')"><UInput v-model="form.joinDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.status')"><USelect v-model="form.status" :options="[{label:t('status.active'),value:'active'},{label:t('status.inactive'),value:'inactive'},{label:t('status.resigned'),value:'resigned'}]" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDeleteConfirm" :title="$t('teachers.delete_teacher')" :loading="deleting" @confirm="deleteTeacher" @cancel="showDeleteConfirm=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn, Teacher } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const loading = ref(false)
const statsLoading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const teachers = ref<Record<string, unknown>[]>([])
const editing = ref(false)
const editId = ref<string | null>(null)
const showDialog = ref(false)
const showDeleteConfirm = ref(false)
const deleteTarget = ref<Record<string, unknown> | null>(null)
const stats = reactive({ total: 0, active: 0, certified: 0 })

const columns: TableColumn[] = [
  { key: 'fullName', label: 'teachers.name', sortable: true },
  { key: 'nip', label: 'teachers.nip' },
  { key: 'gender', label: 'common.gender' },
  { key: 'education', label: 'teachers.education' },
  { key: 'subjects', label: 'teachers.subjects' },
  { key: 'phone', label: 'common.phone' },
  { key: 'email', label: 'common.email' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.active'), value: 'active' }, { label: t('status.inactive'), value: 'inactive' }] },
  { key: 'subject', label: 'teachers.subjects', type: 'text' as const },
]

const form = reactive({ nip: '', fullName: '', gender: 'male', birthPlace: '', birthDate: '', education: '', major: '', phone: '', email: '', address: '', joinDate: '', status: 'active' })

const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get<{total:number;active:number;certified:number}>('/teachers/stats')) } catch { /* ignore */ } finally { statsLoading.value = false } }
const fetchTeachers = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { teachers.value = await api.paginate('/teachers', { ...filters }).then(r => r.data) } catch { /* ignore */ } finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchTeachers(filters)

const openAddDialog = () => { editing.value = false; editId.value = null; Object.assign(form, { nip: '', fullName: '', gender: 'male', birthPlace: '', birthDate: '', education: '', major: '', phone: '', email: '', address: '', joinDate: '', status: 'active' }); showDialog.value = true }
const editTeacher = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false; editId.value = null }
const viewTeacher = (row: Record<string, unknown>) => navigateTo(`/teachers/${row.id}`)

const saveTeacher = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/teachers/${editId.value}`, form); toast.add({ title: t('teachers.teacher_updated'), color: 'success' }) }
    else { await api.post('/teachers', form); toast.add({ title: t('teachers.teacher_created'), color: 'success' }) }
    closeDialog(); fetchTeachers(); fetchStats()
  } catch { /* handled by api */ } finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDeleteConfirm.value = true }
const deleteTeacher = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/teachers/${deleteTarget.value.id}`); toast.add({ title: t('teachers.teacher_deleted'), color: 'success' }); showDeleteConfirm.value = false; fetchTeachers(); fetchStats() } catch { /* handled */ } finally { deleting.value = false } }

onMounted(() => { fetchStats(); fetchTeachers() })
</script>

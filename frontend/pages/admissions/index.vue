<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('admissions.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('admissions.subtitle') }}</p></div>
      <UButton v-if="permissions.can('admissions.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('admissions.new_applicant') }}</UButton>
    </div>
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <StatCard :label="$t('admissions.total')" :value="stats.total" icon="i-heroicons-users" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('admissions.pending')" :value="stats.pending" icon="i-heroicons-clock" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('admissions.accepted')" :value="stats.accepted" icon="i-heroicons-check-circle" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('admissions.enrolled')" :value="stats.enrolled" icon="i-heroicons-academic-cap" color="purple" :loading="statsLoading" />
    </div>
    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />
    <DataTable :columns="columns" :rows="applicants" :loading="loading" :empty-title="$t('admissions.no_applicants')" :show-export="false">
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'pending'" /></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton v-if="(row.status as string) === 'pending'" color="emerald" variant="ghost" size="xs" icon="i-heroicons-check" @click="updateStatus(row, 'accepted')" />
          <UButton v-if="(row.status as string) === 'pending'" color="red" variant="ghost" size="xs" icon="i-heroicons-x-mark" @click="updateStatus(row, 'rejected')" />
          <UButton v-if="(row.status as string) === 'accepted'" color="indigo" variant="ghost" size="xs" icon="i-heroicons-academic-cap" @click="enrollStudent(row)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>
    <FormDialog v-model="showForm" :title="$t('admissions.new_applicant')" :loading="saving" @submit="saveApplicant" @cancel="showForm=false">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <UFormGroup :label="$t('students.full_name')" required><UInput v-model="form.fullName" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.gender')"><USelect v-model="form.gender" :options="[{label:t('common.male'),value:'male'},{label:t('common.female'),value:'female'}]" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('students.birth_place')"><UInput v-model="form.birthPlace" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('students.birth_date')"><UInput v-model="form.birthDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('admissions.applied_grade')"><USelect v-model="form.appliedGrade" :options="gradeOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('admissions.previous_school')"><UInput v-model="form.previousSchool" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('students.parent_name')" required class="sm:col-span-2"><UInput v-model="form.parentName" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('students.parent_phone')"><UInput v-model="form.parentPhone" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('students.parent_email')"><UInput v-model="form.parentEmail" type="email" color="gray" /></UFormGroup>
      </div>
    </FormDialog>
    <ConfirmDialog v-model="showDelete" :title="$t('admissions.delete_applicant')" :loading="deleting" @confirm="deleteApplicant" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast()
const loading = ref(false); const statsLoading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const applicants = ref<Record<string, unknown>[]>([]); const deleteTarget = ref<Record<string, unknown> | null>(null)
const stats = reactive({ total: 0, pending: 0, accepted: 0, enrolled: 0 })
const columns: TableColumn[] = [
  { key: 'registrationNumber', label: 'admissions.registration_number' },
  { key: 'fullName', label: 'students.name', sortable: true },
  { key: 'appliedGrade', label: 'admissions.applied_grade' },
  { key: 'parentName', label: 'students.parent_name' },
  { key: 'parentPhone', label: 'students.parent_phone' },
  { key: 'status', label: 'common.status', type: 'status' },
  { key: 'registrationDate', label: 'admissions.registration_date', type: 'date' },
]
const filterFields = [{ key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.pending'), value: 'pending' }, { label: t('status.accepted'), value: 'accepted' }, { label: t('status.rejected'), value: 'rejected' }, { label: t('status.enrolled'), value: 'enrolled' }] }]
const gradeOpts = Array.from({ length: 12 }, (_, i) => ({ label: `${t('academic.grade')} ${i+1}`, value: i + 1 }))
const form = reactive({ fullName: '', gender: 'male', birthPlace: '', birthDate: '', appliedGrade: 1, previousSchool: '', parentName: '', parentPhone: '', parentEmail: '' })

const fetchApplicants = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { applicants.value = await api.paginate('/admissions', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/admissions/stats')) } catch {} finally { statsLoading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchApplicants(filters)
const openAdd = () => { Object.assign(form, { fullName: '', gender: 'male', birthPlace: '', birthDate: '', appliedGrade: 1, previousSchool: '', parentName: '', parentPhone: '', parentEmail: '' }); showForm.value = true }
const saveApplicant = async () => { saving.value = true; try { await api.post('/admissions', form); toast.add({ title: t('admissions.applicant_created'), color: 'success' }); showForm.value = false; fetchApplicants(); fetchStats() } catch {} finally { saving.value = false } }
const updateStatus = async (row: Record<string, unknown>, status: string) => { try { await api.patch(`/admissions/${row.id}/status`, { status }); toast.add({ title: t('admissions.status_updated'), color: 'success' }); fetchApplicants(); fetchStats() } catch {} }
const enrollStudent = async (row: Record<string, unknown>) => { try { await api.post(`/admissions/${row.id}/enroll`); toast.add({ title: t('admissions.student_enrolled'), color: 'success' }); fetchApplicants(); fetchStats() } catch {} }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteApplicant = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/admissions/${deleteTarget.value.id}`); toast.add({ title: t('admissions.applicant_deleted'), color: 'success' }); showDelete.value = false; fetchApplicants(); fetchStats() } catch {} finally { deleting.value = false } }
onMounted(() => { fetchApplicants(); fetchStats() })
</script>

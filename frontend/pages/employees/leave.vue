<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('employees.leave_management') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('employees.leave_subtitle') }}</p></div>
      <UButton v-if="permissions.can('employees.leave.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openRequest">{{ $t('employees.request_leave') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <StatCard :label="$t('employees.pending_requests')" :value="stats.pending" icon="i-heroicons-clock" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('employees.approved')" :value="stats.approved" icon="i-heroicons-check-circle" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('employees.rejected')" :value="stats.rejected" icon="i-heroicons-x-circle" color="red" :loading="statsLoading" />
    </div>

    <DataFilter :filter-fields="filterFields" :searchable="true" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="leaveRequests" :loading="loading" :empty-title="$t('employees.no_leave_requests')" :show-export="false">
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'pending'" /></template>
      <template #cell-dateRange="{ row }">{{ row.startDate }} - {{ row.endDate }}</template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1" v-if="(row.status as string) === 'pending' && permissions.can('employees.leave.approve')">
          <UButton color="emerald" variant="ghost" size="xs" icon="i-heroicons-check" @click="approveRequest(row as Record<string, unknown>)" />
          <UButton color="red" variant="ghost" size="xs" icon="i-heroicons-x-mark" @click="rejectRequest(row as Record<string, unknown>)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showRequestForm" :title="$t('employees.request_leave')" :loading="saving" @submit="submitRequest" @cancel="showRequestForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('employees.leave_type')" required>
          <USelect v-model="leaveForm.type" :options="leaveTypeOptions" color="gray" />
        </UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('common.start_date')" required><UInput v-model="leaveForm.startDate" type="date" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('common.end_date')" required><UInput v-model="leaveForm.endDate" type="date" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('common.reason')" required><UTextarea v-model="leaveForm.reason" color="gray" :rows="3" /></UFormGroup>
        <FileUpload accept=".pdf,.jpg,.jpeg,.png" :multiple="false" accept-hint="Upload supporting documents" @files-selected="handleDocSelected" />
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showApproveConfirm" :title="$t('employees.approve_leave')" variant="info" :loading="processing" @confirm="doApprove" @cancel="showApproveConfirm=false" />
    <ConfirmDialog v-model="showRejectConfirm" :title="$t('employees.reject_leave')" variant="danger" :loading="processing" @confirm="doReject" @cancel="showRejectConfirm=false" />
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
const processing = ref(false)
const showRequestForm = ref(false)
const showApproveConfirm = ref(false)
const showRejectConfirm = ref(false)
const leaveRequests = ref<Record<string, unknown>[]>([])
const actionTarget = ref<Record<string, unknown> | null>(null)
const docFile = ref<File | null>(null)

const stats = reactive({ pending: 0, approved: 0, rejected: 0 })

const columns: TableColumn[] = [
  { key: 'employeeName', label: 'employees.name', sortable: true },
  { key: 'type', label: 'employees.leave_type' },
  { key: 'dateRange', label: 'employees.leave_period' },
  { key: 'days', label: 'employees.days', type: 'number' },
  { key: 'reason', label: 'common.reason' },
  { key: 'status', label: 'common.status', type: 'status' },
  { key: 'createdAt', label: 'common.created_at', type: 'date' },
]
const filterFields = [
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.pending'), value: 'pending' }, { label: t('status.approved'), value: 'approved' }, { label: t('status.rejected'), value: 'rejected' }] },
  { key: 'type', label: 'employees.leave_type', type: 'select' as const, options: [{ label: t('employees.annual'), value: 'annual' }, { label: t('employees.sick'), value: 'sick' }, { label: t('employees.maternity'), value: 'maternity' }, { label: t('employees.personal'), value: 'personal' }] },
]
const leaveTypeOptions = [
  { label: t('employees.annual_leave'), value: 'annual' }, { label: t('employees.sick_leave'), value: 'sick' },
  { label: t('employees.maternity_leave'), value: 'maternity' }, { label: t('employees.personal_leave'), value: 'personal' },
  { label: t('employees.marriage_leave'), value: 'marriage' }, { label: t('employees.bereavement_leave'), value: 'bereavement' },
]

const leaveForm = reactive({ type: 'annual', startDate: '', endDate: '', reason: '' })

const handleDocSelected = (files: File[]) => { if (files.length > 0) docFile.value = files[0] }

const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get<{pending:number;approved:number;rejected:number}>('/employees/leave/stats')) } catch {} finally { statsLoading.value = false } }
const fetchLeaveRequests = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { leaveRequests.value = await api.paginate('/employees/leave', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchLeaveRequests(filters)

const openRequest = () => { leaveForm.type = 'annual'; leaveForm.startDate = ''; leaveForm.endDate = ''; leaveForm.reason = ''; showRequestForm.value = true }

const submitRequest = async () => {
  saving.value = true
  try {
    const formData = new FormData(); formData.append('type', leaveForm.type); formData.append('startDate', leaveForm.startDate); formData.append('endDate', leaveForm.endDate); formData.append('reason', leaveForm.reason)
    if (docFile.value) formData.append('document', docFile.value)
    await api.upload('/employees/leave', formData)
    toast.add({ title: t('employees.leave_requested'), color: 'success' }); showRequestForm.value = false; fetchLeaveRequests(); fetchStats()
  } catch {} finally { saving.value = false }
}

const approveRequest = (row: Record<string, unknown>) => { actionTarget.value = row; showApproveConfirm.value = true }
const rejectRequest = (row: Record<string, unknown>) => { actionTarget.value = row; showRejectConfirm.value = true }

const doApprove = async () => { if (!actionTarget.value) return; processing.value = true; try { await api.patch(`/employees/leave/${actionTarget.value.id}/approve`); toast.add({ title: t('employees.leave_approved'), color: 'success' }); showApproveConfirm.value = false; fetchLeaveRequests(); fetchStats() } catch {} finally { processing.value = false } }
const doReject = async () => { if (!actionTarget.value) return; processing.value = true; try { await api.patch(`/employees/leave/${actionTarget.value.id}/reject`); toast.add({ title: t('employees.leave_rejected'), color: 'info' }); showRejectConfirm.value = false; fetchLeaveRequests(); fetchStats() } catch {} finally { processing.value = false } }

onMounted(() => { fetchStats(); fetchLeaveRequests() })
</script>

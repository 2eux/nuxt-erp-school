<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('employees.attendance') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('employees.attendance_subtitle') }}</p></div>
      <div class="flex items-center gap-3">
        <UInput v-model="startDate" type="date" color="gray" size="sm" @change="fetchAttendance" />
        <span class="text-gray-400 text-sm">-</span>
        <UInput v-model="endDate" type="date" color="gray" size="sm" @change="fetchAttendance" />
        <UButton color="primary" size="sm" icon="i-heroicons-plus" @click="openBulkEntry">{{ $t('attendance.bulk_entry') }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <StatCard :label="$t('attendance.present')" :value="summary.present" icon="i-heroicons-check-circle" color="emerald" :loading="loading" />
      <StatCard :label="$t('attendance.late')" :value="summary.late" icon="i-heroicons-clock" color="amber" :loading="loading" />
      <StatCard :label="$t('attendance.absent')" :value="summary.absent" icon="i-heroicons-x-circle" color="red" :loading="loading" />
      <StatCard :label="$t('attendance.on_leave')" :value="summary.onLeave" icon="i-heroicons-calendar" color="purple" :loading="loading" />
    </div>

    <DataTable :columns="columns" :rows="records" :loading="loading" :empty-title="$t('attendance.no_records')" :show-filters="false">
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'present'" /></template>
      <template #cell-checkIn="{ row }">
        <span v-if="!row.checkIn" class="text-sm text-amber-600 dark:text-amber-400">{{ $t('attendance.not_checked_in') }}</span>
        <span v-else class="text-sm text-gray-900 dark:text-white">{{ row.checkIn }}</span>
      </template>
      <template #cell-checkOut="{ row }">
        <span v-if="!row.checkOut" class="text-sm text-gray-400">-</span>
        <span v-else class="text-sm text-gray-900 dark:text-white">{{ row.checkOut }}</span>
      </template>
    </DataTable>

    <FormDialog v-model="showBulkDialog" :title="$t('attendance.bulk_entry')" :loading="bulkSaving" @submit="submitBulkEntry" @cancel="showBulkDialog=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('common.date')" required><UInput v-model="bulkDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('attendance.status')"><USelect v-model="bulkStatus" :options="attendanceStatusOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.department')"><USelect v-model="bulkDepartment" :options="departmentOptions" color="gray" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const toast = useToast()
const { $dayjs } = useNuxtApp()
const loading = ref(false)
const bulkSaving = ref(false)
const showBulkDialog = ref(false)
const records = ref<Record<string, unknown>[]>([])
const startDate = ref($dayjs().startOf('month').format('YYYY-MM-DD'))
const endDate = ref($dayjs().format('YYYY-MM-DD'))
const bulkDate = ref($dayjs().format('YYYY-MM-DD'))
const bulkStatus = ref('present')
const bulkDepartment = ref('')

const summary = reactive({ present: 0, late: 0, absent: 0, onLeave: 0 })

const columns: TableColumn[] = [
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'employeeName', label: 'employees.name' },
  { key: 'nip', label: 'employees.nip' },
  { key: 'department', label: 'employees.department' },
  { key: 'checkIn', label: 'attendance.check_in' },
  { key: 'checkOut', label: 'attendance.check_out' },
  { key: 'status', label: 'common.status', type: 'status' },
  { key: 'note', label: 'common.note' },
]
const attendanceStatusOptions = [
  { label: t('attendance.present'), value: 'present' },
  { label: t('attendance.absent'), value: 'absent' },
  { label: t('attendance.late'), value: 'late' },
  { label: t('attendance.permission'), value: 'permission' },
  { label: t('attendance.sick'), value: 'sick' },
]
const departmentOptions = [
  { label: 'Administrasi', value: 'Administrasi' }, { label: 'Keuangan', value: 'Keuangan' },
  { label: 'Kesiswaan', value: 'Kesiswaan' }, { label: 'TU', value: 'TU' },
  { label: 'Kebersihan', value: 'Kebersihan' }, { label: 'Keamanan', value: 'Keamanan' },
  { label: 'Dapur', value: 'Dapur' }, { label: 'Perpustakaan', value: 'Perpustakaan' },
]

const fetchAttendance = async () => {
  loading.value = true
  try {
    const [recordsRes, summaryRes] = await Promise.all([
      api.get('/employees/attendance', { startDate: startDate.value, endDate: endDate.value }),
      api.get<{ present: number; late: number; absent: number; onLeave: number }>('/employees/attendance/summary', { startDate: startDate.value, endDate: endDate.value }),
    ])
    records.value = recordsRes; Object.assign(summary, summaryRes)
  } catch {/* ignore */}
  finally { loading.value = false }
}

const openBulkEntry = () => { bulkDate.value = $dayjs().format('YYYY-MM-DD'); bulkStatus.value = 'present'; bulkDepartment.value = ''; showBulkDialog.value = true }

const submitBulkEntry = async () => {
  bulkSaving.value = true
  try {
    await api.post('/employees/attendance/bulk', { date: bulkDate.value, status: bulkStatus.value, department: bulkDepartment.value || undefined })
    toast.add({ title: t('attendance.bulk_saved'), color: 'success' })
    showBulkDialog.value = false; fetchAttendance()
  } catch {/* handled */} finally { bulkSaving.value = false }
}

onMounted(() => fetchAttendance())
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('attendance.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('attendance.subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedClass" :options="classOptions" :placeholder="$t('academic.select_class')" color="gray" size="sm" class="w-36" @change="fetchAttendance" />
        <UInput v-model="selectedDate" type="date" color="gray" size="sm" @change="fetchAttendance" />
        <UButton color="primary" size="sm" icon="i-heroicons-plus" @click="openBulkEntry">{{ $t('attendance.bulk_entry') }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-2 sm:grid-cols-5 gap-4">
      <StatCard :label="$t('attendance.present')" :value="stats.present" icon="i-heroicons-check-circle" color="emerald" :loading="loading" />
      <StatCard :label="$t('attendance.late')" :value="stats.late" icon="i-heroicons-clock" color="amber" :loading="loading" />
      <StatCard :label="$t('attendance.absent')" :value="stats.absent" icon="i-heroicons-x-circle" color="red" :loading="loading" />
      <StatCard :label="$t('attendance.sick')" :value="stats.sick" icon="i-heroicons-heart" color="purple" :loading="loading" />
      <StatCard :label="$t('attendance.rate')" :value="`${attendanceRate}%`" icon="i-heroicons-chart-bar" color="blue" :loading="loading" />
    </div>

    <DataTable :columns="columns" :rows="attendanceRecords" :loading="loading" :empty-title="$t('attendance.no_records')" :show-export="true" @export="handleExport">
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'present'" /></template>
      <template #cell-studentName="{ row }">
        <span class="text-sm font-medium text-gray-900 dark:text-white">{{ row.studentName }}</span>
        <p class="text-xs text-gray-500">{{ row.className }}</p>
      </template>
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editStatus(row as Record<string, unknown>)" />
      </template>
    </DataTable>

    <FormDialog v-model="showBulkDialog" :title="$t('attendance.bulk_entry')" :loading="bulkSaving" @submit="submitBulkEntry" @cancel="showBulkDialog=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('common.date')" required><UInput v-model="bulkDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.class')" required><USelect v-model="bulkClass" :options="classOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('attendance.status')" required><USelect v-model="bulkStatus" :options="statusOptions" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <FormDialog v-model="showEditStatus" :title="$t('attendance.edit_status')" :loading="saving" @submit="updateStatus" @cancel="showEditStatus=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('attendance.status')" required><USelect v-model="editForm.status" :options="statusOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.note')"><UTextarea v-model="editForm.note" color="gray" :rows="2" /></UFormGroup>
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
const saving = ref(false)
const bulkSaving = ref(false)
const showBulkDialog = ref(false)
const showEditStatus = ref(false)
const attendanceRecords = ref<Record<string, unknown>[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const selectedClass = ref('')
const selectedDate = ref($dayjs().format('YYYY-MM-DD'))
const bulkDate = ref($dayjs().format('YYYY-MM-DD'))
const bulkClass = ref('')
const bulkStatus = ref('present')

const stats = reactive({ present: 0, late: 0, absent: 0, sick: 0, permission: 0, total: 0 })
const attendanceRate = computed(() => {
  if (stats.total === 0) return 0
  return Math.round(((stats.present + stats.late) / stats.total) * 100)
})

const editTarget = ref<Record<string, unknown> | null>(null)
const editForm = reactive({ status: 'present', note: '' })

const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name', sortable: true },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'status', label: 'common.status', type: 'status' },
  { key: 'checkInTime', label: 'attendance.check_in' },
  { key: 'checkOutTime', label: 'attendance.check_out' },
  { key: 'note', label: 'common.note' },
]
const statusOptions = [
  { label: t('attendance.present'), value: 'present' }, { label: t('attendance.absent'), value: 'absent' },
  { label: t('attendance.late'), value: 'late' }, { label: t('attendance.sick'), value: 'sick' },
  { label: t('attendance.permission'), value: 'permission' }, { label: t('attendance.holiday'), value: 'holiday' },
]

const fetchAttendance = async () => {
  if (!selectedClass.value) return
  loading.value = true
  try {
    const [records, statsRes] = await Promise.all([
      api.get('/attendance', { classId: selectedClass.value, date: selectedDate.value }),
      api.get('/attendance/stats', { classId: selectedClass.value, date: selectedDate.value }),
    ])
    attendanceRecords.value = records; Object.assign(stats, statsRes)
  } catch {} finally { loading.value = false }
}

const openBulkEntry = () => { bulkDate.value = selectedDate.value; bulkClass.value = selectedClass.value; bulkStatus.value = 'present'; showBulkDialog.value = true }

const submitBulkEntry = async () => {
  bulkSaving.value = true
  try { await api.post('/attendance/bulk', { date: bulkDate.value, classId: bulkClass.value, status: bulkStatus.value }); toast.add({ title: t('attendance.bulk_saved'), color: 'success' }); showBulkDialog.value = false; fetchAttendance() }
  catch {} finally { bulkSaving.value = false }
}

const editStatus = (row: Record<string, unknown>) => { editTarget.value = row; editForm.status = (row.status as string) || 'present'; editForm.note = (row.note as string) || ''; showEditStatus.value = true }

const updateStatus = async () => {
  if (!editTarget.value) return; saving.value = true
  try { await api.patch(`/attendance/${editTarget.value.id}`, editForm); toast.add({ title: t('attendance.status_updated'), color: 'success' }); showEditStatus.value = false; fetchAttendance() }
  catch {} finally { saving.value = false }
}

const handleExport = async (format: string) => {
  try {
    const res = await api.get(`/attendance/export?classId=${selectedClass.value}&date=${selectedDate.value}&format=${format}`)
    if (typeof res === 'object' && (res as Record<string, unknown>).url) {
      window.open((res as Record<string, unknown>).url as string, '_blank')
    }
  } catch {}
}

const fetchClasses = async () => {
  try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); if (classOptions.value.length > 0) { selectedClass.value = classOptions.value[0].value; fetchAttendance() } } catch {}
}

onMounted(() => fetchClasses())
</script>

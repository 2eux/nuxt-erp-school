<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('employees.payroll') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('employees.payroll_subtitle') }}</p></div>
      <UButton v-if="permissions.can('employees.payroll.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openGeneratePeriod">{{ $t('employees.generate_payroll') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('employees.current_period')" :value="currentPeriod" icon="i-heroicons-calendar" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('employees.total_salary')" :value="formatCurrency(stats.totalSalary)" icon="i-heroicons-currency-dollar" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('employees.processed')" :value="stats.processed" icon="i-heroicons-check-circle" color="amber" :loading="statsLoading" :suffix="`/${stats.totalEmployees}`" />
      <StatCard :label="$t('employees.paid')" :value="stats.paid" icon="i-heroicons-banknotes" color="purple" :loading="statsLoading" :suffix="`/${stats.totalEmployees}`" />
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="payrollRecords" :loading="loading" :empty-title="$t('employees.no_payroll')" :show-export="false" :selectable="true" :selected-ids="selectedIds" @update:selected-ids="selectedIds = $event">
      <template #cell-netSalary="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.netSalary as number) }}</span></template>
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'draft'" /></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-document-text" @click="viewSlip(row as Record<string, unknown>)" />
          <UButton v-if="(row.status as string) === 'draft'" color="gray" variant="ghost" size="xs" icon="i-heroicons-check-circle" @click="approvePayroll(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showPeriodForm" :title="$t('employees.generate_payroll')" :loading="generating" @submit="generatePayroll" @cancel="showPeriodForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('employees.month')" required><USelect v-model="periodForm.month" :options="monthOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.year')" required><UInput v-model.number="periodForm.year" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.department')"><USelect v-model="periodForm.department" :options="deptOptions" color="gray" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const { $dayjs } = useNuxtApp()
const loading = ref(false)
const statsLoading = ref(false)
const generating = ref(false)
const showPeriodForm = ref(false)
const payrollRecords = ref<Record<string, unknown>[]>([])
const selectedIds = ref<string[]>([])

const stats = reactive({ totalSalary: 0, processed: 0, paid: 0, totalEmployees: 0 })
const currentPeriod = computed(() => {
  const now = $dayjs()
  return `${now.format('MMMM')} ${now.format('YYYY')}`
})

const columns: TableColumn[] = [
  { key: 'employeeName', label: 'employees.name', sortable: true },
  { key: 'nip', label: 'employees.nip' },
  { key: 'month', label: 'employees.month' },
  { key: 'year', label: 'employees.year' },
  { key: 'baseSalary', label: 'employees.base_salary', type: 'currency' },
  { key: 'allowances', label: 'employees.allowances', type: 'currency' },
  { key: 'deductions', label: 'employees.deductions', type: 'currency' },
  { key: 'netSalary', label: 'employees.net_salary', type: 'currency' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.draft'), value: 'draft' }, { label: t('status.approved'), value: 'approved' }, { label: t('status.paid'), value: 'paid' }] },
  { key: 'month', label: 'employees.month', type: 'select' as const, options: Array.from({ length: 12 }, (_, i) => ({ label: $dayjs().month(i).format('MMMM'), value: String(i + 1) })) },
]
const monthOptions = Array.from({ length: 12 }, (_, i) => ({ label: $dayjs().month(i).format('MMMM'), value: String(i + 1) }))
const deptOptions = [{ label: 'Semua', value: '' }, { label: 'Guru', value: 'teacher' }, { label: 'Staf', value: 'staff' }]

const periodForm = reactive({ month: String($dayjs().month() + 1), year: $dayjs().year(), department: '' })

const formatCurrency = (value: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value)

const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get<{totalSalary:number;processed:number;paid:number;totalEmployees:number}>('/employees/payroll/stats')) } catch {} finally { statsLoading.value = false } }
const fetchPayroll = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { payrollRecords.value = await api.paginate('/employees/payroll', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchPayroll(filters)

const openGeneratePeriod = () => { periodForm.month = String($dayjs().month() + 1); periodForm.year = $dayjs().year(); periodForm.department = ''; showPeriodForm.value = true }

const generatePayroll = async () => {
  generating.value = true
  try {
    await api.post('/employees/payroll/generate', periodForm)
    toast.add({ title: t('employees.payroll_generated'), color: 'success' })
    showPeriodForm.value = false; fetchPayroll(); fetchStats()
  } catch {} finally { generating.value = false }
}

const approvePayroll = async (row: Record<string, unknown>) => {
  try {
    await api.patch(`/employees/payroll/${row.id}/approve`)
    toast.add({ title: t('employees.payroll_approved'), color: 'success' }); fetchPayroll(); fetchStats()
  } catch {}
}

const viewSlip = (row: Record<string, unknown>) => {
  window.open(`/api/v1/employees/payroll/${row.id}/slip-pdf`, '_blank')
}

onMounted(() => { fetchStats(); fetchPayroll() })
</script>

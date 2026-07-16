<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.payroll') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.payroll_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <UButton v-if="permissions.can('finance.payroll.process')" color="primary" size="sm" icon="i-heroicons-calculator" @click="openProcess">{{ $t('finance.process_payroll') }}</UButton>
        <UButton color="gray" variant="outline" size="sm" icon="i-heroicons-arrow-down-tray" @click="exportPayroll">{{ $t('common.export') }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('finance.total_payroll')" :value="formatCurrency(stats.total)" icon="i-heroicons-currency-dollar" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('finance.bpjs')" :value="formatCurrency(stats.bpjs)" icon="i-heroicons-heart" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('finance.pph21')" :value="formatCurrency(stats.pph21)" icon="i-heroicons-document-text" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('finance.net_total')" :value="formatCurrency(stats.netTotal)" icon="i-heroicons-banknotes" color="purple" :loading="statsLoading" />
    </div>

    <div class="card">
      <div class="card-header"><h3 class="card-title">{{ $t('finance.payroll_periods') }}</h3></div>
      <DataTable :columns="columns" :rows="periods" :loading="loading" :empty-title="$t('finance.no_payroll_periods')" :show-export="false">
        <template #cell-totalSalary="{ row }"><span class="text-sm font-mono">{{ formatCurrency(row.totalSalary as number) }}</span></template>
        <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'draft'" /></template>
        <template #item-actions="{ row }">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewDetail(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-printer" @click="printSlip(row as Record<string, unknown>)" />
        </template>
      </DataTable>
    </div>

    <FormDialog v-model="showProcess" :title="$t('finance.process_payroll')" :loading="processing" @submit="doProcess" @cancel="showProcess=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('employees.month')"><USelect v-model="processForm.month" :options="monthOpts" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('employees.year')"><UInput v-model.number="processForm.year" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('finance.include_bpjs')"><UToggle v-model="processForm.includeBpjs" /></UFormGroup>
        <UFormGroup :label="$t('finance.include_pph21')"><UToggle v-model="processForm.includePph21" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const statsLoading = ref(false); const processing = ref(false)
const showProcess = ref(false)
const periods = ref<Record<string, unknown>[]>([])
const stats = reactive({ total: 0, bpjs: 0, pph21: 0, netTotal: 0 })
const columns: TableColumn[] = [
  { key: 'month', label: 'employees.month' },
  { key: 'year', label: 'employees.year' },
  { key: 'totalEmployees', label: 'employees.total_employees', type: 'number' },
  { key: 'totalSalary', label: 'finance.total_salary', type: 'currency' },
  { key: 'status', label: 'common.status', type: 'status' },
  { key: 'processedAt', label: 'finance.processed_at', type: 'date' },
]
const processForm = reactive({ month: String($dayjs().month() + 1), year: $dayjs().year(), includeBpjs: true, includePph21: true })
const monthOpts = Array.from({ length: 12 }, (_, i) => ({ label: $dayjs().month(i).format('MMMM'), value: String(i + 1) }))
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

const fetchData = async () => {
  loading.value = true; statsLoading.value = true
  try { periods.value = await api.paginate('/finance/payroll').then(r => r.data) } catch {} finally { loading.value = false }
  try { Object.assign(stats, await api.get('/finance/payroll/stats')) } catch {} finally { statsLoading.value = false }
}
const openProcess = () => { processForm.month = String($dayjs().month() + 1); processForm.year = $dayjs().year(); showProcess.value = true }
const doProcess = async () => { processing.value = true; try { await api.post('/finance/payroll/process', processForm); toast.add({ title: t('finance.payroll_processed'), color: 'success' }); showProcess.value = false; fetchData() } catch {} finally { processing.value = false } }
const viewDetail = (row: Record<string, unknown>) => { navigateTo(`/finance/payroll/${row.id}`) }
const printSlip = (row: Record<string, unknown>) => { window.open(`/api/v1/finance/payroll/${row.id}/slips`, '_blank') }
const exportPayroll = () => { window.open('/api/v1/finance/payroll/export', '_blank') }
onMounted(() => fetchData())
</script>

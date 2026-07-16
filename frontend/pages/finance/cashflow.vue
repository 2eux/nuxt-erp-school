<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.cashflow') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.cashflow_subtitle') }}</p></div>
      <UButton v-if="permissions.can('finance.cashflow.create')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('finance.add_transaction') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <StatCard :label="$t('finance.income')" :value="formatCurrency(stats.income)" icon="i-heroicons-arrow-trending-up" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('finance.expense')" :value="formatCurrency(stats.expense)" icon="i-heroicons-arrow-trending-down" color="red" :loading="statsLoading" />
      <StatCard :label="$t('finance.net')" :value="formatCurrency(stats.net)" icon="i-heroicons-scale" color="blue" :loading="statsLoading" />
    </div>

    <div class="card">
      <div class="card-header"><h3 class="card-title">{{ $t('finance.monthly_summary') }}</h3></div>
      <div class="h-72"><ApexChart type="bar" height="100%" :options="cashflowChart.options" :series="cashflowChart.series" /></div>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="transactions" :loading="loading" :empty-title="$t('finance.no_transactions')" :show-export="true" @export="handleExport">
      <template #cell-amount="{ row }">
        <span class="text-sm font-mono" :class="(row.type as string) === 'in' ? 'text-emerald-600' : 'text-red-600'">
          {{ (row.type as string) === 'in' ? '+' : '-' }}{{ formatCurrency(row.amount as number) }}
        </span>
      </template>
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editTransaction(row as Record<string, unknown>)" />
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('finance.edit_transaction') : $t('finance.add_transaction')" :loading="saving" @submit="saveTransaction" @cancel="showForm=false">
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('finance.type')" required><USelect v-model="form.type" :options="[{label:t('finance.cash_in'),value:'in'},{label:t('finance.cash_out'),value:'out'}]" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('finance.amount')" required><UInput v-model.number="form.amount" type="number" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('finance.category')"><UInput v-model="form.category" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.date')"><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.description')" required><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('finance.delete_transaction')" :loading="deleting" @confirm="deleteTransaction" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const statsLoading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const transactions = ref<Record<string, unknown>[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)

const stats = reactive({ income: 0, expense: 0, net: 0 })
const columns: TableColumn[] = [
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'type', label: 'finance.type' },
  { key: 'category', label: 'finance.category' },
  { key: 'amount', label: 'finance.amount', type: 'currency' },
  { key: 'description', label: 'common.description' },
]
const filterFields = [
  { key: 'type', label: 'finance.type', type: 'select' as const, options: [{ label: t('finance.cash_in'), value: 'in' }, { label: t('finance.cash_out'), value: 'out' }] },
]
const form = reactive({ type: 'in', amount: 0, category: '', date: $dayjs().format('YYYY-MM-DD'), description: '' })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

const cashflowChart = computed(() => ({
  series: [{ name: t('finance.income'), data: [45, 52, 48, 58, 52, 62, 55, 68, 60, 55, 65, 58] }, { name: t('finance.expense'), data: [30, 35, 32, 40, 38, 42, 40, 45, 42, 38, 48, 42] }],
  options: { chart: { type: 'bar' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669', '#ef4444'], plotOptions: { bar: { borderRadius: 4 } }, xaxis: { categories: ['Jul','Aug','Sep','Oct','Nov','Dec','Jan','Feb','Mar','Apr','May','Jun'] }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))

const fetchTransactions = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { transactions.value = await api.paginate('/finance/cashflow', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/finance/cashflow/stats')) } catch {} finally { statsLoading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchTransactions(filters)

const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { type: 'in', amount: 0, category: '', date: $dayjs().format('YYYY-MM-DD'), description: '' }); showForm.value = true }
const editTransaction = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveTransaction = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/finance/cashflow/${editId.value}`, form); toast.add({ title: t('finance.transaction_updated'), color: 'success' }) } else { await api.post('/finance/cashflow', form); toast.add({ title: t('finance.transaction_created'), color: 'success' }) } showForm.value = false; fetchTransactions(); fetchStats() } catch {} finally { saving.value = false } }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteTransaction = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/finance/cashflow/${deleteTarget.value.id}`); toast.add({ title: t('finance.transaction_deleted'), color: 'success' }); showDelete.value = false; fetchTransactions(); fetchStats() } catch {} finally { deleting.value = false } }
const handleExport = (format: string) => { window.open(`/api/v1/finance/cashflow/export?format=${format}`, '_blank') }
onMounted(() => { fetchTransactions(); fetchStats() })
</script>

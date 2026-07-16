<template>
  <div class="space-y-6">
    <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.subtitle') }}</p></div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('finance.total_receivables')" :value="formatCurrency(stats.totalReceivables)" icon="i-heroicons-currency-dollar" color="blue" :loading="loading" />
      <StatCard :label="$t('finance.collected')" :value="formatCurrency(stats.collected)" icon="i-heroicons-banknotes" color="emerald" :loading="loading" />
      <StatCard :label="$t('finance.overdue')" :value="formatCurrency(stats.overdue)" icon="i-heroicons-exclamation-triangle" color="red" :loading="loading" />
      <StatCard :label="$t('finance.collection_rate')" :value="`${stats.collectionRate}%`" icon="i-heroicons-chart-bar" color="amber" :loading="loading" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('finance.revenue_trend') }}</h3></div>
        <div class="h-72"><ApexChart type="bar" height="100%" :options="revenueChart.options" :series="revenueChart.series" /></div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('finance.quick_links') }}</h3></div>
        <div class="grid grid-cols-2 gap-3">
          <NuxtLink v-for="link in quickLinks" :key="link.to" :to="link.to" class="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
            <UIcon :name="link.icon" class="w-5 h-5 text-brand-600 dark:text-brand-400" />
            <span class="text-sm font-medium text-gray-900 dark:text-white">{{ link.label }}</span>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi()
const loading = ref(false)
const stats = reactive({ totalReceivables: 0, collected: 0, overdue: 0, collectionRate: 0 })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(v)
const quickLinks = [
  { to: '/finance/invoices', label: t('finance.invoices'), icon: 'i-heroicons-document-text' },
  { to: '/finance/payments', label: t('finance.payments'), icon: 'i-heroicons-banknotes' },
  { to: '/finance/fee-types', label: t('finance.fee_types'), icon: 'i-heroicons-tag' },
  { to: '/finance/journals', label: t('finance.journals'), icon: 'i-heroicons-book-open' },
  { to: '/finance/ledger', label: t('finance.ledger'), icon: 'i-heroicons-table-cells' },
  { to: '/finance/budget', label: t('finance.budget'), icon: 'i-heroicons-calculator' },
  { to: '/finance/payroll', label: t('employees.payroll'), icon: 'i-heroicons-user-group' },
  { to: '/finance/cashflow', label: t('finance.cashflow'), icon: 'i-heroicons-arrows-right-left' },
]
const revenueChart = computed(() => ({
  series: [{ name: t('finance.collected'), data: [45, 52, 48, 58, 52, 62, 55, 68, 60, 55, 65, 58] }],
  options: { chart: { type: 'bar' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669'], plotOptions: { bar: { borderRadius: 4 } }, xaxis: { categories: ['Jul','Aug','Sep','Oct','Nov','Dec','Jan','Feb','Mar','Apr','May','Jun'] }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))
onMounted(async () => { loading.value = true; try { Object.assign(stats, await api.get('/finance/dashboard-stats')) } catch {} finally { loading.value = false } })
</script>

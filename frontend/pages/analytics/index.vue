<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('analytics.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('analytics.subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <UInput v-model="dateRange.start" type="date" color="gray" size="sm" />
        <span class="text-gray-400">-</span>
        <UInput v-model="dateRange.end" type="date" color="gray" size="sm" />
        <UButton color="gray" variant="outline" size="sm" icon="i-heroicons-arrow-down-tray" @click="exportReport">{{ $t('common.export') }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('analytics.enrollment_rate')" :value="`${stats.enrollmentRate}%`" icon="i-heroicons-arrow-trending-up" color="emerald" :loading="loading" :trend="stats.enrollmentTrend" />
      <StatCard :label="$t('analytics.attendance_rate')" :value="`${stats.attendanceRate}%`" icon="i-heroicons-check-badge" color="blue" :loading="loading" />
      <StatCard :label="$t('analytics.revenue')" :value="formatCurrency(stats.revenue)" icon="i-heroicons-currency-dollar" color="amber" :loading="loading" :trend="stats.revenueTrend" />
      <StatCard :label="$t('analytics.tahfidz_progress')" :value="`${stats.tahfidzRate}%`" icon="i-heroicons-book-open" color="purple" :loading="loading" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('analytics.enrollment_trend') }}</h3></div>
        <div class="h-80"><ApexChart type="line" height="100%" :options="enrollmentChart.options" :series="enrollmentChart.series" /></div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('analytics.revenue_trend') }}</h3></div>
        <div class="h-80"><ApexChart type="bar" height="100%" :options="revenueChart.options" :series="revenueChart.series" /></div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('analytics.attendance_by_class') }}</h3></div>
        <div class="h-80"><ApexChart type="bar" height="100%" :options="attendanceChart.options" :series="attendanceChart.series" /></div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('analytics.student_performance') }}</h3></div>
        <div class="h-80"><ApexChart type="radialBar" height="100%" :options="performanceChart.options" :series="performanceChart.series" /></div>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div class="card"><div class="card-header"><h3 class="card-title">{{ $t('analytics.prayer_compliance') }}</h3></div>
        <div class="h-64"><ApexChart type="donut" height="100%" :options="prayerChart.options" :series="prayerChart.series" /></div></div>
      <div class="card"><div class="card-header"><h3 class="card-title">{{ $t('analytics.fee_collection') }}</h3></div>
        <div class="h-64"><ApexChart type="donut" height="100%" :options="feeChart.options" :series="feeChart.series" /></div></div>
      <div class="card"><div class="card-header"><h3 class="card-title">{{ $t('analytics.student_gender') }}</h3></div>
        <div class="h-64"><ApexChart type="donut" height="100%" :options="genderChart.options" :series="genderChart.series" /></div></div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const { $dayjs } = useNuxtApp()
const loading = ref(false)
const stats = reactive({ enrollmentRate: 0, attendanceRate: 0, revenue: 0, tahfidzRate: 0, enrollmentTrend: 0, revenueTrend: 0 })
const dateRange = reactive({ start: $dayjs().subtract(1, 'year').format('YYYY-MM-DD'), end: $dayjs().format('YYYY-MM-DD') })

const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', notation: 'compact', maximumFractionDigits: 1 }).format(v)

const enrollmentChart = computed(() => ({
  series: [{ name: t('analytics.students'), data: [320, 340, 355, 370, 385, 400, 415, 430, 445, 460, 475, 490] }],
  options: { chart: { type: 'line' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669'], stroke: { curve: 'smooth' as const, width: 2 }, xaxis: { categories: ['Jul','Aug','Sep','Oct','Nov','Dec','Jan','Feb','Mar','Apr','May','Jun'] }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))
const revenueChart = computed(() => ({
  series: [{ name: t('analytics.collected'), data: [45, 52, 48, 58, 52, 62, 55, 68, 60, 55, 65, 58] }, { name: t('analytics.expected'), data: [50, 55, 52, 62, 58, 65, 60, 70, 65, 60, 68, 62] }],
  options: { chart: { type: 'bar' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669', '#6366f1'], plotOptions: { bar: { borderRadius: 4 } }, xaxis: { categories: ['Jul','Aug','Sep','Oct','Nov','Dec','Jan','Feb','Mar','Apr','May','Jun'] }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))
const attendanceChart = computed(() => ({
  series: [{ name: t('attendance.present'), data: [92, 88, 95, 90, 87, 93] }],
  options: { chart: { type: 'bar' as const, background: 'transparent' }, colors: ['#059669'], xaxis: { categories: ['1A','1B','2A','2B','3A','3B'] }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))
const performanceChart = computed(() => ({
  series: [78],
  options: { chart: { type: 'radialBar' as const, background: 'transparent' }, colors: ['#059669'], plotOptions: { radialBar: { dataLabels: { name: { show: true }, value: { show: true } } } }, labels: [t('analytics.average_gpa')] },
}))
const prayerChart = computed(() => ({ series: [78, 12, 7, 3], options: { chart: { type: 'donut' as const, background: 'transparent' }, labels: [t('attendance.ontime'),t('attendance.late'),t('attendance.absent'),t('attendance.excused')], colors: ['#059669','#f59e0b','#ef4444','#6366f1'] } }))
const feeChart = computed(() => ({ series: [75, 15, 10], options: { chart: { type: 'donut' as const, background: 'transparent' }, labels: [t('finance.paid'), t('finance.partial'), t('finance.unpaid')], colors: ['#059669','#f59e0b','#ef4444'] } }))
const genderChart = computed(() => ({ series: [52, 48], options: { chart: { type: 'donut' as const, background: 'transparent' }, labels: [t('common.male'), t('common.female')], colors: ['#6366f1','#ec4899'] } }))

const exportReport = () => { window.open(`/api/v1/analytics/export?start=${dateRange.start}&end=${dateRange.end}`, '_blank') }
const fetchData = async () => { loading.value = true; try { Object.assign(stats, await api.get('/analytics/stats', { startDate: dateRange.start, endDate: dateRange.end })) } catch {} finally { loading.value = false } }
onMounted(() => fetchData())
</script>

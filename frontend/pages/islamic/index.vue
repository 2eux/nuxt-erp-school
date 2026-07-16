<template>
  <div class="space-y-6">
    <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('islamic.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('islamic.subtitle') }}</p></div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('islamic.total_hafiz')" :value="stats.totalHafiz" icon="i-heroicons-book-open" color="emerald" :loading="loading" />
      <StatCard :label="$t('islamic.avg_prayer')" :value="`${stats.avgPrayer}%`" icon="i-heroicons-clock" color="blue" :loading="loading" />
      <StatCard :label="$t('islamic.active_halaqah')" :value="stats.activeHalaqah" icon="i-heroicons-user-group" color="amber" :loading="loading" />
      <StatCard :label="$t('islamic.ziswaf_collected')" :value="formatCurrency(stats.ziswafCollected)" icon="i-heroicons-currency-dollar" color="purple" :loading="loading" />
    </div>

    <div class="card">
      <div class="card-header"><h3 class="card-title">{{ $t('islamic.modules') }}</h3></div>
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
        <NuxtLink v-for="mod in modules" :key="mod.to" :to="mod.to" class="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 hover:border-brand-300 dark:hover:border-brand-700 transition-all text-center">
          <div class="w-12 h-12 rounded-full flex items-center justify-center" :class="mod.bg">
            <UIcon :name="mod.icon" class="w-6 h-6" :class="mod.iconColor" />
          </div>
          <span class="text-sm font-medium text-gray-900 dark:text-white">{{ mod.label }}</span>
        </NuxtLink>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('islamic.tahfidz_overview') }}</h3></div>
        <div class="h-72">
          <ApexChart type="bar" height="100%" :options="tahfidzOverviewChart.options" :series="tahfidzOverviewChart.series" />
        </div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('islamic.prayer_compliance') }}</h3></div>
        <div class="h-72">
          <ApexChart type="radar" height="100%" :options="prayerComplianceChart.options" :series="prayerComplianceChart.series" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const loading = ref(false)
const stats = reactive({ totalHafiz: 0, avgPrayer: 0, activeHalaqah: 0, ziswafCollected: 0 })

const modules = [
  { to: '/islamic/tahfidz', label: t('tahfidz.title'), icon: 'i-heroicons-book-open', bg: 'bg-emerald-100 dark:bg-emerald-900/30', iconColor: 'text-emerald-600 dark:text-emerald-400' },
  { to: '/islamic/tahfidz-progress', label: t('islamic.tahfidz_progress'), icon: 'i-heroicons-chart-bar', bg: 'bg-emerald-100 dark:bg-emerald-900/30', iconColor: 'text-emerald-600 dark:text-emerald-400' },
  { to: '/islamic/tasmi', label: t('islamic.tasmi'), icon: 'i-heroicons-microphone', bg: 'bg-blue-100 dark:bg-blue-900/30', iconColor: 'text-blue-600 dark:text-blue-400' },
  { to: '/islamic/mutabaah', label: t('mutabaah.title'), icon: 'i-heroicons-clipboard-document-check', bg: 'bg-amber-100 dark:bg-amber-900/30', iconColor: 'text-amber-600 dark:text-amber-400' },
  { to: '/islamic/prayer-attendance', label: t('islamic.prayer_attendance'), icon: 'i-heroicons-clock', bg: 'bg-sky-100 dark:bg-sky-900/30', iconColor: 'text-sky-600 dark:text-sky-400' },
  { to: '/islamic/islamic-character', label: t('islamic.islamic_character'), icon: 'i-heroicons-star', bg: 'bg-purple-100 dark:bg-purple-900/30', iconColor: 'text-purple-600 dark:text-purple-400' },
  { to: '/islamic/halaqah', label: t('halaqah.title'), icon: 'i-heroicons-user-group', bg: 'bg-rose-100 dark:bg-rose-900/30', iconColor: 'text-rose-600 dark:text-rose-400' },
  { to: '/islamic/quranic-competencies', label: t('islamic.quranic_competencies'), icon: 'i-heroicons-academic-cap', bg: 'bg-indigo-100 dark:bg-indigo-900/30', iconColor: 'text-indigo-600 dark:text-indigo-400' },
  { to: '/islamic/islamic-events', label: t('islamic.islamic_events'), icon: 'i-heroicons-calendar', bg: 'bg-orange-100 dark:bg-orange-900/30', iconColor: 'text-orange-600 dark:text-orange-400' },
  { to: '/islamic/ziswaf', label: t('islamic.ziswaf'), icon: 'i-heroicons-heart', bg: 'bg-red-100 dark:bg-red-900/30', iconColor: 'text-red-600 dark:text-red-400' },
]

const tahfidzOverviewChart = computed(() => ({
  series: [{ name: t('tahfidz.students'), data: [12, 18, 15, 22, 19, 25, 21, 28, 24, 30] }],
  options: { chart: { type: 'bar' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669'], plotOptions: { bar: { borderRadius: 4 } }, xaxis: { categories: Array.from({ length: 10 }, (_, i) => `${t('tahfidz.juz')} ${i + 21}`) }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))
const prayerComplianceChart = computed(() => ({
  series: [{ name: t('islamic.compliance'), data: [85, 78, 72, 90, 82] }],
  options: { chart: { type: 'radar' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669'], xaxis: { categories: [t('prayer.fajr'), t('prayer.dhuhr'), t('prayer.asr'), t('prayer.maghrib'), t('prayer.isha')] }, yaxis: { show: false, min: 0, max: 100 } },
}))

const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

onMounted(async () => {
  loading.value = true
  try { Object.assign(stats, await api.get<{totalHafiz:number;avgPrayer:number;activeHalaqah:number;ziswafCollected:number}>('/islamic/stats')) } catch {} finally { loading.value = false }
})
</script>

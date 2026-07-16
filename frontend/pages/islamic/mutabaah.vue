<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('mutabaah.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('mutabaah.subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedClass" :options="classOptions" :placeholder="$t('academic.select_class')" color="gray" size="sm" class="w-36" @change="fetchMutabaah" />
        <UInput v-model="selectedDate" type="date" color="gray" size="sm" @change="fetchMutabaah" />
        <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openBulkEntry">{{ $t('mutabaah.bulk_entry') }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="card lg:col-span-2">
        <div class="card-header"><h3 class="card-title">{{ $t('mutabaah.weekly_summary') }}</h3></div>
        <div class="h-72">
          <ApexChart type="bar" height="100%" :options="weeklyChart.options" :series="weeklyChart.series" />
        </div>
      </div>
      <div class="card">
        <div class="card-header"><h3 class="card-title">{{ $t('mutabaah.prayer_stats') }}</h3></div>
        <div class="h-72">
          <ApexChart type="radialBar" height="100%" :options="prayerStatsChart.options" :series="prayerStatsChart.series" />
        </div>
      </div>
    </div>

    <DataTable :columns="columns" :rows="records" :loading="loading" :empty-title="$t('mutabaah.no_records')" :show-export="false">
      <template #cell-studentName="{ row }"><span class="text-sm font-medium text-gray-900 dark:text-white">{{ row.studentName }}</span></template>
      <template #cell-fajr="{ row }"><UIcon :name="row.fajr ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.fajr ? 'text-emerald-500' : 'text-red-500'" /></template>
      <template #cell-dhuhr="{ row }"><UIcon :name="row.dhuhr ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.dhuhr ? 'text-emerald-500' : 'text-red-500'" /></template>
      <template #cell-asr="{ row }"><UIcon :name="row.asr ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.asr ? 'text-emerald-500' : 'text-red-500'" /></template>
      <template #cell-maghrib="{ row }"><UIcon :name="row.maghrib ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.maghrib ? 'text-emerald-500' : 'text-red-500'" /></template>
      <template #cell-isha="{ row }"><UIcon :name="row.isha ? 'i-heroicons-check-circle' : 'i-heroicons-x-circle'" class="w-5 h-5" :class="row.isha ? 'text-emerald-500' : 'text-red-500'" /></template>
    </DataTable>

    <FormDialog v-model="showBulkDialog" :title="$t('mutabaah.bulk_entry')" :loading="bulkSaving" @submit="submitBulkEntry" @cancel="showBulkDialog=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('common.date')" required><UInput v-model="bulkDate" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.class')" required><USelect v-model="bulkClass" :options="classOptions" color="gray" /></UFormGroup>
        <div class="grid grid-cols-3 gap-3">
          <UFormGroup :label="$t('prayer.fajr')"><UToggle v-model="bulkForm.fajr" /></UFormGroup>
          <UFormGroup :label="$t('prayer.dhuhr')"><UToggle v-model="bulkForm.dhuhr" /></UFormGroup>
          <UFormGroup :label="$t('prayer.asr')"><UToggle v-model="bulkForm.asr" /></UFormGroup>
          <UFormGroup :label="$t('prayer.maghrib')"><UToggle v-model="bulkForm.maghrib" /></UFormGroup>
          <UFormGroup :label="$t('prayer.isha')"><UToggle v-model="bulkForm.isha" /></UFormGroup>
          <UFormGroup :label="$t('prayer.tahajjud')"><UToggle v-model="bulkForm.tahajjud" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('quran.pages')"><UInput v-model.number="bulkForm.quranPages" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('prayer.dhuha')"><UToggle v-model="bulkForm.dhuha" /></UFormGroup>
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
const bulkSaving = ref(false)
const showBulkDialog = ref(false)
const records = ref<Record<string, unknown>[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const selectedClass = ref('')
const selectedDate = ref($dayjs().format('YYYY-MM-DD'))
const bulkDate = ref($dayjs().format('YYYY-MM-DD'))
const bulkClass = ref('')

const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name' },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'fajr', label: 'prayer.fajr' }, { key: 'dhuhr', label: 'prayer.dhuhr' },
  { key: 'asr', label: 'prayer.asr' }, { key: 'maghrib', label: 'prayer.maghrib' },
  { key: 'isha', label: 'prayer.isha' }, { key: 'tahajjud', label: 'prayer.tahajjud' },
  { key: 'dhuha', label: 'prayer.dhuha' }, { key: 'quranPages', label: 'quran.pages', type: 'number' },
]

const bulkForm = reactive({ fajr: true, dhuhr: true, asr: true, maghrib: true, isha: true, tahajjud: false, dhuha: false, quranPages: 0 })

const weeklyChart = computed(() => ({
  series: [{ name: t('mutabaah.completed'), data: [45, 52, 48, 55, 58, 50, 42] }, { name: t('mutabaah.missed'), data: [5, 3, 2, 3, 4, 5, 3] }],
  options: { chart: { type: 'bar' as const, stacked: true, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669', '#ef4444'], xaxis: { categories: [t('day.monday'), t('day.tuesday'), t('day.wednesday'), t('day.thursday'), t('day.friday'), t('day.saturday'), t('day.sunday')] }, grid: { borderColor: '#e5e7eb', strokeDashArray: 4 } },
}))
const prayerStatsChart = computed(() => ({
  series: [85],
  options: { chart: { type: 'radialBar' as const, background: 'transparent' }, colors: ['#059669'], plotOptions: { radialBar: { dataLabels: { name: { show: true }, value: { show: true } } } }, labels: [t('mutabaah.prayer_compliance')] },
}))

const fetchMutabaah = async () => { if (!selectedClass.value) return; loading.value = true; try { records.value = await api.get('/mutabaah', { classId: selectedClass.value, date: selectedDate.value }) } catch {} finally { loading.value = false } }

const openBulkEntry = () => { bulkDate.value = selectedDate.value; bulkClass.value = selectedClass.value; showBulkDialog.value = true }

const submitBulkEntry = async () => { bulkSaving.value = true; try { await api.post('/mutabaah/bulk', { date: bulkDate.value, classId: bulkClass.value, ...bulkForm }); toast.add({ title: t('mutabaah.entry_saved'), color: 'success' }); showBulkDialog.value = false; fetchMutabaah() } catch {} finally { bulkSaving.value = false } }

const fetchOptions = async () => { try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); if (classOptions.value.length > 0) { selectedClass.value = classOptions.value[0].value; fetchMutabaah() } } catch {} }

onMounted(() => fetchOptions())
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('calendar.title') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('calendar.subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <UButton color="gray" variant="ghost" size="sm" icon="i-heroicons-chevron-left" @click="prevMonth" />
        <span class="text-sm font-semibold text-gray-900 dark:text-white">{{ currentMonthName }} {{ currentYear }}</span>
        <UButton color="gray" variant="ghost" size="sm" icon="i-heroicons-chevron-right" @click="nextMonth" />
        <UButton color="gray" variant="outline" size="sm" @click="goToday">{{ $t('calendar.today') }}</UButton>
      </div>
    </div>

    <div class="card overflow-hidden !p-0">
      <div class="grid grid-cols-7">
        <div v-for="day in dayLabels" :key="day" class="px-2 py-3 text-center text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-800/50 border-b border-gray-200 dark:border-gray-700">{{ day }}</div>
        <div v-for="(day, idx) in calendarDays" :key="idx" class="min-h-[100px] p-2 border-b border-r border-gray-100 dark:border-gray-750 hover:bg-gray-50 dark:hover:bg-gray-800/30 cursor-pointer transition-colors"
          :class="{ 'text-gray-300 dark:text-gray-600': !day.isCurrentMonth, 'bg-brand-50 dark:bg-brand-900/10': day.isToday }"
          @click="selectDate(day.date)">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium" :class="{ 'bg-brand-500 text-white w-7 h-7 rounded-full flex items-center justify-center': day.isToday }">{{ day.dayNumber }}</span>
          </div>
          <div class="mt-1 space-y-0.5">
            <div v-for="event in getEventsForDate(day.date)" :key="event.id" class="truncate text-[10px] px-1 py-0.5 rounded" :class="event.colorClass" :title="event.title">
              {{ event.title }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="selectedDateEvents.length > 0" class="card">
      <div class="card-header"><h3 class="card-title">{{ $dayjs(selectedDate).format('DD MMMM YYYY') }}</h3></div>
      <div class="space-y-2">
        <div v-for="event in selectedDateEvents" :key="event.id" class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50">
          <div class="w-2 h-2 rounded-full" :class="event.dotColor"></div>
          <div class="flex-1"><p class="text-sm font-medium text-gray-900 dark:text-white">{{ event.title }}</p><p class="text-xs text-gray-500">{{ event.time }} - {{ event.description }}</p></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const { $dayjs } = useNuxtApp()
const loading = ref(false)
const currentDate = ref($dayjs()); const selectedDate = ref($dayjs().format('YYYY-MM-DD'))
const events = ref<{ id: string; title: string; date: string; time: string; description: string; type: string; colorClass: string; dotColor: string }[]>([])

const currentMonthName = computed(() => currentDate.value.format('MMMM'))
const currentYear = computed(() => currentDate.value.year())
const dayLabels = computed(() => {
  const locale = $dayjs.locale()
  const days = []; for (let i = 0; i < 7; i++) days.push($dayjs().day(i).format('dd')); return days
})
const calendarDays = computed(() => {
  const start = currentDate.value.startOf('month').startOf('week')
  const end = currentDate.value.endOf('month').endOf('week')
  const days = []; let d = start
  while (d.isBefore(end) || d.isSame(end, 'day')) {
    days.push({ date: d.format('YYYY-MM-DD'), dayNumber: d.date(), isCurrentMonth: d.month() === currentDate.value.month(), isToday: d.format('YYYY-MM-DD') === $dayjs().format('YYYY-MM-DD') })
    d = d.add(1, 'day')
  }
  return days
})
const selectedDateEvents = computed(() => events.value.filter(e => e.date === selectedDate.value))

const getEventsForDate = (date: string) => events.value.filter(e => e.date === date).slice(0, 3)

const prevMonth = () => { currentDate.value = currentDate.value.subtract(1, 'month'); fetchEvents() }
const nextMonth = () => { currentDate.value = currentDate.value.add(1, 'month'); fetchEvents() }
const goToday = () => { currentDate.value = $dayjs(); selectedDate.value = $dayjs().format('YYYY-MM-DD') }
const selectDate = (date: string) => { selectedDate.value = date }

const fetchEvents = async () => { loading.value = true; try { events.value = await api.get('/calendar/events', { month: currentDate.value.format('YYYY-MM'), year: currentDate.value.format('YYYY') }) } catch {} finally { loading.value = false } }
onMounted(() => fetchEvents())
</script>

import dayjs from 'dayjs'
import localizedFormat from 'dayjs/plugin/localizedFormat'
import relativeTime from 'dayjs/plugin/relativeTime'
import utc from 'dayjs/plugin/utc'
import timezone from 'dayjs/plugin/timezone'
import isToday from 'dayjs/plugin/isToday'
import isBetween from 'dayjs/plugin/isBetween'
import weekOfYear from 'dayjs/plugin/weekOfYear'
import 'dayjs/locale/id'
import 'dayjs/locale/en'

declare module '#app' {
  interface NuxtApp {
    $dayjs: typeof dayjs
  }
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $dayjs: typeof dayjs
  }
}

export default defineNuxtPlugin(() => {
  dayjs.extend(localizedFormat)
  dayjs.extend(relativeTime)
  dayjs.extend(utc)
  dayjs.extend(timezone)
  dayjs.extend(isToday)
  dayjs.extend(isBetween)
  dayjs.extend(weekOfYear)

  const { locale } = useI18n()
  watch(locale, (newLocale) => {
    dayjs.locale(newLocale)
  }, { immediate: true })

  dayjs.locale(locale.value || 'id')

  return {
    provide: {
      dayjs,
    },
  }
})

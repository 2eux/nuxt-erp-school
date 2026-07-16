<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('settings.title') }}</h1>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('settings.subtitle') }}</p>
    </div>

    <UTabs
      :items="tabs"
      :ui="{ list: { tab: { size: 'sm' } } }"
      class="w-full"
    >
      <template #item="{ item }">
        <div class="space-y-6 pt-4">
          <template v-if="item.key === 'school'">
            <div class="card">
              <div class="card-header">
                <h3 class="card-title">{{ $t('settings.school_profile') }}</h3>
              </div>

              <div v-if="loading" class="space-y-4">
                <LoadingSkeleton type="detail" />
              </div>

              <div v-else class="space-y-4">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <UFormGroup :label="$t('school.name')" required>
                    <UInput v-model="schoolForm.name" color="gray" />
                  </UFormGroup>
                  <UFormGroup :label="$t('school.npsn')">
                    <UInput v-model="schoolForm.npsn" color="gray" />
                  </UFormGroup>
                  <UFormGroup :label="$t('school.address')" class="sm:col-span-2">
                    <UInput v-model="schoolForm.address" color="gray" />
                  </UFormGroup>
                  <UFormGroup :label="$t('school.city')">
                    <UInput v-model="schoolForm.city" color="gray" />
                  </UFormGroup>
                  <UFormGroup :label="$t('school.province')">
                    <UInput v-model="schoolForm.province" color="gray" />
                  </UFormGroup>
                  <UFormGroup :label="$t('school.phone')">
                    <UInput v-model="schoolForm.phone" color="gray" />
                  </UFormGroup>
                  <UFormGroup :label="$t('school.email')">
                    <UInput v-model="schoolForm.email" color="gray" />
                  </UFormGroup>
                  <UFormGroup :label="$t('school.principal')">
                    <UInput v-model="schoolForm.principal" color="gray" />
                  </UFormGroup>
                  <UFormGroup :label="$t('school.accreditation')">
                    <USelect
                      v-model="schoolForm.accreditation"
                      :options="['A', 'B', 'C', 'Not Accredited']"
                      color="gray"
                    />
                  </UFormGroup>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    {{ $t('school.logo') }}
                  </label>
                  <FileUpload
                    accept="image/*"
                    :accept-hint="$t('upload.image_hint')"
                    @files-selected="handleLogoSelected"
                  />
                </div>

                <div class="flex justify-end">
                  <UButton color="primary" :loading="saving" @click="saveSchool">
                    {{ $t('common.save_changes') }}
                  </UButton>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="item.key === 'academic'">
            <div class="space-y-4">
              <div class="card">
                <div class="card-header">
                  <div>
                    <h3 class="card-title">{{ $t('settings.academic_years') }}</h3>
                    <p class="text-sm text-gray-500 mt-1">{{ $t('settings.academic_years_description') }}</p>
                  </div>
                  <UButton color="primary" size="sm" icon="i-heroicons-plus">
                    {{ $t('common.add') }}
                  </UButton>
                </div>

                <DataTable
                  :columns="academicYearColumns"
                  :rows="academicYears"
                  :loading="loading"
                  :empty-title="$t('settings.no_academic_years')"
                >
                  <template #cell-isActive="{ row }">
                    <StatusBadge :status="row.isActive ? 'active' : 'inactive'" />
                  </template>
                </DataTable>
              </div>

              <div class="card">
                <div class="card-header">
                  <h3 class="card-title">{{ $t('settings.grading_system') }}</h3>
                </div>

                <div class="flex items-center gap-3 mb-4">
                  <UInput
                    v-model.number="gradingSystem.passingScore"
                    type="number"
                    :label="$t('settings.passing_score')"
                    color="gray"
                    class="w-32"
                  />
                  <UInput
                    v-model.number="gradingSystem.maxScore"
                    type="number"
                    :label="$t('settings.max_score')"
                    color="gray"
                    class="w-32"
                  />
                </div>

                <div class="flex justify-end">
                  <UButton color="primary" size="sm">{{ $t('common.save') }}</UButton>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="item.key === 'notifications'">
            <div class="card">
              <div class="card-header">
                <h3 class="card-title">{{ $t('settings.notification_preferences') }}</h3>
              </div>

              <div class="space-y-4">
                <div v-for="pref in notificationPrefs" :key="pref.key" class="flex items-center justify-between py-2">
                  <div>
                    <p class="text-sm font-medium text-gray-900 dark:text-white">{{ $t(pref.label) }}</p>
                    <p class="text-xs text-gray-500 dark:text-gray-400">{{ $t(pref.description) }}</p>
                  </div>
                  <UToggle v-model="pref.enabled" />
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="item.key === 'theme'">
            <div class="card">
              <div class="card-header">
                <h3 class="card-title">{{ $t('settings.theme_settings') }}</h3>
              </div>

              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    {{ $t('theme.mode') }}
                  </label>
                  <USelect
                    v-model="themeMode"
                    :options="['system', 'light', 'dark']"
                    color="gray"
                    @change="handleThemeChange"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    {{ $t('theme.primary_color') }}
                  </label>
                  <div class="flex items-center gap-2">
                    <button
                      v-for="color in primaryColors"
                      :key="color.value"
                      class="w-8 h-8 rounded-full border-2 transition-all"
                      :class="[
                        color.bg,
                        themeStore.themePreference.primaryColor === color.value
                          ? 'ring-2 ring-offset-2 dark:ring-offset-gray-800 ring-' + color.value + '-500 border-white'
                          : 'border-transparent'
                      ]"
                      :title="color.label"
                      @click="themeStore.setPrimaryColor(color.value)"
                    />
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </template>
    </UTabs>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'

definePageMeta({
  middleware: ['auth'],
})

const { t } = useI18n()
const themeStore = useTheme()

const loading = ref(false)
const saving = ref(false)

const tabs = computed(() => [
  { key: 'school', label: t('settings.school') },
  { key: 'academic', label: t('settings.academic') },
  { key: 'notifications', label: t('settings.notifications') },
  { key: 'theme', label: t('theme.title') },
])

const schoolForm = reactive({
  name: '',
  npsn: '',
  address: '',
  city: '',
  province: '',
  phone: '',
  email: '',
  principal: '',
  accreditation: '',
})

const academicYearColumns: TableColumn[] = [
  { key: 'name', label: 'academic_year.name' },
  { key: 'startDate', label: 'academic_year.start_date', type: 'date' },
  { key: 'endDate', label: 'academic_year.end_date', type: 'date' },
  { key: 'isActive', label: 'common.status', type: 'status' },
]

const academicYears = ref<Record<string, unknown>[]>([])

const gradingSystem = reactive({
  passingScore: 75,
  maxScore: 100,
})

const notificationPrefs = reactive([
  { key: 'email', label: 'settings.email_notifications', description: 'settings.email_notifications_desc', enabled: true },
  { key: 'push', label: 'settings.push_notifications', description: 'settings.push_notifications_desc', enabled: true },
  { key: 'sms', label: 'settings.sms_notifications', description: 'settings.sms_notifications_desc', enabled: false },
])

const primaryColors = [
  { value: 'emerald', label: 'Emerald', bg: 'bg-emerald-500' },
  { value: 'blue', label: 'Blue', bg: 'bg-blue-500' },
  { value: 'indigo', label: 'Indigo', bg: 'bg-indigo-500' },
  { value: 'purple', label: 'Purple', bg: 'bg-purple-500' },
  { value: 'rose', label: 'Rose', bg: 'bg-rose-500' },
]

const themeMode = ref(themeStore.currentMode.value)

const handleThemeChange = () => {
  themeStore.setMode(themeMode.value as 'light' | 'dark' | 'system')
}

const handleLogoSelected = (files: File[]) => {
  // handle logo upload
}

const saveSchool = async () => {
  saving.value = true
  try {
    // save school profile
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('settings.school_profile') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('settings.school_profile_subtitle') }}</p></div>

    <div class="card">
      <div class="card-header"><h3 class="card-title">{{ $t('settings.basic_info') }}</h3></div>
      <div v-if="loading"><LoadingSkeleton type="detail" /></div>
      <div v-else class="space-y-6">
        <div class="flex items-center gap-6">
          <div class="w-24 h-24 rounded-xl bg-gray-100 dark:bg-gray-800 flex items-center justify-center overflow-hidden">
            <img v-if="form.logo" :src="form.logo" class="w-full h-full object-cover" />
            <UIcon v-else name="i-heroicons-building-office" class="w-10 h-10 text-gray-400" />
          </div>
          <div>
            <UButton color="gray" variant="outline" size="sm" @click="triggerLogoUpload">{{ $t('settings.upload_logo') }}</UButton>
            <p class="text-xs text-gray-500 mt-1">{{ $t('settings.logo_hint') }}</p>
          </div>
          <input ref="logoInput" type="file" accept="image/*" class="hidden" @change="handleLogoChange" />
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <UFormGroup :label="$t('school.name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.npsn')"><UInput v-model="form.npsn" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.address')" class="sm:col-span-2"><UTextarea v-model="form.address" color="gray" :rows="2" /></UFormGroup>
          <UFormGroup :label="$t('school.city')"><UInput v-model="form.city" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.province')"><UInput v-model="form.province" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.postal_code')"><UInput v-model="form.postalCode" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.phone')"><UInput v-model="form.phone" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.email')"><UInput v-model="form.email" type="email" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.website')"><UInput v-model="form.website" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.principal')"><UInput v-model="form.principal" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.established')"><UInput v-model="form.established" type="date" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.accreditation')"><USelect v-model="form.accreditation" :options="['A','B','C','Not Accredited']" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('school.curriculum')"><UInput v-model="form.curriculum" color="gray" /></UFormGroup>
        </div>
        <div class="flex justify-end pt-4 border-t border-gray-200 dark:border-gray-700">
          <UButton color="primary" :loading="saving" @click="saveSchool">{{ $t('common.save_changes') }}</UButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const toast = useToast()
const loading = ref(false); const saving = ref(false)
const logoInput = ref<HTMLInputElement>()

const form = reactive({ name: '', npsn: '', logo: '', address: '', city: '', province: '', postalCode: '', phone: '', email: '', website: '', principal: '', established: '', accreditation: '', curriculum: '' })

const fetchSchool = async () => { loading.value = true; try { const s = await api.get<Record<string, unknown>>('/schools/current'); Object.assign(form, s) } catch {} finally { loading.value = false } }
const triggerLogoUpload = () => logoInput.value?.click()
const handleLogoChange = (e: Event) => { const file = (e.target as HTMLInputElement).files?.[0]; if (file) { const fd = new FormData(); fd.append('logo', file); api.upload('/schools/current/logo', fd); toast.add({ title: t('settings.logo_updated'), color: 'success' }) } }
const saveSchool = async () => { saving.value = true; try { await api.put('/schools/current', form); toast.add({ title: t('settings.school_updated'), color: 'success' }) } catch {} finally { saving.value = false } }
onMounted(() => fetchSchool())
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('students.enroll_student') }}</h1>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('students.enroll_description') }}</p>
    </div>

    <div class="card">
      <div class="flex items-center gap-2 mb-6">
        <div v-for="(step, i) in steps" :key="step.key" class="flex items-center gap-2">
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold transition-all"
            :class="currentStep >= i ? 'bg-brand-500 text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-500'"
          >
            {{ i + 1 }}
          </div>
          <span class="text-sm font-medium" :class="currentStep >= i ? 'text-gray-900 dark:text-white' : 'text-gray-400'">
            {{ step.label }}
          </span>
          <div v-if="i < steps.length - 1" class="w-12 h-px bg-gray-300 dark:bg-gray-600" />
        </div>
      </div>

      <div v-if="currentStep === 0">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <UFormGroup :label="$t('students.full_name')" required>
            <UInput v-model="form.fullName" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.nick_name')">
            <UInput v-model="form.nickName" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.nis')">
            <UInput v-model="form.nis" color="gray" readonly placeholder="Auto-generated" />
          </UFormGroup>
          <UFormGroup :label="$t('students.nisn')">
            <UInput v-model="form.nisn" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('common.gender')" required>
            <USelect v-model="form.gender" :options="genderOptions" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.birth_place')">
            <UInput v-model="form.birthPlace" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.birth_date')" required>
            <UInput v-model="form.birthDate" type="date" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.religion')">
            <USelect v-model="form.religion" :options="religionOptions" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.enrollment_date')">
            <UInput v-model="form.enrollmentDate" type="date" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.class')" required>
            <USelect v-model="form.classId" :options="classOptions" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.address')" class="sm:col-span-2">
            <UTextarea v-model="form.address" color="gray" :rows="2" />
          </UFormGroup>
          <UFormGroup :label="$t('students.city')">
            <UInput v-model="form.city" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.postal_code')">
            <UInput v-model="form.postalCode" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.phone')">
            <UInput v-model="form.phone" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.email')">
            <UInput v-model="form.email" type="email" color="gray" />
          </UFormGroup>
        </div>
      </div>

      <div v-else-if="currentStep === 1">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <UFormGroup :label="$t('students.father_name')" required class="sm:col-span-2">
            <UInput v-model="form.fatherName" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.father_phone')" required>
            <UInput v-model="form.fatherPhone" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.father_occupation')">
            <UInput v-model="form.fatherOccupation" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.mother_name')" required class="sm:col-span-2">
            <UInput v-model="form.motherName" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.mother_phone')" required>
            <UInput v-model="form.motherPhone" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.mother_occupation')">
            <UInput v-model="form.motherOccupation" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.guardian_name')">
            <UInput v-model="form.guardianName" color="gray" />
          </UFormGroup>
          <UFormGroup :label="$t('students.guardian_phone')">
            <UInput v-model="form.guardianPhone" color="gray" />
          </UFormGroup>
        </div>
      </div>

      <div v-else-if="currentStep === 2">
        <div class="space-y-6">
          <div>
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-2">{{ $t('students.required_documents') }}</h3>
            <div class="space-y-3">
              <div v-for="doc in requiredDocuments" :key="doc.key" class="flex items-center justify-between p-3 rounded-lg border border-gray-200 dark:border-gray-700">
                <div class="flex items-center gap-3">
                  <UIcon name="i-heroicons-document-text" class="w-5 h-5 text-gray-400" />
                  <div><p class="text-sm text-gray-900 dark:text-white">{{ doc.label }}</p><p class="text-xs text-gray-500">{{ doc.hint }}</p></div>
                </div>
                <div class="flex items-center gap-2">
                  <UIcon v-if="docFiles[doc.key]" name="i-heroicons-check-circle" class="w-5 h-5 text-emerald-500" />
                  <UButton color="gray" variant="outline" size="xs" @click="selectDocFile(doc.key)">
                    {{ docFiles[doc.key] ? $t('common.change') : $t('upload.upload') }}
                  </UButton>
                </div>
              </div>
            </div>
          </div>
          <div>
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-2">{{ $t('students.photo') }}</h3>
            <div class="flex items-center gap-4">
              <UAvatar :src="form.photo || undefined" size="xl" />
              <FileUpload
                ref="photoUploadRef"
                accept="image/*"
                :multiple="false"
                :accept-hint="$t('upload.image_hint')"
                @files-selected="handlePhotoSelected"
              />
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="currentStep === 3">
        <div class="space-y-4">
          <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50">
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-3">{{ $t('students.personal_info') }}</h3>
            <div class="grid grid-cols-2 gap-3">
              <div><span class="text-xs text-gray-500">{{ $t('students.nis') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.nis || '-' }}</p></div>
              <div><span class="text-xs text-gray-500">{{ $t('students.nisn') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.nisn || '-' }}</p></div>
              <div><span class="text-xs text-gray-500">{{ $t('students.full_name') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.fullName }}</p></div>
              <div><span class="text-xs text-gray-500">{{ $t('students.class') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ selectedClassName }}</p></div>
              <div><span class="text-xs text-gray-500">{{ $t('students.birth_info') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.birthPlace }}, {{ form.birthDate || '-' }}</p></div>
              <div><span class="text-xs text-gray-500">{{ $t('students.address') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.address }}, {{ form.city }}</p></div>
            </div>
          </div>
          <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50">
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-3">{{ $t('students.parent_info') }}</h3>
            <div class="grid grid-cols-2 gap-3">
              <div><span class="text-xs text-gray-500">{{ $t('students.father_name') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.fatherName }}</p></div>
              <div><span class="text-xs text-gray-500">{{ $t('students.father_phone') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.fatherPhone }}</p></div>
              <div><span class="text-xs text-gray-500">{{ $t('students.mother_name') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.motherName }}</p></div>
              <div><span class="text-xs text-gray-500">{{ $t('students.mother_phone') }}</span><p class="text-sm text-gray-900 dark:text-white">{{ form.motherPhone }}</p></div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-between mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
        <UButton
          v-if="currentStep > 0"
          color="gray"
          variant="ghost"
          icon="i-heroicons-arrow-left"
          @click="currentStep--"
        >
          {{ $t('common.previous') }}
        </UButton>
        <div v-else />
        <div class="flex items-center gap-3">
          <UButton color="gray" variant="ghost" @click="navigateTo('/students')">
            {{ $t('common.cancel') }}
          </UButton>
          <UButton
            v-if="currentStep < steps.length - 1"
            color="primary"
            icon="i-heroicons-arrow-right"
            trailing
            :disabled="!canProceedFrom(currentStep)"
            @click="currentStep++"
          >
            {{ $t('common.next') }}
          </UButton>
          <UButton
            v-else
            color="primary"
            icon="i-heroicons-check"
            :loading="submitting"
            @click="submitEnrollment"
          >
            {{ $t('students.complete_enrollment') }}
          </UButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

const { t } = useI18n()
const api = useApi()
const toast = useToast()

const currentStep = ref(0)
const submitting = ref(false)

const classOptions = ref<{ label: string; value: string }[]>([])
const photoUploadRef = ref()
const docFiles = reactive<Record<string, File | null>>({})

const steps = [
  { key: 'personal', label: t('students.personal_info') },
  { key: 'parent', label: t('students.parent_info') },
  { key: 'documents', label: t('documents.title') },
  { key: 'review', label: t('common.review') },
]

const requiredDocuments = [
  { key: 'birth_certificate', label: t('students.birth_certificate'), hint: 'PDF, JPG - Max 5MB' },
  { key: 'family_card', label: t('students.family_card'), hint: 'PDF, JPG - Max 5MB' },
  { key: 'previous_report', label: t('students.previous_report_card'), hint: 'PDF, JPG - Max 5MB' },
  { key: 'transfer_certificate', label: t('students.transfer_certificate'), hint: 'PDF, JPG - Max 5MB (optional)' },
]

const genderOptions = [
  { label: t('common.male'), value: 'male' },
  { label: t('common.female'), value: 'female' },
]

const religionOptions = [
  { label: 'Islam', value: 'Islam' },
  { label: 'Kristen', value: 'Kristen' },
  { label: 'Katolik', value: 'Katolik' },
  { label: 'Hindu', value: 'Hindu' },
  { label: 'Buddha', value: 'Buddha' },
  { label: 'Konghucu', value: 'Konghucu' },
]

const form = reactive({
  nis: '',
  nisn: '',
  fullName: '',
  nickName: '',
  gender: 'male' as string,
  birthPlace: '',
  birthDate: '',
  religion: 'Islam',
  address: '',
  city: '',
  postalCode: '',
  phone: '',
  email: '',
  classId: '',
  enrollmentDate: new Date().toISOString().split('T')[0],
  fatherName: '',
  fatherPhone: '',
  fatherOccupation: '',
  motherName: '',
  motherPhone: '',
  motherOccupation: '',
  guardianName: '',
  guardianPhone: '',
  photo: '',
  status: 'active' as string,
})

const selectedClassName = computed(() => {
  const found = classOptions.value.find(c => c.value === form.classId)
  return found?.label || '-'
})

const canProceedFrom = (step: number): boolean => {
  if (step === 0) {
    return !!form.fullName && !!form.gender && !!form.birthDate && !!form.classId
  }
  if (step === 1) {
    return !!form.fatherName && !!form.fatherPhone && !!form.motherName && !!form.motherPhone
  }
  return true
}

const handlePhotoSelected = (files: File[]) => {
  if (files.length > 0) {
    const reader = new FileReader()
    reader.onload = (e) => {
      form.photo = e.target?.result as string
    }
    reader.readAsDataURL(files[0])
  }
}

const selectDocFile = (key: string) => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.pdf,.jpg,.jpeg,.png'
  input.onchange = (e) => {
    const target = e.target as HTMLInputElement
    if (target.files?.[0]) {
      docFiles[key] = target.files[0]
    }
  }
  input.click()
}

const fetchClasses = async () => {
  try {
    const res = await api.get<{ id: string; name: string }[]>('/classes', { limit: 100 })
    classOptions.value = res.map((c: { id: string; name: string }) => ({ label: c.name, value: c.id }))
  } catch { /* ignore */ }
}

const submitEnrollment = async () => {
  submitting.value = true
  try {
    const formData = new FormData()
    const payload = { ...form }
    delete (payload as Record<string, unknown>).photo
    formData.append('data', JSON.stringify(payload))

    if (photoUploadRef.value?.files?.[0]) {
      formData.append('photo', photoUploadRef.value.files[0])
    }
    for (const [key, file] of Object.entries(docFiles)) {
      if (file) {
        formData.append(`documents[${key}]`, file)
      }
    }

    await api.upload('/students/enroll', formData)
    toast.add({ title: t('students.enrollment_success'), color: 'success' })
    navigateTo('/students')
  } catch { /* handled by api */ }
  finally { submitting.value = false }
}

onMounted(async () => {
  await fetchClasses()
  try {
    const nisRes = await api.get<{ nis: string }>('/students/generate-nis')
    form.nis = nisRes.nis
  } catch { /* ignore */ }
})
</script>

<template>
  <div class="space-y-6">
    <div class="text-center">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white">
        {{ $t('auth.reset_password') }}
      </h2>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
        {{ $t('auth.reset_password_subtitle') }}
      </p>
    </div>

    <form class="space-y-4" @submit.prevent="handleSubmit">
      <div>
        <UInput
          v-model="form.password"
          :type="showPassword ? 'text' : 'password'"
          :placeholder="$t('auth.new_password')"
          icon="i-heroicons-lock-closed"
          size="lg"
          color="gray"
          :disabled="loading || success"
          autocomplete="new-password"
          autofocus
        />
        <p v-if="errors.password" class="text-xs text-red-500 mt-1">{{ errors.password }}</p>
      </div>

      <div>
        <UInput
          v-model="form.passwordConfirmation"
          :type="showPassword ? 'text' : 'password'"
          :placeholder="$t('auth.confirm_password')"
          icon="i-heroicons-lock-closed"
          size="lg"
          color="gray"
          :disabled="loading || success"
          autocomplete="new-password"
        />
        <p v-if="errors.passwordConfirmation" class="text-xs text-red-500 mt-1">{{ errors.passwordConfirmation }}</p>
      </div>

      <label class="flex items-center gap-2 cursor-pointer">
        <UCheckbox v-model="showPassword" size="sm" />
        <span class="text-sm text-gray-600 dark:text-gray-400">
          {{ $t('auth.show_password') }}
        </span>
      </label>

      <UButton
        type="submit"
        color="primary"
        size="lg"
        block
        :loading="loading"
        :disabled="loading || success"
      >
        {{ $t('auth.reset_password') }}
      </UButton>
    </form>

    <div v-if="errorMessage" class="p-3 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-800">
      <p class="text-sm text-red-600 dark:text-red-400 text-center">
        {{ errorMessage }}
      </p>
    </div>

    <div v-if="success" class="p-3 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-100 dark:border-emerald-800">
      <div class="flex items-center gap-2 mb-2">
        <UIcon name="i-heroicons-check-circle" class="w-5 h-5 text-emerald-500" />
        <p class="text-sm font-medium text-emerald-700 dark:text-emerald-400">
          {{ $t('auth.password_reset_success') }}
        </p>
      </div>
    </div>

    <div class="text-center">
      <NuxtLink
        to="/auth/login"
        class="text-sm font-medium text-brand-600 hover:text-brand-500 transition-colors"
      >
        {{ $t('auth.back_to_login') }}
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { z } from 'zod'

definePageMeta({
  layout: 'auth',
  middleware: ['auth'],
})

const route = useRoute()

const resetSchema = z.object({
  password: z.string().min(8, 'Password must be at least 8 characters'),
  passwordConfirmation: z.string().min(1, 'Confirm your password'),
}).refine(data => data.password === data.passwordConfirmation, {
  message: 'Passwords do not match',
  path: ['passwordConfirmation'],
})

const form = reactive({
  password: '',
  passwordConfirmation: '',
})

const errors = reactive({
  password: '',
  passwordConfirmation: '',
})

const showPassword = ref(false)
const loading = ref(false)
const success = ref(false)
const errorMessage = ref('')

const auth = useAuth()

const handleSubmit = async () => {
  errors.password = ''
  errors.passwordConfirmation = ''
  errorMessage.value = ''

  const result = resetSchema.safeParse(form)
  if (!result.success) {
    for (const issue of result.error.issues) {
      const field = issue.path[0] as string
      if (field in errors) {
        errors[field as keyof typeof errors] = issue.message
      }
    }
    return
  }

  const token = route.query.token as string
  if (!token) {
    errorMessage.value = 'Invalid or expired reset token.'
    return
  }

  loading.value = true
  try {
    await auth.resetPassword(token, form.password, form.passwordConfirmation)
    success.value = true
    setTimeout(() => {
      navigateTo('/auth/login')
    }, 3000)
  } catch (err: unknown) {
    const error = err as { statusCode?: number; message?: string }
    errorMessage.value = error.message || 'An error occurred. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

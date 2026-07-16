<template>
  <div class="space-y-6">
    <div class="text-center">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white">
        {{ $t('auth.forgot_password') }}
      </h2>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
        {{ $t('auth.forgot_password_subtitle') }}
      </p>
    </div>

    <form class="space-y-4" @submit.prevent="handleSubmit">
      <div>
        <UInput
          v-model="email"
          type="email"
          :placeholder="$t('auth.email')"
          icon="i-heroicons-envelope"
          size="lg"
          color="gray"
          :disabled="loading || sent"
          autofocus
        />
        <p v-if="errorMessage" class="text-xs text-red-500 mt-1">{{ errorMessage }}</p>
      </div>

      <UButton
        type="submit"
        color="primary"
        size="lg"
        block
        :loading="loading"
        :disabled="loading || sent"
      >
        {{ sent ? $t('auth.reset_link_sent') : $t('auth.send_reset_link') }}
      </UButton>
    </form>

    <div v-if="sent" class="p-3 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-100 dark:border-emerald-800">
      <div class="flex items-center gap-2 mb-2">
        <UIcon name="i-heroicons-check-circle" class="w-5 h-5 text-emerald-500" />
        <p class="text-sm font-medium text-emerald-700 dark:text-emerald-400">
          {{ $t('auth.check_email') }}
        </p>
      </div>
      <p class="text-sm text-emerald-600 dark:text-emerald-400">
        {{ $t('auth.check_email_description') }}
      </p>
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

const emailSchema = z.object({
  email: z.string().email('Valid email is required').min(1, 'Email is required'),
})

const email = ref('')
const loading = ref(false)
const sent = ref(false)
const errorMessage = ref('')

const auth = useAuth()

const handleSubmit = async () => {
  errorMessage.value = ''

  const result = emailSchema.safeParse({ email: email.value })
  if (!result.success) {
    errorMessage.value = result.error.issues[0].message
    return
  }

  loading.value = true
  try {
    await auth.forgotPassword(email.value)
    sent.value = true
  } catch (err: unknown) {
    const error = err as { statusCode?: number; message?: string }
    errorMessage.value = error.message || 'An error occurred. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

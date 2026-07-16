<template>
  <div class="space-y-6">
    <div class="text-center">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white">
        {{ $t('auth.login') }}
      </h2>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
        {{ $t('auth.login_subtitle') }}
      </p>

      <div class="mt-3 p-3 rounded-lg bg-brand-50 dark:bg-brand-900/20 border border-brand-100 dark:border-brand-800">
        <p class="text-xs text-brand-700 dark:text-brand-400 font-arabic text-lg leading-relaxed">
          &#1576;&#1616;&#1587;&#1618;&#1605;&#1616; &#1575;&#1604;&#1604;&#1617;&#1614;&#1607;&#1616; &#1575;&#1604;&#1585;&#1617;&#1614;&#1581;&#1618;&#1605;&#1614;&#1600;&#1606;&#1616; &#1575;&#1604;&#1585;&#1617;&#1614;&#1581;&#1616;&#1610;&#1605;&#1616;
        </p>
      </div>
    </div>

    <form class="space-y-4" @submit.prevent="handleLogin">
      <div>
        <UInput
          v-model="form.email"
          type="text"
          :placeholder="$t('auth.email_phone')"
          icon="i-heroicons-envelope"
          size="lg"
          color="gray"
          :disabled="loading"
          autocomplete="email"
          autofocus
        />
        <p v-if="errors.email" class="text-xs text-red-500 mt-1">{{ errors.email }}</p>
      </div>

      <div>
        <UInput
          v-model="form.password"
          :type="showPassword ? 'text' : 'password'"
          :placeholder="$t('auth.password')"
          icon="i-heroicons-lock-closed"
          size="lg"
          color="gray"
          :disabled="loading"
          autocomplete="current-password"
          :ui="{ trailing: { padding: { lg: 'pe-1' } } }"
        >
          <template #trailing>
            <UButton
              :icon="showPassword ? 'i-heroicons-eye-slash' : 'i-heroicons-eye'"
              color="gray"
              variant="link"
              :padded="false"
              @click="showPassword = !showPassword"
            />
          </template>
        </UInput>
        <p v-if="errors.password" class="text-xs text-red-500 mt-1">{{ errors.password }}</p>
      </div>

      <div class="flex items-center justify-between">
        <label class="flex items-center gap-2 cursor-pointer">
          <UCheckbox v-model="form.rememberMe" size="sm" />
          <span class="text-sm text-gray-600 dark:text-gray-400">
            {{ $t('auth.remember_me') }}
          </span>
        </label>

        <NuxtLink
          to="/auth/forgot-password"
          class="text-sm font-medium text-brand-600 hover:text-brand-500 transition-colors"
        >
          {{ $t('auth.forgot_password') }}
        </NuxtLink>
      </div>

      <UButton
        type="submit"
        color="primary"
        size="lg"
        block
        :loading="loading"
        :disabled="loading"
      >
        {{ $t('auth.login') }}
      </UButton>
    </form>

    <div v-if="errorMessage" class="p-3 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-800">
      <p class="text-sm text-red-600 dark:text-red-400 text-center">
        {{ errorMessage }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { z } from 'zod'

definePageMeta({
  layout: 'auth',
  middleware: ['auth'],
})

const loginSchema = z.object({
  email: z.string().min(1, 'Email or phone is required'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
})

const form = reactive({
  email: '',
  password: '',
  rememberMe: false,
})

const errors = reactive({
  email: '',
  password: '',
})

const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')

const auth = useAuth()
const route = useRoute()
const router = useRouter()

const handleLogin = async () => {
  errors.email = ''
  errors.password = ''
  errorMessage.value = ''

  const result = loginSchema.safeParse(form)
  if (!result.success) {
    for (const issue of result.error.issues) {
      const field = issue.path[0] as string
      if (field in errors) {
        errors[field as keyof typeof errors] = issue.message
      }
    }
    return
  }

  loading.value = true
  try {
    await auth.login({
      email: form.email,
      password: form.password,
      rememberMe: form.rememberMe,
    })

    const redirect = (route.query.redirect as string) || '/dashboard'
    await router.push(redirect)
  } catch (err: unknown) {
    const error = err as { statusCode?: number; message?: string }
    if (error.statusCode === 401) {
      errorMessage.value = 'Invalid email or password'
    } else {
      errorMessage.value = 'An error occurred. Please try again.'
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (import.meta.client) {
    const savedEmail = localStorage.getItem('remembered_email')
    if (savedEmail) {
      form.email = savedEmail
      form.rememberMe = true
    }
  }
})
</script>

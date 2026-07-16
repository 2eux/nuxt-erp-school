<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('parents.title') }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('parents.subtitle') }}</p>
      </div>
    </div>
    <DataFilter :filter-fields="filterFields" :searchable="true" @apply="handleFilter" />
    <DataTable :columns="columns" :rows="parents" :loading="loading" :empty-title="$t('parents.no_parents')" :show-export="false">
      <template #cell-fullName="{ row }">
        <div class="flex items-center gap-3">
          <UAvatar :src="(row.photo as string) || undefined" size="sm" />
          <span class="text-sm font-medium text-gray-900 dark:text-white">{{ row.fullName }}</span>
        </div>
      </template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" :to="`/students/${row.studentId}`" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-envelope" @click="sendMessage(row as Record<string, unknown>)" />
        </div>
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const toast = useToast()
const loading = ref(false)
const parents = ref<Record<string, unknown>[]>([])

const columns: TableColumn[] = [
  { key: 'fullName', label: 'common.name', sortable: true },
  { key: 'relationship', label: 'parents.relationship' },
  { key: 'phone', label: 'common.phone' },
  { key: 'email', label: 'common.email' },
  { key: 'childrenCount', label: 'parents.children_count' },
  { key: 'childrenNames', label: 'parents.children' },
]
const filterFields = [
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.active'), value: 'active' }, { label: t('status.inactive'), value: 'inactive' }] },
]
const handleFilter = (filters: Record<string, unknown>) => { fetchParents(filters) }

const sendMessage = (row: Record<string, unknown>) => {
  navigateTo({ path: '/messages', query: { to: row.id as string, type: 'parent' } })
}

const fetchParents = async (filters: Record<string, unknown> = {}) => {
  loading.value = true
  try { parents.value = await api.paginate('/parents', filters).then(r => r.data) } catch { /* ignore */ }
  finally { loading.value = false }
}
onMounted(() => fetchParents())
</script>

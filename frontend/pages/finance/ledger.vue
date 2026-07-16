<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.ledger') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.ledger_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedAccount" :options="accountOptions" :placeholder="$t('finance.select_account')" color="gray" size="sm" class="w-48" @change="fetchLedger" />
        <UInput v-model="startDate" type="date" color="gray" size="sm" @change="fetchLedger" />
        <UInput v-model="endDate" type="date" color="gray" size="sm" @change="fetchLedger" />
        <UButton color="gray" variant="outline" size="sm" icon="i-heroicons-arrow-down-tray" @click="exportLedger">{{ $t('common.export') }}</UButton>
      </div>
    </div>

    <DataTable :columns="columns" :rows="transactions" :loading="loading" :empty-title="$t('finance.no_transactions')" :show-export="false">
      <template #cell-amount="{ row }">
        <span class="text-sm font-mono" :class="(row.type as string) === 'debit' ? 'text-emerald-600' : 'text-red-600'">
          {{ (row.type as string) === 'debit' ? '+' : '-' }}{{ formatCurrency(row.amount as number) }}
        </span>
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const { $dayjs } = useNuxtApp()
const loading = ref(false)
const transactions = ref<Record<string, unknown>[]>([])
const accountOptions = ref<{ label: string; value: string }[]>([])
const selectedAccount = ref('')
const startDate = ref($dayjs().startOf('month').format('YYYY-MM-DD'))
const endDate = ref($dayjs().format('YYYY-MM-DD'))

const columns: TableColumn[] = [
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'journalNumber', label: 'finance.journal_number' },
  { key: 'description', label: 'common.description' },
  { key: 'type', label: 'finance.type' },
  { key: 'amount', label: 'finance.amount' },
  { key: 'balance', label: 'finance.balance' },
]
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

const fetchLedger = async () => {
  if (!selectedAccount.value) return; loading.value = true
  try { transactions.value = await api.get('/finance/ledger', { accountId: selectedAccount.value, startDate: startDate.value, endDate: endDate.value }) }
  catch {} finally { loading.value = false }
}
const exportLedger = () => { window.open(`/api/v1/finance/ledger/export?accountId=${selectedAccount.value}&startDate=${startDate.value}&endDate=${endDate.value}`, '_blank') }

const fetchAccounts = async () => {
  try { accountOptions.value = (await api.get<{id:string;name:string}[]>('/finance/accounts')).map(a => ({ label: a.name, value: a.id })); if (accountOptions.value.length > 0) { selectedAccount.value = accountOptions.value[0].value; fetchLedger() } }
  catch {}
}
onMounted(() => fetchAccounts())
</script>

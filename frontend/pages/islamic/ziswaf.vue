<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('islamic.ziswaf') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('islamic.ziswaf_subtitle') }}</p></div>
      <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('islamic.add_record') }}</UButton>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard :label="$t('islamic.zakat')" :value="formatCurrency(stats.zakat)" icon="i-heroicons-heart" color="emerald" :loading="statsLoading" />
      <StatCard :label="$t('islamic.infaq')" :value="formatCurrency(stats.infaq)" icon="i-heroicons-gift" color="blue" :loading="statsLoading" />
      <StatCard :label="$t('islamic.shadaqah')" :value="formatCurrency(stats.shadaqah)" icon="i-heroicons-hand-raised" color="amber" :loading="statsLoading" />
      <StatCard :label="$t('islamic.waqf')" :value="formatCurrency(stats.waqf)" icon="i-heroicons-building-library" color="purple" :loading="statsLoading" />
    </div>

    <UTabs :items="tabs">
      <template #item="{ item }">
        <div class="pt-4">
          <DataTable :columns="ziswafColumns" :rows="filteredRecords(item.key)" :loading="loading" :empty-title="$t('islamic.no_ziswaf_records')" :show-export="false">
            <template #cell-amount="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.amount as number) }}</span></template>
            <template #cell-type="{ row }"><StatusBadge :status="(row.type as string) || 'zakat'" /></template>
            <template #item-actions="{ row }">
              <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editRecord(row as Record<string, unknown>)" />
              <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
            </template>
          </DataTable>
        </div>
      </template>
    </UTabs>

    <FormDialog v-model="showForm" :title="editing ? $t('islamic.edit_record') : $t('islamic.add_record')" :loading="saving" @submit="saveRecord" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('islamic.type')" required><USelect v-model="form.type" :options="typeOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('finance.amount')" required><UInput v-model.number="form.amount" type="number" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.date')"><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('islamic.donor')"><USelect v-model="form.donorId" :options="donorOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('islamic.recipient')"><USelect v-model="form.recipientId" :options="recipientOptions" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.note')"><UTextarea v-model="form.note" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('islamic.delete_record')" :loading="deleting" @confirm="deleteRecord" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const statsLoading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const records = ref<Record<string, unknown>[]>([])
const donorOptions = ref<{ label: string; value: string }[]>([])
const recipientOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)
const stats = reactive({ zakat: 0, infaq: 0, shadaqah: 0, waqf: 0 })

const tabs = computed(() => [
  { key: 'all', label: t('common.all') },
  { key: 'zakat', label: t('islamic.zakat') },
  { key: 'infaq', label: t('islamic.infaq') },
  { key: 'shadaqah', label: t('islamic.shadaqah') },
  { key: 'waqf', label: t('islamic.waqf') },
])

const ziswafColumns: TableColumn[] = [
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'type', label: 'islamic.type' },
  { key: 'donorName', label: 'islamic.donor' },
  { key: 'recipientName', label: 'islamic.recipient' },
  { key: 'amount', label: 'finance.amount', type: 'currency' },
  { key: 'note', label: 'common.note' },
]
const typeOptions = [
  { label: t('islamic.zakat'), value: 'zakat' }, { label: t('islamic.infaq'), value: 'infaq' },
  { label: t('islamic.shadaqah'), value: 'shadaqah' }, { label: t('islamic.waqf'), value: 'waqf' },
]
const form = reactive({ type: 'zakat', amount: 0, date: $dayjs().format('YYYY-MM-DD'), donorId: '', recipientId: '', note: '' })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)

const filteredRecords = (type: string) => type === 'all' ? records.value : records.value.filter(r => r.type === type)

const fetchRecords = async () => { loading.value = true; try { records.value = await api.paginate('/ziswaf').then(r => r.data) } catch {} finally { loading.value = false } }
const fetchStats = async () => { statsLoading.value = true; try { Object.assign(stats, await api.get('/ziswaf/stats')) } catch {} finally { statsLoading.value = false } }
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { type: 'zakat', amount: 0, date: $dayjs().format('YYYY-MM-DD'), donorId: '', recipientId: '', note: '' }); showForm.value = true }
const editRecord = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveRecord = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/ziswaf/${editId.value}`, form); toast.add({ title: t('islamic.record_updated'), color: 'success' }) } else { await api.post('/ziswaf', form); toast.add({ title: t('islamic.record_created'), color: 'success' }) } showForm.value = false; fetchRecords(); fetchStats() } catch {} finally { saving.value = false } }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteRecord = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/ziswaf/${deleteTarget.value.id}`); toast.add({ title: t('islamic.record_deleted'), color: 'success' }); showDelete.value = false; fetchRecords(); fetchStats() } catch {} finally { deleting.value = false } }
const fetchOptions = async () => { try { donorOptions.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })) } catch {} }
onMounted(() => { fetchRecords(); fetchStats(); fetchOptions() })
</script>

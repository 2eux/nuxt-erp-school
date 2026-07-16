<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('finance.fee_types') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('finance.fee_types_subtitle') }}</p></div>
      <UButton v-if="permissions.can('finance.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('finance.add_fee_type') }}</UButton>
    </div>
    <DataTable :columns="columns" :rows="feeTypes" :loading="loading" :empty-title="$t('finance.no_fee_types')" :show-export="false">
      <template #cell-amount="{ row }"><span class="text-sm font-mono text-gray-900 dark:text-white">{{ formatCurrency(row.amount as number) }}</span></template>
      <template #item-actions="{ row }">
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editFee(row as Record<string, unknown>)" />
        <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
      </template>
    </DataTable>
    <FormDialog v-model="showForm" :title="editing ? $t('finance.edit_fee_type') : $t('finance.add_fee_type')" :loading="saving" @submit="saveFeeType" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('finance.name')" required><UInput v-model="form.name" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('finance.category')" required><USelect v-model="form.category" :options="[{label:t('finance.spp'),value:'spp'},{label:t('finance.registration'),value:'registration'},{label:t('finance.development'),value:'development'},{label:t('finance.uniform'),value:'uniform'},{label:t('finance.book'),value:'book'},{label:t('finance.exam'),value:'exam'},{label:t('finance.event'),value:'event'}]" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('finance.amount')" required><UInput v-model.number="form.amount" type="number" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('common.description')"><UTextarea v-model="form.description" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>
    <ConfirmDialog v-model="showDelete" :title="$t('finance.delete_fee_type')" :loading="deleting" @confirm="deleteFee" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast()
const loading = ref(false); const saving = ref(false); const deleting = ref(false)
const showForm = ref(false); const showDelete = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const feeTypes = ref<Record<string, unknown>[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)
const columns: TableColumn[] = [
  { key: 'name', label: 'finance.name', sortable: true },
  { key: 'category', label: 'finance.category' },
  { key: 'amount', label: 'finance.amount', type: 'currency' },
  { key: 'description', label: 'common.description' },
]
const form = reactive({ name: '', category: 'spp', amount: 0, description: '' })
const formatCurrency = (v: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)
const fetchFeeTypes = async () => { loading.value = true; try { feeTypes.value = await api.paginate('/finance/fee-types').then(r => r.data) } catch {} finally { loading.value = false } }
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { name: '', category: 'spp', amount: 0, description: '' }); showForm.value = true }
const editFee = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveFeeType = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/finance/fee-types/${editId.value}`, form); toast.add({ title: t('finance.fee_type_updated'), color: 'success' }) } else { await api.post('/finance/fee-types', form); toast.add({ title: t('finance.fee_type_created'), color: 'success' }) } showForm.value = false; fetchFeeTypes() } catch {} finally { saving.value = false } }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteFee = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/finance/fee-types/${deleteTarget.value.id}`); toast.add({ title: t('finance.fee_type_deleted'), color: 'success' }); showDelete.value = false; fetchFeeTypes() } catch {} finally { deleting.value = false } }
onMounted(() => fetchFeeTypes())
</script>

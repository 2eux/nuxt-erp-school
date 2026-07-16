<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('islamic.islamic_character') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('islamic.islamic_character_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedClass" :options="classOptions" :placeholder="$t('academic.select_class')" color="gray" size="sm" class="w-36" @change="fetchData" />
        <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('islamic.add_assessment') }}</UButton>
      </div>
    </div>

    <div class="card">
      <div class="card-header"><h3 class="card-title">{{ $t('islamic.character_assessment') }}</h3></div>
      <div class="h-72">
        <ApexChart type="radar" height="100%" :options="characterChart.options" :series="characterChart.series" />
      </div>
    </div>

    <DataTable :columns="columns" :rows="assessments" :loading="loading" :empty-title="$t('islamic.no_assessments')" :show-export="false">
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editAssessment(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('islamic.edit_assessment') : $t('islamic.add_assessment')" :loading="saving" @submit="saveAssessment" @cancel="showForm=false">
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('tahfidz.student')" required><USelect v-model="form.studentId" :options="studentOptions" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('common.date')"><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('islamic.adab')" required><div class="flex items-center gap-2"><input v-model.number="form.adab" type="range" min="1" max="5" class="w-full" /><span class="text-sm font-semibold">{{ form.adab }}/5</span></div></UFormGroup>
        <UFormGroup :label="$t('islamic.akhlaq')" required><div class="flex items-center gap-2"><input v-model.number="form.akhlaq" type="range" min="1" max="5" class="w-full" /><span class="text-sm font-semibold">{{ form.akhlaq }}/5</span></div></UFormGroup>
        <UFormGroup :label="$t('islamic.ibadah')" required><div class="flex items-center gap-2"><input v-model.number="form.ibadah" type="range" min="1" max="5" class="w-full" /><span class="text-sm font-semibold">{{ form.ibadah }}/5</span></div></UFormGroup>
        <UFormGroup :label="$t('islamic.discipline')" required><div class="flex items-center gap-2"><input v-model.number="form.discipline" type="range" min="1" max="5" class="w-full" /><span class="text-sm font-semibold">{{ form.discipline }}/5</span></div></UFormGroup>
        <UFormGroup :label="$t('islamic.social')" required><div class="flex items-center gap-2"><input v-model.number="form.social" type="range" min="1" max="5" class="w-full" /><span class="text-sm font-semibold">{{ form.social }}/5</span></div></UFormGroup>
        <UFormGroup :label="$t('common.note')"><UTextarea v-model="form.note" color="gray" :rows="3" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('islamic.delete_assessment')" :loading="deleting" @confirm="deleteItem" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const { $dayjs } = useNuxtApp()
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const showForm = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const assessments = ref<Record<string, unknown>[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const studentOptions = ref<{ label: string; value: string }[]>([])
const selectedClass = ref('')
const deleteTarget = ref<Record<string, unknown> | null>(null)

const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name' },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'adab', label: 'islamic.adab', type: 'number' },
  { key: 'akhlaq', label: 'islamic.akhlaq', type: 'number' },
  { key: 'ibadah', label: 'islamic.ibadah', type: 'number' },
  { key: 'discipline', label: 'islamic.discipline', type: 'number' },
  { key: 'social', label: 'islamic.social', type: 'number' },
]

const form = reactive({ studentId: '', date: $dayjs().format('YYYY-MM-DD'), adab: 3, akhlaq: 3, ibadah: 3, discipline: 3, social: 3, note: '' })

const characterChart = computed(() => ({
  series: [{ name: t('islamic.score'), data: [4.2, 3.8, 4.5, 3.9, 4.1] }],
  options: { chart: { type: 'radar' as const, toolbar: { show: false }, background: 'transparent' }, colors: ['#059669'], xaxis: { categories: [t('islamic.adab'), t('islamic.akhlaq'), t('islamic.ibadah'), t('islamic.discipline'), t('islamic.social')] }, yaxis: { show: false, min: 0, max: 5 } },
}))

const fetchData = async () => { if (!selectedClass.value) return; loading.value = true; try { assessments.value = await api.get('/islamic-character', { classId: selectedClass.value }) } catch {} finally { loading.value = false } }
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { studentId: '', date: $dayjs().format('YYYY-MM-DD'), adab: 3, akhlaq: 3, ibadah: 3, discipline: 3, social: 3, note: '' }); showForm.value = true }
const editAssessment = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveAssessment = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/islamic-character/${editId.value}`, form); toast.add({ title: t('islamic.assessment_updated'), color: 'success' }) } else { await api.post('/islamic-character', form); toast.add({ title: t('islamic.assessment_created'), color: 'success' }) } showForm.value = false; fetchData() } catch {} finally { saving.value = false } }
const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteItem = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/islamic-character/${deleteTarget.value.id}`); toast.add({ title: t('islamic.assessment_deleted'), color: 'success' }); showDelete.value = false; fetchData() } catch {} finally { deleting.value = false } }
const fetchOptions = async () => { try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })) } catch {}; try { studentOptions.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })) } catch {} }

onMounted(() => fetchOptions())
</script>

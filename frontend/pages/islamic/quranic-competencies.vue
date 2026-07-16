<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('islamic.quranic_competencies') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('islamic.quranic_competencies_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedClass" :options="classOptions" :placeholder="$t('academic.select_class')" color="gray" size="sm" class="w-36" @change="fetchData" />
        <UButton v-if="permissions.can('islamic.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('islamic.add_score') }}</UButton>
      </div>
    </div>

    <DataTable :columns="columns" :rows="records" :loading="loading" :empty-title="$t('islamic.no_competencies')" :show-export="true" @export="exportData">
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editRecord(row as Record<string, unknown>)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showForm" :title="editing ? $t('islamic.edit_score') : $t('islamic.add_score')" :loading="saving" @submit="saveRecord" @cancel="showForm=false">
      <div class="space-y-4">
        <UFormGroup :label="$t('tahfidz.student')" required><USelect v-model="form.studentId" :options="studentOptions" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('islamic.tajwid')"><UInput v-model.number="form.tajwid" type="number" min="0" max="100" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('islamic.tilawah')"><UInput v-model.number="form.tilawah" type="number" min="0" max="100" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('islamic.tahfidz_score')"><UInput v-model.number="form.tahfidz" type="number" min="0" max="100" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('islamic.tafsir')"><UInput v-model.number="form.tafsir" type="number" min="0" max="100" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('islamic.date')"><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.note')"><UTextarea v-model="form.note" color="gray" :rows="2" /></UFormGroup>
      </div>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n(); const api = useApi(); const permissions = usePermissions(); const toast = useToast(); const { $dayjs } = useNuxtApp()
const loading = ref(false); const saving = ref(false); const showForm = ref(false)
const editing = ref(false); const editId = ref<string | null>(null)
const records = ref<Record<string, unknown>[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const studentOptions = ref<{ label: string; value: string }[]>([])
const selectedClass = ref('')

const columns: TableColumn[] = [
  { key: 'studentName', label: 'students.name' },
  { key: 'tajwid', label: 'islamic.tajwid', type: 'number' },
  { key: 'tilawah', label: 'islamic.tilawah', type: 'number' },
  { key: 'tahfidz', label: 'islamic.tahfidz_score', type: 'number' },
  { key: 'tafsir', label: 'islamic.tafsir', type: 'number' },
  { key: 'average', label: 'academic.average', type: 'number' },
  { key: 'grade', label: 'academic.grade' },
]
const form = reactive({ studentId: '', tajwid: 0, tilawah: 0, tahfidz: 0, tafsir: 0, date: $dayjs().format('YYYY-MM-DD'), note: '' })

const fetchData = async () => { if (!selectedClass.value) return; loading.value = true; try { records.value = await api.get('/quranic-competencies', { classId: selectedClass.value }) } catch {} finally { loading.value = false } }
const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { studentId: '', tajwid: 0, tilawah: 0, tahfidz: 0, tafsir: 0, date: $dayjs().format('YYYY-MM-DD'), note: '' }); showForm.value = true }
const editRecord = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showForm.value = true }
const saveRecord = async () => { saving.value = true; try { if (editing.value && editId.value) { await api.put(`/quranic-competencies/${editId.value}`, form); toast.add({ title: t('islamic.score_updated'), color: 'success' }) } else { await api.post('/quranic-competencies', form); toast.add({ title: t('islamic.score_created'), color: 'success' }) } showForm.value = false; fetchData() } catch {} finally { saving.value = false } }
const exportData = () => { /* export CSV */ }
const fetchOptions = async () => { try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })) } catch {} try { studentOptions.value = (await api.get<{id:string;fullName:string}[]>('/students')).map(s => ({ label: s.fullName, value: s.id })) } catch {} }
onMounted(() => fetchOptions())
</script>

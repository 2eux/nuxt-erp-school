<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.lesson_plans') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.lesson_plans_subtitle') }}</p></div>
      <UButton v-if="permissions.can('academic.lesson_plans.manage')" color="primary" size="sm" icon="i-heroicons-plus" @click="openAdd">{{ $t('academic.create_lesson_plan') }}</UButton>
    </div>

    <DataFilter :filter-fields="filterFields" @apply="handleFilter" />

    <DataTable :columns="columns" :rows="lessonPlans" :loading="loading" :empty-title="$t('academic.no_lesson_plans')" :show-export="false">
      <template #cell-status="{ row }"><StatusBadge :status="(row.status as string) || 'draft'" /></template>
      <template #item-actions="{ row }">
        <div class="flex items-center justify-end gap-1">
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-eye" @click="viewLessonPlan(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-pencil-square" @click="editLessonPlan(row as Record<string, unknown>)" />
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-trash" @click="confirmDelete(row)" />
        </div>
      </template>
    </DataTable>

    <FormDialog v-model="showDialog" :title="editing ? $t('academic.edit_lesson_plan') : $t('academic.create_lesson_plan')" :loading="saving" @submit="saveLessonPlan" @cancel="closeDialog">
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('subjects.title')" required><USelect v-model="form.subjectId" :options="subjectOptions" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('academic.class')" required><USelect v-model="form.classId" :options="classOptions" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('academic.topic')" required><UInput v-model="form.topic" color="gray" /></UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="$t('academic.week')"><UInput v-model.number="form.week" type="number" color="gray" /></UFormGroup>
          <UFormGroup :label="$t('academic.date')"><UInput v-model="form.date" type="date" color="gray" /></UFormGroup>
        </div>
        <UFormGroup :label="$t('academic.objectives')" required><RichEditor v-model="form.objectives" /></UFormGroup>
        <UFormGroup :label="$t('academic.materials')" required><RichEditor v-model="form.materials" /></UFormGroup>
        <UFormGroup :label="$t('academic.activities')" required><RichEditor v-model="form.activities" /></UFormGroup>
        <UFormGroup :label="$t('academic.assessment')"><RichEditor v-model="form.assessment" /></UFormGroup>
        <UFormGroup :label="$t('academic.media')"><UInput v-model="form.media" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('academic.source')"><UInput v-model="form.source" color="gray" /></UFormGroup>
        <UFormGroup :label="$t('common.status')"><USelect v-model="form.status" :options="[{label:t('status.draft'),value:'draft'},{label:t('status.approved'),value:'approved'}]" color="gray" /></UFormGroup>
      </div>
    </FormDialog>

    <ConfirmDialog v-model="showDelete" :title="$t('academic.delete_lesson_plan')" :loading="deleting" @confirm="deleteItem" @cancel="showDelete=false" />
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '~/types'
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const permissions = usePermissions()
const toast = useToast()
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const showDialog = ref(false)
const showDelete = ref(false)
const editing = ref(false)
const editId = ref<string | null>(null)
const lessonPlans = ref<Record<string, unknown>[]>([])
const classOptions = ref<{ label: string; value: string }[]>([])
const subjectOptions = ref<{ label: string; value: string }[]>([])
const deleteTarget = ref<Record<string, unknown> | null>(null)

const columns: TableColumn[] = [
  { key: 'topic', label: 'academic.topic', sortable: true },
  { key: 'subjectName', label: 'subjects.title' },
  { key: 'className', label: 'academic.class' },
  { key: 'week', label: 'academic.week', type: 'number' },
  { key: 'date', label: 'common.date', type: 'date' },
  { key: 'status', label: 'common.status', type: 'status' },
]
const filterFields = [
  { key: 'classId', label: 'academic.class', type: 'select' as const, options: [] },
  { key: 'subjectId', label: 'subjects.title', type: 'select' as const, options: [] },
  { key: 'status', label: 'common.status', type: 'select' as const, options: [{ label: t('status.draft'), value: 'draft' }, { label: t('status.approved'), value: 'approved' }] },
]

const form = reactive({ topic: '', subjectId: '', classId: '', week: 1, date: '', objectives: '', materials: '', activities: '', assessment: '', media: '', source: '', status: 'draft' })

const fetchLessonPlans = async (filters: Record<string, unknown> = {}) => { loading.value = true; try { lessonPlans.value = await api.paginate('/lesson-plans', { ...filters }).then(r => r.data) } catch {} finally { loading.value = false } }
const handleFilter = (filters: Record<string, unknown>) => fetchLessonPlans(filters)

const openAdd = () => { editing.value = false; editId.value = null; Object.assign(form, { topic: '', subjectId: '', classId: '', week: 1, date: '', objectives: '', materials: '', activities: '', assessment: '', media: '', source: '', status: 'draft' }); showDialog.value = true }
const editLessonPlan = (row: Record<string, unknown>) => { editing.value = true; editId.value = row.id as string; Object.assign(form, row); showDialog.value = true }
const closeDialog = () => { showDialog.value = false; editing.value = false }
const viewLessonPlan = (row: Record<string, unknown>) => { navigateTo(`/academic/lesson-plans/${row.id}`) }

const saveLessonPlan = async () => {
  saving.value = true
  try {
    if (editing.value && editId.value) { await api.put(`/lesson-plans/${editId.value}`, form); toast.add({ title: t('academic.lesson_plan_updated'), color: 'success' }) }
    else { await api.post('/lesson-plans', form); toast.add({ title: t('academic.lesson_plan_created'), color: 'success' }) }
    closeDialog(); fetchLessonPlans()
  } catch {} finally { saving.value = false }
}

const confirmDelete = (row: Record<string, unknown>) => { deleteTarget.value = row; showDelete.value = true }
const deleteItem = async () => { if (!deleteTarget.value) return; deleting.value = true; try { await api.delete(`/lesson-plans/${deleteTarget.value.id}`); toast.add({ title: t('academic.lesson_plan_deleted'), color: 'success' }); showDelete.value = false; fetchLessonPlans() } catch {} finally { deleting.value = false } }

const fetchOptions = async () => {
  try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })); filterFields[0].options = classOptions.value } catch {}
  try { subjectOptions.value = (await api.get<{id:string;name:string}[]>('/subjects')).map(s => ({ label: s.name, value: s.id })); filterFields[1].options = subjectOptions.value } catch {}
}

onMounted(() => { fetchLessonPlans(); fetchOptions() })
</script>

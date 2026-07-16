<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div><h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ $t('academic.gradebook') }}</h1><p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ $t('academic.gradebook_subtitle') }}</p></div>
      <div class="flex items-center gap-2">
        <USelect v-model="selectedClass" :options="classOptions" :placeholder="$t('academic.select_class')" color="gray" size="sm" class="w-36" @change="fetchGradebook" />
        <USelect v-model="selectedSubject" :options="subjectOptions" :placeholder="$t('subjects.title')" color="gray" size="sm" class="w-36" @change="fetchGradebook" />
        <UButton color="gray" variant="outline" size="sm" icon="i-heroicons-arrow-down-tray" @click="exportGradebook">{{ $t('common.export') }}</UButton>
      </div>
    </div>

    <div class="card overflow-hidden !p-0">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead><tr class="bg-gray-50 dark:bg-gray-800/50">
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">No</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">{{ $t('students.name') }}</th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">NIS</th>
            <th v-for="comp in scoreComponents" :key="comp.key" class="px-4 py-3 text-center text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">
              {{ comp.label }}<br/><span class="font-normal">({{ comp.weight }}%)</span>
            </th>
            <th class="px-4 py-3 text-center text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">{{ $t('academic.total') }}</th>
            <th class="px-4 py-3 text-center text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">{{ $t('academic.grade') }}</th>
          </tr></thead>
          <tbody>
            <template v-if="loading">
              <tr v-for="i in 8" :key="'sk-'+i"><td class="px-4 py-3" :colspan="5 + scoreComponents.length"><div class="h-4 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" /></td></tr>
            </template>
            <template v-else-if="gradebookData.length === 0">
              <tr><td :colspan="5 + scoreComponents.length" class="px-6 py-16"><EmptyState :title="$t('academic.no_gradebook_data')" icon="i-heroicons-table-cells" /></td></tr>
            </template>
            <template v-else>
              <tr v-for="(student, idx) in gradebookData" :key="student.id" class="border-t border-gray-100 dark:border-gray-750 hover:bg-gray-50 dark:hover:bg-gray-800/50">
                <td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">{{ idx + 1 }}</td>
                <td class="px-4 py-2 text-sm font-medium text-gray-900 dark:text-white">{{ student.studentName }}</td>
                <td class="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">{{ student.nis }}</td>
                <td v-for="comp in scoreComponents" :key="comp.key" class="px-2 py-2 text-center">
                  <UInput
                    v-model.number="student.scores[comp.key]"
                    type="number"
                    size="xs"
                    color="gray"
                    class="w-16 text-center mx-auto"
                    :min="0"
                    :max="100"
                    @blur="calculateGrade(student)"
                  />
                </td>
                <td class="px-4 py-2 text-center text-sm font-semibold text-gray-900 dark:text-white">{{ student.total }}</td>
                <td class="px-4 py-2 text-center">
                  <span class="inline-flex px-2 py-0.5 rounded-full text-xs font-bold"
                    :class="gradeColor(student.grade)">{{ student.grade }}</span>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <div class="flex justify-end gap-3" v-if="gradebookData.length > 0">
      <UButton color="gray" variant="ghost" @click="resetChanges">{{ $t('common.cancel') }}</UButton>
      <UButton color="primary" :loading="saving" @click="saveGradebook">{{ $t('common.save_grades') }}</UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
const { t } = useI18n()
const api = useApi()
const toast = useToast()
const loading = ref(false)
const saving = ref(false)
const selectedClass = ref('')
const selectedSubject = ref('')
const classOptions = ref<{ label: string; value: string }[]>([])
const subjectOptions = ref<{ label: string; value: string }[]>([])
const gradebookData = ref<{ id: string; studentId: string; studentName: string; nis: string; scores: Record<string, number>; total: number; grade: string }[]>([])
const originalData = ref<typeof gradebookData.value>([])

const scoreComponents = [
  { key: 'daily', label: t('grades.daily'), weight: 20 },
  { key: 'assignment', label: t('grades.assignment'), weight: 20 },
  { key: 'midterm', label: t('grades.midterm'), weight: 25 },
  { key: 'final', label: t('grades.final'), weight: 35 },
]

const calculateGrade = (student: typeof gradebookData.value[0]) => {
  let total = 0
  for (const comp of scoreComponents) {
    total += (student.scores[comp.key] || 0) * comp.weight / 100
  }
  student.total = Math.round(total)
  student.grade = getGradeLetter(student.total)
}

const getGradeLetter = (score: number): string => {
  if (score >= 90) return 'A'
  if (score >= 80) return 'B'
  if (score >= 70) return 'C'
  if (score >= 60) return 'D'
  return 'E'
}

const gradeColor = (grade: string): string => {
  const map: Record<string, string> = { A: 'bg-emerald-100 text-emerald-800', B: 'bg-blue-100 text-blue-800', C: 'bg-amber-100 text-amber-800', D: 'bg-orange-100 text-orange-800', E: 'bg-red-100 text-red-800' }
  return map[grade] || 'bg-gray-100 text-gray-800'
}

const fetchGradebook = async () => {
  if (!selectedClass.value || !selectedSubject.value) return
  loading.value = true
  try {
    const res = await api.get<{ id: string; studentId: string; studentName: string; nis: string; scores: Record<string, number> }[]>('/gradebook', { classId: selectedClass.value, subjectId: selectedSubject.value })
    gradebookData.value = res.map(r => ({
      ...r,
      scores: r.scores || {},
      total: 0,
      grade: 'E',
    }))
    gradebookData.value.forEach(s => calculateGrade(s))
    originalData.value = JSON.parse(JSON.stringify(gradebookData.value))
  } catch {} finally { loading.value = false }
}

const resetChanges = () => {
  gradebookData.value = JSON.parse(JSON.stringify(originalData.value))
}

const saveGradebook = async () => {
  saving.value = true
  try {
    const grades = gradebookData.value.map(s => ({
      studentId: s.studentId,
      subjectId: selectedSubject.value,
      classId: selectedClass.value,
      scores: s.scores,
      total: s.total,
      grade: s.grade,
    }))
    await api.post('/gradebook', { grades })
    toast.add({ title: t('academic.grades_saved'), color: 'success' })
  } catch {} finally { saving.value = false }
}

const exportGradebook = () => {
  const headers = ['No', 'Name', 'NIS', ...scoreComponents.map(c => c.label), 'Total', 'Grade']
  const rows = gradebookData.value.map((s, i) => [i + 1, s.studentName, s.nis, ...scoreComponents.map(c => s.scores[c.key] || 0), s.total, s.grade])
  const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a'); a.href = url; a.download = `gradebook_${selectedClass.value}_${selectedSubject.value}.csv`; a.click()
  URL.revokeObjectURL(url)
}

const fetchOptions = async () => {
  try { classOptions.value = (await api.get<{id:string;name:string}[]>('/classes')).map(c => ({ label: c.name, value: c.id })) } catch {}
  try { subjectOptions.value = (await api.get<{id:string;name:string}[]>('/subjects')).map(s => ({ label: s.name, value: s.id })) } catch {}
}

onMounted(() => fetchOptions())
</script>

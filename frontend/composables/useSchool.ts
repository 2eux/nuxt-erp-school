import type { School, AcademicYear } from '~/types'

export const useSchool = () => {
  const schoolStore = useSchoolStore()
  const api = useApi()

  const currentSchool = computed(() => schoolStore.currentSchool)
  const currentAcademicYear = computed(() => schoolStore.currentAcademicYear)
  const academicYears = computed(() => schoolStore.academicYears)
  const schools = computed(() => schoolStore.schools)

  const fetchSchools = async (): Promise<void> => {
    try {
      const response = await api.get<School[]>('/schools')
      schoolStore.setSchools(response)
    } catch {
      // handled by useApi
    }
  }

  const setCurrentSchool = (school: School): void => {
    schoolStore.setCurrentSchool(school)
    if (import.meta.client) {
      localStorage.setItem('current_school_id', school.id)
    }
  }

  const fetchAcademicYears = async (schoolId: string): Promise<void> => {
    try {
      const response = await api.get<AcademicYear[]>(`/schools/${schoolId}/academic-years`)
      schoolStore.setAcademicYears(response)

      const activeYear = response.find(y => y.isActive)
      if (activeYear) {
        schoolStore.setCurrentAcademicYear(activeYear)
      }
    } catch {
      // handled by useApi
    }
  }

  const setCurrentAcademicYear = (year: AcademicYear): void => {
    schoolStore.setCurrentAcademicYear(year)
    if (import.meta.client) {
      localStorage.setItem('current_academic_year_id', year.id)
    }
  }

  const initializeSchoolContext = async (): Promise<void> => {
    const savedSchoolId = import.meta.client ? localStorage.getItem('current_school_id') : null
    const savedYearId = import.meta.client ? localStorage.getItem('current_academic_year_id') : null

    if (schools.value.length === 0) {
      await fetchSchools()
    }

    if (savedSchoolId) {
      const school = schools.value.find(s => s.id === savedSchoolId)
      if (school) {
        schoolStore.setCurrentSchool(school)
        await fetchAcademicYears(school.id)
      }
    } else if (schools.value.length > 0) {
      schoolStore.setCurrentSchool(schools.value[0])
      await fetchAcademicYears(schools.value[0].id)
    }

    if (savedYearId && academicYears.value.length > 0) {
      const year = academicYears.value.find(y => y.id === savedYearId)
      if (year) {
        schoolStore.setCurrentAcademicYear(year)
      }
    }
  }

  const schoolOptions = computed(() =>
    schools.value.map(s => ({ label: s.name, value: s.id }))
  )

  const academicYearOptions = computed(() =>
    academicYears.value.map(y => ({ label: y.name, value: y.id }))
  )

  const isSchoolContextReady = computed(() =>
    !!currentSchool.value && !!currentAcademicYear.value
  )

  return {
    currentSchool,
    currentAcademicYear,
    academicYears,
    schools,
    fetchSchools,
    setCurrentSchool,
    fetchAcademicYears,
    setCurrentAcademicYear,
    initializeSchoolContext,
    schoolOptions,
    academicYearOptions,
    isSchoolContextReady,
  }
}

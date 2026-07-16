import { defineStore } from 'pinia'
import type { School, AcademicYear } from '~/types'

interface SchoolState {
  currentSchool: School | null
  currentAcademicYear: AcademicYear | null
  schools: School[]
  academicYears: AcademicYear[]
}

export const useSchoolStore = defineStore('school', {
  state: (): SchoolState => ({
    currentSchool: null,
    currentAcademicYear: null,
    schools: [],
    academicYears: [],
  }),

  getters: {
    schoolId: (state) => state.currentSchool?.id || null,
    academicYearId: (state) => state.currentAcademicYear?.id || null,
    schoolName: (state) => state.currentSchool?.name || '',
    academicYearName: (state) => state.currentAcademicYear?.name || '',
    activeAcademicYears: (state) => state.academicYears.filter(y => y.isActive),
    schoolCount: (state) => state.schools.length,
  },

  actions: {
    setSchools(schools: School[]) {
      this.schools = schools
    },

    setCurrentSchool(school: School) {
      this.currentSchool = school
    },

    setAcademicYears(years: AcademicYear[]) {
      this.academicYears = years
    },

    setCurrentAcademicYear(year: AcademicYear) {
      this.currentAcademicYear = year
    },

    clearSchool() {
      this.currentSchool = null
      this.currentAcademicYear = null
      this.schools = []
      this.academicYears = []
    },
  },
})

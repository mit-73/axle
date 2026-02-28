import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Project as ProtoProject } from '@axle/contracts/bff/v1/projects_pb'
import { projectsClient } from '@/api/bff'

// Local projection of the proto Project message for reactive state.
export interface Project {
  id: string
  name: string
  description: string
  status: string
}

export const useProjectsStore = defineStore('projects', () => {
  const projects = ref<Project[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchProjects(page = 1, pageSize = 20): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const res = await projectsClient.listProjects({ page, pageSize })
      projects.value = res.projects.map((p: ProtoProject) => ({
        id: p.id,
        name: p.name,
        description: p.description,
        status: String(p.status),
      }))
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  return { projects, loading, error, fetchProjects }
})

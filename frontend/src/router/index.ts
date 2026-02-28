import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // Root → redirect to projects
    { path: '/', redirect: '/projects' },

    // Projects
    {
      path: '/projects',
      component: () => import('@/views/projects/ProjectList.vue'),
      name: 'project-list',
    },
    {
      path: '/projects/new',
      component: () => import('@/views/projects/ProjectCreate.vue'),
      name: 'project-create',
    },
    {
      path: '/projects/:id',
      component: () => import('@/views/projects/ProjectDetail.vue'),
      name: 'project-detail',
    },
    {
      path: '/projects/:id/board',
      component: () => import('@/views/projects/KanbanBoard.vue'),
      name: 'kanban-board',
    },
    {
      path: '/projects/:id/backlog',
      component: () => import('@/views/projects/Backlog.vue'),
      name: 'backlog',
    },
    {
      path: '/projects/:id/epics',
      component: () => import('@/views/projects/EpicList.vue'),
      name: 'epic-list',
    },
    {
      path: '/projects/:id/insights',
      component: () => import('@/views/projects/Insights.vue'),
      name: 'insights',
    },

    // Settings
    {
      path: '/settings',
      component: () => import('@/views/settings/Settings.vue'),
      name: 'settings',
    },
    {
      path: '/settings/ai-providers',
      component: () => import('@/views/settings/AIProviders.vue'),
      name: 'ai-providers',
    },
  ],
})

export default router

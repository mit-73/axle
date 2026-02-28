<script setup lang="ts">
import { h, computed } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NLayout, NLayoutSider, NLayoutContent, NMenu } from 'naive-ui'
import type { MenuOption } from 'naive-ui'

const router = useRouter()
const { t } = useI18n()

const menuOptions: MenuOption[] = [
  {
    label: () => h(RouterLink, { to: '/projects' }, () => t('nav.projects')),
    key: 'projects',
  },
  {
    label: () => h(RouterLink, { to: '/settings' }, () => t('nav.settings')),
    key: 'settings',
    children: [
      {
        label: () => h(RouterLink, { to: '/settings/ai-providers' }, () => t('nav.aiProviders')),
        key: 'ai-providers',
      },
    ],
  },
]

const activeKey = computed(() => {
  const path = router.currentRoute.value.path
  if (path.startsWith('/projects')) return 'projects'
  if (path.startsWith('/settings/ai-providers')) return 'ai-providers'
  if (path.startsWith('/settings')) return 'settings'
  return ''
})
</script>

<template>
  <n-layout style="height: 100vh" has-sider>
    <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="220" show-trigger>
      <div style="padding: 16px 20px; font-weight: 700; font-size: 18px">Axle</div>
      <n-menu :options="menuOptions" :value="activeKey" />
    </n-layout-sider>

    <n-layout-content content-style="padding: 24px;">
      <RouterView />
    </n-layout-content>
  </n-layout>
</template>

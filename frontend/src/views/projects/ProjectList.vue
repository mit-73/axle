<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <div>
    <h1>{{ t('pages.projectList') }}</h1>
    <!-- TODO: render project list from useProjectsStore -->
  </div>
</template>

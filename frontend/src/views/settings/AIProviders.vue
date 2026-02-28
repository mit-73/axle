<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <div>
    <h1>{{ t('pages.aiProviders') }}</h1>
    <!-- TODO: LiteLLM provider configuration -->
  </div>
</template>

<template>
  <Toolbar class="app-nav">
    <template #start>
      <div class="nav-brand">
        <img src="@/assets/logo.png" alt="SmiteTrainer logo" class="nav-logo" />
        <span class="nav-title">SmiteTrainer</span>
      </div>
    </template>
  </Toolbar>
</template>

<script lang="ts">
  import { Options, Vue } from 'vue-class-component';
  import Toolbar from 'primevue/toolbar';

  @Options({
    components: { Toolbar },
  })
  export default class AppNavbar extends Vue {}
</script>

<style scoped>
  :deep(.p-toolbar) {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
    height: 64px;
    padding: 0 2rem;
    background: #0a1628 !important;
    border: none !important;
    border-bottom: 1px solid rgba(77, 159, 255, 0.15) !important;
    border-radius: 0 !important;
    box-shadow: 0 2px 16px rgba(0, 0, 0, 0.4) !important;
  }

  .nav-brand {
    display: flex;
    align-items: center;
    gap: 0.65rem;
  }

  .nav-logo {
    width: 38px;
    height: 38px;
    object-fit: contain;
    mix-blend-mode: lighten;
    filter: drop-shadow(0 0 6px rgba(77, 159, 255, 0.6));
  }

  .nav-title {
    font-weight: 700;
    font-size: 1.25rem;
    letter-spacing: 0.12em;
    color: #ffffff;
    text-shadow: 0 0 16px rgba(77, 159, 255, 0.4);
    user-select: none;
  }
</style>

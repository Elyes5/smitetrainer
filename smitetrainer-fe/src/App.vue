<template>
  <div id="app-root">
    <AppNavbar />
    <main class="app-main">
      <router-view />
    </main>
  </div>
</template>

<script lang="ts">
  import { Options, Vue } from 'vue-class-component';
  import AppNavbar from '@/components/AppNavbar.vue';

  @Options({
    components: { AppNavbar },
  })
  export default class App extends Vue {}
</script>

<style>
  #app-root {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: var(--color-bg);
  }
</style>

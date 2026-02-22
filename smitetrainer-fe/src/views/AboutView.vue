<template>
  <div class="about-page">
    <div class="about-shell">
      <h1>About SmiteTrainer</h1>
      <p>
        SmiteTrainer is a reaction-focused practice tool for Smite timing and dragon secure decision
        windows.
      </p>
      <p>
        The backend uses Riot Match-V5 timeline data to estimate dragon HP over time and can be
        queried from the Training page.
      </p>
      <div class="about-actions">
        <button class="about-btn" @click="$router.push('/training')">Open Training</button>
        <button class="about-btn about-btn--ghost" @click="$router.push('/')">Back Home</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  import { defineComponent } from 'vue';

  export default defineComponent({
    name: 'AboutView',
  });
</script>

<style scoped>
  .about-page {
    min-height: 100vh;
    padding: 6rem 1.2rem 2rem;
    display: grid;
    place-items: center;
    background: radial-gradient(900px 400px at 0% 100%, rgba(77, 159, 255, 0.1), transparent 65%),
      radial-gradient(700px 350px at 100% 0%, rgba(0, 212, 255, 0.08), transparent 60%), #050d1a;
  }

  .about-shell {
    width: min(760px, 100%);
    background: rgba(10, 22, 40, 0.8);
    border: 1px solid rgba(77, 159, 255, 0.2);
    border-radius: 14px;
    padding: 1.4rem;
  }

  h1 {
    color: #fff;
    letter-spacing: 0.08em;
    margin-bottom: 0.9rem;
  }

  p {
    color: var(--color-text-secondary);
    line-height: 1.6;
    margin-bottom: 0.7rem;
  }

  .about-actions {
    margin-top: 1.2rem;
    display: flex;
    gap: 0.7rem;
  }

  .about-btn {
    border: 1px solid #fff;
    background: #fff;
    color: #0f2040;
    border-radius: 8px;
    height: 2.5rem;
    padding: 0 0.9rem;
    cursor: pointer;
    font-weight: 700;
    letter-spacing: 0.06em;
  }

  .about-btn--ghost {
    border-color: rgba(77, 159, 255, 0.3);
    background: transparent;
    color: var(--color-blue-bright);
  }

  @media (max-width: 560px) {
    .about-actions {
      flex-direction: column;
    }
  }
</style>

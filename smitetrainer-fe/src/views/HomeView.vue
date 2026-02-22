<template>
  <div class="home-container">
    <div class="bg-orb bg-orb--1"></div>
    <div class="bg-orb bg-orb--2"></div>
    <div class="bg-orb bg-orb--3"></div>

    <div class="content-wrapper">
      <!-- Logo -->
      <div class="hero-section">
        <img src="@/assets/logo.png" alt="SmiteTrainer" class="logo-img" />
        <div class="title-divider">
          <span class="divider-line"></span>
          <span class="divider-dot"></span>
          <span class="divider-line"></span>
        </div>
      </div>

      <!-- Navigation Buttons -->
      <div class="button-group">
        <AppButton
          label="MULTIPLAYER"
          variant="primary"
          icon="⊹"
          :show-arrow="true"
          @click="onMultiplayer"
        />
        <AppButton
          label="TRAINING"
          variant="primary"
          icon="◎"
          :show-arrow="true"
          @click="onTraining"
        />
        <AppButton label="ABOUT" variant="primary" icon="◈" :show-arrow="true" @click="onAbout" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  import { Options, Vue } from 'vue-class-component';
  import AppButton from '@/components/AppButton.vue';

  @Options({
    components: { AppButton },
  })
  export default class HomeView extends Vue {
    onMultiplayer() {
      this.$router.push('/multiplayer');
    }
    onTraining() {
      this.$router.push('/training');
    }
    onAbout() {
      this.$router.push('/about');
    }
  }
</script>

<style scoped>
  .home-container {
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: var(--color-bg);
    overflow: hidden;
  }

  /* ─── Orbs ─── */
  .bg-orb {
    position: absolute;
    border-radius: 50%;
    filter: blur(90px);
    animation: drift 14s ease-in-out infinite;
    pointer-events: none;
  }

  .bg-orb--1 {
    width: 520px;
    height: 520px;
    background: radial-gradient(circle, rgba(26, 58, 110, 0.55) 0%, transparent 70%);
    top: -160px;
    right: -120px;
    animation-delay: 0s;
  }

  .bg-orb--2 {
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, rgba(77, 159, 255, 0.1) 0%, transparent 70%);
    bottom: -100px;
    left: -80px;
    animation-delay: -5s;
  }

  .bg-orb--3 {
    width: 280px;
    height: 280px;
    background: radial-gradient(circle, rgba(0, 212, 255, 0.07) 0%, transparent 70%);
    top: 55%;
    left: 50%;
    transform: translate(-50%, -50%);
    animation-delay: -9s;
  }

  @keyframes drift {
    0%,
    100% {
      transform: translate(0, 0) scale(1);
    }
    33% {
      transform: translate(18px, -28px) scale(1.04);
    }
    66% {
      transform: translate(-14px, 14px) scale(0.97);
    }
  }

  /* ─── Grid ─── */
  .grid-overlay {
    position: absolute;
    inset: 0;
    background-image: linear-gradient(rgba(77, 159, 255, 0.035) 1px, transparent 1px),
      linear-gradient(90deg, rgba(77, 159, 255, 0.035) 1px, transparent 1px);
    background-size: 44px 44px;
    pointer-events: none;
  }

  /* ─── Content ─── */
  .content-wrapper {
    position: relative;
    z-index: 10;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-2xl);
    padding: var(--space-2xl) var(--space-xl);
    width: 100%;
    max-width: 420px;
  }

  /* ─── Hero / Logo ─── */
  .hero-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-sm);
    animation: fadeDown 0.7s ease both;
  }

  @keyframes fadeDown {
    from {
      opacity: 0;
      transform: translateY(-18px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .logo-img {
    /* The logo has its own black background so we blend it into the dark bg */
    width: 260px;
    height: 260px;
    object-fit: contain;
    /* mix-blend-mode: lighten makes the black parts of the logo transparent */
    mix-blend-mode: lighten;
    filter: drop-shadow(0 0 24px rgba(77, 159, 255, 0.5))
      drop-shadow(0 0 48px rgba(0, 212, 255, 0.25));
    animation: logoPulse 4s ease-in-out infinite;
  }

  @keyframes logoPulse {
    0%,
    100% {
      filter: drop-shadow(0 0 24px rgba(77, 159, 255, 0.5))
        drop-shadow(0 0 48px rgba(0, 212, 255, 0.25));
    }
    50% {
      filter: drop-shadow(0 0 32px rgba(77, 159, 255, 0.75))
        drop-shadow(0 0 64px rgba(0, 212, 255, 0.4));
    }
  }

  .title-divider {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    width: 100%;
    margin-top: var(--space-sm);
  }

  .divider-line {
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(77, 159, 255, 0.35), transparent);
  }

  .divider-dot {
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background: var(--color-blue-bright);
    box-shadow: 0 0 8px var(--color-blue-bright);
  }

  /* ─── Buttons ─── */
  .button-group {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    width: 100%;
    animation: fadeUp 0.7s ease 0.15s both;
  }

  @keyframes fadeUp {
    from {
      opacity: 0;
      transform: translateY(18px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>

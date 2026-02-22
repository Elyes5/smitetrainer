<template>
  <button
    :class="['app-btn', `app-btn--${variant}`, { 'app-btn--full': fullWidth }]"
    :disabled="disabled"
    v-bind="$attrs"
    @click="$emit('click', $event)"
  >
    <span v-if="icon" class="app-btn__icon">{{ icon }}</span>
    <span class="app-btn__label">
      <slot>{{ label }}</slot>
    </span>
    <span v-if="showArrow" class="app-btn__arrow">→</span>
  </button>
</template>

<script lang="ts">
  import { defineComponent } from 'vue';

  export default defineComponent({
    name: 'AppButton',

    inheritAttrs: false,

    props: {
      label: {
        type: String,
        default: '',
      },
      variant: {
        type: String as () => 'primary' | 'secondary' | 'ghost' | 'danger',
        default: 'primary',
        validator: (v: string) => ['primary', 'secondary', 'ghost', 'danger'].includes(v),
      },
      icon: {
        type: String,
        default: '',
      },
      showArrow: {
        type: Boolean,
        default: true,
      },
      fullWidth: {
        type: Boolean,
        default: true,
      },
      disabled: {
        type: Boolean,
        default: false,
      },
    },

    emits: ['click'],
  });
</script>

<style scoped>
  /* ─── Base ─── */
  .app-btn {
    position: relative;
    display: inline-flex;
    align-items: center;
    gap: 0.85rem;
    padding: 0.95rem 1.5rem;
    border: 1px solid transparent;
    border-radius: var(--radius-md, 10px);
    background: transparent;
    cursor: pointer;
    font-weight: 700;
    font-size: 1rem;
    letter-spacing: 0.12em;
    text-align: left;
    overflow: hidden;
    transition: background 0.22s ease, border-color 0.22s ease, box-shadow 0.22s ease,
      color 0.22s ease, transform 0.12s ease;
  }

  .app-btn--full {
    width: 100%;
  }

  .app-btn:disabled {
    opacity: 0.38;
    cursor: not-allowed;
    pointer-events: none;
  }

  /* Shimmer sweep */
  .app-btn::after {
    content: '';
    position: absolute;
    top: 0;
    left: -120%;
    width: 60%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.18), transparent);
    transition: left 0.5s ease;
    pointer-events: none;
  }

  .app-btn:not(:disabled):hover::after {
    left: 160%;
  }

  .app-btn:not(:disabled):active {
    transform: scale(0.975);
  }

  /* ─── Inner elements ─── */
  .app-btn__icon {
    flex-shrink: 0;
    font-size: 1.05rem;
    transition: transform 0.22s ease;
  }

  .app-btn:not(:disabled):hover .app-btn__icon {
    transform: rotate(25deg) scale(1.15);
  }

  .app-btn__label {
    flex: 1;
  }

  .app-btn__arrow {
    margin-left: auto;
    opacity: 0;
    transform: translateX(-10px);
    transition: opacity 0.22s ease, transform 0.22s ease;
    font-size: 0.95rem;
  }

  .app-btn:not(:disabled):hover .app-btn__arrow {
    opacity: 1;
    transform: translateX(0);
  }

  /* ─── Variant: primary ─── */
  /* White background, blue text — high-emphasis action */
  .app-btn--primary {
    background: #ffffff;
    border-color: #ffffff;
    color: var(--color-btn-text, #0f2040);
    box-shadow: 0 2px 16px rgba(255, 255, 255, 0.12), 0 1px 4px rgba(0, 0, 0, 0.25);
  }

  .app-btn--primary .app-btn__icon,
  .app-btn--primary .app-btn__arrow {
    color: var(--color-blue-bright, #4d9fff);
  }

  .app-btn--primary:not(:disabled):hover {
    background: #e8f2ff;
    border-color: #e8f2ff;
    color: var(--color-blue, #2451a0);
    box-shadow: 0 4px 28px rgba(255, 255, 255, 0.2), 0 0 0 3px rgba(77, 159, 255, 0.18);
  }

  /* ─── Variant: secondary ─── */
  /* Frosted white surface, softer emphasis */
  .app-btn--secondary {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.18);
    color: var(--color-text-secondary, #c2d9f5);
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.07);
  }

  .app-btn--secondary .app-btn__icon {
    color: var(--color-blue-bright, #4d9fff);
    opacity: 0.7;
  }

  .app-btn--secondary:not(:disabled):hover {
    background: rgba(255, 255, 255, 0.14);
    border-color: rgba(255, 255, 255, 0.35);
    color: #ffffff;
    box-shadow: 0 4px 20px rgba(77, 159, 255, 0.12), inset 0 1px 0 rgba(255, 255, 255, 0.1);
  }

  .app-btn--secondary:not(:disabled):hover .app-btn__icon {
    opacity: 1;
  }

  /* ─── Variant: ghost ─── */
  /* Borderline invisible, low emphasis */
  .app-btn--ghost {
    background: transparent;
    border-color: rgba(255, 255, 255, 0.12);
    color: var(--color-text-muted, #5a88b8);
  }

  .app-btn--ghost .app-btn__icon {
    font-size: 0.9rem;
    color: var(--color-text-muted, #5a88b8);
  }

  .app-btn--ghost:not(:disabled):hover {
    background: rgba(77, 159, 255, 0.06);
    border-color: rgba(255, 255, 255, 0.25);
    color: var(--color-text-secondary, #c2d9f5);
  }

  /* ─── Variant: danger ─── */
  /* White bg, red text — destructive actions */
  .app-btn--danger {
    background: #ffffff;
    border-color: #ffffff;
    color: #c0392b;
    box-shadow: 0 2px 16px rgba(255, 255, 255, 0.1);
  }

  .app-btn--danger:not(:disabled):hover {
    background: #fff0ef;
    border-color: #ffcdd1;
    color: #96201a;
    box-shadow: 0 4px 24px rgba(255, 80, 60, 0.15);
  }
</style>

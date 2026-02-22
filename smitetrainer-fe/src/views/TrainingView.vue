<template>
  <div class="training-page">
    <div class="back-glow back-glow--left"></div>
    <div class="back-glow back-glow--right"></div>

    <div class="training-shell">
      <section class="control-panel">
        <header class="panel-head">
          <h1>Dragon HP Analyzer</h1>
          <p>
            Pulls match timeline from your Go backend and renders the estimated dragon HP curve for
            the selected dragon kill.
          </p>
          <div class="api-chip font-mono">API: {{ apiBaseUrl }}</div>
        </header>

        <form class="query-form" @submit.prevent="runLookup">
          <label>
            <span>Match ID</span>
            <input
              v-model.trim="matchId"
              type="text"
              placeholder="EUW1_1234567890"
              required
              autocomplete="off"
            />
          </label>

          <label>
            <span>Dragon Index</span>
            <input v-model.number="dragonIndex" type="number" min="1" step="1" required />
          </label>

          <label>
            <span>Tick (ms)</span>
            <input v-model.number="tickMs" type="number" min="100" step="100" />
          </label>

          <label>
            <span>Window (ms)</span>
            <input v-model.number="windowMs" type="number" min="1000" step="500" />
          </label>

          <div class="form-actions">
            <AppButton
              label="FETCH HP SERIES"
              icon="⛏"
              variant="primary"
              :show-arrow="true"
              type="submit"
              :disabled="loading"
              :full-width="false"
            />
            <AppButton
              label="RESET"
              icon="↺"
              variant="ghost"
              :show-arrow="false"
              type="button"
              :disabled="loading"
              :full-width="false"
              @click="resetResult"
            />
          </div>
        </form>
      </section>

      <section class="result-panel">
        <div v-if="loading" class="state-card">
          <div class="spinner"></div>
          <p>Fetching match and timeline data...</p>
        </div>

        <div v-else-if="errorMessage" class="state-card state-card--error">
          <h2>Request failed</h2>
          <p>{{ errorMessage }}</p>
        </div>

        <div v-else-if="result" class="result-content">
          <div class="summary-grid">
            <article class="summary-card">
              <h3>Dragon</h3>
              <p>{{ result.dragonType }}</p>
              <small class="font-mono">Index #{{ result.dragonIndex }}</small>
            </article>
            <article class="summary-card">
              <h3>Kill Time</h3>
              <p>{{ formatGameTime(result.killMs) }}</p>
              <small class="font-mono">{{ result.killMs }} ms</small>
            </article>
            <article class="summary-card">
              <h3>Estimated HP0</h3>
              <p>{{ result.hpModel.hp0.toLocaleString() }}</p>
              <small class="font-mono">Level {{ result.hpModel.levelEstimate }}</small>
            </article>
            <article class="summary-card">
              <h3>Model</h3>
              <p>{{ result.hpModel.model }}</p>
              <small class="font-mono">tick {{ result.hpModel.tickMs }} ms</small>
            </article>
          </div>

          <div class="chart-card">
            <div class="chart-head">
              <h3>HP Curve (Estimated)</h3>
              <span class="font-mono">{{ result.series.length }} points</span>
            </div>
            <svg
              class="hp-chart"
              viewBox="0 0 1000 240"
              preserveAspectRatio="none"
              role="img"
              aria-label="Dragon HP percent over time"
            >
              <defs>
                <linearGradient id="lineGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                  <stop offset="0%" stop-color="#00d4ff" />
                  <stop offset="100%" stop-color="#4d9fff" />
                </linearGradient>
              </defs>
              <line x1="0" y1="20" x2="1000" y2="20" class="grid-line" />
              <line x1="0" y1="120" x2="1000" y2="120" class="grid-line" />
              <line x1="0" y1="220" x2="1000" y2="220" class="grid-line" />
              <polyline :points="seriesPath" class="hp-line" />
            </svg>
          </div>

          <div class="table-card">
            <h3>Series Preview</h3>
            <table>
              <thead>
                <tr>
                  <th>tMs</th>
                  <th>HP</th>
                  <th>HP %</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="point in seriesPreview" :key="point.tMs">
                  <td class="font-mono">{{ point.tMs }}</td>
                  <td>{{ point.hp.toLocaleString() }}</td>
                  <td>{{ (point.hpPct * 100).toFixed(2) }}%</td>
                </tr>
              </tbody>
            </table>
          </div>

          <details class="raw-json">
            <summary>Raw JSON</summary>
            <pre>{{ prettyResponse }}</pre>
          </details>
        </div>

        <div v-else class="state-card">
          <h2>Ready</h2>
          <p>Fill the form and fetch a match to load dragon HP data.</p>
        </div>
      </section>
    </div>
  </div>
</template>

<script lang="ts">
  import { defineComponent } from 'vue';
  import AppButton from '@/components/AppButton.vue';
  import {
    fetchDragonHP,
    getApiBaseUrl,
    type DragonHPResponse,
    type HPPoint,
  } from '@/services/dragonApi';

  export default defineComponent({
    name: 'TrainingView',
    components: { AppButton },
    data() {
      return {
        matchId: '',
        dragonIndex: 1,
        tickMs: 500,
        windowMs: 20000,
        loading: false,
        errorMessage: '',
        result: null as DragonHPResponse | null,
      };
    },
    computed: {
      apiBaseUrl(): string {
        return getApiBaseUrl();
      },
      seriesPath(): string {
        if (!this.result || this.result.series.length < 2) return '';

        const points = this.result.series;
        const lastIndex = points.length - 1;

        return points
          .map((point, index) => {
            const x = (index / lastIndex) * 1000;
            const y = 220 - point.hpPct * 200;
            return `${x.toFixed(2)},${y.toFixed(2)}`;
          })
          .join(' ');
      },
      seriesPreview(): HPPoint[] {
        if (!this.result) return [];

        const series = this.result.series;
        if (series.length <= 16) return series;

        const head = series.slice(0, 8);
        const tail = series.slice(series.length - 8);
        return [...head, ...tail];
      },
      prettyResponse(): string {
        if (!this.result) return '';
        return JSON.stringify(this.result, null, 2);
      },
    },
    methods: {
      async runLookup() {
        if (this.loading) return;
        if (!this.matchId.trim()) {
          this.errorMessage = 'Match ID is required.';
          this.result = null;
          return;
        }
        if (this.dragonIndex < 1) {
          this.errorMessage = 'Dragon index must be 1 or higher.';
          this.result = null;
          return;
        }

        this.loading = true;
        this.errorMessage = '';

        try {
          const payload = await fetchDragonHP({
            matchId: this.matchId.trim(),
            dragonIndex: this.dragonIndex,
            tickMs: this.tickMs,
            windowMs: this.windowMs,
          });
          this.result = payload;
        } catch (err: unknown) {
          this.result = null;
          this.errorMessage =
            err instanceof Error ? err.message : 'Unknown error while fetching data.';
        } finally {
          this.loading = false;
        }
      },
      resetResult() {
        this.result = null;
        this.errorMessage = '';
      },
      formatGameTime(ms: number): string {
        const totalSec = Math.floor(ms / 1000);
        const min = Math.floor(totalSec / 60);
        const sec = totalSec % 60;
        return `${String(min).padStart(2, '0')}:${String(sec).padStart(2, '0')}`;
      },
    },
  });
</script>

<style scoped>
  .training-page {
    position: relative;
    min-height: 100vh;
    padding: 6rem 1.5rem 2rem;
    background: radial-gradient(1200px 600px at 0% 0%, rgba(77, 159, 255, 0.12), transparent 60%),
      radial-gradient(1000px 500px at 100% 100%, rgba(0, 212, 255, 0.07), transparent 60%), #050d1a;
    overflow: hidden;
  }

  .back-glow {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    pointer-events: none;
  }

  .back-glow--left {
    width: 340px;
    height: 340px;
    top: -120px;
    left: -100px;
    background: rgba(77, 159, 255, 0.28);
  }

  .back-glow--right {
    width: 320px;
    height: 320px;
    right: -120px;
    bottom: -80px;
    background: rgba(0, 212, 255, 0.2);
  }

  .training-shell {
    position: relative;
    z-index: 10;
    width: min(1180px, 100%);
    margin: 0 auto;
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 1rem;
  }

  .control-panel,
  .result-panel {
    background: rgba(10, 22, 40, 0.78);
    border: 1px solid rgba(77, 159, 255, 0.2);
    border-radius: 14px;
    backdrop-filter: blur(6px);
  }

  .control-panel {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .panel-head h1 {
    font-size: 1.25rem;
    letter-spacing: 0.06em;
    color: #fff;
  }

  .panel-head p {
    margin-top: 0.45rem;
    color: var(--color-text-secondary);
    line-height: 1.45;
    font-size: 0.9rem;
  }

  .api-chip {
    margin-top: 0.75rem;
    padding: 0.4rem 0.55rem;
    background: rgba(77, 159, 255, 0.1);
    border: 1px solid rgba(77, 159, 255, 0.28);
    border-radius: 6px;
    color: var(--color-blue-bright);
    font-size: 0.7rem;
    letter-spacing: 0.08em;
    word-break: break-all;
  }

  .query-form {
    display: grid;
    gap: 0.8rem;
  }

  .query-form label {
    display: grid;
    gap: 0.3rem;
  }

  .query-form span {
    color: var(--color-text-secondary);
    font-size: 0.82rem;
    letter-spacing: 0.06em;
  }

  .query-form input {
    height: 2.5rem;
    border-radius: 8px;
    border: 1px solid rgba(77, 159, 255, 0.25);
    background: rgba(2, 8, 16, 0.75);
    color: #fff;
    padding: 0 0.7rem;
    outline: none;
    font-size: 0.95rem;
  }

  .query-form input:focus {
    border-color: rgba(77, 159, 255, 0.7);
    box-shadow: 0 0 0 3px rgba(77, 159, 255, 0.15);
  }

  .form-actions {
    margin-top: 0.45rem;
    display: flex;
    gap: 0.6rem;
  }

  .result-panel {
    padding: 1rem;
    min-height: 500px;
  }

  .state-card {
    min-height: 460px;
    border: 1px dashed rgba(77, 159, 255, 0.28);
    border-radius: 10px;
    background: rgba(2, 8, 16, 0.4);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 0.75rem;
    text-align: center;
    color: var(--color-text-secondary);
    padding: 1.2rem;
  }

  .state-card h2 {
    color: #fff;
    letter-spacing: 0.07em;
  }

  .state-card--error {
    border-color: rgba(255, 100, 100, 0.45);
    background: rgba(48, 8, 8, 0.35);
  }

  .spinner {
    width: 34px;
    height: 34px;
    border-radius: 50%;
    border: 3px solid rgba(77, 159, 255, 0.2);
    border-top-color: var(--color-blue-bright);
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .result-content {
    display: grid;
    gap: 1rem;
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 0.7rem;
  }

  .summary-card {
    border: 1px solid rgba(77, 159, 255, 0.22);
    border-radius: 10px;
    padding: 0.75rem;
    background: rgba(2, 8, 16, 0.55);
  }

  .summary-card h3 {
    color: var(--color-text-muted);
    font-size: 0.72rem;
    letter-spacing: 0.12em;
  }

  .summary-card p {
    margin-top: 0.35rem;
    color: #fff;
    font-size: 1.1rem;
    font-weight: 700;
  }

  .summary-card small {
    margin-top: 0.3rem;
    display: block;
    color: var(--color-blue-bright);
    opacity: 0.85;
    font-size: 0.68rem;
  }

  .chart-card,
  .table-card,
  .raw-json {
    border: 1px solid rgba(77, 159, 255, 0.22);
    border-radius: 10px;
    background: rgba(2, 8, 16, 0.55);
    padding: 0.85rem;
  }

  .chart-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.55rem;
  }

  .chart-head h3,
  .table-card h3 {
    color: #fff;
    font-size: 0.95rem;
    letter-spacing: 0.08em;
  }

  .chart-head span {
    color: var(--color-blue-bright);
    font-size: 0.72rem;
    letter-spacing: 0.1em;
  }

  .hp-chart {
    width: 100%;
    height: 220px;
    background: rgba(0, 0, 0, 0.25);
    border: 1px solid rgba(77, 159, 255, 0.15);
    border-radius: 8px;
  }

  .grid-line {
    stroke: rgba(77, 159, 255, 0.2);
    stroke-width: 1;
  }

  .hp-line {
    fill: none;
    stroke: url(#lineGrad);
    stroke-width: 4;
    stroke-linejoin: round;
    stroke-linecap: round;
    filter: drop-shadow(0 0 6px rgba(77, 159, 255, 0.4));
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 0.45rem 0.5rem;
    border-bottom: 1px solid rgba(77, 159, 255, 0.14);
    text-align: left;
    font-size: 0.86rem;
    color: var(--color-text-secondary);
  }

  th {
    color: var(--color-blue-bright);
    letter-spacing: 0.08em;
    font-size: 0.72rem;
  }

  .raw-json summary {
    cursor: pointer;
    color: var(--color-blue-bright);
    font-size: 0.86rem;
  }

  .raw-json pre {
    margin-top: 0.7rem;
    max-height: 260px;
    overflow: auto;
    background: rgba(0, 0, 0, 0.35);
    border-radius: 8px;
    padding: 0.75rem;
    color: #d8ebff;
    font-size: 0.76rem;
    line-height: 1.45;
  }

  @media (max-width: 980px) {
    .training-shell {
      grid-template-columns: 1fr;
    }

    .summary-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }

  @media (max-width: 560px) {
    .summary-grid {
      grid-template-columns: 1fr;
    }

    .form-actions {
      flex-direction: column;
    }
  }
</style>

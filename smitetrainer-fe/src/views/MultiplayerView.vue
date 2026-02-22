<template>
  <div class="arena" :class="{ 'arena--shake': shaking }">
    <div class="arena-bg">
      <div class="energy-grid"></div>
      <div class="vignette"></div>
      <div class="bg-beam bg-beam--left"></div>
      <div class="bg-beam bg-beam--right"></div>
    </div>

    <!-- Header -->
    <div class="arena-header">
      <div
        class="player-nameplate nameplate--left"
        :class="{ 'nameplate--win': roundResult === 'p1', 'nameplate--lose': roundResult === 'p2' }"
      >
        <span class="nameplate__tag font-mono">P1</span>
        <span class="nameplate__name">SUMMONER ONE</span>
      </div>
      <div class="header-center">
        <div class="score-board">
          <span class="score-num" :class="{ 'score-num--hot': score.p1 > score.p2 }">{{
            score.p1
          }}</span>
          <div class="score-pips">
            <span
              v-for="r in 5"
              :key="r"
              class="score-pip"
              :class="{
                'score-pip--p1': roundHistory[r - 1] === 'p1',
                'score-pip--p2': roundHistory[r - 1] === 'p2',
                'score-pip--active': r === currentRound && !matchOver,
              }"
            >
            </span>
          </div>
          <span class="score-num" :class="{ 'score-num--hot': score.p2 > score.p1 }">{{
            score.p2
          }}</span>
        </div>
        <div class="match-meta font-mono">BEST OF 5 · ROUND {{ currentRound }}</div>
      </div>
      <div
        class="player-nameplate nameplate--right"
        :class="{ 'nameplate--win': roundResult === 'p2', 'nameplate--lose': roundResult === 'p1' }"
      >
        <span class="nameplate__name">SUMMONER TWO</span>
        <span class="nameplate__tag font-mono">P2</span>
      </div>
    </div>

    <div class="stage-divider">
      <div class="divider-line"></div>
      <div class="divider-vs">VS</div>
      <div class="divider-line"></div>
    </div>

    <!-- Main Stage -->
    <div class="arena-stage">
      <!-- P1 Card -->
      <div
        class="player-card player-card--left"
        :class="{
          'player-card--fired': player1Ready,
          'player-card--win': roundResult === 'p1',
          'player-card--lose': roundResult === 'p2',
        }"
      >
        <div class="player-card__bg"></div>

        <div class="champ-avatar">
          <div class="champ-avatar__img">⚔</div>
        </div>

        <div class="player-stats">
          <div class="stat-row">
            <span class="stat-label font-mono">RANK</span>
            <span class="stat-value">Diamond II</span>
          </div>
          <div class="stat-row">
            <span class="stat-label font-mono">WIN RATE</span>
            <span class="stat-value stat-value--good">63%</span>
          </div>
          <div class="stat-row">
            <span class="stat-label font-mono">SMITE ACC</span>
            <span class="stat-value stat-value--good">91%</span>
          </div>
          <div class="stat-row">
            <span class="stat-label font-mono">REACTION</span>
            <span class="stat-value">{{ player1Ready ? playerTimings[1] + 'ms' : '—' }}</span>
          </div>
        </div>

        <button
          class="smite-trigger"
          :class="{
            'smite-trigger--fired': player1Ready,
            'smite-trigger--locked': !canSmite && !player1Ready,
          }"
          :disabled="player1Ready || !canSmite"
          @click="fireSmite(1)"
        >
          <div class="trigger-inner">
            <span class="trigger-icon">⚡</span>
            <span class="trigger-label">{{ player1Ready ? 'SMITED' : 'SMITE' }}</span>
            <span class="trigger-key font-mono">Q</span>
          </div>
          <div class="trigger-shockwave" v-if="player1Ready"></div>
        </button>
      </div>

      <!-- Center -->
      <div class="center-stage">
        <transition name="slam">
          <div v-if="gameState === 'countdown'" class="countdown-display">
            <div class="countdown-num" :key="countdownValue">
              {{ countdownValue === 0 ? 'GO' : countdownValue }}
            </div>
          </div>
        </transition>

        <div class="clip-frame" :class="{ 'clip-frame--live': gameState === 'playing' }">
          <div class="clip-inner">
            <div v-if="gameState !== 'playing' && gameState !== 'countdown'" class="clip-idle">
              <div class="clip-idle__icon">▶</div>
              <div class="clip-idle__text font-mono">CLIP LOADS HERE</div>
            </div>
            <div v-if="gameState === 'playing'" class="clip-playing">
              <div class="clip-playing__label font-mono">● LIVE</div>
            </div>
          </div>
          <div class="danger-bar" v-if="gameState === 'playing'">
            <div
              class="danger-bar__fill"
              :class="{ 'danger-bar__fill--critical': timerPercent < 25 }"
              :style="{ width: timerPercent + '%' }"
            ></div>
          </div>
          <div class="clip-corner clip-corner--tl"></div>
          <div class="clip-corner clip-corner--tr"></div>
          <div class="clip-corner clip-corner--bl"></div>
          <div class="clip-corner clip-corner--br"></div>
        </div>

        <div class="timer-row" v-if="gameState === 'playing' || gameState === 'countdown'">
          <div class="timer-track">
            <div
              class="timer-fill"
              :class="{ 'timer-fill--critical': timerPercent < 25 }"
              :style="{ width: timerPercent + '%' }"
            ></div>
          </div>
          <span class="timer-val font-mono">{{ formattedTimer }}</span>
        </div>

        <div class="cta-zone">
          <button v-if="gameState === 'waiting'" class="start-btn" @click="startRound">
            <span class="start-btn__text">START ROUND</span>
            <span class="start-btn__sub font-mono">PRESS TO BEGIN</span>
          </button>
        </div>
      </div>

      <!-- P2 Card -->
      <div
        class="player-card player-card--right"
        :class="{
          'player-card--fired': player2Ready,
          'player-card--win': roundResult === 'p2',
          'player-card--lose': roundResult === 'p1',
        }"
      >
        <div class="player-card__bg"></div>

        <div class="champ-avatar">
          <div class="champ-avatar__img">🛡</div>
        </div>

        <div class="player-stats">
          <div class="stat-row">
            <span class="stat-label font-mono">RANK</span>
            <span class="stat-value">Platinum I</span>
          </div>
          <div class="stat-row">
            <span class="stat-label font-mono">WIN RATE</span>
            <span class="stat-value stat-value--warn">51%</span>
          </div>
          <div class="stat-row">
            <span class="stat-label font-mono">SMITE ACC</span>
            <span class="stat-value stat-value--warn">78%</span>
          </div>
          <div class="stat-row">
            <span class="stat-label font-mono">REACTION</span>
            <span class="stat-value">{{ player2Ready ? playerTimings[2] + 'ms' : '—' }}</span>
          </div>
        </div>

        <button
          class="smite-trigger"
          :class="{
            'smite-trigger--fired': player2Ready,
            'smite-trigger--locked': !canSmite && !player2Ready,
          }"
          :disabled="player2Ready || !canSmite"
          @click="fireSmite(2)"
        >
          <div class="trigger-inner">
            <span class="trigger-icon">⚡</span>
            <span class="trigger-label">{{ player2Ready ? 'SMITED' : 'SMITE' }}</span>
            <span class="trigger-key font-mono">P</span>
          </div>
          <div class="trigger-shockwave" v-if="player2Ready"></div>
        </button>
      </div>
    </div>

    <!-- Round overlay -->
    <transition name="slam-overlay">
      <div v-if="showRoundOverlay" class="round-overlay">
        <div class="round-overlay__slash"></div>
        <div class="round-overlay__content">
          <div class="round-overlay__winner font-mono">
            {{ lastRoundWinner === 'p1' ? 'PLAYER 1' : 'PLAYER 2' }}
          </div>
          <div class="round-overlay__label">WINS THE ROUND</div>
          <div class="round-overlay__score font-mono">
            {{ score.p1 }}<span>—</span>{{ score.p2 }}
          </div>
          <div class="round-overlay__next font-mono" v-if="!matchOver">NEXT ROUND IN...</div>
        </div>
      </div>
    </transition>

    <!-- Match over -->
    <transition name="slam-overlay">
      <div v-if="matchOver" class="match-overlay">
        <div class="match-overlay__flare match-overlay__flare--left"></div>
        <div class="match-overlay__flare match-overlay__flare--right"></div>
        <div class="match-overlay__content">
          <div class="match-overlay__eyebrow font-mono">MATCH COMPLETE</div>
          <div class="match-overlay__winner">{{ score.p1 >= 3 ? 'PLAYER 1' : 'PLAYER 2' }}</div>
          <div class="match-overlay__crown">👑</div>
          <div class="match-overlay__score font-mono">{{ score.p1 }} — {{ score.p2 }}</div>
          <button class="start-btn start-btn--rematch" @click="resetMatch">
            <span class="start-btn__text">REMATCH</span>
            <span class="start-btn__sub font-mono">PLAY AGAIN</span>
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script lang="ts">
  import { Options, Vue } from 'vue-class-component';

  type GameState = 'waiting' | 'countdown' | 'playing' | 'finished';
  type Winner = 'p1' | 'p2' | null;

  @Options({})
  export default class MultiplayerView extends Vue {
    gameState: GameState = 'waiting';
    currentRound = 1;
    score = { p1: 0, p2: 0 };
    roundHistory: Array<'p1' | 'p2' | null> = [null, null, null, null, null];
    timer = 0;
    timerMax = 0;
    timerInterval: number | null = null;
    smiteWindowStart = 0;
    countdownValue = 3;
    countdownInterval: number | null = null;
    player1Ready = false;
    player2Ready = false;
    playerTimings: Record<number, number> = {};
    roundResult: Winner = null;
    lastRoundWinner: Winner = null;
    showRoundOverlay = false;
    matchOver = false;
    shaking = false;

    get canSmite(): boolean {
      return this.gameState === 'playing';
    }

    get formattedTimer(): string {
      const s = Math.floor(this.timer / 1000);
      const ms = Math.floor((this.timer % 1000) / 10);
      return `${String(s).padStart(2, '0')}:${String(ms).padStart(2, '0')}`;
    }

    get timerPercent(): number {
      if (this.timerMax === 0) return 100;
      return (this.timer / this.timerMax) * 100;
    }

    shake() {
      this.shaking = true;
      setTimeout(() => {
        this.shaking = false;
      }, 400);
    }

    startRound() {
      this.gameState = 'countdown';
      this.countdownValue = 3;
      this.countdownInterval = window.setInterval(() => {
        this.countdownValue--;
        if (this.countdownValue < 0) {
          clearInterval(this.countdownInterval!);
          this.beginSmiteWindow();
        }
      }, 1000);
    }

    beginSmiteWindow() {
      this.gameState = 'playing';
      const duration = Math.floor(Math.random() * 5000) + 4000;
      this.timerMax = duration;
      this.timer = duration;
      this.smiteWindowStart = Date.now();
      this.timerInterval = window.setInterval(() => {
        this.timer -= 50;
        if (this.timer <= 0) {
          clearInterval(this.timerInterval!);
          this.finishRound();
        }
      }, 50);
    }

    fireSmite(player: 1 | 2) {
      if (!this.canSmite) return;
      this.shake();
      const elapsed = Date.now() - this.smiteWindowStart;
      this.playerTimings[player] = elapsed;
      if (player === 1) this.player1Ready = true;
      if (player === 2) this.player2Ready = true;
      if (this.player1Ready && this.player2Ready) {
        clearInterval(this.timerInterval!);
        this.finishRound();
      }
    }

    finishRound() {
      this.gameState = 'finished';
      let winner: 'p1' | 'p2';
      if (this.player1Ready && !this.player2Ready) winner = 'p1';
      else if (this.player2Ready && !this.player1Ready) winner = 'p2';
      else winner = this.playerTimings[1] <= this.playerTimings[2] ? 'p1' : 'p2';
      this.roundResult = winner;
      this.lastRoundWinner = winner;
      this.score[winner]++;
      this.roundHistory[this.currentRound - 1] = winner;
      this.showRoundOverlay = true;
      if (this.score.p1 >= 3 || this.score.p2 >= 3) {
        setTimeout(() => {
          this.showRoundOverlay = false;
          this.matchOver = true;
        }, 2800);
      } else {
        setTimeout(() => {
          this.showRoundOverlay = false;
          this.currentRound++;
          this.resetRound();
        }, 2800);
      }
    }

    resetRound() {
      this.gameState = 'waiting';
      this.timer = 0;
      this.timerMax = 0;
      this.player1Ready = false;
      this.player2Ready = false;
      this.playerTimings = {};
      this.roundResult = null;
      if (this.timerInterval) clearInterval(this.timerInterval);
      if (this.countdownInterval) clearInterval(this.countdownInterval);
    }

    resetMatch() {
      this.matchOver = false;
      this.currentRound = 1;
      this.score = { p1: 0, p2: 0 };
      this.roundHistory = [null, null, null, null, null];
      this.lastRoundWinner = null;
      this.resetRound();
    }

    beforeUnmount() {
      if (this.timerInterval) clearInterval(this.timerInterval);
      if (this.countdownInterval) clearInterval(this.countdownInterval);
    }
  }
</script>

<style scoped>
  .arena {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    position: relative;
    overflow: hidden;
    background: #020810;
  }

  .arena--shake {
    animation: screenShake 0.35s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
  }

  @keyframes screenShake {
    0% {
      transform: translate(0, 0) rotate(0);
    }
    15% {
      transform: translate(-5px, 3px) rotate(-0.3deg);
    }
    30% {
      transform: translate(5px, -3px) rotate(0.3deg);
    }
    45% {
      transform: translate(-4px, 2px) rotate(-0.2deg);
    }
    60% {
      transform: translate(4px, -2px) rotate(0.2deg);
    }
    75% {
      transform: translate(-2px, 1px);
    }
    100% {
      transform: translate(0, 0);
    }
  }

  /* ── Background ── */
  .arena-bg {
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 0;
  }

  .energy-grid {
    position: absolute;
    inset: 0;
    background-image: linear-gradient(rgba(77, 159, 255, 0.04) 1px, transparent 1px),
      linear-gradient(90deg, rgba(77, 159, 255, 0.04) 1px, transparent 1px);
    background-size: 60px 60px;
  }

  .vignette {
    position: absolute;
    inset: 0;
    background: radial-gradient(ellipse at center, transparent 40%, rgba(2, 8, 16, 0.9) 100%);
  }

  .bg-beam {
    position: absolute;
    top: 0;
    bottom: 0;
    width: 40%;
    opacity: 0.07;
  }
  .bg-beam--left {
    left: 0;
    background: linear-gradient(90deg, rgba(77, 159, 255, 0.6), transparent);
  }
  .bg-beam--right {
    right: 0;
    background: linear-gradient(270deg, rgba(77, 159, 255, 0.6), transparent);
  }

  /* ── Header ── */
  .arena-header {
    position: relative;
    z-index: 10;
    display: grid;
    grid-template-columns: 1fr auto 1fr;
    align-items: center;
    padding: 1rem 2rem;
    background: rgba(2, 8, 16, 0.8);
    border-bottom: 1px solid rgba(77, 159, 255, 0.12);
  }

  .player-nameplate {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    border: 1px solid rgba(77, 159, 255, 0.08);
    background: rgba(10, 22, 40, 0.5);
    transition: border-color 0.3s, box-shadow 0.3s, opacity 0.3s;
  }

  .nameplate--right {
    justify-content: flex-end;
  }
  .nameplate--win {
    border-color: rgba(80, 200, 120, 0.5) !important;
    box-shadow: 0 0 20px rgba(80, 200, 120, 0.1);
  }
  .nameplate--lose {
    opacity: 0.4;
  }
  .nameplate__tag {
    font-size: 0.6rem;
    letter-spacing: 0.25em;
    color: var(--color-blue-bright);
    opacity: 0.6;
  }
  .nameplate__name {
    font-size: 0.9rem;
    font-weight: 700;
    letter-spacing: 0.1em;
    color: var(--color-text-primary);
  }

  .header-center {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.4rem;
  }
  .score-board {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .score-num {
    font-size: 2.5rem;
    font-weight: 700;
    color: rgba(255, 255, 255, 0.2);
    letter-spacing: 0.05em;
    transition: color 0.3s, text-shadow 0.3s;
  }
  .score-num--hot {
    color: #fff;
    text-shadow: 0 0 20px rgba(77, 159, 255, 0.7), 0 0 40px rgba(77, 159, 255, 0.3);
  }

  .score-pips {
    display: flex;
    gap: 0.3rem;
  }
  .score-pip {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(77, 159, 255, 0.15);
    transition: all 0.3s;
  }
  .score-pip--active {
    background: rgba(77, 159, 255, 0.2);
    border-color: var(--color-blue-bright);
    animation: pipPulse 1s ease infinite;
  }
  .score-pip--p1 {
    background: rgba(80, 200, 120, 0.4);
    border-color: #50c878;
    box-shadow: 0 0 6px rgba(80, 200, 120, 0.3);
  }
  .score-pip--p2 {
    background: rgba(255, 80, 80, 0.4);
    border-color: #ff5050;
    box-shadow: 0 0 6px rgba(255, 80, 80, 0.3);
  }

  @keyframes pipPulse {
    0%,
    100% {
      box-shadow: 0 0 4px rgba(77, 159, 255, 0.4);
    }
    50% {
      box-shadow: 0 0 10px rgba(77, 159, 255, 0.8);
    }
  }

  .match-meta {
    font-size: 0.58rem;
    letter-spacing: 0.25em;
    color: var(--color-text-muted);
  }

  /* ── Divider ── */
  .stage-divider {
    position: relative;
    z-index: 10;
    display: flex;
    align-items: center;
    padding: 0 2rem;
  }
  .divider-line {
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(77, 159, 255, 0.2), transparent);
  }
  .divider-vs {
    padding: 0 1rem;
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.4em;
    color: rgba(77, 159, 255, 0.25);
  }

  /* ── Stage ── */
  .arena-stage {
    position: relative;
    z-index: 10;
    display: grid;
    grid-template-columns: 200px 1fr 200px;
    flex: 1;
    padding: 1.5rem 1.5rem 2rem;
    gap: 1.5rem;
    align-items: center;
  }

  /* ── Player Cards ── */
  .player-card {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 1.25rem 1rem;
    border-radius: 12px;
    border: 1px solid rgba(77, 159, 255, 0.08);
    overflow: hidden;
    transition: border-color 0.3s, box-shadow 0.3s, opacity 0.3s;
  }

  .player-card__bg {
    position: absolute;
    inset: 0;
    background: rgba(10, 22, 40, 0.5);
    transition: background 0.3s;
  }

  .player-card--fired .player-card__bg {
    background: rgba(77, 159, 255, 0.04);
  }
  .player-card--win {
    border-color: rgba(80, 200, 120, 0.4) !important;
    box-shadow: inset 0 0 30px rgba(80, 200, 120, 0.04), 0 0 24px rgba(80, 200, 120, 0.08);
  }
  .player-card--lose {
    opacity: 0.45;
  }

  /* ── Champion Avatar ── */
  .champ-avatar {
    position: relative;
    z-index: 2;
    width: 72px;
    height: 72px;
    border-radius: 50%;
    background: rgba(77, 159, 255, 0.08);
    border: 1px solid rgba(77, 159, 255, 0.2);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
  .champ-avatar__img {
    font-size: 2rem;
    line-height: 1;
  }
  .player-card--win .champ-avatar {
    border-color: rgba(80, 200, 120, 0.5);
    box-shadow: 0 0 20px rgba(80, 200, 120, 0.2);
  }

  /* ── Stats ── */
  .player-stats {
    position: relative;
    z-index: 2;
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .stat-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.28rem 0.5rem;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid rgba(77, 159, 255, 0.05);
  }

  .stat-label {
    font-size: 0.52rem;
    letter-spacing: 0.18em;
    color: var(--color-text-muted);
  }
  .stat-value {
    font-size: 0.72rem;
    font-weight: 700;
    letter-spacing: 0.05em;
    color: var(--color-text-secondary);
  }
  .stat-value--good {
    color: #50c878;
  }
  .stat-value--warn {
    color: #f0a030;
  }

  /* ── Smite Trigger ── */
  .smite-trigger {
    position: relative;
    z-index: 2;
    width: 110px;
    height: 110px;
    border-radius: 50%;
    border: 2px solid rgba(255, 255, 255, 0.8);
    background: #ffffff;
    cursor: pointer;
    transition: all 0.15s ease;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
    overflow: visible;
    flex-shrink: 0;
  }

  .smite-trigger:not(.smite-trigger--locked):not(.smite-trigger--fired):hover {
    transform: scale(1.07);
    box-shadow: 0 0 0 6px rgba(255, 255, 255, 0.1), 0 0 40px rgba(255, 255, 255, 0.25),
      0 8px 32px rgba(0, 0, 0, 0.5);
  }

  .smite-trigger:not(.smite-trigger--locked):not(.smite-trigger--fired):active {
    transform: scale(0.93);
  }

  .smite-trigger--locked {
    opacity: 0.2;
    cursor: not-allowed;
    border-color: rgba(255, 255, 255, 0.2);
    background: rgba(255, 255, 255, 0.05);
  }

  .smite-trigger--fired {
    background: rgba(77, 159, 255, 0.08) !important;
    border-color: var(--color-blue-bright) !important;
    box-shadow: 0 0 0 4px rgba(77, 159, 255, 0.15), 0 0 40px rgba(77, 159, 255, 0.3),
      inset 0 0 20px rgba(77, 159, 255, 0.05);
    cursor: default;
    animation: firedPulse 2s ease infinite;
  }

  @keyframes firedPulse {
    0%,
    100% {
      box-shadow: 0 0 0 4px rgba(77, 159, 255, 0.15), 0 0 40px rgba(77, 159, 255, 0.3);
    }
    50% {
      box-shadow: 0 0 0 8px rgba(77, 159, 255, 0.08), 0 0 60px rgba(77, 159, 255, 0.5);
    }
  }

  .trigger-inner {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    gap: 0.15rem;
  }

  .trigger-icon {
    font-size: 1.4rem;
    transition: all 0.15s;
    color: #0f2040;
  }
  .smite-trigger--fired .trigger-icon {
    color: var(--color-blue-bright);
  }

  .trigger-label {
    font-size: 0.8rem;
    font-weight: 700;
    letter-spacing: 0.15em;
    color: #0f2040;
    transition: color 0.15s;
  }
  .smite-trigger--fired .trigger-label {
    color: var(--color-blue-bright);
  }

  .trigger-key {
    font-size: 0.5rem;
    letter-spacing: 0.15em;
    color: #5a88b8;
    border: 1px solid rgba(90, 136, 184, 0.3);
    padding: 0.1rem 0.3rem;
    border-radius: 3px;
  }
  .smite-trigger--fired .trigger-key {
    color: var(--color-cyan);
    border-color: var(--color-cyan);
  }

  .trigger-shockwave {
    position: absolute;
    inset: -2px;
    border-radius: 50%;
    border: 3px solid rgba(77, 159, 255, 0.6);
    animation: shockwave 0.6s ease-out forwards;
    pointer-events: none;
  }

  @keyframes shockwave {
    0% {
      transform: scale(1);
      opacity: 1;
    }
    100% {
      transform: scale(1.8);
      opacity: 0;
    }
  }

  /* ── Center ── */
  .center-stage {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .countdown-display {
    position: absolute;
    inset: 0;
    z-index: 20;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(2, 8, 16, 0.7);
    backdrop-filter: blur(4px);
  }

  .countdown-num {
    font-size: 10rem;
    font-weight: 700;
    color: #fff;
    line-height: 1;
    text-shadow: 0 0 30px rgba(77, 159, 255, 1), 0 0 80px rgba(77, 159, 255, 0.5),
      0 0 120px rgba(77, 159, 255, 0.2);
    animation: countSlam 0.9s cubic-bezier(0.22, 1, 0.36, 1) both;
  }

  @keyframes countSlam {
    0% {
      transform: scale(2.5);
      opacity: 0;
      filter: blur(20px);
    }
    40% {
      transform: scale(0.9);
      opacity: 1;
      filter: blur(0);
    }
    70% {
      transform: scale(1.03);
    }
    85% {
      transform: scale(0.98);
      opacity: 1;
    }
    100% {
      transform: scale(0.85);
      opacity: 0;
    }
  }

  /* ── Clip ── */
  .clip-frame {
    position: relative;
    width: 100%;
    aspect-ratio: 16/9;
    background: #010609;
    border-radius: 10px;
    overflow: hidden;
    border: 1px solid rgba(77, 159, 255, 0.1);
    transition: border-color 0.3s, box-shadow 0.3s;
  }

  .clip-frame--live {
    border-color: rgba(77, 159, 255, 0.3);
    box-shadow: 0 0 30px rgba(77, 159, 255, 0.08);
  }
  .clip-inner {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .clip-idle {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }
  .clip-idle__icon {
    font-size: 2.5rem;
    color: rgba(77, 159, 255, 0.15);
  }
  .clip-idle__text {
    font-size: 0.55rem;
    letter-spacing: 0.25em;
    color: rgba(77, 159, 255, 0.2);
  }
  .clip-playing {
    position: absolute;
    top: 0.75rem;
    left: 0.75rem;
  }
  .clip-playing__label {
    font-size: 0.6rem;
    letter-spacing: 0.2em;
    color: #ff4444;
  }

  .danger-bar {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: rgba(255, 255, 255, 0.05);
  }
  .danger-bar__fill {
    height: 100%;
    background: linear-gradient(90deg, var(--color-blue-bright), var(--color-cyan));
    box-shadow: 0 0 10px rgba(77, 159, 255, 0.7);
    transition: width 0.05s linear, background 0.3s;
  }
  .danger-bar__fill--critical {
    background: linear-gradient(90deg, #ff4444, #ff8800) !important;
    box-shadow: 0 0 12px rgba(255, 80, 0, 0.8) !important;
    animation: dangerFlash 0.3s ease infinite;
  }

  @keyframes dangerFlash {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.6;
    }
  }

  .clip-corner {
    position: absolute;
    width: 16px;
    height: 16px;
    border-color: rgba(77, 159, 255, 0.3);
    border-style: solid;
  }
  .clip-corner--tl {
    top: 8px;
    left: 8px;
    border-width: 1px 0 0 1px;
  }
  .clip-corner--tr {
    top: 8px;
    right: 8px;
    border-width: 1px 1px 0 0;
  }
  .clip-corner--bl {
    bottom: 8px;
    left: 8px;
    border-width: 0 0 1px 1px;
  }
  .clip-corner--br {
    bottom: 8px;
    right: 8px;
    border-width: 0 1px 1px 0;
  }

  /* ── Timer ── */
  .timer-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    width: 100%;
  }
  .timer-track {
    flex: 1;
    height: 3px;
    border-radius: 2px;
    background: rgba(255, 255, 255, 0.05);
    overflow: hidden;
  }
  .timer-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--color-blue), var(--color-blue-bright));
    transition: width 0.05s linear;
    border-radius: 2px;
  }
  .timer-fill--critical {
    background: linear-gradient(90deg, #cc2200, #ff5500) !important;
    animation: dangerFlash 0.3s ease infinite;
  }
  .timer-val {
    font-size: 0.75rem;
    font-weight: 700;
    color: var(--color-text-muted);
    letter-spacing: 0.1em;
    flex-shrink: 0;
    width: 3rem;
    text-align: right;
  }

  /* ── CTA ── */
  .cta-zone {
    display: flex;
    justify-content: center;
  }

  .start-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.2rem;
    padding: 1rem 3rem;
    background: #ffffff;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    box-shadow: 0 4px 24px rgba(255, 255, 255, 0.15), 0 2px 8px rgba(0, 0, 0, 0.4);
    transition: all 0.18s ease;
  }
  .start-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 32px rgba(255, 255, 255, 0.25);
  }
  .start-btn:active {
    transform: scale(0.97);
  }
  .start-btn--rematch {
    background: rgba(77, 159, 255, 0.1);
    border: 1px solid rgba(77, 159, 255, 0.3);
    box-shadow: 0 0 20px rgba(77, 159, 255, 0.15);
  }
  .start-btn__text {
    font-size: 1rem;
    font-weight: 700;
    letter-spacing: 0.2em;
    color: #0f2040;
  }
  .start-btn--rematch .start-btn__text {
    color: var(--color-blue-bright);
  }
  .start-btn__sub {
    font-size: 0.55rem;
    letter-spacing: 0.2em;
    color: #5a88b8;
  }

  /* ── Round Overlay ── */
  .round-overlay {
    position: fixed;
    inset: 0;
    z-index: 200;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(2, 8, 16, 0.85);
    backdrop-filter: blur(8px);
  }
  .round-overlay__slash {
    position: absolute;
    inset: 0;
    background: linear-gradient(
      105deg,
      transparent 48%,
      rgba(77, 159, 255, 0.06) 49%,
      rgba(77, 159, 255, 0.06) 51%,
      transparent 52%
    );
  }
  .round-overlay__content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    animation: slamIn 0.4s cubic-bezier(0.22, 1, 0.36, 1) both;
  }
  .round-overlay__winner {
    font-size: 0.7rem;
    letter-spacing: 0.5em;
    color: var(--color-blue-bright);
  }
  .round-overlay__label {
    font-size: 3.5rem;
    font-weight: 700;
    letter-spacing: 0.1em;
    color: #fff;
    text-shadow: 0 0 40px rgba(255, 255, 255, 0.3);
    line-height: 1;
  }
  .round-overlay__score {
    font-size: 2rem;
    font-weight: 700;
    color: var(--color-text-muted);
    letter-spacing: 0.3em;
  }
  .round-overlay__score span {
    color: rgba(77, 159, 255, 0.3);
    margin: 0 0.5rem;
  }
  .round-overlay__next {
    font-size: 0.6rem;
    letter-spacing: 0.25em;
    color: var(--color-text-muted);
    margin-top: 0.5rem;
  }

  @keyframes slamIn {
    0% {
      transform: scaleY(0.3) scaleX(1.2);
      opacity: 0;
      filter: blur(10px);
    }
    60% {
      transform: scaleY(1.05) scaleX(0.98);
      opacity: 1;
      filter: blur(0);
    }
    100% {
      transform: scale(1);
    }
  }

  /* ── Match Overlay ── */
  .match-overlay {
    position: fixed;
    inset: 0;
    z-index: 300;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(2, 8, 16, 0.92);
    backdrop-filter: blur(12px);
  }
  .match-overlay__flare {
    position: absolute;
    top: 0;
    bottom: 0;
    width: 35%;
    opacity: 0.06;
  }
  .match-overlay__flare--left {
    left: 0;
    background: linear-gradient(90deg, rgba(77, 159, 255, 1), transparent);
  }
  .match-overlay__flare--right {
    right: 0;
    background: linear-gradient(270deg, rgba(77, 159, 255, 1), transparent);
  }
  .match-overlay__content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.25rem;
    animation: slamIn 0.5s cubic-bezier(0.22, 1, 0.36, 1) both;
  }
  .match-overlay__eyebrow {
    font-size: 0.65rem;
    letter-spacing: 0.5em;
    color: var(--color-text-muted);
  }
  .match-overlay__winner {
    font-size: 4.5rem;
    font-weight: 700;
    letter-spacing: 0.15em;
    color: #fff;
    text-shadow: 0 0 30px rgba(77, 159, 255, 0.8), 0 0 60px rgba(77, 159, 255, 0.4),
      0 0 100px rgba(77, 159, 255, 0.2);
    line-height: 1;
  }
  .match-overlay__crown {
    font-size: 3rem;
    animation: crownFloat 2s ease infinite;
  }
  @keyframes crownFloat {
    0%,
    100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-6px);
    }
  }
  .match-overlay__score {
    font-size: 1.2rem;
    letter-spacing: 0.3em;
    color: var(--color-blue-bright);
    padding: 0.4rem 1.5rem;
    border: 1px solid rgba(77, 159, 255, 0.2);
    border-radius: 4px;
    background: rgba(77, 159, 255, 0.06);
  }

  /* ── Transitions ── */
  .slam-overlay-enter-active {
    transition: opacity 0.3s ease;
  }
  .slam-overlay-leave-active {
    transition: opacity 0.4s ease;
  }
  .slam-overlay-enter-from,
  .slam-overlay-leave-to {
    opacity: 0;
  }
  .slam-enter-active,
  .slam-leave-active {
    transition: opacity 0.2s;
  }
  .slam-enter-from,
  .slam-leave-to {
    opacity: 0;
  }
</style>

export interface HPPoint {
  tMs: number;
  hp: number;
  hpPct: number;
}

export interface FightMarkers {
  championKillsInWindow: number;
  championKillsNearPit: number;
  distanceThreshold: number;
}

export interface DragonHPResponse {
  matchId: string;
  dragonIndex: number;
  dragonType: string;
  killMs: number;
  startMs: number;
  hpModel: {
    isEstimated: boolean;
    patchBasis: string;
    levelEstimate: number;
    hp0: number;
    tickMs: number;
    windowMs: number;
    burstMs: number;
    burstStartHpPct: number;
    model: string;
    sourceGameLength: number;
  };
  series: HPPoint[];
  markers: FightMarkers;
}

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export interface FetchDragonHPParams {
  matchId: string;
  dragonIndex: number;
  tickMs?: number;
  windowMs?: number;
}

export function getApiBaseUrl(): string {
  const configured = (process.env.VUE_APP_API_BASE_URL || '').trim();
  const base = configured || 'http://localhost:8080';
  return base.replace(/\/+$/, '');
}

export async function fetchDragonHP(params: FetchDragonHPParams): Promise<DragonHPResponse> {
  const { matchId, dragonIndex, tickMs, windowMs } = params;

  const routeMatchID = encodeURIComponent(matchId.trim());
  const url = new URL(`${getApiBaseUrl()}/v1/matches/${routeMatchID}/dragons/${dragonIndex}/hp`);

  if (tickMs && tickMs > 0) {
    url.searchParams.set('tickMs', String(tickMs));
  }
  if (windowMs && windowMs > 0) {
    url.searchParams.set('windowMs', String(windowMs));
  }

  const response = await fetch(url.toString(), {
    method: 'GET',
    headers: {
      Accept: 'application/json',
    },
  });

  let payload: unknown = null;
  try {
    payload = await response.json();
  } catch (err) {
    payload = null;
  }

  if (!response.ok) {
    const errorMessage =
      typeof payload === 'object' &&
      payload !== null &&
      'error' in payload &&
      typeof (payload as { error?: unknown }).error === 'string'
        ? ((payload as { error: string }).error as string)
        : `Request failed with status ${response.status}`;
    const message = errorMessage;
    throw new ApiError(message, response.status);
  }

  return payload as DragonHPResponse;
}

import { useEffect, useMemo, useState } from 'react'

type HealthState = 'loading' | 'ready' | 'error'

type HealthResponse = {
  service: string
  status: string
  message: string
  checkedAt: string
}

const navItems = [
  'Niches',
  'Trends',
  'Weekly Plan',
  'Drafts',
  'Review',
  'Uploads',
  'Analytics',
  'Settings',
]

const queues = [
  {
    label: 'Niches',
    title: 'No niches defined',
    detail: 'Create the first channel niche to start topic discovery.',
  },
  {
    label: 'Trends',
    title: 'Signals pending',
    detail: 'Trend and search providers will appear here after setup.',
  },
  {
    label: 'Weekly Plan',
    title: 'No plan generated',
    detail: 'The weekly queue starts after niche candidates are available.',
  },
  {
    label: 'Review',
    title: 'No drafts waiting',
    detail: 'Human approval stays required before upload workflows begin.',
  },
]

function App() {
  const [health, setHealth] = useState<HealthResponse | null>(null)
  const [healthState, setHealthState] = useState<HealthState>('loading')

  useEffect(() => {
    const controller = new AbortController()

    async function loadHealth() {
      try {
        setHealthState('loading')

        const response = await fetch('/api/health', {
          headers: { Accept: 'application/json' },
          signal: controller.signal,
        })

        if (!response.ok) {
          throw new Error(`Health request failed with ${response.status}`)
        }

        const payload = (await response.json()) as HealthResponse
        setHealth(payload)
        setHealthState('ready')
      } catch (error) {
        if (!controller.signal.aborted) {
          console.error(error)
          setHealth(null)
          setHealthState('error')
        }
      }
    }

    void loadHealth()

    return () => controller.abort()
  }, [])

  const checkedAt = useMemo(() => {
    if (!health?.checkedAt) {
      return 'Waiting for API'
    }

    return new Intl.DateTimeFormat(undefined, {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(new Date(health.checkedAt))
  }, [health?.checkedAt])

  return (
    <main className="app-shell">
      <aside className="sidebar" aria-label="Primary navigation">
        <div className="brand-block">
          <p className="eyebrow">AutoTube</p>
          <h1>Production Desk</h1>
        </div>
        <nav>
          {navItems.map((item) => (
            <a href="/" key={item}>
              {item}
            </a>
          ))}
        </nav>
      </aside>

      <section className="workspace" aria-labelledby="workspace-title">
        <div className="workspace-header">
          <div>
            <p className="eyebrow">Milestone 1</p>
            <h2 id="workspace-title">Project foundation</h2>
          </div>
          <StatusBadge healthState={healthState} />
        </div>

        <section className="status-band" aria-label="Local API status">
          <div>
            <span>API status</span>
            <strong>{healthState === 'ready' ? health?.service : 'autotube-api'}</strong>
            <p>{statusMessage(healthState, health)}</p>
          </div>
          <dl>
            <div>
              <dt>Checked</dt>
              <dd>{checkedAt}</dd>
            </div>
            <div>
              <dt>Publishing</dt>
              <dd>Approval gated</dd>
            </div>
            <div>
              <dt>Cadence</dt>
              <dd>2-3 videos/week</dd>
            </div>
          </dl>
        </section>

        <section className="work-grid" aria-label="Milestone work areas">
          {queues.map((queue) => (
            <article key={queue.label}>
              <span>{queue.label}</span>
              <strong>{queue.title}</strong>
              <p>{queue.detail}</p>
            </article>
          ))}
        </section>
      </section>
    </main>
  )
}

function StatusBadge({ healthState }: { healthState: HealthState }) {
  const label =
    healthState === 'ready' ? 'API online' : healthState === 'loading' ? 'Checking API' : 'API offline'

  return <span className={`status-pill status-pill--${healthState}`}>{label}</span>
}

function statusMessage(healthState: HealthState, health: HealthResponse | null) {
  if (healthState === 'loading') {
    return 'Checking the local Go API through the dashboard proxy.'
  }

  if (healthState === 'error') {
    return 'Start the API locally and refresh the dashboard status.'
  }

  return health?.message ?? 'AutoTube API is ready.'
}

export default App

import { useEffect, useState } from 'react'
import { fetchHealth, type ApiHealth } from './health'

type Connection = { status: 'loading' } | { status: 'available'; health: ApiHealth } | { status: 'failed' }
type Page = 'overview' | 'connection'

export default function App() {
  const [page, setPage] = useState<Page>('overview')
  const [attempt, setAttempt] = useState(0)
  const [connection, setConnection] = useState<Connection>({ status: 'loading' })
  useEffect(() => {
    const controller = new AbortController()
    setConnection({ status: 'loading' })
    void fetchHealth(controller.signal).then((health) => {
      if (!controller.signal.aborted) setConnection({ status: 'available', health })
    }).catch(() => {
      if (!controller.signal.aborted) setConnection({ status: 'failed' })
    })
    return () => { controller.abort() }
  }, [attempt])

  return (
    <div className="workspace">
      <header>
        <a className="brand" href="#overview" onClick={() => { setPage('overview') }}>AutoTube <span>v2</span></a>
        <p>Local operations</p>
      </header>
      <nav aria-label="Operations">
        <button aria-current={page === 'overview' ? 'page' : undefined} onClick={() => { setPage('overview') }}>Overview</button>
        <button aria-current={page === 'connection' ? 'page' : undefined} onClick={() => { setPage('connection') }}>API connection</button>
      </nav>
      <main id="overview">
        <p className="eyebrow">Operations</p>
        <h1>{page === 'overview' ? 'Workspace overview' : 'API connection'}</h1>
        {page === 'overview' && <p>The first milestone is an inspectable, evidence-backed topic backlog. Workflow controls will appear here as they become available.</p>}
        <section className="connection" aria-labelledby="connection-title" aria-busy={connection.status === 'loading'}>
          <h2 id="connection-title">Local API</h2>
          <div role="status" aria-live="polite">
            {connection.status === 'loading' && <p className="status loading">Checking API availability…</p>}
            {connection.status === 'failed' && <><p className="status failed">API unavailable</p><p>Could not confirm a healthy API response. Start the local API and try again.</p></>}
            {connection.status === 'available' && <><p className="status available">API available at last check</p><p>Checked <time dateTime={connection.health.checkedAt}>{new Date(connection.health.checkedAt).toLocaleTimeString()}</time></p></>}
          </div>
          <button className="primary" disabled={connection.status === 'loading'} onClick={() => { setAttempt((value) => value + 1) }}>
            {connection.status === 'failed' ? 'Retry connection' : 'Check again'}
          </button>
          <p className="note">This checks the API process only. Database, worker and provider readiness are not included.</p>
          {page === 'connection' && <dl><dt>Health endpoint</dt><dd><a href="/api/health">/api/health</a></dd><dt>Service</dt><dd>{connection.status === 'available' ? connection.health.service : 'Not confirmed'}</dd></dl>}
        </section>
        {page === 'overview' && <p className="note">Discovery, video production and publishing are not enabled.</p>}
      </main>
    </div>
  )
}

const navItems = ['Niches', 'Trends', 'Planner', 'Drafts', 'Review']

function App() {
  return (
    <main className="app-shell">
      <aside className="sidebar" aria-label="Primary navigation">
        <div>
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
          <span className="status-pill">Local first</span>
        </div>

        <div className="summary-grid">
          <article>
            <span>Cadence</span>
            <strong>2-3 videos/week</strong>
          </article>
          <article>
            <span>Review gate</span>
            <strong>Human approval</strong>
          </article>
          <article>
            <span>Asset policy</span>
            <strong>Original only</strong>
          </article>
        </div>
      </section>
    </main>
  )
}

export default App

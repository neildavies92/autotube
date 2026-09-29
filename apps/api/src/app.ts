import Fastify from 'fastify'

// Constructing the app has no network, database or provider side effects.
export function buildApp({ logger = false } = {}) {
  const app = Fastify({ logger })
  app.get('/health', () => ({
    service: 'autotube-api',
    status: 'ok',
    message: 'AutoTube API foundation is running.',
    checkedAt: new Date().toISOString(),
  }))
  return app
}

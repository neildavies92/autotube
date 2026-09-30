import { buildApp } from './app.js'
import { ConfigurationError, readConfig } from './config.js'

async function main() {
  let config: ReturnType<typeof readConfig>
  try {
    config = readConfig(process.env)
  } catch (error) {
    console.error(JSON.stringify({ level: 'error', message: error instanceof ConfigurationError ? error.message : 'Invalid configuration.' }))
    process.exitCode = 1
    return
  }
  const app = buildApp({ logger: true, logLevel: config.logLevel })
  let closing = false
  async function shutdown() {
    if (closing) return
    closing = true
    const deadline = setTimeout(() => {
      app.log.error('API shutdown timed out')
      process.exit(1)
    }, 10_000)
    deadline.unref()
    try {
      await app.close()
      app.log.info('API stopped')
    } catch {
      app.log.error('API shutdown failed')
      process.exitCode = 1
    } finally {
      clearTimeout(deadline)
    }
  }
  process.once('SIGINT', () => { void shutdown() })
  process.once('SIGTERM', () => { void shutdown() })
  try {
    await app.listen({ host: config.host, port: config.port })
  } catch (error) {
    const safeCodes = new Set(['EADDRINUSE', 'EACCES', 'EPERM', 'EADDRNOTAVAIL', 'ENOTFOUND'])
    const candidate = error && typeof error === 'object' && 'code' in error ? error.code : undefined
    const code = typeof candidate === 'string' && safeCodes.has(candidate) ? candidate : 'STARTUP_FAILED'
    app.log.error({ code }, 'API startup failed')
    await shutdown()
    process.exitCode = 1
  }
}

await main()

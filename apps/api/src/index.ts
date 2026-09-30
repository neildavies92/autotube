import { buildApp } from './app.js'
import { readConfig } from './config.js'

const app = buildApp({ logger: true })

async function shutdown() {
  try {
    await app.close()
  } catch (error) {
    app.log.error(error)
    process.exitCode = 1
  }
}

try {
  const config = readConfig(process.env)
  await app.listen(config)
  process.once('SIGINT', () => { void shutdown() })
  process.once('SIGTERM', () => { void shutdown() })
} catch (error) {
  app.log.error(error)
  await app.close()
  process.exitCode = 1
}

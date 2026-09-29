// Keep a separate process entrypoint; durable jobs arrive in SWA-49.
const idle = setInterval(() => {}, 60_000)
console.info('AutoTube worker started (idle; no jobs configured).')

function shutdown() {
  clearInterval(idle)
  console.info('AutoTube worker stopped.')
}

process.once('SIGINT', shutdown)
process.once('SIGTERM', shutdown)

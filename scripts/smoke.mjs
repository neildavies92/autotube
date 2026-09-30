import assert from 'node:assert/strict'
import { spawn } from 'node:child_process'
import { fileURLToPath } from 'node:url'

const root = fileURLToPath(new URL('../', import.meta.url))
const children = []
// Keep verification independent of developer/provider credentials.
const env = { PATH: process.env.PATH, SYSTEMROOT: process.env.SYSTEMROOT, NODE_ENV: 'test', AUTOTUBE_HOST: '127.0.0.1', AUTOTUBE_PORT: '8080' }

function start(args, extraEnv = {}, cwd = root) {
  const child = spawn(process.execPath, args, { cwd, env: { ...env, ...extraEnv }, stdio: ['ignore', 'pipe', 'pipe'] })
  const state = { child, output: '', exited: false, code: null, signal: null }
  child.stdout.on('data', (data) => { state.output += String(data) })
  child.stderr.on('data', (data) => { state.output += String(data) })
  child.on('error', (error) => { state.output += error.message; state.exited = true })
  child.on('exit', (code, signal) => { Object.assign(state, { exited: true, code, signal }) })
  children.push(state)
  return state
}

async function until(predicate, description, state) {
  const deadline = Date.now() + 15_000
  while (!predicate()) {
    if (Date.now() > deadline || state?.exited) throw new Error(`${description}\n${state?.output ?? ''}`)
    await new Promise((resolve) => setTimeout(resolve, 50))
  }
}

async function stop(state, signal = 'SIGTERM', vite = false) {
  state.child.kill(signal)
  await until(() => state.exited, 'Process did not shut down', state)
  if (vite) {
    // Vite uses the conventional 128 + signal exit status. Our apps must exit 0.
    assert.ok(state.code === 0 || state.code === 143 || state.signal === signal, state.output)
  } else {
    assert.equal(state.code, 0, state.output)
    assert.equal(state.signal, null, 'Expected handled shutdown, not signal termination')
  }
}

async function get(url) {
  const response = await fetch(url, { signal: AbortSignal.timeout(5_000) })
  assert.equal(response.status, 200, String(url))
  return response
}

try {
  const invalid = start(['apps/api/dist/index.js'], { AUTOTUBE_PORT: 'invalid' })
  await until(() => invalid.exited, 'Invalid configuration did not exit')
  assert.equal(invalid.code, 1)
  assert.match(invalid.output, /Invalid configuration: AUTOTUBE_PORT/)

  const api = start(['apps/api/dist/index.js'])
  await until(() => api.output.includes('Server listening at'), 'API failed to start (port 8080 must be free)', api)
  assert.equal((await (await get('http://127.0.0.1:8080/health')).json()).status, 'ok')

  const conflict = start(['apps/api/dist/index.js'])
  await until(() => conflict.exited, 'Port conflict did not exit')
  assert.equal(conflict.code, 1)
  assert.match(conflict.output, /EADDRINUSE/)

  const worker = start(['apps/worker/dist/index.js'])
  await until(() => worker.output.includes('worker started'), 'Worker failed to start', worker)
  assert.equal(worker.exited, false)

  const vite = fileURLToPath(new URL('../node_modules/vite/bin/vite.js', import.meta.url))
  const webRoot = fileURLToPath(new URL('../apps/web', import.meta.url))
  const web = start([vite], {}, webRoot)
  await until(() => web.output.includes('http://127.0.0.1:5173'), 'Web dev server failed to start (port 5173 must be free)', web)
  assert.match(await (await get('http://127.0.0.1:5173')).text(), /<title>AutoTube<\/title>/)
  assert.match(await (await get('http://127.0.0.1:5173/src/App.tsx')).text(), /Development foundation/)
  assert.equal((await (await get('http://127.0.0.1:5173/api/health')).json()).service, 'autotube-api')
  await stop(web, 'SIGTERM', true)

  const preview = start([vite, 'preview'], {}, webRoot)
  await until(() => preview.output.includes('http://127.0.0.1:4173'), 'Web preview failed to start (port 4173 must be free)', preview)
  const html = await (await get('http://127.0.0.1:4173')).text()
  const asset = html.match(/src="([^"]+\.js)"/)?.[1]
  assert.ok(asset, 'Built HTML must reference a JavaScript asset')
  await get(new URL(asset, 'http://127.0.0.1:4173'))
  await stop(preview, 'SIGTERM', true)
  await stop(worker, 'SIGINT')
  assert.match(worker.output, /worker stopped/)
  await stop(api)
  console.info('Smoke passed: API, worker, web dev/proxy, built web, invalid config, port conflict and graceful shutdown.')
} finally {
  for (const state of children) {
    if (!state.exited) state.child.kill('SIGTERM')
  }
  await new Promise((resolve) => setTimeout(resolve, 250))
  for (const state of children) {
    if (!state.exited) state.child.kill('SIGKILL')
  }
}

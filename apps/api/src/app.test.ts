import { Writable } from 'node:stream'
import { z } from 'zod'
import { afterEach, describe, expect, it } from 'vitest'
import { buildApp } from './app.js'
import { readConfig } from './config.js'

describe('credential-free API foundation', () => {
  const apps: ReturnType<typeof buildApp>[] = []
  afterEach(async () => { await Promise.all(apps.splice(0).map((app) => app.close())) })

  it('serves the retained health contract without opening a port or configuring providers', async () => {
    const app = buildApp()
    apps.push(app)
    const response = await app.inject({ method: 'GET', url: '/health' })
    expect(response.statusCode).toBe(200)
    expect(response.headers['content-type']).toContain('application/json')
    const body: unknown = response.json()
    expect(body).toEqual({ service: 'autotube-api', status: 'ok', message: expect.any(String), checkedAt: expect.stringMatching(/^\d{4}-\d{2}-\d{2}T/) })
  })

  it('does not expose retired research routes or accept health mutations', async () => {
    const app = buildApp()
    apps.push(app)
    for (const request of [{ method: 'GET', url: '/research/youtube?q=test' }, { method: 'POST', url: '/health' }] as const) {
      expect((await app.inject(request)).statusCode).toBe(404)
    }
  })

  it('starts with local defaults and ignores unrelated provider/database settings', () => {
    expect(readConfig({ DATABASE_URL: 'unused', YOUTUBE_API_KEY: 'unused' })).toEqual({ host: '127.0.0.1', port: 8080, logLevel: 'info' })
    expect(readConfig({ AUTOTUBE_HOST: 'localhost', AUTOTUBE_PORT: '9090' })).toEqual({ host: 'localhost', port: 9090, logLevel: 'info' })
  })

  it.each(['', '0', '-1', '65536', '3.5', 'not-a-port'])('rejects invalid port %j before startup', (port) => {
    expect(() => readConfig({ AUTOTUBE_PORT: port })).toThrow('Invalid configuration: AUTOTUBE_PORT')
  })

  it('rejects blank hosts without exposing input values', () => {
    expect(() => readConfig({ AUTOTUBE_HOST: ' ', AUTOTUBE_PORT: 'secret-value' })).toThrow('Invalid configuration: AUTOTUBE_HOST, AUTOTUBE_PORT')
  })
})


describe('HTTP boundary and safe structured logging', () => {
  const apps: ReturnType<typeof buildApp>[] = []
  afterEach(async () => { await Promise.all(apps.splice(0).map((app) => app.close())) })
  const correlationId = '55ad41a2-9088-4120-8cdb-531ad18f3fe4'

  it('validates query input, returns one error envelope and correlates responses', async () => {
    const app = buildApp()
    apps.push(app)
    for (const [url, status, code] of [['/health?token=private-value', 400, 'INVALID_REQUEST'], ['/missing', 404, 'NOT_FOUND']] as const) {
      const response = await app.inject({ url, headers: { 'x-request-id': correlationId } })
      expect(response.statusCode).toBe(status)
      expect(response.headers['x-request-id']).toBe(correlationId)
      expect(response.json()).toEqual({ error: { code, message: expect.any(String), requestId: correlationId } })
      expect(response.body).not.toContain('private-value')
    }
    const health = await app.inject('/health')
    expect(health.headers['x-request-id']).toMatch(/^[a-f0-9-]{36}$/)
  })

  it.each(['untrusted-request-id', 'x'.repeat(256), 'line\nbreak'])('rejects unsafe correlation input without echoing it', async (requestId) => {
    const app = buildApp()
    apps.push(app)
    const response = await app.inject({ url: '/health', headers: { 'x-request-id': requestId } })
    expect(response.statusCode).toBe(400)
    expect(response.headers['x-request-id']).toMatch(/^[a-f0-9-]{36}$/)
    expect(response.headers['x-request-id']).not.toBe(requestId)
    expect(response.json().error.code).toBe('INVALID_REQUEST')
    expect(response.body).not.toContain(requestId)
  })

  it('contains unexpected errors and omits secret request/error values from logs', async () => {
    let logs = ''
    const logStream = new Writable({ write(chunk, _encoding, done) { logs += String(chunk); done() } })
    const app = buildApp({ logger: true, logStream })
    apps.push(app)
    // Test-only failing handlers exercise the shared boundary, not a product endpoint.
    app.post('/failure', () => { throw new Error('exception-secret') })
    app.get('/broken-contract', () => z.literal('ok').parse('contract-secret'))
    const response = await app.inject({
      method: 'POST', url: '/failure?token=query-secret',
      headers: { authorization: 'Bearer auth-secret', cookie: 'session=cookie-secret', 'x-api-key': 'header-secret', 'x-request-id': correlationId },
      payload: { password: 'body-secret' },
    })
    expect(response.statusCode).toBe(500)
    expect(response.json()).toEqual({ error: { code: 'INTERNAL_ERROR', message: 'Internal server error.', requestId: correlationId } })
    const malformed = await app.inject({ method: 'POST', url: '/failure', headers: { 'content-type': 'application/json' }, payload: '{"password":"parser-secret"' })
    expect(malformed.statusCode).toBe(400)
    expect(malformed.json().error.code).toBe('BAD_REQUEST')
    const brokenContract = await app.inject('/broken-contract')
    expect(brokenContract.statusCode).toBe(500)
    expect(brokenContract.json().error.code).toBe('INTERNAL_ERROR')
    const malformedUrl = await app.inject({ url: '/%E0%A4%A?token=malformed-url-secret', headers: { 'x-request-id': correlationId } })
    expect(malformedUrl.statusCode).toBe(400)
    expect(malformedUrl.headers['x-request-id']).toBe(correlationId)
    expect(malformedUrl.json()).toEqual({ error: { code: 'BAD_REQUEST', message: 'Invalid request.', requestId: correlationId } })
    expect(malformedUrl.body).not.toContain('malformed-url-secret')
    await app.inject('/path-secret?token=unknown-secret')
    app.log.info({ authorization: 'direct-secret', headers: { cookie: 'direct-cookie-secret' }, err: new Error('logged-error-secret') }, 'redaction check')
    for (const secret of ['exception-secret', 'query-secret', 'auth-secret', 'cookie-secret', 'header-secret', 'body-secret', 'parser-secret', 'path-secret', 'unknown-secret', 'direct-secret', 'direct-cookie-secret', 'logged-error-secret', 'malformed-url-secret', 'contract-secret']) {
      expect(logs).not.toContain(secret)
      expect(response.body).not.toContain(secret)
    }
    const entries = logs.trim().split('\n').map((line) => JSON.parse(line) as Record<string, unknown>)
    expect(entries.some((entry) => entry['reqId'] === correlationId && entry['msg'] === 'request failed')).toBe(true)
    expect(entries.some((entry) => entry['statusCode'] === 500 && typeof entry['durationMs'] === 'number')).toBe(true)
    expect(logs).toContain('[REDACTED]')
  })

  it('validates log level without printing environment values', () => {
    expect(readConfig({ AUTOTUBE_LOG_LEVEL: 'warn' }).logLevel).toBe('warn')
    expect(() => readConfig({ AUTOTUBE_LOG_LEVEL: 'secret' })).toThrow('Invalid configuration: AUTOTUBE_LOG_LEVEL')
  })
})

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
    expect(readConfig({ DATABASE_URL: 'unused', YOUTUBE_API_KEY: 'unused' })).toEqual({ host: '127.0.0.1', port: 8080 })
    expect(readConfig({ AUTOTUBE_HOST: 'localhost', AUTOTUBE_PORT: '9090' })).toEqual({ host: 'localhost', port: 9090 })
  })

  it.each(['', '0', '-1', '65536', '3.5', 'not-a-port'])('rejects invalid port %j before startup', (port) => {
    expect(() => readConfig({ AUTOTUBE_PORT: port })).toThrow('Invalid configuration: AUTOTUBE_PORT')
  })

  it('rejects blank hosts without exposing input values', () => {
    expect(() => readConfig({ AUTOTUBE_HOST: ' ', AUTOTUBE_PORT: 'secret-value' })).toThrow('Invalid configuration: AUTOTUBE_HOST, AUTOTUBE_PORT')
  })
})

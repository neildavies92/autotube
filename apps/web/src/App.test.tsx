// @vitest-environment jsdom
import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import App from './App'

const health = { service: 'autotube-api', status: 'ok', message: 'Healthy', checkedAt: '2026-09-30T12:00:00.000Z' }
afterEach(() => { cleanup(); vi.unstubAllGlobals() })

describe('operational API connection', () => {
  it('shows loading until a validated response arrives and provides working navigation', async () => {
    let resolve: (response: Response) => void = () => { throw new Error('Request has not started') }
    const fetchMock = vi.fn(() => new Promise<Response>((done) => { resolve = done }))
    vi.stubGlobal('fetch', fetchMock)
    render(<App />)
    expect(screen.getByRole('status').textContent).toContain('Checking API availability')
    expect(screen.getByRole<HTMLButtonElement>('button', { name: 'Check again' }).disabled).toBe(true)
    await act(async () => { resolve(Response.json(health)) })
    expect(await screen.findByText('API available at last check')).toBeTruthy()
    expect(fetchMock).toHaveBeenCalledWith('/api/health', expect.objectContaining({ cache: 'no-store' }))
    fireEvent.click(screen.getByRole('button', { name: 'API connection' }))
    expect(screen.getByRole('heading', { level: 1 }).textContent).toBe('API connection')
    expect(screen.getByRole('link', { name: '/api/health' }).getAttribute('href')).toBe('/api/health')
  })

  it.each(['network', 'http', 'contract', 'json'])('shows failure for %s errors and recovers on explicit retry', async (failure) => {
    const fetchMock = vi.fn()
    if (failure === 'network') fetchMock.mockRejectedValueOnce(new Error('private connection details'))
    else if (failure === 'http') fetchMock.mockResolvedValueOnce(Response.json(health, { status: 503 }))
    else if (failure === 'json') fetchMock.mockResolvedValueOnce(new Response('not-json'))
    else fetchMock.mockResolvedValueOnce(Response.json({ ...health, status: 'broken' }))
    fetchMock.mockResolvedValueOnce(Response.json(health))
    vi.stubGlobal('fetch', fetchMock)
    render(<App />)
    expect(await screen.findByText('API unavailable')).toBeTruthy()
    expect(screen.queryByText(/private connection details/)).toBeNull()
    fireEvent.click(screen.getByRole('button', { name: 'Retry connection' }))
    expect(await screen.findByText('API available at last check')).toBeTruthy()
    expect(fetchMock).toHaveBeenCalledTimes(2)
  })

  it('aborts pending requests when the shell unmounts', () => {
    let signal: AbortSignal | undefined
    vi.stubGlobal('fetch', vi.fn((_url: string, init: RequestInit) => {
      signal = init.signal ?? undefined
      return new Promise<Response>(() => {})
    }))
    const { unmount } = render(<App />)
    expect(signal?.aborted).toBe(false)
    unmount()
    expect(signal?.aborted).toBe(true)
  })
})

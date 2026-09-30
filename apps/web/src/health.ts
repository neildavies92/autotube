import { z } from 'zod'

const healthSchema = z.object({
  service: z.literal('autotube-api'), status: z.literal('ok'),
  message: z.string(), checkedAt: z.iso.datetime(),
})
export type ApiHealth = z.infer<typeof healthSchema>

export async function fetchHealth(signal: AbortSignal): Promise<ApiHealth> {
  const response = await fetch('/api/health', {
    signal: AbortSignal.any([signal, AbortSignal.timeout(5_000)]),
    headers: { accept: 'application/json' }, cache: 'no-store',
  })
  if (!response.ok) throw new Error('API health request failed.')
  const body: unknown = await response.json()
  return healthSchema.parse(body)
}

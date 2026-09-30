import { z } from 'zod'

const configSchema = z.object({
  AUTOTUBE_HOST: z.string().trim().min(1).default('127.0.0.1'),
  AUTOTUBE_PORT: z.coerce.number().int().min(1).max(65535).default(8080),
  AUTOTUBE_LOG_LEVEL: z.enum(['fatal', 'error', 'warn', 'info', 'debug', 'trace', 'silent']).default('info'),
})

export class ConfigurationError extends Error {
  constructor(readonly fields: string[]) {
    super(`Invalid configuration: ${fields.join(', ')}`)
  }
}

export function readConfig(env: NodeJS.ProcessEnv) {
  const result = configSchema.safeParse(env)
  if (!result.success) {
    // Report field names only, never environment values.
    throw new ConfigurationError(result.error.issues.map((issue) => issue.path.join('.')))
  }
  return { host: result.data.AUTOTUBE_HOST, port: result.data.AUTOTUBE_PORT, logLevel: result.data.AUTOTUBE_LOG_LEVEL }
}

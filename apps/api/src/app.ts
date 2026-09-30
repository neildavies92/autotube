import { randomUUID } from 'node:crypto'
import type { Writable } from 'node:stream'
import Fastify, { LogController, type FastifyReply } from 'fastify'
import { z } from 'zod'

class InvalidRequestError extends Error {}

const requestIdSchema = z.uuid()
const healthQuerySchema = z.strictObject({})
const healthResponseSchema = z.strictObject({
  service: z.literal('autotube-api'), status: z.literal('ok'),
  message: z.string(), checkedAt: z.iso.datetime(),
})

type AppOptions = {
  logger?: boolean
  logLevel?: 'fatal' | 'error' | 'warn' | 'info' | 'debug' | 'trace' | 'silent'
  logStream?: Writable
}

// Constructing the app has no network, database or provider side effects.
export function buildApp({ logger = false, logLevel = 'info', logStream }: AppOptions = {}) {
  const app = Fastify({
    logger: logger ? {
      level: logLevel,
      ...(logStream ? { stream: logStream } : {}),
      redact: {
        paths: ['authorization', 'cookie', 'password', 'token', 'apiKey', 'secret',
          'headers.authorization', 'headers.cookie', 'req.headers', 'req.body', 'req.query', 'res.headers'],
        censor: '[REDACTED]',
      },
      // Log allowlisted metadata only: even malformed URLs and exception messages
      // can contain credentials. Do not serialize request/response payloads.
      serializers: {
        req: (request: { method: string }) => ({ method: request.method }),
        res: (response: { statusCode: number }) => ({ statusCode: response.statusCode }),
        err: () => ({ type: 'Error', message: 'Error details omitted', stack: '' }),
      },
    } : false,
    logController: new LogController({ disableRequestLogging: true }),
    // Router failures run before onRequest/errorHandler (e.g. invalid URL encoding).
    frameworkErrors: (error, request, reply: FastifyReply) => {
      const statusCode = error.statusCode && error.statusCode >= 400 && error.statusCode < 500 ? error.statusCode : 500
      const code = statusCode >= 500 ? 'INTERNAL_ERROR' : 'BAD_REQUEST'
      request.log.warn({ code, statusCode }, 'request rejected before routing')
      void reply.header('x-request-id', request.id).code(statusCode).send({ error: {
        code, message: statusCode >= 500 ? 'Internal server error.' : 'Invalid request.', requestId: request.id,
      } })
    },
    requestIdHeader: false,
    genReqId: (request) => {
      const parsed = requestIdSchema.safeParse(request.headers['x-request-id'])
      return parsed.success ? parsed.data : randomUUID()
    },
  })

  app.addHook('onRequest', async (request, reply) => {
    reply.header('x-request-id', request.id)
    request.log.info({ method: request.method }, 'request started')
    const supplied = request.headers['x-request-id']
    if (supplied !== undefined && !requestIdSchema.safeParse(supplied).success) {
      return reply.code(400).send({ error: { code: 'INVALID_REQUEST', message: 'Invalid request.', requestId: request.id } })
    }
  })
  app.addHook('onResponse', async (request, reply) => {
    request.log.info({ statusCode: reply.statusCode, durationMs: reply.elapsedTime }, 'request completed')
  })
  app.setErrorHandler((error: Error & { statusCode?: number }, request, reply) => {
    const invalid = error instanceof InvalidRequestError
    const statusCode = invalid ? 400 : error.statusCode && error.statusCode >= 400 && error.statusCode < 500 ? error.statusCode : 500
    const code = invalid ? 'INVALID_REQUEST' : statusCode >= 500 ? 'INTERNAL_ERROR' : 'BAD_REQUEST'
    const message = statusCode >= 500 ? 'Internal server error.' : 'Invalid request.'
    // Never log raw exceptions: parser/validation errors can echo user input.
    request.log[statusCode >= 500 ? 'error' : 'warn']({ code, statusCode }, 'request failed')
    void reply.code(statusCode).send({ error: { code, message, requestId: request.id } })
  })
  app.setNotFoundHandler((request, reply) => {
    void reply.code(404).send({ error: { code: 'NOT_FOUND', message: 'Route not found.', requestId: request.id } })
  })
  app.get('/health', (request) => {
    if (!healthQuerySchema.safeParse(request.query).success) throw new InvalidRequestError()
    return healthResponseSchema.parse({
      service: 'autotube-api', status: 'ok',
      message: 'AutoTube API foundation is running.', checkedAt: new Date().toISOString(),
    })
  })
  return app
}

import { pgSchema } from 'drizzle-orm/pg-core'

// Domain tables arrive in their owning issues. Graphile Worker owns a separate schema.
export const appSchema = pgSchema('app')

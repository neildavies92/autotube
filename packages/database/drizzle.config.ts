import { defineConfig } from 'drizzle-kit'

// Generation is offline. Applying migrations uses the validated, finally-closed client.
export default defineConfig({
  dialect: 'postgresql',
  schema: './src/schema.ts',
  out: './migrations',
  strict: true,
  verbose: false,
})

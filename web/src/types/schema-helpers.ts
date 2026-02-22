type JsonPrimitive = string | number | boolean | null

function isJsonPrimitive(value: unknown): value is JsonPrimitive {
  return (
    value === null ||
    typeof value === 'string' ||
    typeof value === 'number' ||
    typeof value === 'boolean'
  )
}

export function isObjectRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

export function expectObjectRecord(value: unknown, schemaName: string): Record<string, unknown> {
  if (!isObjectRecord(value)) {
    throw new Error(`${schemaName} response must be an object`)
  }
  return value
}

export function expectStringField(
  payload: Record<string, unknown>,
  fieldName: string,
  schemaName: string,
): string {
  const value = payload[fieldName]
  if (typeof value !== 'string') {
    throw new Error(`${schemaName}.${fieldName} must be a string`)
  }
  return value
}

export function expectOptionalStringField(
  payload: Record<string, unknown>,
  fieldName: string,
  schemaName: string,
): string | undefined {
  const value = payload[fieldName]
  if (value === undefined) {
    return undefined
  }
  if (typeof value !== 'string') {
    throw new Error(`${schemaName}.${fieldName} must be a string when provided`)
  }
  return value
}

export function expectBooleanField(
  payload: Record<string, unknown>,
  fieldName: string,
  schemaName: string,
): boolean {
  const value = payload[fieldName]
  if (typeof value !== 'boolean') {
    throw new Error(`${schemaName}.${fieldName} must be a boolean`)
  }
  return value
}

export function expectFiniteNumberField(
  payload: Record<string, unknown>,
  fieldName: string,
  schemaName: string,
): number {
  const value = payload[fieldName]
  if (typeof value !== 'number' || !Number.isFinite(value)) {
    throw new Error(`${schemaName}.${fieldName} must be a finite number`)
  }
  return value
}

export function expectIntegerField(
  payload: Record<string, unknown>,
  fieldName: string,
  schemaName: string,
): number {
  const value = expectFiniteNumberField(payload, fieldName, schemaName)
  if (!Number.isInteger(value)) {
    throw new Error(`${schemaName}.${fieldName} must be an integer`)
  }
  return value
}

export function ensureJsonValue(value: unknown, fieldPath: string): unknown {
  if (isJsonPrimitive(value)) {
    return value
  }

  if (Array.isArray(value)) {
    return value.map((item, index) => ensureJsonValue(item, `${fieldPath}[${index}]`))
  }

  if (isObjectRecord(value)) {
    const normalized: Record<string, unknown> = {}
    for (const [key, nested] of Object.entries(value)) {
      normalized[key] = ensureJsonValue(nested, `${fieldPath}.${key}`)
    }
    return normalized
  }

  throw new Error(`${fieldPath} must contain JSON-serializable values`)
}

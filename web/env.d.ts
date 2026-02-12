/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string
  readonly VITE_API_MODE?: 'real' | 'mock'
  readonly VITE_MOCK_AUTH?: 'true' | 'false'
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

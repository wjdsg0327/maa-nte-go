const BASE_URL = '/api'

export async function api(url, options = {}) {
  const res = await fetch(BASE_URL + url, {
    headers: { 'Content-Type': 'application/json' },
    ...options
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || '请求失败')
  return data
}

export const pipelineApi = {
  list: () => api('/pipelines'),
  get: (name) => api(`/pipelines/${name}`),
  create: (name, content) => api('/pipelines', { method: 'POST', body: JSON.stringify({ name, content }) }),
  update: (name, content) => api(`/pipelines/${name}`, { method: 'PUT', body: JSON.stringify({ content }) }),
  delete: (name) => api(`/pipelines/${name}`, { method: 'DELETE' })
}

export const taskApi = {
  execute: (task, node) => api('/tasks', { method: 'POST', body: JSON.stringify({ task, node }) }),
  run: (task, node) => api('/tasks/run', { method: 'POST', body: JSON.stringify({ task, node }) }),
  status: () => api('/tasks/status'),
  stop: () => api('/tasks/stop', { method: 'POST' })
}

export const windowApi = {
  list: () => api('/windows'),
  connect: (handle) => api('/windows/connect', { method: 'POST', body: JSON.stringify({ handle }) })
}

export const appApi = {
  config: () => api('/config')
}

export const resourceApi = {
  images: () => api('/images'),
  ocrModels: () => api('/models/ocr'),
  detectModels: () => api('/models/detect'),
  nodes: () => api('/nodes'),
  screenshot: () => api('/screenshot'),
  reload: () => api('/resources/reload', { method: 'POST' })
}

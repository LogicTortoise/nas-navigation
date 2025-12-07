import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000
})

// 响应拦截器
api.interceptors.response.use(
  response => response.data,
  error => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

// 分类API
export const categoryApi = {
  getAll: (withLinks = false) => api.get('/categories', { params: { with_links: withLinks } }),
  getById: (id) => api.get(`/categories/${id}`),
  create: (data) => api.post('/categories', data),
  update: (id, data) => api.put(`/categories/${id}`, data),
  delete: (id) => api.delete(`/categories/${id}`)
}

// 链接API
export const linkApi = {
  getAll: (categoryId) => api.get('/links', { params: categoryId ? { category_id: categoryId } : {} }),
  search: (query) => api.get('/links/search', { params: { q: query } }),
  getById: (id) => api.get(`/links/${id}`),
  create: (data) => api.post('/links', data),
  update: (id, data) => api.put(`/links/${id}`, data),
  delete: (id) => api.delete(`/links/${id}`)
}

// 脚本API
export const scriptApi = {
  getAll: () => api.get('/scripts'),
  getById: (id) => api.get(`/scripts/${id}`),
  create: (data) => api.post('/scripts', data),
  update: (id, data) => api.put(`/scripts/${id}`, data),
  delete: (id) => api.delete(`/scripts/${id}`),
  execute: (id) => api.post(`/scripts/${id}/execute`)
}

// 设置API
export const settingApi = {
  getAll: () => api.get('/settings'),
  update: (data) => api.put('/settings', data)
}

// 系统API
export const systemApi = {
  getInfo: () => api.get('/system/info'),
  restart: (target = '', confirm = true) => api.post('/system/restart', { target, confirm }),
  restartAll: () => api.post('/system/restart-all'),
  quickAction: (action) => api.post(`/system/quick/${action}`)
}

export default api

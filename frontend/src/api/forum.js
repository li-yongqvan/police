import api from './index'

export function getBoards() {
  return api.get('/boards')
}

export function getPosts(params) {
  return api.get('/posts', { params })
}

export function getPost(id) {
  return api.get(`/posts/${id}`)
}

export function createPost(data) {
  if (data instanceof FormData) {
    return api.post('/posts', data, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  }
  return api.post('/posts', data)
}

export function updatePost(id, data) {
  return api.put(`/posts/${id}`, data)
}

export function deletePost(id) {
  return api.delete(`/posts/${id}`)
}

export function getComments(postId, page = 1) {
  return api.get(`/posts/${postId}/comments`, { params: { page } })
}

export function createComment(postId, data) {
  return api.post(`/posts/${postId}/comments`, data)
}

export function likePost(postId) {
  return api.post(`/posts/${postId}/like`)
}

export function collectPost(postId) {
  return api.post(`/posts/${postId}/collect`)
}

export function uploadAttachment(formData) {
  return api.post('/attachments/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function getAttachment(id) {
  return api.get(`/attachments/${id}`)
}

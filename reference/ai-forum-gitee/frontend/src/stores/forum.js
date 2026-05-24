import { defineStore } from 'pinia'
import { getBoards, getPosts, getPost, createPost, updatePost, deletePost, getComments, createComment, likePost, collectPost } from '../api/forum'

export const useForumStore = defineStore('forum', {
  state: () => ({
    boards: [],
    posts: [],
    currentPost: null,
    comments: [],
    totalPosts: 0,
    totalComments: 0,
    loading: false,
  }),
  actions: {
    async loadBoards() {
      try {
        const { data } = await getBoards()
        this.boards = data
      } catch (e) {
        console.error('Failed to load boards:', e)
      }
    },
    async loadPosts(boardId, page = 1) {
      this.loading = true
      try {
        const params = { page }
        if (boardId) params.board_id = boardId
        const { data } = await getPosts(params)
        this.posts = data.posts
        this.totalPosts = data.total
      } catch (e) {
        console.error('Failed to load posts:', e)
      } finally {
        this.loading = false
      }
    },
    async loadPost(id) {
      this.loading = true
      try {
        const { data } = await getPost(id)
        this.currentPost = data
      } catch (e) {
        console.error('Failed to load post:', e)
      } finally {
        this.loading = false
      }
    },
    async createPost(data) {
      const { data: post } = await createPost(data)
      return post
    },
    async updatePost(id, data) {
      const { data: post } = await updatePost(id, data)
      return post
    },
    async deletePost(id) {
      await deletePost(id)
    },
    async loadComments(postId, page = 1) {
      try {
        const { data } = await getComments(postId, page)
        this.comments = data.comments
        this.totalComments = data.total
      } catch (e) {
        console.error('Failed to load comments:', e)
      }
    },
    async createComment(postId, content) {
      const { data } = await createComment(postId, { content })
      this.comments.push(data)
      this.totalComments++
      return data
    },
    async likePost(postId) {
      const { data } = await likePost(postId)
      if (this.currentPost) {
        this.currentPost.like_count = data.like_count
        this.currentPost.user_liked = data.liked
      }
      return data
    },
    async collectPost(postId) {
      const { data } = await collectPost(postId)
      if (this.currentPost) {
        this.currentPost.user_collected = data.collected
      }
      return data
    },
  },
})

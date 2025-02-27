import { createRouter, createWebHistory } from 'vue-router'
import PlaylistView from '../views/PlaylistView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'playlist',
      component: PlaylistView,
    },
    {
      path: '/upload',
      name: 'upload-song',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/UploadSongView.vue'),
    },
  ],
})

export default router

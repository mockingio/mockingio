import { createRouter, createWebHistory } from 'vue-router'
import MockView from '../views/MockView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: MockView
    }
  ]
})

export default router

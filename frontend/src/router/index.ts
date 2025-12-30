import { createRouter, createWebHistory } from 'vue-router'
import Login from 'views/Login.vue'
import Dashboard from './views/Dashboard.vue'

const routes = [
  { path: '/login', component: Login },
  { path: '/dashboard', component: Dashboard },
  { path: '/', redirect: '/login' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

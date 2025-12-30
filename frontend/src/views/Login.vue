<template>
  <div class="login-container">
    <h1>Login</h1>
    <form @submit.prevent="login">
      <input v-model="username" placeholder="Username" required />
      <input type="password" v-model="password" placeholder="Password" required />
      <button type="submit">Login</button>
    </form>
    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import request from '@/utils/request'
import { useRouter } from 'vue-router'

const username = ref('')
const password = ref('')
const error = ref('')
const router = useRouter()

const login = async () => {
  try {
    const res = await request.post('/api/login', { username: username.value, password: password.value })
    if (res.data.success) {
      router.push('/dashboard')
    } else {
      error.value = 'Login failed'
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Login error'
  }
}
</script>

<style scoped>
.login-container { max-width: 400px; margin: 50px auto; }
input { display: block; margin: 10px 0; width: 100%; padding: 8px; }
button { padding: 8px 16px; }
</style>

<template>
  <el-form @submit.prevent="login">
    <el-input v-model="form.username" placeholder="Username" />
    <el-input v-model="form.password" type="password" placeholder="Password" />
    <el-button type="primary" @click="login">Login</el-button>
  </el-form>
</template>

<script setup>
import request from '@/utils/request'

const form = reactive({
  username: '',
  password: ''
})

const login = async () => {
  await request.post('/login', form)
  location.href = '/'
}
</script>

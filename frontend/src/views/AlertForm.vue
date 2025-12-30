<template>
  <form @submit.prevent="sendAlert">
    <h2>Send Alert</h2>
    <input v-model="token" placeholder="Token" required />
    <input v-model="alertname" placeholder="Alert Name" />
    <input v-model="severity" placeholder="Severity" />
    <input v-model="instance" placeholder="Instance" />
    <textarea v-model="message" placeholder="Message"></textarea>
    <select v-model="status">
      <option value="firing">firing</option>
      <option value="resolved">resolved</option>
    </select>
    <button type="submit">Send</button>
  </form>
</template>

<script setup>
import { ref } from 'vue'
import request from '@/utils/request'
import { defineEmits } from 'vue'

const emit = defineEmits(['sent'])

const token = ref('')
const alertname = ref('')
const severity = ref('')
const instance = ref('')
const message = ref('')
const status = ref('firing')

const sendAlert = async () => {
  try {
    await request.post('/api/alert/send', {
      token: token.value,
      alertname: alertname.value,
      severity: severity.value,
      instance: instance.value,
      message: message.value,
      status: status.value
    })
    alertname.value = severity.value = instance.value = message.value = ''
    status.value = 'firing'
    emit('sent')
  } catch (e) {
    console.error('Send alert failed', e)
  }
}
</script>

<style scoped>
form { margin: 20px 0; }
input, textarea, select { display: block; margin: 5px 0; width: 100%; padding: 5px; }
button { margin-top: 10px; padding: 5px 10px; }
</style>

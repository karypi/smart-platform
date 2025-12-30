<template>
  <div>
    <h1>Dashboard</h1>
    <p>Welcome!</p>
    <AlertForm @sent="loadHistory"/>
    <h2>Sent Alerts</h2>
    <ul>
      <li v-for="alert in alerts" :key="alert.ID">
        {{ alert.AlertName }} - {{ alert.Severity }} - {{ alert.Instance }} - {{ alert.Message }} - {{ alert.Status }} - {{ alert.CreatedAt }}
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import AlertForm from '@/components/AlertForm.vue'

const alerts = ref([])

const loadHistory = async () => {
  try {
    const res = await request.get('/api/alert/history')
    alerts.value = res.data
  } catch (e) {
    console.error('Failed to load alerts', e)
  }
}

onMounted(() => {
  loadHistory()
})
</script>

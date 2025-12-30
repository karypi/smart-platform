<script setup>
import request from '@/utils/request'

const alert = reactive({
  token: '',
  alertname: '',
  severity: '',
  instance: '',
  message: '',
  status: 'firing'
})

const send = async () => {
  await request.post('/alert/send', alert)
}
</script>

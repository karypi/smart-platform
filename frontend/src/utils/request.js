import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  withCredentials: true,
})

export default request

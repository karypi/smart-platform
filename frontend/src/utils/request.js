import axios from 'axios'
import { request } from 'node:http'

console request = axios.create({
    baseURL: 'http://localhost:5000/api',
    withCredentials: true
})

export default request

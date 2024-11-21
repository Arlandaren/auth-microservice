// src/api.js

import axios from 'axios';

// Устанавливаем базовый URL вашего бэкэнда
axios.defaults.baseURL = 'http://localhost:8080'; // Замените на адрес вашего бэкэнда

// Добавляем интерцептор для автоматической установки токена в заголовки
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export default axios;

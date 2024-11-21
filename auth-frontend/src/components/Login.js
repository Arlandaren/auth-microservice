// src/components/Login.js

import React, { useState } from 'react';
import { TextField, Button, Container, Typography, Box, CircularProgress } from '@mui/material';
import axios from '../api';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import {jwtDecode} from 'jwt-decode';

function Login() {
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [clientId, setClientId] = useState('');
  const [loading, setLoading] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();

  // Извлекаем redirect_uri из параметров URL
  const params = new URLSearchParams(location.search);
  const redirectUri = params.get('redirect_uri');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await axios.post('/v1/auth/login', {
        name,
        password,
        clientId: parseInt(clientId),
      });

      const token = response.data.token;

      if (token) {
        localStorage.setItem('token', token);

        // Получаем роль пользователя из токена
        let userRole = null;
        try {
          const decodedToken = jwtDecode(token);
          userRole = decodedToken.role || null;
        } catch (error) {
          console.error('Ошибка декодирования токена:', error);
        }

        const allowedRoles = ['Supreme', 'Client_Supreme'];

        if (!allowedRoles.includes(userRole)) {
          // Пользователь имеет другую роль, перенаправляем на клиентское приложение
          const redirectURL = new URL(redirectUri || 'http://localhost:3002');
          redirectURL.hash = `token=${token}`;
          window.location.href = redirectURL.toString();
        } else {
          // Пользователь имеет роль 'Supreme' или 'Client_Supreme'
          enqueueSnackbar('Вы успешно вошли в систему!', { variant: 'success' });
          navigate('/dashboard');
        }
      } else {
        enqueueSnackbar('Не удалось получить токен после входа', { variant: 'error' });
      }
    } catch (error) {
      console.error('Ошибка при входе:', error);
      enqueueSnackbar('Неверные учетные данные', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  };


  return (
    <Container maxWidth="xs">
      <Box sx={{ mt: 8 }}>
        <Typography component="h1" variant="h5" align="center">
          Вход
        </Typography>
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            label="Имя пользователя"
            value={name}
            onChange={(e) => setName(e.target.value)}
            autoFocus
          />
          <TextField
            margin="normal"
            required
            fullWidth
            label="Пароль"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            label="ID клиента"
            type="number"
            value={clientId}
            onChange={(e) => setClientId(e.target.value)}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            disabled={loading}
            sx={{ mt: 3, mb: 2 }}
          >
            {loading ? <CircularProgress size={24} /> : 'Войти'}
          </Button>
          <Typography align="center">
            Нет аккаунта? <Link to="/register">Зарегистрируйтесь</Link>
          </Typography>
        </Box>
      </Box>
    </Container>
  );
}

export default Login;

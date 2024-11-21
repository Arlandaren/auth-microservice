// src/components/Register.js

import React, { useState } from 'react';
import {
  TextField,
  Button,
  Container,
  Typography,
  Box,
  CircularProgress,
} from '@mui/material';
import axios from '../api';
import { useLocation, Link, useNavigate } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import {jwtDecode} from 'jwt-decode';

function Register() {
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [clientId, setClientId] = useState('');
  const [role, setRole] = useState('');
  const [loading, setLoading] = useState(false);
  const { enqueueSnackbar } = useSnackbar();
  const location = useLocation();
  const navigate = useNavigate();

  // Извлекаем redirect_uri из параметров URL
  const params = new URLSearchParams(location.search);
  const redirectUri = params.get('redirect_uri');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await axios.post('/v1/auth/register', {
        name,
        password,
        clientId: parseInt(clientId),
        role,
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
          enqueueSnackbar('Вы успешно зарегистрировались!', { variant: 'success' });
          navigate('/dashboard');
        }
      } else {
        enqueueSnackbar('Не удалось получить токен после регистрации', { variant: 'error' });
      }
    } catch (error) {
      console.error('Ошибка при регистрации:', error);
      enqueueSnackbar('Не удалось зарегистрироваться', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  };

  return (
    // Ваш JSX для формы регистрации
    <Container maxWidth="xs">
      <Box sx={{ mt: 8 }}>
        <Typography component="h1" variant="h5" align="center">
          Регистрация
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
          <TextField
            margin="normal"
            required
            fullWidth
            label="Роль"
            value={role}
            onChange={(e) => setRole(e.target.value)}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            disabled={loading}
            sx={{ mt: 3, mb: 2 }}
          >
            {loading ? <CircularProgress size={24} /> : 'Зарегистрироваться'}
          </Button>
          <Typography align="center">
            Уже есть аккаунт? <Link to="/login">Войдите</Link>
          </Typography>
        </Box>
      </Box>
    </Container>
  );
}

export default Register;

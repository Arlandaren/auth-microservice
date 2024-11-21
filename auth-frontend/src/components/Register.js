// src/components/Register.js

import React, { useState } from 'react';
import { TextField, Button, Container, Typography, Box, CircularProgress } from '@mui/material';
import axios from '../api';
import { useNavigate, Link } from 'react-router-dom';
import { useSnackbar } from 'notistack';

function Register() {
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [clientId, setClientId] = useState('');
  const [role, setRole] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      await axios.post('/v1/auth/register', {
        name,
        password,
        clientId: parseInt(clientId),
        role,
      });

      enqueueSnackbar('Регистрация успешна! Теперь вы можете войти.', { variant: 'success' });
      navigate('/login');
    } catch (error) {
      console.error('Ошибка при регистрации:', error);
      enqueueSnackbar('Не удалось зарегистрироваться', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  };

  return (
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

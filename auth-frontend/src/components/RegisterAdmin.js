// src/components/RegisterAdmin.js

import React, { useState } from 'react';
import {
  TextField,
  Button,
  Container,
  Typography,
  Box,
  CircularProgress,
} from '@mui/material';
import axios from '../api'; // Исправлено имя файла на 'api' с маленькой буквы
import { Link, useNavigate } from 'react-router-dom';
import { useSnackbar } from 'notistack';

function RegisterAdmin() {
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [clientId, setClientId] = useState('');
  const [role, setRole] = useState('');
  const [loading, setLoading] = useState(false);
  const { enqueueSnackbar } = useSnackbar();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      await axios.post('/v1/auth/register/admin', {
        name,
        password,
        clientId: parseInt(clientId),
        role,
      });

      enqueueSnackbar('Администратор успешно зарегистрирован!', { variant: 'success' });
      setName('');
      setPassword('');
      setClientId('');
      setRole('');
      navigate('/dashboard');
    } catch (error) {
      console.error('Ошибка при регистрации администратора:', error);
      enqueueSnackbar('Не удалось зарегистрировать администратора', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 8 }}>
        <Typography component="h1" variant="h5" align="center">
          Регистрация администратора
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
            variant="contained"
            disabled={loading}
            sx={{ mt: 3 }}
          >
            {loading ? <CircularProgress size={24} /> : 'Зарегистрировать администратора'}
          </Button>
        </Box>
        <Box sx={{ mt: 2 }}>
          <Button component={Link} to="/dashboard" variant="outlined">
            Вернуться на дашборд
          </Button>
        </Box>
      </Box>
    </Container>
  );
}

export default RegisterAdmin;

// src/components/RegisterClient.js

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
import { Link, useNavigate } from 'react-router-dom';
import { useSnackbar } from 'notistack';

function RegisterClient() {
  const [name, setName] = useState('');
  const [roles, setRoles] = useState('');
  const [loading, setLoading] = useState(false);
  const { enqueueSnackbar } = useSnackbar();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      await axios.post('/v1/auth/register/client', {
        name,
        roles: roles.split(',').map((role) => role.trim()),
      });

      enqueueSnackbar('Клиент успешно зарегистрирован!', { variant: 'success' });
      setName('');
      setRoles('');
      // Перенаправить на дашборд после успешной регистрации
      navigate('/dashboard');
    } catch (error) {
      console.error('Ошибка при регистрации клиента:', error);
      enqueueSnackbar('Не удалось зарегистрировать клиента', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 8 }}>
        <Typography component="h1" variant="h5" align="center">
          Регистрация клиента
        </Typography>
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            label="Название клиента"
            value={name}
            onChange={(e) => setName(e.target.value)}
            autoFocus
          />
          <TextField
            margin="normal"
            required
            fullWidth
            label="Роли (через запятую)"
            value={roles}
            onChange={(e) => setRoles(e.target.value)}
          />
          <Button
            type="submit"
            variant="contained"
            disabled={loading}
            sx={{ mt: 3 }}
          >
            {loading ? <CircularProgress size={24} /> : 'Зарегистрировать клиента'}
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

export default RegisterClient;

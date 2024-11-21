// src/components/Dashboard.js

import React from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Button, Container, Typography, Box } from '@mui/material';

function Dashboard() {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <Container maxWidth="md">
      <Box sx={{ mt: 8, textAlign: 'center' }}>
        <Typography component="h1" variant="h4">
          Добро пожаловать в панель управления!
        </Typography>
        <Box sx={{ mt: 4 }}>
          <Button component={Link} to="/register-client" variant="contained" sx={{ mr: 2 }}>
            Регистрация клиента
          </Button>
          <Button component={Link} to="/register-admin" variant="contained" sx={{ mr: 2 }}>
            Регистрация администратора
          </Button>
          <Button variant="outlined" color="secondary" onClick={handleLogout}>
            Выйти
          </Button>
        </Box>
      </Box>
    </Container>
  );
}

export default Dashboard;

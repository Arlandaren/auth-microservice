// src/hooks/useRoleCheck.js

import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';

function useRoleCheck(allowedRoles) {
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('token');
    let userRole = null;

    if (token) {
      try {
        const decodedToken = jwtDecode(token);
        userRole = decodedToken.role || null;

        const currentTime = Math.floor(Date.now() / 1000);
        if (decodedToken.exp && decodedToken.exp < currentTime) {
          console.error('Срок действия токена истёк');
          localStorage.removeItem('token');
          navigate('/login');
        }
      } catch (error) {
        console.error('Ошибка декодирования токена:', error);
        localStorage.removeItem('token');
        navigate('/login');
      }
    } else {
      navigate('/login');
    }

    if (!allowedRoles.includes(userRole)) {
      window.location.href = 'http://localhost:3001'; // Замените на URL вашего клиентского приложения
    }
  }, [allowedRoles, navigate]);
}

export default useRoleCheck;

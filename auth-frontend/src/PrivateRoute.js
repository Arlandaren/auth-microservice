// src/PrivateRoute.js

import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import {jwtDecode} from 'jwt-decode'; // Обратите внимание на правильный импорт

function PrivateRoute({ roles }) {
  const token = localStorage.getItem('token');
  let isAuthenticated = !!token;
  let userRole = null;

  if (token) {
    try {
      const decodedToken = jwtDecode(token);
      userRole = decodedToken.role || null;

      // Проверяем, не истёк ли срок действия токена
      const currentTime = Math.floor(Date.now() / 1000);
      if (decodedToken.exp && decodedToken.exp < currentTime) {
        console.error('Срок действия токена истёк');
        isAuthenticated = false;
        localStorage.removeItem('token');
      }
    } catch (error) {
      console.error('Ошибка декодирования токена:', error);
      isAuthenticated = false;
      localStorage.removeItem('token');
    }
  } else {
    console.log('Токен отсутствует');
  }

  if (!isAuthenticated) {
    console.log('Пользователь не аутентифицирован');
    return <Navigate to="/login" replace />;
  }

  // Если пользователь не имеет роль 'Supreme' или 'Client_Supreme', перенаправляем его на клиентское приложение
  const allowedRoles = ['Supreme', 'Client_Supreme'];
  if (!allowedRoles.includes(userRole)) {
    console.log('Пользователь не имеет достаточных прав, перенаправляем на клиентское приложение');
    window.location.href = 'http://localhost:3002'; // Замените на URL вашего клиентского приложения
    return null; // Возвращаем null, так как перенаправление уже произошло
  }

  // Проверяем, если требуются определенные роли для доступа к этому маршруту
  if (roles && !roles.includes(userRole)) {
    console.log('У пользователя нет необходимой роли для этого маршрута');
    console.log('Роль пользователя:', userRole);
    console.log('Требуемые роли:', roles);
    return <Navigate to="/unauthorized" replace />;
  }

  return <Outlet />;
}

export default PrivateRoute;

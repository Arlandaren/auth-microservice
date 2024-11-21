  // src/PrivateRoute.js

  import React from 'react';
  import { Navigate, Outlet } from 'react-router-dom';
  import { jwtDecode } from 'jwt-decode';

  function PrivateRoute({ roles }) {
    const token = localStorage.getItem('token');
    let isAuthenticated = !!token;
    let userRole = null;

    if (token) {
      try {
        const decodedToken = jwtDecode(token);
        userRole = decodedToken.role || null;

        // Проверяем, не истёк ли срок действия токена (опционально)
        const currentTime = Math.floor(Date.now() / 1000);
        if (decodedToken.exp && decodedToken.exp < currentTime) {
          console.error('Срок действия токена истёк');
          isAuthenticated = false;
        }
      } catch (error) {
        console.error('Ошибка декодирования токена:', error);
        isAuthenticated = false;
      }
    } else {
      console.log('Токен отсутствует');
    }

    if (!isAuthenticated) {
      console.log('Пользователь не аутентифицирован');
      return <Navigate to="/login" replace />;
    }

    // Проверяем, есть ли у пользователя необходимая роль
    if (roles && !roles.includes(userRole)) {
      console.log('У пользователя нет необходимой роли');
      console.log('Роль пользователя:', userRole);
      console.log('Требуемые роли:', roles);
      return <Navigate to="/unauthorized" replace />;
    }

    return <Outlet />;
  }

  export default PrivateRoute;

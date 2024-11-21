// src/App.js

import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './components/Login';
import Register from './components/Register';
import RegisterClient from './components/RegisterClient';
import RegisterAdmin from './components/RegisterAdmin';
import Dashboard from './components/Dashboard';
import PrivateRoute from './PrivateRoute';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          {/*Публичные маршруты*/}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          {/*Приватные маршруты*/}
          <Route element={<PrivateRoute />}>
            <Route path="/dashboard" element={<Dashboard />} />
          </Route>

          <Route element={<PrivateRoute roles={['Supreme', 'Client_Supreme', 'Admin']} />}>
            <Route path="/register-client" element={<RegisterClient />} />
          </Route>

          <Route element={<PrivateRoute roles={['Supreme', 'Client_Supreme']} />}>
            <Route path="/register-admin" element={<RegisterAdmin />} />
          </Route>

          {/*Перенаправление по умолчанию*/}
          <Route path="*" element={<Navigate to="/login" replace />} />*
        </Routes>
      </div>
    </Router>
  );
}

export default App;

import { useEffect } from 'react';
import { Route, Routes } from 'react-router-dom';
import './App.css';
import Dashboard from './components/Dashboard';
import Login from './components/Login';
import { getToken, useToken } from './api/auth';

function App() {
  const { token, setToken, setRefreshToken } = useToken();

  useEffect(() => {
    const intervalId = setInterval(async () => {
      if (!getToken()) {
        setToken('');
      }
    }, 1000);

    return () => clearInterval(intervalId);
  });

  if (!token) {
    return <Login setToken={setToken} setRefreshToken={setRefreshToken} />;
  }

  return (
    <div className="app">
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route
          path="/login"
          element={<Login setToken={setToken} setRefreshToken={setRefreshToken} />}
        />
      </Routes>
    </div>
  );
}

export default App;

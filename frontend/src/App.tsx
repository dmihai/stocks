import { Route, Routes } from 'react-router-dom';
import './App.css';
import Dashboard from './components/Dashboard';
import Login from './components/Login';
import { useToken } from './api/auth';
import { AuthContext } from './components/AuthContext';

function App() {
  const { token, refreshToken, setToken, setRefreshToken } = useToken();

  let app = (
    <div className="App">
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </div>
  );

  if (!token) {
    app = <Login />;
  }

  return (
    <AuthContext.Provider value={{ token, refreshToken, setToken, setRefreshToken }}>
      {app}
    </AuthContext.Provider>
  );
}

export default App;

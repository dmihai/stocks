import { Route, Routes } from 'react-router-dom';
import './App.css';
import Dashboard from './components/Dashboard';
import Login from './components/Login';
import { useToken } from './api/auth';

function App() {
  const { token, setToken, setRefreshToken } = useToken();

  if (!token) {
    return <Login setToken={setToken} setRefreshToken={setRefreshToken} />;
  }

  return (
    <div className="App">
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

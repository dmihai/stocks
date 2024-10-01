import { Route, Routes } from 'react-router-dom';
import './App.css';
import Dashboard from './components/Dashboard';
import Login from './components/Login';
import useToken from './api/auth';
import { AuthContext } from './components/AuthContext';

function App() {
  const { token, setToken } = useToken();

  if (!token) {
    return <Login setToken={setToken} />;
  }

  return (
    <div className="App">
      <AuthContext.Provider value={{ token }}>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/login" element={<Login setToken={setToken} />} />
        </Routes>
      </AuthContext.Provider>
    </div>
  );
}

export default App;

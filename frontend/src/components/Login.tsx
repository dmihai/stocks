import { useState, useEffect } from 'react';
import './Login.css';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = () => {
    console.log(`The username you entered was: ${username}`);
  };

  return (
    <div className="container">
      <div id="login-row" className="row justify-content-center align-items-center">
        <div id="login-column" className="col-md-6">
          <div id="login-box" className="col-md-12">
            <form id="login-form" onSubmit={handleSubmit} autoComplete="off" action="#">
              <h3 className="text-center text-info">Login</h3>
              <div className="mb-3">
                <label htmlFor="username" className="form-label">
                  Username:
                </label>
                <input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="form-control"
                />
              </div>
              <div className="mb-3">
                <label htmlFor="password" className="form-label">
                  Password:
                </label>
                <input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="form-control"
                />
              </div>
              <div className="mb-3 form-check">
                <input id="remember-me" type="checkbox" className="form-check-input" />
                <label htmlFor="remember-me" className="form-check-label">
                  Remember me
                </label>
              </div>
              <div className="mb-3">
                <input type="submit" className="btn btn-primary btn-md" value="Submit" />
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Login;

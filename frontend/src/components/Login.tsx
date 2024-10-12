import { useState } from 'react';
import { login } from '../api/api';

type LoginProps = {
  setToken: (token: string) => void;
  setRefreshToken: (refreshToken: string) => void;
};

function Login(props: LoginProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    const result = await login(username, password);
    props.setToken(result.token);
    props.setRefreshToken(result.refreshToken);
  };

  return (
    <div className="container">
      <div className="row justify-content-center mt-5">
        <div className="col-md-6">
          <div className="card border-primary bg-light">
            <h5 className="card-header bg-primary text-light">Login</h5>
            <div className="card-body">
              <form onSubmit={handleSubmit} autoComplete="off">
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
                  <input
                    type="submit"
                    className="btn btn-primary btn-md"
                    value="Submit"
                  />
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Login;

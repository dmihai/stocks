import { useState } from 'react';

const tokenKey = "token";

export default function useToken() {
  const getToken = () => {
    return sessionStorage.getItem(tokenKey);
  };

  const [token, setToken] = useState(getToken());

  const saveToken = (token: string) => {
    sessionStorage.setItem(tokenKey, token);
    setToken(token);
  };

  return {
    setToken: saveToken,
    token
  }
}

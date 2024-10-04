import { useState } from 'react';

const tokenKey = 'token';
const refreshTokenKey = 'refreshToken';

export function useToken() {
  const [token, setToken] = useState(getToken());
  const [refreshToken, setRefreshToken] = useState(getRefreshToken());

  return {
    token,
    refreshToken,
    setToken: (token: string) => {
      saveToken(token);
      setToken(token);
    },
    setRefreshToken: (refreshToken: string) => {
      saveRefreshToken(refreshToken);
      setRefreshToken(refreshToken);
    },
  };
}

export function getToken() {
  return localStorage.getItem(tokenKey);
}

export function saveToken(token: string) {
  localStorage.setItem(tokenKey, token);
}

export function getRefreshToken() {
  return localStorage.getItem(refreshTokenKey);
}

export function saveRefreshToken(refreshToken: string) {
  localStorage.setItem(refreshTokenKey, refreshToken);
}

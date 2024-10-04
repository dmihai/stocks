import { createContext } from 'react';

export type AuthPayload = {
  token: string | null;
  refreshToken: string | null;
  setToken: (token: string) => void;
  setRefreshToken: (refreshToken: string) => void;
};

export const AuthContext = createContext<AuthPayload>({
  token: '',
  refreshToken: '',
  setToken: () => {},
  setRefreshToken: () => {},
});

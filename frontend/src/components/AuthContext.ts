import { createContext } from 'react';

type AuthContextPayload = {
    token: string;
}

export const AuthContext = createContext<AuthContextPayload>({
    token: '',
});

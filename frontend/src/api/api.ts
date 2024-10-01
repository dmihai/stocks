import axios from 'axios';

export type TopGainer = {
    symbol: string;
    percentChanged: number;
    yesterday: {
        open: number;
        high: number;
        low: number;
        close: number;
        volume: number;
    }
    current: {
        price: number;
        volume: number;
    }
    lastUpdated: string;
};

type Login = {
    token: string;
}

const api = axios.create({
  baseURL: process.env.REACT_APP_API_ADDRESS,
});

function auth(token: string) {
    return {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    };
};

export async function getTopGainers(token: string) {
    try {
        const response = await api.get<TopGainer[]>('top-gainers', auth(token));
        return response.data;
    } catch (error) {
        console.error(error);
        return [];
    }
}

export async function login(username: string, password: string) {
    try {
        const response = await api.post<Login>('login', {}, {
            auth: {
              username: username,
              password: password,
            }
          });
        return response.data.token;
    } catch (error) {
        if (!axios.isAxiosError(error) || error.status != 401) {
          console.error(error);
        }

        return "";
    }
}

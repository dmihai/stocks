import axios from 'axios';
import { getRefreshToken, getToken, saveRefreshToken, saveToken } from './auth';

export type TopGainer = {
  symbol: string;
  percentChanged: number;
  yesterday: {
    open: number;
    high: number;
    low: number;
    close: number;
    volume: number;
  };
  current: {
    price: number;
    volume: number;
  };
  lastUpdated: string;
};

type Login = {
  token: string;
  refreshToken: string;
};

const api = axios.create({
  baseURL: process.env.REACT_APP_API_ADDRESS,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config) => {
    if (config.url && !config.url.includes('api/')) {
      return config;
    }

    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    if (
      error.response.status === 401 &&
      !originalRequest._retry &&
      originalRequest.url.includes('api/')
    ) {
      originalRequest._retry = true;
      try {
        const response = await exchange();
        const { token, refreshToken: newRefreshToken } = response;

        saveToken(token);
        saveRefreshToken(newRefreshToken);

        api.defaults.headers.Authorization = `Bearer ${token}`;
        return api(originalRequest);
      } catch (refreshError) {
        console.error('Token refresh failed:', refreshError);
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(error);
  },
);

export async function getTopGainers() {
  try {
    const response = await api.get<TopGainer[]>('api/top-gainers');
    return response.data;
  } catch (error) {
    if (!axios.isAxiosError(error) || error.status !== 401) {
      console.error(error);
    }
    return [];
  }
}

export async function login(username: string, password: string) {
  try {
    const basicAuth = {
      auth: {
        username: username,
        password: password,
      },
    };

    const response = await api.post<Login>('login', {}, basicAuth);
    return response.data;
  } catch (error) {
    if (!axios.isAxiosError(error) || error.status !== 401) {
      console.error(error);
    }

    return {
      token: '',
      refreshToken: '',
    };
  }
}

export async function exchange() {
  try {
    const refreshToken = getRefreshToken();
    const bearerAuth = {
      headers: {
        Authorization: `Bearer ${refreshToken}`,
      },
    };

    const response = await api.post<Login>('exchange', {}, bearerAuth);
    return response.data;
  } catch (error) {
    if (!axios.isAxiosError(error) || error.status !== 401) {
      console.error(error);
    }

    return {
      token: '',
      refreshToken: '',
    };
  }
}

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

export async function getTopGainers() {
    try {
        const url = process.env.REACT_APP_API_ADDRESS + 'top-gainers';
        const response = await axios.get<TopGainer[]>(url);
        return response.data;
    } catch (error) {
        console.error(error);
        return [];
    }
}

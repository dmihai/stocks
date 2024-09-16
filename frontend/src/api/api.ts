import axios from 'axios';

export type TopGainer = {
    symbol: string;
    percentChange: number;
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
};

const server = 'http://localhost:3100/';

export async function getTopGainers() {
    try {
        const response = await axios.get<TopGainer[]>(server + 'top-gainers');
        return response.data;
    } catch (error) {
        console.error(error);
        return [];
    }
}

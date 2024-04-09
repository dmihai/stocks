import logging

import pandas as pd
import yfinance as yf

from utils.Database import Database

provider = 'yahoo'

class Importer:
    def __init__(self, db: Database, startdate: str):
        self._db = db
        self._startdate = startdate
    
    def populate(self, symbol: str):
        self._db.update_symbol_history(symbol, provider, None, None)
    
        stock = yf.Ticker(symbol)

        try:
            df = stock.history(period='max', interval='1d', start=self._startdate)
            if df.size == 0:
                raise Exception('empty history')

            df.insert(loc=0, column='date', value=df.index.strftime('%Y-%m-%d'))
            df.insert(loc=0, column='symbol', value=symbol)
            history = df[['symbol', 'date', 'Open', 'High', 'Low', 'Close', 'Volume', 'Dividends', 'Stock Splits']]
            history_list = history.values.tolist()
                
            self._db.insert_history(history_list)
            self._db.update_symbol_history(symbol, provider, history_list[0][1], history_list[-1][1])

            logging.info(f"Imported {len(history_list)} entries for symbol {symbol} between {history_list[0][1]} and {history_list[-1][1]}")
        except Exception as e:
            logging.warning(f"Failed to import price history for symbol {symbol}: {e}")
    
    def update(self, items: list[tuple[str,str]]):
        symbols = [item[0] for item in items]

        startdates = [item[1] for item in items]
        startdates.sort()
        startdate = startdates[0]

        columns = ['Ticker', 'Formatted Date', 'Open', 'High', 'Low', 'Close', 'Volume', 'Dividends', 'Stock Splits']

        try:
            df = yf.download(symbols, period='max', interval='1d', start=startdate, actions=True, progress=False, group_by='ticker')

            if len(df) > 0:
                df = pd.concat([gp.xs(idx, level=0, axis=1).assign(Ticker=idx) for idx, gp in df.groupby(level=0, axis=1)])
                df=df.reset_index()
                df=df.dropna(subset=['Open', 'High', 'Low', 'Close'], how='all')

                df['Formatted Date'] = df['Date'].dt.strftime('%Y-%m-%d')

                df = df[columns]
                history = df.values.tolist()
                
                self._db.insert_history(history)
            else:
                df = pd.DataFrame(columns=columns)

            for symbol in symbols:
                df_symbol = df[df['Ticker'] == symbol]
                if len(df_symbol) > 0:
                    self._db.update_symbol_history(symbol, provider, None, df_symbol['Formatted Date'].max())
                else:
                    self._db.update_symbol_history(symbol, provider, None, None)

            logging.info(f"Updated {len(df)} entries for symbols {symbols}")
        except Exception as e:
            logging.warning(f"Failed to update price history for symbols {symbols}: {e}")

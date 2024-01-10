import logging

import yfinance as yf

from utils.Database import Database


class Importer:
    def __init__(self, db: Database, startdate: str):
        self._db = db
        self._startdate = startdate
    
    def populate(self, symbol: str):
        self._db.update_symbol_history(symbol, 'yahoo', None, None)
    
        stock = yf.Ticker(symbol)

        try:
            df = stock.history(period='max', interval='1d', start=self._startdate)
            if df.size == 0:
                raise Exception('empty history')

            df.insert(loc=0, column='date', value=df.index.strftime('%Y-%m-%d'))
            df.insert(loc=0, column='symbol', value=symbol)
            history = df[['symbol', 'date', 'Open', 'High', 'Low', 'Close', 'Volume', 'Dividends', 'Stock Splits']].values.tolist()
                
            self._db.insert_history(history)
            self._db.update_symbol_history(symbol, 'yahoo', history[0][1], history[-1][1])

            logging.info(f"Imported {len(history)} entries for symbol {symbol} between {history[0][1]} and {history[-1][1]}")
        except Exception as e:
            logging.warning(f"Failed to import price history for symbol {symbol}: {e}")

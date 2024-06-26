from abc import ABC, abstractmethod
from datetime import date


class Database(ABC):
    @abstractmethod
    def __del__(self):
        pass

    @abstractmethod
    def connect(self, config):
        pass

    def get_next_unpopulated_symbol(self):
        cursor = self._conn.cursor()

        query = self._transform_query("SELECT symbol, history_start, history_end\
            FROM stocks\
            WHERE history IS NULL AND symbol NOT LIKE '%-%' AND last_update IS NULL\
            LIMIT 1")
        cursor.execute(query)

        result = cursor.fetchone()

        if result is None:
            return None

        return {
            "symbol": result[0],
            "start": result[1],
            "end": result[2]
        }

    def get_next_symbols_for_update(self, limit: int):
        cursor = self._conn.cursor()

        query = self._transform_query("SELECT symbol, history_end\
            FROM stocks\
            WHERE status='Active' AND history_start IS NOT NULL AND last_update<%s\
            ORDER BY history_end\
            LIMIT %s")
        cursor.execute(query, (self.today(), limit))

        return cursor.fetchall()

    def insert_history(self, history):
        cursor = self._conn.cursor()

        query = self._transform_query("REPLACE INTO history\
            (symbol, date, open, high, low, close, volume, dividends, splits)\
            VALUES(%s, %s, %s, %s, %s, %s, %s, %s, %s)")
        cursor.executemany(query, history)
        self._conn.commit()

    def insert_stocks(self, stocks):
        cursor = self._conn.cursor()

        query = self._transform_query("INSERT IGNORE INTO stocks\
            (symbol, name, exchange, asset_type, ipo_date, delisting_date, status)\
            VALUES (%s, %s, %s, %s, %s, %s, %s)")
        cursor.executemany(query, stocks)
        self._conn.commit()

        return cursor.rowcount

    def update_symbol_history(self, symbol, history, history_start, history_end):
        cursor = self._conn.cursor()

        update_fields = "history = %s, last_update = %s"
        params = (history, self.today())
        
        if history_start is not None:
            update_fields += ", history_start = %s"
            params += (history_start,)
        
        if history_end is not None:
            update_fields += ", history_end = %s"
            params += (history_end,)
        
        query = self._transform_query(f"UPDATE stocks\
            SET {update_fields}\
            WHERE symbol = %s")
        cursor.execute(query, params + (symbol,))
        
        self._conn.commit()

    def today(self):
        return str(date.today())

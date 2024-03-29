import sqlite3

from datetime import date

from utils.Database import Database


class DatabaseSQLite(Database):
    def __del__(self):
        self._conn.close()

    def connect(self, db_file):
        self._conn = sqlite3.connect(db_file)

    def get_next_unpopulated_symbol(self):
        cursor = self._conn.cursor()

        query = "SELECT symbol, history_start, history_end\
            FROM stocks\
            WHERE history is null AND symbol NOT LIKE '%-%' AND last_update<?\
            LIMIT 1"
        cursor.execute(query, (self.today(),))

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

        query = "SELECT symbol, history_end\
            FROM stocks\
            WHERE status='Active' AND history_start IS NOT NULL AND last_update<?\
            ORDER BY history_end\
            LIMIT ?"
        cursor.execute(query, (self.today(), limit))

        return cursor.fetchall()

    def insert_history(self, history):
        cursor = self._conn.cursor()

        cursor.executemany(
            "REPLACE INTO history (symbol, date, open, high, low, close, volume, dividends, splits) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)", history)
        self._conn.commit()

    def insert_stocks(self, stocks):
        cursor = self._conn.cursor()

        cursor.executemany("""INSERT OR IGNORE INTO stocks
            (symbol, name, exchange, asset_type, ipo_date, delisting_date, status)
            VALUES (?, ?, ?, ?, ?, ?, ?)""", stocks)
        self._conn.commit()

        return cursor.rowcount

    def update_symbol_history(self, symbol, history, history_start, history_end):
        cursor = self._conn.cursor()

        update_fields = "history = ?, last_update = ?"
        params = (history, self.today())
        
        if history_start is not None:
            update_fields += ", history_start = ?"
            params += (history_start,)
        
        if history_end is not None:
            update_fields += ", history_end = ?"
            params += (history_end,)
        
        query = f"UPDATE stocks\
            SET {update_fields}\
            WHERE symbol = ?"
        cursor.execute(query, params + (symbol,))
        
        self._conn.commit()

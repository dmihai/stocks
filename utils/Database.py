import sqlite3

from datetime import date


class Database:
    def __del__(self):
        self._conn.close()

    def connect(self, db_file):
        self._conn = sqlite3.connect(db_file)

    def get_next_unpopulated_symbol(self):
        cursor = self._conn.cursor()

        query = "SELECT symbol, history_start, history_end\
            FROM stocks\
            WHERE history is null AND symbol NOT LIKE '%-%'\
            LIMIT 1"
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

        query = "SELECT symbol, history_end\
            FROM stocks\
            WHERE status='Active' AND last_update<?\
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

        query = "UPDATE stocks\
            SET history = ?, history_start = ?, history_end = ?, last_update = ?\
            WHERE symbol = ?"
        cursor.execute(query, (history, history_start,
                       history_end, self.today(), symbol))
        self._conn.commit()

    def today(self):
        return str(date.today())

import sqlite3

class Database:
    def __del__(self):
        self._conn.close()

    def connect(self, db_file):
        self._conn = sqlite3.connect(db_file)
    
    def get_next_symbol(self):
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

    def insert_history(self, history):
        cursor = self._conn.cursor()

        cursor.executemany("INSERT INTO history (symbol, date, open, high, low, close, volume, dividends, splits) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)", history)
        self._conn.commit()
    
    def update_symbol_history(self, symbol, history, history_start, history_end):
        cursor = self._conn.cursor()

        query = "UPDATE stocks\
            SET history = ?, history_start = ?, history_end = ?\
            WHERE symbol = ?"
        cursor.execute(query, (history, history_start, history_end, symbol))
        self._conn.commit()

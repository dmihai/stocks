import mysql.connector

from utils.Database import Database


class DatabaseMysql(Database):
    def __del__(self):
        self._conn.close()

    def connect(self, config):
        self._conn = mysql.connector.connect(
            host=config['host'],
            user=config['user'],
            password=config['password'],
            database=config['database']
        )
    
    def _transform_query(self, query):
        return query

    def populate_history_month(self):
        cursor = self._conn.cursor()

        query = "INSERT INTO history_month\
                (symbol, date, open, high, low, close, volume)\
            SELECT\
                symbol,\
                DATE_SUB(date, INTERVAL (DAY(date)-1) DAY) month,\
                (SELECT open\
                    FROM history\
                    WHERE symbol=h.symbol AND YEAR(date)=YEAR(month) AND MONTH(date)=MONTH(month)\
                    ORDER BY date ASC\
                    LIMIT 1),\
                MAX(high),\
                MIN(low),\
                (SELECT close\
                    FROM history\
                    WHERE symbol=h.symbol AND YEAR(date)=YEAR(month) AND MONTH(date)=MONTH(month)\
                    ORDER BY date DESC\
                    LIMIT 1),\
                SUM(volume)\
            FROM history h\
            GROUP BY symbol, DATE_SUB(date, INTERVAL (day(date)-1) DAY)"
        cursor.execute(query)

        self._conn.commit()

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
    
    def _insert_stocks_query(self):
        return "INSERT IGNORE INTO stocks\
            (symbol, name, exchange, asset_type, ipo_date, delisting_date, status)\
            VALUES (?, ?, ?, ?, ?, ?, ?)"

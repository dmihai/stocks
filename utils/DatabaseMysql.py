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

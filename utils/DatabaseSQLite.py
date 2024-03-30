import sqlite3

from utils.Database import Database


class DatabaseSQLite(Database):
    def __del__(self):
        self._conn.close()

    def connect(self, config):
        self._conn = sqlite3.connect(config["db_file"])

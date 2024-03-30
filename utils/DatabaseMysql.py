import mysql.connector

from utils.Database import Database


class DatabaseMysql(Database):
    def __del__(self):
        pass

    def connect(self, config):
        pass

    def get_next_unpopulated_symbol(self):
        pass

    def get_next_symbols_for_update(self, limit: int):
        pass

    def insert_history(self, history):
        pass

    def insert_stocks(self, stocks):
        pass

    def update_symbol_history(self, symbol, history, history_start, history_end):
        pass

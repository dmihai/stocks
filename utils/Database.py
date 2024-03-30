from abc import ABC, abstractmethod
from datetime import date


class Database(ABC):
    @abstractmethod
    def __del__(self):
        pass

    @abstractmethod
    def connect(self, config):
        pass

    @abstractmethod
    def get_next_unpopulated_symbol(self):
        pass

    @abstractmethod
    def get_next_symbols_for_update(self, limit: int):
        pass

    @abstractmethod
    def insert_history(self, history):
        pass

    @abstractmethod
    def insert_stocks(self, stocks):
        pass

    @abstractmethod
    def update_symbol_history(self, symbol, history, history_start, history_end):
        pass

    def today(self):
        return str(date.today())

import argparse
from datetime import date
import json
import logging
import signal
import time

import yfinance as yf

from utils.Database import Database

default_db_file = 'data/stocks.sqlite3'
default_history_start = '1995-01-01'
wait_time = 2


def get_config(file):
    f = open(file)
    config = json.load(f)
    f.close()

    return config

def config_logging(config):
    filename = config['filename'].replace("%(today)s", str(date.today())) if config['output'] == 'file' else None
    level = getattr(logging, config['level'].upper())
    
    logging.basicConfig(filename=filename, format="%(asctime)s - %(levelname)s - %(message)s", level=level)

def stop_script(signum=None, frame=None):
    global is_running
    is_running = False


config = get_config('config.json')
config_logging(config['logging'])

parser = argparse.ArgumentParser()
parser.add_argument('--dbfile', type=str, required=False, default=default_db_file,
                    help='Choose the sqlite3 file')
parser.add_argument('--startdate', type=str, required=False, default=default_history_start,
                    help='Start date for stock history (2015-12-30)')
args = parser.parse_args()

db = Database()
db.connect(args.dbfile)

is_running = True

signal.signal(signal.SIGINT, stop_script)
signal.signal(signal.SIGTERM, stop_script)

while is_running:
    res = db.get_next_symbol()

    if res is None:
        logging.info('No more symbols to import.')
        stop_script()
    
    db.update_symbol_history(res['symbol'], 'yahoo', None, None)
    
    stock = yf.Ticker(res['symbol'])

    try:
        df = stock.history(period='max', interval='1d', start=args.startdate)
        if df.size == 0:
            raise Exception('empty history')

        df.insert(loc=0, column='date', value=df.index.strftime('%Y-%m-%d'))
        df.insert(loc=0, column='symbol', value=res['symbol'])
        history = df[['symbol', 'date', 'Open', 'High', 'Low', 'Close', 'Volume', 'Dividends', 'Stock Splits']].values.tolist()
            
        db.insert_history(history)
        db.update_symbol_history(res['symbol'], 'yahoo', history[0][1], history[-1][1])

        logging.info(f"Imported {len(history)} entries for symbol {res['symbol']} between {history[0][1]} and {history[-1][1]}")
    except Exception as e:
        logging.warning(f"Failed to import price history for symbol {res['symbol']}: {e}")

    time.sleep(wait_time)

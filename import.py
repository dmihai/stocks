import argparse
from datetime import date
import json
import logging
import pandas as pd
import signal
import time

from utils.Database import Database
from utils.Importer import Importer

default_db_file = 'data/stocks.sqlite3'
default_history_start = '1995-01-01'
wait_time = 2
stocks_list_url = "https://www.alphavantage.co/query?function=LISTING_STATUS&apikey=demo"


def get_config(file):
    f = open(file)
    config = json.load(f)
    f.close()

    return config


def config_logging(config):
    filename = config['filename'].replace("%(today)s", str(
        date.today())) if config['output'] == 'file' else None
    level = getattr(logging, config['level'].upper())

    logging.basicConfig(
        filename=filename, format="%(asctime)s - %(levelname)s - %(message)s", level=level)


def stop_script(signum=None, frame=None):
    global is_running
    is_running = False


config = get_config('config.json')
config_logging(config['logging'])

parser = argparse.ArgumentParser()
parser.add_argument('--operation', type=str, required=False, default="populate",
                    help='Operation to perform (populate, update, new_stocks)')
parser.add_argument('--dbfile', type=str, required=False, default=default_db_file,
                    help='Choose the sqlite3 file')
parser.add_argument('--startdate', type=str, required=False, default=default_history_start,
                    help='Start date for stock history (2015-12-30)')
args = parser.parse_args()

if args.operation not in {'populate', 'update', 'new_stocks'}:
    logging.error('Invalid operation')
    exit()

db = Database()
db.connect(args.dbfile)

importer = Importer(db, args.startdate)

is_running = True

signal.signal(signal.SIGINT, stop_script)
signal.signal(signal.SIGTERM, stop_script)

while is_running:
    if args.operation == 'populate':
        res = db.get_next_unpopulated_symbol()

        if res is None:
            logging.info('No more symbols to import.')
            break

        importer.populate(res['symbol'])

        time.sleep(wait_time)

    elif args.operation == 'update':
        items = db.get_next_symbols_for_update(2)

        if len(items) == 0:
            logging.info('No more symbols to update.')
            break

        importer.update(items)

        exit()

    elif args.operation == 'new_stocks':
        data = pd.read_csv(stocks_list_url)
        data = data.query('symbol.notnull()')
        count = db.insert_stocks(data.values.tolist())

        logging.info(f"Inserted {count} new stocks.")

        exit()

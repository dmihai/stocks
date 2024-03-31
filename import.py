import argparse
from datetime import date
import json
import logging
import pandas as pd
import numpy as np
import signal
import time

from utils.DatabaseSQLite import DatabaseSQLite
from utils.DatabaseMysql import DatabaseMysql
from utils.Importer import Importer

default_history_start = '1995-01-01'
wait_time = 1
update_group_count = 3


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


def init_database(storage):
    db = None
    if storage['provider'] == 'sqlite':
        db = DatabaseSQLite()
    elif storage['provider'] == 'mysql':
        db = DatabaseMysql()
    
    if db is not None:
        db.connect(storage['config'])

    return db


def stop_script(signum=None, frame=None):
    global is_running
    is_running = False


config = get_config('config.json')
config_logging(config['logging'])

parser = argparse.ArgumentParser()
parser.add_argument('--operation', type=str, required=False, default="populate",
                    help='Operation to perform (populate, update, new_stocks)')
parser.add_argument('--startdate', type=str, required=False, default=default_history_start,
                    help='Start date for stock history (2015-12-30)')
args = parser.parse_args()

if args.operation not in {'populate', 'update', 'new_stocks'}:
    logging.error('Invalid operation')
    exit()

db = init_database(config['storage'])

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
        items = db.get_next_symbols_for_update(update_group_count)
        if len(items) == 0:
            logging.info('No more symbols to update.')
            break

        importer.update(items)

        time.sleep(wait_time)

    elif args.operation == 'new_stocks':
        alphavantage_url = f"{config['alphavantage']['url']}?apikey={config['alphavantage']['apikey']}"

        new_stocks_url = f"{alphavantage_url}&function=LISTING_STATUS"
        data = pd.read_csv(new_stocks_url, na_values='null', keep_default_na=False)
        data = data.replace(np.nan, None)
        logging.info(f"Downloaded {len(data)} active stocks.")

        count = db.insert_stocks(data.values.tolist())
        logging.info(f"Inserted {count} new active stocks in the DB.")

        today = str(date.today())
        delisted_stocks_url = f"{alphavantage_url}&function=LISTING_STATUS&date={today}&state=delisted"
        data = pd.read_csv(delisted_stocks_url, na_values='null', keep_default_na=False)
        data = data.replace(np.nan, None)
        logging.info(f"Downloaded {len(data)} delisted stocks.")

        count = db.insert_stocks(data.values.tolist())
        logging.info(f"Inserted {count} new delisted stocks in the DB.")

        break

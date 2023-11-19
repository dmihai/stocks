CREATE TABLE stocks (
    symbol         TEXT (15) PRIMARY KEY NOT NULL,
    name           TEXT,
    exchange       TEXT,
    asset_type     TEXT,
    ipo_date       TEXT,
    delisting_date TEXT,
    status         TEXT,
    history        TEXT,
    history_start  TEXT,
    history_end    TEXT
);

CREATE TABLE history (
    symbol    TEXT (10),
    date      TEXT (10),
    open      REAL,
    high      REAL,
    low       REAL,
    close     REAL,
    volume    INTEGER,
    dividends REAL,
    splits    REAL
);

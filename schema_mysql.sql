CREATE TABLE `stocks` (
  `symbol` varchar(15) NOT NULL,
  `name` varchar(200) NOT NULL,
  `exchange` varchar(10) NOT NULL,
  `asset_type` varchar(10) NOT NULL,
  `ipo_date` date DEFAULT NULL,
  `delisting_date` date DEFAULT NULL,
  `status` varchar(10) NOT NULL,
  `history` varchar(10) DEFAULT NULL,
  `history_start` date DEFAULT NULL,
  `history_end` date DEFAULT NULL,
  `last_update` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE `stocks`
  ADD PRIMARY KEY (`symbol`);

CREATE TABLE `history` (
  `symbol` char(15) NOT NULL,
  `date` date NOT NULL,
  `open` double NOT NULL,
  `high` double NOT NULL,
  `low` double NOT NULL,
  `close` double NOT NULL,
  `volume` int NOT NULL,
  `dividends` double NOT NULL,
  `splits` double NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE `history`
  ADD PRIMARY KEY (`symbol`,`date`);

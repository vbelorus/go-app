#!/bin/bash
set -e

clickhouse client -n <<-EOSQL
	CREATE DATABASE IF NOT EXISTS app;
	create table if not exists app.events
  (
      client_time DateTime,
      device_id   String,
      device_os   String,
      session     String,
      sequence    int,
      event       String,
      param_int   int,
      param_str   String,
      server_time DateTime,
      client_id   String
  ) engine = MergeTree() order by device_id;
EOSQL
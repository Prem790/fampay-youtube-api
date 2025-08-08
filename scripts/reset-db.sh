#!/bin/bash

psql -U postgres -d fampay_youtube -c "DROP TABLE IF EXISTS videos;"
psql -U postgres -d fampay_youtube -f migrations/001_create_videos_table.sql

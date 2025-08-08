#!/bin/bash

psql -U postgres -d fampay_youtube -c "INSERT INTO videos (id, title, description, published_at) VALUES
('sample1', 'Test Video 1', 'Description 1', NOW()),
('sample2', 'Test Video 2', 'Description 2', NOW());"

#!/usr/bin/env bash

docker compose -f env/docker-compose.yml --env-file .env.local up -d

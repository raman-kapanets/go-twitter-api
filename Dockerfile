FROM postgres:11-alpine
COPY scripts/create-multiple-postgresql-databases.sh /docker-entrypoint-initdb.d/

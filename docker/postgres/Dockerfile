FROM postgis/postgis:18-3.6

COPY init-database.sh /docker-entrypoint-initdb.d/10-init-database.sh
RUN chmod +x /docker-entrypoint-initdb.d/10-init-database.sh
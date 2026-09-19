FROM ghcr.io/project-osrm/osrm-backend:v26.9.0-debian

RUN apt-get update \
    && apt-get install -y --no-install-recommends curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN chmod +x ./cmd/prepare.sh

EXPOSE 5000

ENTRYPOINT ["./cmd/prepare.sh"]

CMD ["osrm-routed", "--algorithm", "mld", "--port", "5000", "/data/SanktPetersburg.osrm"]

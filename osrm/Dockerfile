FROM ghcr.io/project-osrm/osrm-backend:v26.9.0-debian

ARG PBF_URL=https://download.bbbike.org/osm/bbbike/SanktPetersburg/SanktPetersburg.osm.pbf

RUN apt-get update \
    && apt-get install -y --no-install-recommends curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN chmod +x ./cmd/prepare.sh \
    && mkdir -p /source \
    && if [ -f ./data/SanktPetersburg.osm.pbf ]; then \
        cp ./data/SanktPetersburg.osm.pbf /source/SanktPetersburg.osm.pbf; \
    else \
        curl --fail --location --retry 3 \
            --output /source/SanktPetersburg.osm.pbf.download \
            "$PBF_URL" \
        && mv /source/SanktPetersburg.osm.pbf.download /source/SanktPetersburg.osm.pbf; \
    fi

EXPOSE 5000

ENTRYPOINT ["./cmd/prepare.sh"]

CMD ["osrm-routed", "--algorithm", "mld", "--port", "5000", "/data/SanktPetersburg.osrm"]

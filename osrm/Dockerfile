FROM ghcr.io/project-osrm/osrm-backend:v26.9.0-debian AS builder

ARG PBF_URL=https://www.osmlab.ru/regions/RU-SPE.pbf

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY . .

RUN mkdir -p /source /graph \
    && if [ -f ./data/SanktPetersburg.osm.pbf ]; then \
    cp ./data/SanktPetersburg.osm.pbf \
    /source/SanktPetersburg.osm.pbf; \
    else \
    curl --fail --location --retry 3 \
    --output /source/SanktPetersburg.osm.pbf \
    "$PBF_URL"; \
    fi

RUN osrm-extract \
    -p /opt/foot.lua \
    --output /graph/SanktPetersburg.osrm \
    /source/SanktPetersburg.osm.pbf \
    && osrm-partition /graph/SanktPetersburg.osrm \
    && osrm-customize /graph/SanktPetersburg.osrm \
    && osrm-routed \
    --algorithm mld \
    --trial=true \
    /graph/SanktPetersburg.osrm


FROM ghcr.io/project-osrm/osrm-backend:v26.9.0-debian

COPY --from=builder /graph /data

EXPOSE 5000

CMD ["osrm-routed", "--algorithm", "mld", "--port", "5000", "/data/SanktPetersburg.osrm"]
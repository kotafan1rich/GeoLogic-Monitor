FROM ghcr.io/project-osrm/osrm-backend:v26.9.0-debian

WORKDIR /app

COPY . .

RUN chmod +x ./cmd/prepare.sh

EXPOSE 5000

ENTRYPOINT ["./cmd/prepare.sh"]

CMD ["osrm-routed", "--algorithm", "mld", "--port", "5000", "/data/SanktPetersburg.osrm"]
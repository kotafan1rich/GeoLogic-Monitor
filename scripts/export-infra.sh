#!/bin/sh

set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
project_dir=$(CDPATH= cd -- "$script_dir/.." && pwd)
output_dir=${1:-"$project_dir/exports/infra"}
tmp_dir=$(mktemp -d)
trap 'rm -r "$tmp_dir"' EXIT HUP INT TERM

cd "$project_dir"

docker compose exec -T postgres sh -c \
	'exec psql -q -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d api' \
	> "$tmp_dir/infra-source.csv" <<'SQL'
COPY (
	SELECT
		o.id::text AS id,
		o.external_id,
		t.slug AS type_slug,
		t.name AS type_name,
		t.max_radius,
		coalesce(nullif(btrim(o.name), ''), t.name) AS name,
		o.address,
		ST_Y(o.location) AS latitude,
		ST_X(o.location) AS longitude,
		o.created_at,
		o.updated_at
	FROM infra_objects o
	JOIN infra_types t ON t.id = o.type_id
	ORDER BY t.slug, o.id
) TO STDOUT WITH (FORMAT CSV, HEADER TRUE, ENCODING 'UTF8');
SQL

python3 - "$tmp_dir" <<'PY'
import csv
import json
import re
import sys
from collections import defaultdict
from pathlib import Path

export_dir = Path(sys.argv[1])
source_csv = export_dir / "infra-source.csv"
full_geojson = export_dir / "infra-all.geojson"
category_dir = export_dir / "by-category"
category_dir.mkdir()

expected_header = [
    "id",
    "external_id",
    "type_slug",
    "type_name",
    "max_radius",
    "name",
    "address",
    "latitude",
    "longitude",
    "created_at",
    "updated_at",
]
with source_csv.open(newline="", encoding="utf-8") as source:
    reader = csv.DictReader(source)
    if reader.fieldnames != expected_header:
        raise SystemExit("invalid infrastructure export header")
    rows = list(reader)

for row in rows:
    latitude = float(row["latitude"])
    longitude = float(row["longitude"])
    max_radius = int(row["max_radius"])
    if not (-90 <= latitude <= 90 and -180 <= longitude <= 180):
        raise SystemExit(f"invalid coordinates: {latitude}, {longitude}")
    if max_radius <= 0:
        raise SystemExit(f'invalid max_radius for category {row["type_slug"]!r}')


def feature(row):
    return {
        "type": "Feature",
        "id": row["id"],
        "geometry": {
            "type": "Point",
            "coordinates": [float(row["longitude"]), float(row["latitude"])],
        },
        "properties": {
            "id": row["id"],
            "external_id": row["external_id"],
            "type_slug": row["type_slug"],
            "type_name": row["type_name"],
            "max_radius": int(row["max_radius"]),
            "name": row["name"],
            "address": row["address"],
            "created_at": row["created_at"],
            "updated_at": row["updated_at"],
        },
    }


def write_geojson(path, source_rows):
    with path.open("w", encoding="utf-8", newline="\n") as target:
        target.write('{"type":"FeatureCollection","features":[')
        for index, row in enumerate(source_rows):
            if index:
                target.write(",")
            json.dump(feature(row), target, ensure_ascii=False, separators=(",", ":"))
        target.write("]}\n")


write_geojson(full_geojson, rows)

category_rows = defaultdict(list)
for row in rows:
    category_rows[row["type_slug"]].append(row)

used_file_slugs = set()
category_files = []
for slug, grouped_rows in sorted(category_rows.items()):
    file_slug = re.sub(r"[^A-Za-z0-9_-]+", "_", slug).strip("_-")
    if not file_slug or file_slug in used_file_slugs:
        raise SystemExit(f"category slug cannot be used as a unique filename: {slug!r}")
    used_file_slugs.add(file_slug)
    path = category_dir / f"{file_slug}.geojson"
    write_geojson(path, grouped_rows)
    category_files.append((slug, path, len(grouped_rows)))

with full_geojson.open(encoding="utf-8") as source:
    full_data = json.load(source)
if full_data.get("type") != "FeatureCollection" or len(full_data.get("features", [])) != len(rows):
    raise SystemExit("full GeoJSON feature count does not match database export")
if any(not isinstance(item["properties"].get("max_radius"), int) for item in full_data["features"]):
    raise SystemExit("max_radius must be an integer in every feature")

category_total = 0
for slug, path, expected_count in category_files:
    with path.open(encoding="utf-8") as source:
        data = json.load(source)
    features = data.get("features", [])
    if data.get("type") != "FeatureCollection" or len(features) != expected_count:
        raise SystemExit(f"{slug}: invalid GeoJSON feature count")
    if any(feature["properties"].get("type_slug") != slug for feature in features):
        raise SystemExit(f"{slug}: category file contains another type")
    category_total += len(features)
if category_total != len(rows):
    raise SystemExit("category exports do not cover the full dataset")

print(f"Exported {len(rows)} objects to GeoJSON across {len(category_files)} categories.")
PY

mkdir -p "$output_dir"
for legacy_file in \
	infra-all-yandex.csv \
	infra-all.gpx \
	infra-all-kepler.csv \
	infra-all-kepler.geojson \
	infra-yandex-sample-10000.csv \
	infra-yandex-sample-10000.gpx \
	infra-counts-by-type.csv
do
	rm -f "$output_dir/$legacy_file"
done
if [ -d "$output_dir/by-category" ]; then
	rm -r "$output_dir/by-category"
fi
mv "$tmp_dir/infra-all.geojson" "$output_dir/infra-all.geojson"
mv "$tmp_dir/by-category" "$output_dir/by-category"

printf 'Saved GeoJSON exports to %s\n' "$output_dir"

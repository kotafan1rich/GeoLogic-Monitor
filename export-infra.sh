#!/bin/sh

set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
output_dir=${1:-"$script_dir/exports/infra"}
tmp_dir=$(mktemp -d)
trap 'rm -r "$tmp_dir"' EXIT HUP INT TERM

cd "$script_dir"

psql_api() {
	docker compose exec -T postgres sh -c 'exec psql -q -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d api'
}

psql_api > "$tmp_dir/infra-all-yandex.csv" <<'SQL'
COPY (
	SELECT
		ST_Y(o.location)::text,
		ST_X(o.location)::text,
		concat(
			'Тип: ', t.name, ' (', t.slug, '); Адрес: ', o.address,
			'; ID: ', o.external_id
		),
		coalesce(nullif(btrim(o.name), ''), t.name),
		''
	FROM infra_objects o
	JOIN infra_types t ON t.id = o.type_id
	ORDER BY t.slug, o.id
) TO STDOUT WITH (FORMAT CSV, ENCODING 'UTF8');
SQL

psql_api > "$tmp_dir/infra-yandex-sample-10000.csv" <<'SQL'
COPY (
	WITH ranked AS (
		SELECT
			o.id,
			o.location,
			o.address,
			o.external_id,
			o.name AS object_name,
			t.slug,
			t.name AS type_name,
			row_number() OVER (
				PARTITION BY t.slug
				ORDER BY md5(o.id::text)
			) AS type_rank
		FROM infra_objects o
		JOIN infra_types t ON t.id = o.type_id
	),
	stats AS (
		SELECT least(10000, count(*)) AS sample_size
		FROM ranked
	),
	guaranteed AS (
		SELECT *
		FROM ranked
		WHERE type_rank = 1
	),
	remainder AS (
		SELECT *
		FROM ranked
		WHERE type_rank > 1
		ORDER BY md5(id::text)
		LIMIT (
			SELECT sample_size - (SELECT count(*) FROM guaranteed)
			FROM stats
		)
	),
	sample AS (
		SELECT * FROM guaranteed
		UNION ALL
		SELECT * FROM remainder
	)
	SELECT
		ST_Y(sample.location)::text,
		ST_X(sample.location)::text,
		concat(
			'Тип: ', sample.type_name, ' (', sample.slug, '); Адрес: ', sample.address,
			'; ID: ', sample.external_id
		),
		coalesce(nullif(btrim(sample.object_name), ''), sample.type_name),
		''
	FROM sample
	ORDER BY sample.slug, sample.id
) TO STDOUT WITH (FORMAT CSV, ENCODING 'UTF8');
SQL

psql_api > "$tmp_dir/infra-counts-by-type.csv" <<'SQL'
COPY (
	SELECT t.slug, t.name, count(*)
	FROM infra_objects o
	JOIN infra_types t ON t.id = o.type_id
	GROUP BY t.slug, t.name
	ORDER BY count(*) DESC, t.slug
) TO STDOUT WITH (FORMAT CSV, HEADER TRUE, ENCODING 'UTF8');
SQL

python3 - "$tmp_dir" <<'PY'
import csv
import sys
from pathlib import Path
from xml.etree import ElementTree
from xml.sax.saxutils import escape

export_dir = Path(sys.argv[1])
full_csv = export_dir / "infra-all-yandex.csv"
sample_csv = export_dir / "infra-yandex-sample-10000.csv"
summary_csv = export_dir / "infra-counts-by-type.csv"
gpx = export_dir / "infra-yandex-sample-10000.gpx"

with sample_csv.open(newline="", encoding="utf-8") as source, gpx.open(
    "w", encoding="utf-8", newline="\n"
) as target:
    target.write('<?xml version="1.0" encoding="UTF-8"?>\n')
    target.write(
        '<gpx version="1.1" creator="GeoLogic Monitor" '
        'xmlns="http://www.topografix.com/GPX/1/1">\n'
    )
    for lat, lon, description, label, _ in csv.reader(source):
        target.write(
            f'  <wpt lat="{lat}" lon="{lon}">'
            f"<name>{escape(label)}</name>"
            f"<desc>{escape(description)}</desc></wpt>\n"
        )
    target.write("</gpx>\n")


def read_yandex_csv(path: Path) -> list[list[str]]:
    with path.open(newline="", encoding="utf-8") as source:
        rows = list(csv.reader(source))
    if any(len(row) != 5 for row in rows):
        raise SystemExit(f"{path.name}: expected exactly five columns")
    for row in rows:
        lat, lon = map(float, row[:2])
        if not (-90 <= lat <= 90 and -180 <= lon <= 180):
            raise SystemExit(f"{path.name}: invalid coordinates {lat}, {lon}")
    return rows


full_rows = read_yandex_csv(full_csv)
sample_rows = read_yandex_csv(sample_csv)
expected_sample_size = min(10000, len(full_rows))
if len(sample_rows) != expected_sample_size:
    raise SystemExit(
        f"sample contains {len(sample_rows)} rows, expected {expected_sample_size}"
    )

with summary_csv.open(newline="", encoding="utf-8") as source:
    summary_rows = list(csv.reader(source))
if not summary_rows or summary_rows[0] != ["slug", "name", "count"]:
    raise SystemExit("invalid category summary header")
if sum(int(row[2]) for row in summary_rows[1:]) != len(full_rows):
    raise SystemExit("category summary total does not match full export")

root = ElementTree.parse(gpx).getroot()
waypoints = root.findall("{http://www.topografix.com/GPX/1/1}wpt")
if len(waypoints) != len(sample_rows):
    raise SystemExit("GPX waypoint count does not match CSV sample")

print(
    f"Exported {len(full_rows)} objects, "
    f"{len(sample_rows)} objects in the Yandex sample, "
    f"{len(summary_rows) - 1} categories."
)
PY

mkdir -p "$output_dir"
for name in \
	infra-all-yandex.csv \
	infra-yandex-sample-10000.csv \
	infra-yandex-sample-10000.gpx \
	infra-counts-by-type.csv
do
	mv "$tmp_dir/$name" "$output_dir/$name"
done

printf 'Saved infrastructure exports to %s\n' "$output_dir"

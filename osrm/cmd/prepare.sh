#!/bin/sh

set -eu

source=../source/SanktPetersburg.osm.pbf
marker=../data/.prepared-v26.9.0-foot

test -f "$source"

if [ ! -f "$marker" ] || [ "$source" -nt "$marker" ]; then
	echo "Building OSRM pedestrian graph..."

	osrm-extract \
		-p /opt/foot.lua \
		--output ../data/SanktPetersburg.osrm \
		"$source"

	osrm-partition ../data/SanktPetersburg.osrm
	osrm-customize ../data/SanktPetersburg.osrm

	osrm-routed \
		--algorithm mld \
		--trial=true \
		../data/SanktPetersburg.osrm

	touch "$marker"
fi

exec "$@"

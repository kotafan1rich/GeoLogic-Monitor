package geo

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

const (
	minEWKBLen = 25
	sridFlag   = 0x20000000
)

type GeoPoint struct {
	Lat float64
	Lon float64
}

func (p GeoPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", p.Lon, p.Lat), nil
}

func (p *GeoPoint) Scan(val any) error {
	var source []byte
	switch v := val.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return fmt.Errorf("%w: %T", ErrInvalidType, val)
	}

	decoded := make([]byte, hex.DecodedLen(len(source)))
	_, err := hex.Decode(decoded, source)
	if err != nil {
		decoded = source
	}

	if len(decoded) < minEWKBLen {
		return fmt.Errorf("%w: %d", ErrInvalidEWKBLen, len(decoded))
	}

	var order binary.ByteOrder = binary.BigEndian
	if decoded[0] == 1 {
		order = binary.LittleEndian
	}

	geomType := order.Uint32(decoded[1:5])

	var offset int
	if (geomType & sridFlag) != 0 {
		offset = 9
	} else {
		offset = 5
	}

	if len(decoded) < offset+16 {
		return ErrInsufficientBytes
	}

	lonBits := order.Uint64(decoded[offset : offset+8])
	latBits := order.Uint64(decoded[offset+8 : offset+16])

	p.Lon = math.Float64frombits(lonBits)
	p.Lat = math.Float64frombits(latBits)

	return nil
}

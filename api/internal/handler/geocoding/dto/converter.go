package dto

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"

func ToResponseList(addresses []domain.Address) []*AddressResponse {
	result := make([]*AddressResponse, 0, len(addresses))
	for _, address := range addresses {
		result = append(result, &AddressResponse{
			Address: address.Address,
			Lat:     address.Lat,
			Lon:     address.Lon,
		})
	}
	return result
}

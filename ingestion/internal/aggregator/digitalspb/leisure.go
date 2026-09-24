package digitalspb

import (
	"context"
	"fmt"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const (
	restaurantEndpoint     = "datasets/143/versions/latest/data/570/"
	hotelEndpoint          = "datasets/132/versions/latest/data/153/"
	cinemaEndpoint         = "datasets/230/versions/latest/data/345/"
	museumEndpoint         = "datasets/139/versions/latest/data/569/"
	exhibitionHallEndpoint = "datasets/130/versions/latest/data/151/"
	theatreEndpoint        = "datasets/145/versions/latest/data/649/"
	pharmacyEndpoint       = "datasets/122/versions/latest/data/142/"
	vetClinicEndpoint      = "datasets/217/versions/latest/data/247/"
	kidsPlaceEndpoint      = "iparent/places/all/"
)

func (c *Client) ParseRestaurantData(ctx context.Context, src a.URL) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseRestaurantData"

	data, err := fetchSpbClassifGate[Restaurant](ctx, c, src, restaurantEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData := mapToInfraObjects[Restaurant](data, "", mapRestaurant)

	return mappedData, nil
}

func (c *Client) ParseHotelData(ctx context.Context, src a.URL) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseHotelData"

	data, err := fetchSpbClassifGate[Hotel](ctx, c, src, hotelEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData := mapToInfraObjects[Hotel](data, "", mapHotel)

	return mappedData, nil
}

func (c *Client) ParseCinemaData(
	ctx context.Context, src a.URL, convert AddressConverter,
) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseCinemaData"

	data, err := fetchSpbClassifGate[Cinema](ctx, c, src, cinemaEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData, err := mapToGeocodedInfraObjects[Cinema](ctx, data, "", convert, mapCinema)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mappedData, nil
}

func (c *Client) ParseMuseumData(ctx context.Context, src a.URL) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseMuseumData"

	data, err := fetchSpbClassifGate[Museum](ctx, c, src, museumEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData := mapToInfraObjects[Museum](data, "", mapMuseum)

	return mappedData, nil
}

func (c *Client) ParseExhibitionHallData(ctx context.Context, src a.URL) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseExhibitionHallData"

	data, err := fetchSpbClassifGate[ExhibitionHall](ctx, c, src, exhibitionHallEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData := mapToInfraObjects[ExhibitionHall](data, "", mapExhibitionHall)

	return mappedData, nil
}

func (c *Client) ParseTheatreData(ctx context.Context, src a.URL) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseTheatreData"

	data, err := fetchSpbClassifGate[Theatre](ctx, c, src, theatreEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData := mapToInfraObjects[Theatre](data, "", mapTheatre)

	return mappedData, nil
}

func (c *Client) ParsePharmacyData(ctx context.Context, src a.URL) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParsePharmacyData"

	data, err := fetchSpbClassifGate[Pharmacy](ctx, c, src, pharmacyEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData := mapToInfraObjects[Pharmacy](data, "", mapPharmacy)

	return mappedData, nil
}

func (c *Client) ParseVetClinicData(
	ctx context.Context, src a.URL, convert AddressConverter,
) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseVetClinicData"

	data, err := fetchSpbClassifGate[VetClinic](ctx, c, src, vetClinicEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData, err := mapToGeocodedInfraObjects[VetClinic](ctx, data, "", convert, mapVetClinic)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mappedData, nil
}

func (c *Client) ParseKidsPlaceData(
	ctx context.Context, src a.URL, convert CoordinatesConverter,
) ([]geoapi.InfraObjectInput, error) {
	const op = "digitalspb.Client.ParseKidsPlaceData"

	data, err := fetchYazzhGate[KidsPlace](ctx, c, src, kidsPlaceEndpoint, op)
	if err != nil {
		return nil, err
	}

	mappedData, err := mapToAddressedInfraObjects[KidsPlace](ctx, data, "", convert, mapKidsPlace)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mappedData, nil
}

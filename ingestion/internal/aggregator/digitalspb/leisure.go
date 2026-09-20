package digitalspb

import (
	"context"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
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

func (c *Client) ParseRestaurantData(ctx context.Context, src a.URL) ([]Restaurant, error) {
	const op = "digitalspb.Client.ParseRestaurantData"

	data, err := fetchSpbClassifGate[Restaurant](ctx, c, src, restaurantEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseHotelData(ctx context.Context, src a.URL) ([]Hotel, error) {
	const op = "digitalspb.Client.ParseHotelData"

	data, err := fetchSpbClassifGate[Hotel](ctx, c, src, hotelEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseCinemaData(ctx context.Context, src a.URL) ([]Cinema, error) {
	const op = "digitalspb.Client.ParseCinemaData"

	data, err := fetchSpbClassifGate[Cinema](ctx, c, src, cinemaEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseMuseumData(ctx context.Context, src a.URL) ([]Museum, error) {
	const op = "digitalspb.Client.ParseMuseumData"

	data, err := fetchSpbClassifGate[Museum](ctx, c, src, museumEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseExhibitionHallData(ctx context.Context, src a.URL) ([]ExhibitionHall, error) {
	const op = "digitalspb.Client.ParseExhibitionHallData"

	data, err := fetchSpbClassifGate[ExhibitionHall](ctx, c, src, exhibitionHallEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseTheatreData(ctx context.Context, src a.URL) ([]Theatre, error) {
	const op = "digitalspb.Client.ParseTheatreData"

	data, err := fetchSpbClassifGate[Theatre](ctx, c, src, theatreEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParsePharmacyData(ctx context.Context, src a.URL) ([]Pharmacy, error) {
	const op = "digitalspb.Client.ParsePharmacyData"

	data, err := fetchSpbClassifGate[Pharmacy](ctx, c, src, pharmacyEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseVetClinicData(ctx context.Context, src a.URL) ([]VetClinic, error) {
	const op = "digitalspb.Client.ParseVetClinicData"

	data, err := fetchSpbClassifGate[VetClinic](ctx, c, src, vetClinicEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ParseKidsPlaceData(ctx context.Context, src a.URL) ([]KidsPlace, error) {
	const op = "digitalspb.Client.ParseKidsPlaceData"

	data, err := fetchYazzhGate[KidsPlace](ctx, c, src, kidsPlaceEndpoint, op)
	if err != nil {
		return nil, err
	}

	return data, nil
}

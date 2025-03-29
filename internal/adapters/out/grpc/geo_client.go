package grpc

import (
	"context"
	"fmt"

	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/pkg/clients/geopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GeoClient struct {
	grpcConnection *grpc.ClientConn
	grpcGeoClient  geopb.GeoClient
}

func NewGeoClient(geoServiceAddress string) (*GeoClient, error) {
	connection, err := grpc.NewClient(
		geoServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC Geo client: %w", err)
	}

	grpcGeoClient := geopb.NewGeoClient(connection)

	return &GeoClient{
		grpcConnection: connection,
		grpcGeoClient:  grpcGeoClient,
	}, nil
}

func (c *GeoClient) GetLocation(ctx context.Context, street string) (*kernel.Location, error) {
	request := &geopb.GetGeolocationRequest{Street: street}
	response, err := c.grpcGeoClient.GetGeolocation(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("request to Geo service failed: %w", err)
	}

	location, err := kernel.NewLocation(
		int(response.Location.X),
		int(response.Location.Y),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create location from Geo response: %w, location from Geo: %v",
			err,
			response.Location,
		)
	}

	return location, nil
}

func (c *GeoClient) Close() error {
	err := c.grpcConnection.Close()
	if err != nil {
		return fmt.Errorf("failed to close Geo client gRPC connection: %w", err)
	}

	return nil
}

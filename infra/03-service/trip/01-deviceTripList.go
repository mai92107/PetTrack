package trip

import (
	"PetTrack/core/model"
	"PetTrack/core/util/formatter"
	"context"
	"fmt"
)

func (s *TripServiceImpl) GetDeviceTrips(ctx context.Context, member model.Claims, deviceId string, pageable model.Pageable) ([]map[string]interface{}, int64, int64, error) {
	trips := []map[string]interface{}{}
	var total int64
	var totalPages int64

	err := validateTripsRequest(deviceId)
	if err != nil {
		return trips, total, totalPages, err
	}

	err = s.commonService.ValidateDeviceOwner(ctx, deviceId, member)
	if err != nil {
		return trips, total, totalPages, err
	}

	tripsData, total, totalPages, err := s.tripRepo.GetDeviceTrips(ctx, deviceId, pageable)
	if err != nil {
		return trips, total, totalPages, err
	}

	for _, v := range tripsData {
		trips = append(trips, map[string]interface{}{
			"uuid":     v.DataRef,
			"time":     formatter.ToLocalTimeShortStr(v.EndTime),
			"duration": formatter.FormatDigits(v.DurationMinutes, 4),
			"distance": formatter.FormatDigits(v.DistanceKM, 4),
		})
	}

	return trips, total, totalPages, nil
}

func validateTripsRequest(deviceId string) error {
	if deviceId == "" {
		return fmt.Errorf("deviceId 參數錯誤")
	}
	return nil
}

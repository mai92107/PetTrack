package trip

import (
	"PetTrack/core/model"
	"PetTrack/core/util/formatter"
	"context"
	"fmt"
)

func (s *TripServiceImpl) GetTripDetail(ctx context.Context, member model.Claims, deviceId string, tripUuid string) (map[string]interface{}, error) {
	trip := map[string]interface{}{}

	err := validateTripDetailRequest(deviceId, tripUuid)
	if err != nil {
		return trip, err
	}
	err = s.commonService.ValidateDeviceOwner(ctx, deviceId, member)
	if err != nil {
		return trip, err
	}

	tripData, err := s.tripRepo.GetTripDetail(ctx, tripUuid)
	if err != nil {
		return trip, err
	}

	trip = map[string]interface{}{
		"tripUuid":  tripData.DataRef,
		"deviceId":  tripData.DeviceID,
		"distance":  tripData.DistanceKM,
		"duration":  tripData.DurationMinutes,
		"startAt":   formatter.ToLocalTimeStr(tripData.StartTime),
		"endAt":     formatter.ToLocalTimeStr(tripData.EndTime),
		"point":     tripData.PointCount,
		"createdAt": formatter.ToLocalTime(tripData.CreatedAt),
		"updatedAt": formatter.ToLocalTime(tripData.UpdatedAt),
	}

	return trip, nil
}

func validateTripDetailRequest(deviceId string, tripUuid string) error {
	if deviceId == "" {
		return fmt.Errorf("deviceId 參數錯誤")
	}
	if tripUuid == "" {
		return fmt.Errorf("tripUuid 參數錯誤")
	}
	return nil
}

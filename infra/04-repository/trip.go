package repo

import (
	"PetTrack/core/model"
	domainRepo "PetTrack/domain/repo"
	"context"
	"fmt"
	"math"
)

func (r *tripRepoImpl) GetDeviceTrips(ctx context.Context, deviceId string, pageable model.Pageable) ([]domainRepo.TripSummary, int64, int64, error) {
	var deviceTrips []domainRepo.TripSummary
	var totalCount int64
	var totalPage int64

	// 查總筆數
	if err := r.read.Model(&domainRepo.TripSummary{}).
		Where("device_id = ?", deviceId).
		Count(&totalCount).Error; err != nil {
		// logafa.Error("統計裝置行程數量失敗 deviceId=%s, error: %+v", deviceId, err)
		return deviceTrips, totalCount, totalPage, fmt.Errorf("統計行程數量失敗")
	}

	// 如果總筆數為 0，直接回傳空陣列
	if totalCount == 0 {
		// logafa.Info("裝置 %s 無任何行程紀錄", deviceId)
		return deviceTrips, totalCount, totalPage, nil
	}

	totalPage = int64(math.Ceil(float64(totalCount) / float64(pageable.Size)))

	// 2. 正式查詢資料（分頁 + 排序）
	err := r.read.WithContext(ctx).
		Where("device_id = ?", deviceId).
		Offset(pageable.Offset()).    // 分頁
		Limit(pageable.Limit()).      // 每頁筆數
		Order(pageable.OrderBySQL()). // 排序
		Find(&deviceTrips).Error

	if err != nil {
		// logafa.Error("查詢裝置行程失敗 deviceId=%s, error: %+v", deviceId, err)
		return deviceTrips, totalCount, totalPage, fmt.Errorf("查詢行程失敗")
	}

	return deviceTrips, totalCount, totalPage, nil
}

func (r *tripRepoImpl) GetTripDetail(ctx context.Context, tripUuid string) (domainRepo.TripSummary, error) {
	var trip domainRepo.TripSummary

	err := r.read.WithContext(ctx).
		Where("data_ref = ?", tripUuid).
		First(&trip).Error
	if err != nil {
		// logafa.Error("查詢裝置行程失敗 data_ref=%s, error: %+v", tripUuid, err)
		return trip, fmt.Errorf("查詢行程失敗")
	}
	return trip, nil
}

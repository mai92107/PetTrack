package trip

import (
	"PetTrack/infra/00-core/util/logafa"
	"context"
	"time"
)

// 計算近 30min 每趟行程資訊
func (s *TripServiceImpl) FlushTripFmMongoToMaria(ctx context.Context, d time.Duration) {
	trips, err := s.tripRepo.ReadTripFromMongo(ctx, d)
	if err != nil {
		logafa.Error("讀取mongo旅程資料錯誤", "error", err)
	}
	err = s.tripRepo.SaveTripToDB(ctx, trips)
	if err != nil {
		logafa.Error("儲存旅程資料至DB錯誤", "error", err)
	}
}

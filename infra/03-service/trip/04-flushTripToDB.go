package trip

import (
	"PetTrack/core/util/logafa"
)

// 計算近 30min 每趟行程資訊
func (s *TripServiceImpl) FlushTripFmMongoToMaria(timeDuration int) {
	trips, err := s.tripRepo.ReadTripFromMongo(timeDuration)
	if err != nil {
		logafa.Error("讀取mongo旅程資料錯誤", "error", err)
	}
	err = s.tripRepo.SaveTripToDB(trips)
	if err != nil {
		logafa.Error("儲存旅程資料至DB錯誤", "error", err)
	}
}

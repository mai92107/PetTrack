package trip

import (
	"PetTrack/infra/00-core/util/logafa"
	domainRepo "PetTrack/domain/repo"
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func (s *TripServiceImpl) SaveGpsFmRdsToMongo(ctx context.Context) {
	logafa.Info("開始執行 GPS DATA 持久化...")

	deviceKeyPattern := "device:*"
	keys, err := s.redisUtil.KeyScan(ctx, deviceKeyPattern)
	if err != nil {
		logafa.Error("取得 redis device key 值發生錯誤", "error", err)
		return
	}

	if len(keys) == 0 {
		logafa.Debug("無裝置資料, 罷工回家睡覺")
		return
	}

	logafa.Debug("取得裝置資料, 開始讀取", "筆數", len(keys))

	end := time.Now().UTC()
	start := end.Add(-5 * time.Minute)

	for _, key := range keys {
		datas, err := s.redisUtil.ZRangeByScore(ctx, key, start.UnixMilli(), end.UnixMilli())
		if err != nil {
			logafa.Error("取得 redis device data 發生錯誤", "key", key, "error", err)
			continue
		}

		if len(datas) == 0 {
			logafa.Debug("讀取到 0 筆資料", "KEY", key)
			continue
		}
		var records []domainRepo.DeviceLocation
		for _, jsonStr := range datas {
			data := domainRepo.GPS{}
			if err := jsoniter.UnmarshalFromString(jsonStr, &data); err != nil {
				logafa.Error("解析 GPS JSON 失敗", "json", jsonStr, "error", err)
				continue
			}
			record := domainRepo.DeviceLocation{
				DeviceID:   data.DeviceId,
				Location:   domainRepo.NewGeoJSONPoint(data.Longitude, data.Latitude),
				RecordedAt: data.RecordTime,
				DataRef:    data.DataRef,
				CreatedAt:  time.Now().UTC(),
			}
			records = append(records, record)
		}

		if err = s.tripRepo.SaveLocationToDB(ctx, records); err != nil {
			logafa.Error("批次寫入資料至 DB 失敗", "error", err)
			continue
		}

		// 只有成功寫入才刪除
		if err := s.redisUtil.ZRemRangeByScore(ctx, key, start.UnixMilli(), end.UnixMilli()); err != nil {
			logafa.Error("⚠️ 刪除 redis 資料失敗", "key", key, "error", err)
			// TODO: 觸發告警或記錄到監控系統
		}
	}
}

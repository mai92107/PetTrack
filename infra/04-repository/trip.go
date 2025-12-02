package repo

import (
	"PetTrack/core/model"
	"PetTrack/core/util/logafa"
	"PetTrack/core/util/parser"
	domainRepo "PetTrack/domain/repo"
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *tripRepoImpl) GetDeviceTrips(deviceId string, pageable model.Pageable) ([]domainRepo.TripSummary, int64, int64, error) {
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
	err := r.read.Where("device_id = ?", deviceId).
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

func (r *tripRepoImpl) GetTripDetail(tripUuid string) (domainRepo.TripSummary, error) {
	var trip domainRepo.TripSummary

	err := r.read.Where("data_ref = ?", tripUuid).
		First(&trip).Error
	if err != nil {
		// logafa.Error("查詢裝置行程失敗 data_ref=%s, error: %+v", tripUuid, err)
		return trip, fmt.Errorf("查詢行程失敗")
	}
	return trip, nil
}

func (r *tripRepoImpl) SaveLocationToDB(records []domainRepo.DeviceLocation) error {
	if len(records) < 1 {
		return fmt.Errorf("無有效紀錄可存入資料庫")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用 BulkWrite 進行 upsert,防止重複資料
	var operations []mongo.WriteModel
	for _, record := range records {
		filter := bson.M{
			"device_id":   record.DeviceID,
			"recorded_at": record.RecordedAt,
		}

		update := bson.M{
			"$setOnInsert": bson.M{
				"device_id":   record.DeviceID,
				"location":    record.Location,
				"recorded_at": record.RecordedAt,
				"data_ref":    record.DataRef,
				"created_at":  record.CreatedAt,
			},
		}

		operation := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, operation)
	}

	logafa.Debug("批次寫入資料至 DB...", " 筆數", len(records))
	collection := r.mongo.Collection("pettrack")
	result, err := collection.BulkWrite(ctx, operations)
	if err != nil {
		logafa.Error("MongoDB 批次寫入失敗: %+v", err)
		return err
	}

	logafa.Debug("資料成功批次寫入 DB", "筆數", result.UpsertedCount)
	return nil
}

func (r *tripRepoImpl) ReadTripFromMongo(timeDuration int) ([]domainRepo.TripSummary, error) {
	ctx := context.Background()
	coll := r.mongo.Collection("pettrack")

	duration := time.Now().UTC().Add(time.Minute * -(time.Duration(timeDuration)))
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "recorded_at", Value: bson.D{{Key: "$gte", Value: duration}}},
			{Key: "location", Value: bson.D{{Key: "$ne", Value: nil}}},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "data_ref", Value: 1},
			{Key: "recorded_at", Value: 1},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$data_ref"},
			{Key: "device_id", Value: bson.D{{Key: "$first", Value: "$device_id"}}},
			{Key: "start_time", Value: bson.D{{Key: "$min", Value: "$recorded_at"}}},
			{Key: "end_time", Value: bson.D{{Key: "$max", Value: "$recorded_at"}}},
			{Key: "coords", Value: bson.D{{Key: "$push", Value: "$location.coordinates"}}}, // [lng, lat]
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "data_ref", Value: "$_id"},
			{Key: "device_id", Value: 1},
			{Key: "start_time", Value: 1},
			{Key: "end_time", Value: 1},
			{Key: "coords", Value: 1},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		logafa.Error("Mongo 資料讀取錯誤", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domainRepo.TripSummary

	for cursor.Next(ctx) {
		rawData := decodeRawData(cursor)
		distance := getDistance(*rawData)
		duration := rawData.EndTime.Sub(rawData.StartTime).Minutes()
		now := time.Now()

		results = append(results, domainRepo.TripSummary{
			DataRef:         rawData.DataRef,
			DeviceID:        rawData.DeviceID,
			StartTime:       parser.ToUtcTime(rawData.StartTime),
			EndTime:         parser.ToUtcTime(rawData.EndTime),
			DurationMinutes: math.Round(duration*100) / 100, // 保留2位
			PointCount:      len(rawData.Coords),
			DistanceKM:      math.Round(distance*1000) / 1000, // 保留3位
			CreatedAt:       parser.ToUtcTime(now),
			UpdatedAt:       parser.ToUtcTime(now),
		})
	}
	return results, nil
}

func (r *tripRepoImpl) SaveTripToDB(results []domainRepo.TripSummary) error {
	tx := r.write.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("開始交易失敗: %w", err)
	}

	// 確保一定會 rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logafa.Error("交易 panic，已 rollback", "error", r)
		}
	}()

	for i, t := range results {
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "data_ref"}},
			DoUpdates: []clause.Assignment{
				{Column: clause.Column{Name: "updated_at"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(updated_at), updated_at)")},
				{Column: clause.Column{Name: "end_time"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(end_time), end_time)")},
				{Column: clause.Column{Name: "distance_km"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(distance_km), distance_km)")},
				{Column: clause.Column{Name: "duration_minutes"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(duration_minutes), duration_minutes)")},
				{Column: clause.Column{Name: "point_count"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(point_count), point_count)")},
			}}).Create(t).Error
		if err != nil {
			logafa.Error("儲存失敗，rollback", "error", err, "筆數", i+1, "原始資料編號", t.DataRef)
			return fmt.Errorf("儲存失敗: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		logafa.Error("交易提交失敗", "error", err)
		return fmt.Errorf("commit 失敗: %w", err)
	}

	logafa.Debug("全部寫入成功！")
	return nil
}

func getDistance(rawData domainRepo.RawTrip) float64 {
	distance := 0.0
	for i := 1; i < len(rawData.Coords); i++ {
		distance += haversine(
			rawData.Coords[i-1][1], rawData.Coords[i-1][0], // lat1, lng1
			rawData.Coords[i][1], rawData.Coords[i][0], // lat2, lng2
		)
	}
	return distance
}

func decodeRawData(cursor *mongo.Cursor) *domainRepo.RawTrip {
	var temp domainRepo.RawTrip
	if err := cursor.Decode(&temp); err != nil {
		log.Printf("decode error: %v", err)
		return nil
	}
	return &temp
}
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // 地球半徑（公尺）
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c / 1000 // 回傳公里
}

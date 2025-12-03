package initMethod

import (
	"PetTrack/infra/00-core/model"
	bun "PetTrack/infra/00-core/model/bunMachines"
	"PetTrack/infra/00-core/util/logafa"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg model.Config) *bun.DB {
	db, err := initDataBase(
		cfg.Machines.MariaDB.Host,
		cfg.Machines.MariaDB.Port,
		cfg.Machines.MariaDB.User,
		cfg.Machines.MariaDB.Password,
		cfg.Machines.MariaDB.Name,
	)
	if err != nil {
		return nil
	}
	return db
}
func initDataBase(host, port, user, password, database string) (*bun.DB, error) {
	// 暫時共用
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logafa.NewGormLogger(),
	})
	if err != nil {
		logafa.Error(" ❌ 無法連接讀取資料庫", "error", err)
		return nil, err
	}

	logafa.Debug(" ✅ 資料庫連接成功")
	initSQLTables(db)
	return &bun.DB{
		Write: db,
		Read:  db,
	}, nil
}

func InitMongo(cfg model.Config) *mongo.Database {
	mongo, err := initMongoDB(
		cfg.Machines.MongoDB.Host,
		cfg.Machines.MongoDB.Port,
		cfg.Machines.MongoDB.User,
		cfg.Machines.MongoDB.Password,
	)
	if err != nil {
		return nil
	}
	return mongo
}

func initMongoDB(host, port, user, password string) (*mongo.Database, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		user, password, host, port)

	clientOptions := options.Client().ApplyURI(uri)

	// 設置執行timeout時間
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: 增加retry
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logafa.Error("無法連接 Mongodb", "error", err)
		return nil, err
	}
	logafa.Debug("✅ 成功連線 MongoDB!")

	// 初始化index
	initMongoIndexes(client)
	return client.Database("pettrack"), nil
}

func initSQLTables(db *gorm.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Panic 保護
	defer func() {
		if r := recover(); r != nil {
			// logafa.Error("初始化 MySQL Tables 失敗 (panic): %v", r)
		}
	}()

	// 定義所有 Table
	tables := map[string]string{
		"member": `
			CREATE TABLE member (
				id         BIGINT AUTO_INCREMENT PRIMARY KEY,
				last_name  VARCHAR(255),
				first_name VARCHAR(255),
				nick_name  VARCHAR(255),
				email      VARCHAR(255),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		"device": `
			CREATE TABLE device (
				uuid CHAR(36) PRIMARY KEY,
				device_id VARCHAR(36) UNIQUE,
				device_type VARCHAR(50),
				create_by_member BIGINT NOT NULL,
				remark CHAR(50),
				CONSTRAINT fk_device_create_by_member 
					FOREIGN KEY (create_by_member) REFERENCES member(id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		"member_device": `
			CREATE TABLE member_device (
				id          BIGINT AUTO_INCREMENT PRIMARY KEY,
				member_id   BIGINT NOT NULL,
				device_id   VARCHAR(36) NOT NULL,
				device_name VARCHAR(255),
				created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				
				CONSTRAINT fk_member_device_member 
					FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
				CONSTRAINT fk_member_device_device 
					FOREIGN KEY (device_id) REFERENCES device(device_id) ON DELETE CASCADE,
				CONSTRAINT uq_member_device UNIQUE (member_id, device_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		"account": `
			CREATE TABLE account (
				uuid            CHAR(36) PRIMARY KEY,
				member_id       BIGINT NOT NULL,
				username        VARCHAR(255) NOT NULL UNIQUE,
				password        VARCHAR(255) NOT NULL,
				email           VARCHAR(255) NOT NULL UNIQUE,
				identity        VARCHAR(50),
				last_login_time DATETIME,
				created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				
				CONSTRAINT fk_account_member 
					FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		"password_history": `
			CREATE TABLE password_history (
				id BIGINT     AUTO_INCREMENT PRIMARY KEY,
				account_uuid  CHAR(36) NOT NULL,
				password      VARCHAR(255) NOT NULL,
				created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				
				CONSTRAINT fk_password_history_account 
					FOREIGN KEY (account_uuid) REFERENCES account(uuid) ON DELETE CASCADE
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		"record_summary": `
			CREATE TABLE trip_summary (
				data_ref         varchar(64) NOT NULL PRIMARY KEY COMMENT '行程唯一編號',
				device_id        varchar(64) NOT NULL COMMENT '裝置/寵物ID',
				start_time       datetime(3) NOT NULL COMMENT '開始時間',
				end_time         datetime(3) NOT NULL COMMENT '結束時間',
				duration_minutes double DEFAULT 0 COMMENT '總耗時(分鐘)',
				point_count      int DEFAULT 0 COMMENT 'GPS點數量',
				distance_km      decimal(10,3) DEFAULT 0.000 COMMENT '總距離(km)',
				created_at       datetime(3),
				updated_at       datetime(3)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}

	// 建立順序（外鍵依賴）
	createOrder := []string{
		"member",
		"device",
		"member_device",
		"account",
		"password_history",
	}

	// 計算開了多少TABLE
	newTable := 0

	for _, tableName := range createOrder {
		sqlStmt, ok := tables[tableName]
		if !ok {
			// logafa.Warn("Table 定義遺失: %s", tableName)
			continue
		}

		// 使用 GORM Raw 查詢 information_schema
		var count int
		err := db.WithContext(ctx).
			Raw(`
				SELECT COUNT(*) 
				FROM information_schema.tables 
				WHERE table_schema = DATABASE() 
				  AND table_name = ?
			`, tableName).
			Scan(&count).Error

		if err != nil {
			// logafa.Error("檢查 Table %s 是否存在失敗: %v", tableName, err)
			continue
		}

		if count > 0 {
			continue
		}

		// 建立 Table
		// logafa.Info("正在建立 Table `%s`...", tableName)
		if err := db.WithContext(ctx).Exec(sqlStmt).Error; err != nil {
			// logafa.Error("建立 Table `%s` 失敗: %v", tableName, err)
			continue
		}
		newTable++
	}
	if newTable == 0 {
		return
	}

	// logafa.Info("💾SQL Tables 初始化完成")
}

func initMongoIndexes(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Panic 處理
	defer func() {
		if r := recover(); r != nil {
			// logafa.Error("初始化Mongo Index 失敗 (panic): %v", r)
			panic(r)
		}
	}()

	collection := client.Database("pettrack").Collection("pettrack")

	// 定義索引：{ name, model }
	type namedIndex struct {
		Name  string
		Model mongo.IndexModel
	}

	indexesToEnsure := []namedIndex{
		{
			Name: "idx_location_2dsphere",
			Model: mongo.IndexModel{
				Keys:    bson.D{{Key: "location", Value: "2dsphere"}},
				Options: options.Index().SetName("idx_location_2dsphere"),
			},
		},
		{
			Name: "idx_device_recorded_unique",
			Model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "device_id", Value: 1},
					{Key: "recorded_at", Value: 1},
				},
				Options: options.Index().SetName("idx_device_recorded_unique").SetUnique(true),
			},
		},
		{
			Name: "idx_recorded_at_desc",
			Model: mongo.IndexModel{
				Keys:    bson.D{{Key: "recorded_at", Value: -1}},
				Options: options.Index().SetName("idx_recorded_at_desc"),
			},
		},
		{
			Name: "idx_device_id",
			Model: mongo.IndexModel{
				Keys:    bson.D{{Key: "device_id", Value: 1}},
				Options: options.Index().SetName("idx_device_id"),
			},
		},
		{
			Name: "idx_created_at_ttl",
			Model: mongo.IndexModel{
				Keys:    bson.D{{Key: "created_at", Value: 1}},
				Options: options.Index().SetName("idx_created_at_ttl").SetExpireAfterSeconds(30 * 24 * 60 * 60),
			},
		},
	}

	// 取得現有索引名稱
	existingNames := make(map[string]bool)
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		// logafa.Error("無法列出現有索引: %v", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&result); err == nil {
			existingNames[result.Name] = true
		}
	}
	if err := cursor.Err(); err != nil {
		// logafa.Error("遍歷索引時發生錯誤: %v", err)
		return
	}

	// 過濾出需要建立的索引
	var toCreate []mongo.IndexModel

	for _, idx := range indexesToEnsure {
		if !existingNames[idx.Name] {
			toCreate = append(toCreate, idx.Model)
		}
	}

	if len(toCreate) == 0 {
		return
	}

	// 建立索引
	_, err = collection.Indexes().CreateMany(ctx, toCreate)
	if err != nil {
		// logafa.Error("建立索引失敗: %v", err)
		return
	}
	// logafa.Info("MongoDB 索引初始化完成")
}

package domainRepo

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Uuid          uuid.UUID `gorm:"type:char(36);primaryKey" json:"uuid"`
	MemberId      int64     `gorm:"not null" json:"memberId"`
	Username      string    `gorm:"type:varchar(255);unique;not null" json:"username"`
	Password      string    `gorm:"type:varchar(255);not null" json:"password"`
	Email         string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Identity      string    `gorm:"type:varchar(50)" json:"identity"`
	LastLoginTime time.Time `gorm:"type:datetime" json:"lastLoginTime"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (a *Account) TableName() string {
	return "account"
}

type PasswordHistory struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountUuid uuid.UUID `gorm:"type:char(36);not null" json:"accountUuid"`
	Password    string    `gorm:"type:varchar(255);not null" json:"password"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (pp *PasswordHistory) TableName() string {
	return "password_history"
}

type Member struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	LastName  string    `gorm:"type:varchar(255)" json:"lastName"`
	FirstName string    `gorm:"type:varchar(255)" json:"firstName"`
	NickName  string    `gorm:"type:varchar(255)" json:"nickName"`
	Email     string    `gorm:"type:varchar(255)" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (m *Member) TableName() string {
	return "member"
}

type MemberDevice struct {
	Id         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberId   int64     `gorm:"not null" json:"memberId"`
	DeviceId   string    `gorm:"type:char(36);not null" json:"deviceId"`
	DeviceName string    `gorm:"type:varchar(255)" json:"deviceName"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (md *MemberDevice) TableName() string {
	return "member_device"
}

type Device struct {
	Uuid           uuid.UUID `gorm:"type:char(36);primaryKey" json:"uuid"`
	DeviceId       string    `gorm:"type:varchar(255)" json:"deviceId"`
	DeviceType     string    `gorm:"type:varchar(50)" json:"deviceType"`
	CreateByMember int64     `gorm:"not null" json:"memberId"`
	Remark         string    `gorm:"type:char(50)" json:"remark"`
}

func (d *Device) TableName() string {
	return "device"
}

type TripSummary struct {
	DataRef         string    `gorm:"column:data_ref;uniqueIndex:uk_data_ref;size:64;not null;comment:'行程唯一編號'" bson:"data_ref"`
	DeviceID        string    `gorm:"column:device_id;index:idx_device_date;size:64;not null;comment:'裝置/寵物ID'" bson:"device_id"`
	StartTime       time.Time `gorm:"column:start_time;index:idx_device_date;not null;comment:'開始時間'" bson:"start_time"`
	EndTime         time.Time `gorm:"column:end_time;not null;comment:'結束時間'" bson:"end_time"`
	DurationMinutes float64   `gorm:"column:duration_minutes;type:double;default:0;comment:'總耗時(分鐘)'" bson:"duration_minutes"`
	PointCount      int       `gorm:"column:point_count;type:int;default:0;comment:'GPS點數量'" bson:"point_count"`
	DistanceKM      float64   `gorm:"column:distance_km;type:decimal(10,3);default:0.000;index:idx_distance;comment:'總距離(km)'" bson:"distance_km"`

	CreatedAt 		time.Time `gorm:"column:created_at"`
	UpdatedAt 		time.Time `gorm:"column:updated_at"`
}

func (d *TripSummary) TableName() string {
	return "trip_summary"
}
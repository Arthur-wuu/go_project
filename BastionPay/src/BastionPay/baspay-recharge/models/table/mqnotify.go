package table

type (
	MqNotify struct {
		Id             *int64  `json:"id"                gorm:"primary_key;column:ID"`
		NotifyType     *int64  `json:"notify_type"      gorm:"column:NOTIFY_TYPE";type:tinyint(4)"`
		NotifyContent  *string `json:"notify_content"   gorm:"column:NOTIFY_CONTENT;type:varchar(1000)"`
		CreateTime     *string `json:"CREATE_TIME"          gorm:"column:CREATE_TIME;type:datetime"`
		LastUpdateTime *int64  `json:"last_update_time" gorm:"column:LAST_UPDATE_TIME;type:varchar(30)"`
	}
)

func (this *RederInfoT) TableName() string {
	return "MQ_NOTIFY"
}

type (
	Content struct {
		UId         *int64  `valid:"optional" json:"uid"`
		RegistTime  *int64  `valid:"optional" json:"regist_time, "`
		Country     *string `valid:"optional" json:"country"`
		Phone       *string `valid:"optional" json:"phone"`
		Channel     *string `valid:"optional" json:"channel,omitempty"`
		CountryCode *string `valid:"optional" json:"phone_district"`
	}
)

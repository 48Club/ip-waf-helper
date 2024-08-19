package types

type IPWaf struct {
	ID uint64 `gorm:"bigint;primaryKey;autoIncrement;column:id" json:"id"`
	IP string `gorm:"varchar(255);unique;column:ip" json:"ip" form:"ip" binding:"required"`
}

func (IPWaf) TableName() string {
	return "ip_waf"
}

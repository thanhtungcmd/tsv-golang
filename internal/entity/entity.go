package entity

type User struct {
	CompanyCode   string  `json:"company_code" gorm:"primaryKey"`
	ID            string  `json:"id" gorm:"primaryKey"`
	UserName      string  `json:"user_name"`
	Password      *string `json:"password"`
	Email         *string `json:"email"`
	Authority     *string `json:"authority"`
	PositionLevel *string `json:"position_level"`
	JobTitle      *string `json:"job_title"`
	UseYN         string  `json:"use_yn" gorm:"default:Y"`
	CreatedAt     *string `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
	CreatedUser   *string `json:"created_user"`
	UpdatedUser   *string `json:"updated_user"`
}

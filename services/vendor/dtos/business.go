package dtos

type RegisterBusiness struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Desc  string `json:"desc"`
}

type UpdateBusiness struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

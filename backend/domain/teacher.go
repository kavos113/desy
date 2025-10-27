package domain

type Teacher struct {
	ID int
	Name string
	Url string
}

type TeacherRepository interface {
	FindByID(id int) (*Teacher, error)	
	Create(teacher *Teacher) error
	Creates(teachers []Teacher) error
	Update(teacher *Teacher) error
	Delete(id int) error
}
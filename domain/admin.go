package domain

type AdminID string

//genconstructor
type Admin struct {
	id   AdminID `required:"" getter:""`
	name string  `required:"" getter:"" setter:"Rename"`
}

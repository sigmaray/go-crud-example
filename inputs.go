package main

type UserInput struct {
	Login    string `validate:"required,min=3"`
	Password string `validate:"required,min=3"`
}

type PageInput struct {
	Slug    string `validate:"required"`
	Content string `validate:"required"`
}

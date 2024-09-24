package test

import (
	"app-server/interfaces/controller"
	"app-server/utils/mock"
)

var (
	userApp    mock.UserAppInterface
	foodApp    mock.FoodAppInterface
	fakeUpload mock.UploadFileInterface
	fakeAuth   mock.AuthInterface
	fakeToken  mock.TokenInterface

	s  = controller.NewUsers(&userApp, &fakeAuth, &fakeToken)                       //We use all mocked data here
	f  = controller.NewFood(&foodApp, &userApp, &fakeUpload, &fakeAuth, &fakeToken) //We use all mocked data here
	au = controller.NewAuthenticate(&userApp, &fakeAuth, &fakeToken)                //We use all mocked data here

)

package model_test

import (
	"testing"

	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/domain/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestModel_NewUser(t *testing.T) {
	name := "Jonathan"
	email := "email@email.com"

	user, err := model.NewUser(name, email);

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(user.ID))
	require.Equal(t, user.Name, name)
	require.Equal(t, user.Email, email)

	_, err = model.NewUser("", "")
	require.NotNil(t, err)
}
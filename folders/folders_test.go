package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllFolders(t *testing.T) {
	t.Run("returns all folders when orgID matches", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		res, err := folders.GetAllFolders(
			&folders.FetchFolderRequest{
				OrgID: orgID,
			})
		assert.NoError(t, err)
		assert.NotNil(t, res)

		for _, folder := range res.Folders {
			assert.Equal(t, orgID, folder.OrgId)
		}
	})

	t.Run("returns no folders when orgID does not match", func(t *testing.T) {
		orgID := uuid.Must(uuid.NewV4())
		res, err := folders.GetAllFolders(
			&folders.FetchFolderRequest{
				OrgID: orgID,
			})
		assert.NoError(t, err)
		assert.Empty(t, res.Folders)
	})

	t.Run("returns error when orgID is nil", func(t *testing.T) {
		res, err := folders.GetAllFolders(
			&folders.FetchFolderRequest{
				OrgID: uuid.Nil,
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("returns error when request is nil", func(t *testing.T) {
		res, err := folders.GetAllFolders(nil)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

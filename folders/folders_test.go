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
}

func Test_GetAllFoldersPagination(t *testing.T) {
	t.Run("returns error when index is less than or equal to 0", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		res, err := folders.GetAllFoldersPaginated(&folders.FetchFoldersPaginatedRequest{
			OrgID: orgID,
			Limit: -1,
			Token: "",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("returns error when orgID is nil", func(t *testing.T) {
		res, err := folders.GetAllFoldersPaginated(&folders.FetchFoldersPaginatedRequest{
			OrgID: uuid.Nil,
			Limit: 5,
			Token: "",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("returns an error when given an invalid base64 token", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		res, err := folders.GetAllFoldersPaginated(&folders.FetchFoldersPaginatedRequest{
			OrgID: orgID,
			Limit: 5,
			Token: "token",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("returns no folders when orgID does not match", func(t *testing.T) {
		orgID := uuid.Must(uuid.NewV4())
		res, err := folders.GetAllFoldersPaginated(&folders.FetchFoldersPaginatedRequest{
			OrgID: orgID,
			Limit: 5,
			Token: "",
		})
		assert.NoError(t, err)
		assert.Empty(t, res.Folders)
	})

	t.Run("fetches first 5 folders", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		res, err := folders.GetAllFoldersPaginated(&folders.FetchFoldersPaginatedRequest{
			OrgID: orgID,
			Limit: 5,
			Token: "",
		})

		assert.NoError(t, err)

		expected, _ := folders.FetchAllFoldersByOrgID(orgID)
		assert.Equal(t, expected[0:5], res.Folders)
	})

	t.Run("fetches first 5 folders then the next 5 folders", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		first, _ := folders.GetAllFoldersPaginated(&folders.FetchFoldersPaginatedRequest{
			OrgID: orgID,
			Limit: 5,
			Token: "",
		})

		second, err := folders.GetAllFoldersPaginated(&folders.FetchFoldersPaginatedRequest{
			OrgID: orgID,
			Limit: 5,
			Token: first.NextToken,
		})

		expected, _ := folders.FetchAllFoldersByOrgID(orgID)
		assert.NoError(t, err)
		assert.Equal(t, expected[5:10], second.Folders)
	})

	t.Run("fetches near the end of the list of folders", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)

		expected, _ := folders.FetchAllFoldersByOrgID(orgID)
		nextToken := folders.EncodeNextIndex(len(expected) - 3)

		res, err := folders.GetAllFoldersPaginated(&folders.FetchFoldersPaginatedRequest{
			OrgID: orgID,
			Limit: 5,
			Token: nextToken,
		})

		assert.NoError(t, err)
		assert.Equal(t, len(expected[len(expected)-3:]), len(res.Folders))
		assert.Equal(t, expected[len(expected)-3:], res.Folders)
	})
}

package folders

import (
	"github.com/gofrs/uuid"
)

// Fetches all of the folders with the matching OrgID
// There are a few things we can do to refactor the code:
//   - Remove any instances of unused variable declarations
//   - There was a bug in the second for loop where the address of &v1
//     never changes, causing the same Folder item to be appended to the list each time
//   - We can remove both for loops as it is unecessary to instantiate multiple lists, since
//     the response from FetchAllFoldersByOrgID has already retrived all of the data
//   - We can also remove the ffr variable instantiated, and instead directly return the FetchFolderResponse
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	r, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}
	return &FetchFolderResponse{Folders: r}, nil
}

func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	folders := GetSampleData()

	resFolder := []*Folder{}
	for _, folder := range folders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}

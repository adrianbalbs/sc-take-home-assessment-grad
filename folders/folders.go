package folders

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// Fetches all of the folders with the matching OrgID
// There are a few things we can do to refactor the code:
//   - Remove any instances of unused variable declarations
//   - There was a bug in the second for loop where the address of &v1
//     never changes, causing the same Folder item to be appended to the list each time
//   - We can remove both for loops as it is unecessary to instantiate multiple lists and re-process
//     the data, since the response from FetchAllFoldersByOrgID has already retrived all of the data
//   - We can also remove the ffr variable instantiated
//     and instead, directly return the FetchFolderResponse
//   - Original variable names are not very descriptive and should be changed, e.g. "r" to "folders"
//   - Do error handling for when the request passed to the function is nil and deal with
//     any propagating errors from FetchAllFoldersByOrgID
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	folders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}
	return &FetchFolderResponse{Folders: folders}, nil
}

// Retrives all the folders and then filters out the folders which do not match
// the given orgID
//   - We can preallocate resFolder to have the size of folders if we know that the size
//     of the returned folders could be at most the size of the sample data retrieved.
//   - We should handle the case where the UUID passed into the function is potentially Nil
//     and return an error accordingly
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	if orgID == uuid.Nil {
		return nil, fmt.Errorf("invalid orgID: Nil UUID")
	}
	folders := GetSampleData()

	resFolder := make([]*Folder, 0, len(folders))
	for _, folder := range folders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}

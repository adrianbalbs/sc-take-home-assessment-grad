package folders

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Fetches all of the folders with the matching OrgID
// There are a few things we can do to refactor the code:
//
//   - Remove any instances of unused variable declarations. The first variable
//     initialised does not get used anywhere, makes the intent of our code unclear,
//     and prevents our code from compiling. Removing this variable will improve the
//     clarity of our code
//
//   - Variable names in the current code like f1, fs, f, r are not very descriptive and make it hard
//     to clearly reason about what our program is doing. We should be declaring variables with more
//     descriptive names.
//
//   - We can remove both for loops as the first for loop is just collecting and
//     dereferencing the addresses of the folders returned by FetchAllFoldersByOrgID, while
//     the second is collecting the addresses of those folders into another slice. This is
//     unecessary repeated work which is already accomplished by FetchAllFoldersByOrgID.
//
//   - Instead of declaring the ffr variable and assigning it's value on two lines, we can instead
//     use the := operator to instantiate and assign the variable at the same time.
//
//   - Functions like FetchAllFoldersByOrgID may return an error,
//     but we are discarding the error result returned which can lead to potential bugs in our code.
//     We should be handling errors by checking if an error exists and return early to propogate the error.
//
//   - Additionally we could have a request passed in that is potentially nil,
//     so we should also handle this case accordingly.
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	if req == nil {
		return nil, errors.New("invalid request: request cannot be nil")
	}
	folders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}
	return &FetchFolderResponse{Folders: folders}, nil
}

// Retrives all the folders and then filters out the folders which do not match
// the given orgID
//
//   - We can preallocate resFolder to have the size of folders if we know that the size
//     of the returned folders could be at most the size of the sample data retrieved.
//
//   - We should handle the case where the UUID passed into the function is potentially Nil
//     and return an error accordingly
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	if orgID == uuid.Nil {
		return nil, errors.New("invalid orgID: Nil UUID")
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

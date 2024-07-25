package folders

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
)

// Explanation of Solution:
//
// I decided to go with a cursor-based approach for implementing pagination. The reason I went
// with this implementation is that compared to offset-pagination, we can directly index on
// the next set of data to be fetched rather than linearly search for the offset that
// we want to start on.
//
// A user sends a request payload with the orgID, a limit for how many folders are to
// be fetched (which should be an integer greater than 0), and a base64 encoded cursor token,
// which is either empty or contains the encoded starting index of the next set of folders. Calling
// the FetchFoldersPaginatedRequest function returns a slice of the Folders from the starting index up
// to the ending offset index (which is just start + limit), and the encoded cursor pointing to the starting
// index of the next set of folders to be retrieved. If the cursor token passed into the function is empty,
// then we default to searching from the beginning of the slice. The function is repeatedly called
// until we reach the end of the slice, where the cursor returned is an empty string.

type FetchFoldersPaginatedRequest struct {
	Cursor string
	Limit  int
	OrgID  uuid.UUID
}

type FetchFoldersPaginatedResponse struct {
	NextCursor string
	Folders    []*Folder
}

func GetAllFoldersPaginated(req *FetchFoldersPaginatedRequest) (*FetchFoldersPaginatedResponse, error) {
	if req == nil {
		return nil, errors.New("invalid request: request cannot be nil")
	}

	if req.Limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	start := 0
	if req.Cursor != "" {
		index, err := DecodeNextCursor(req.Cursor)
		if err != nil {
			return nil, err
		}
		start = index
	}

	folders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	end := start + req.Limit
	if end > len(folders) {
		end = len(folders)
	}

	nextCursor := ""
	if end != len(folders) {
		nextCursor = EncodeNextCursor(end)
	}

	return &FetchFoldersPaginatedResponse{Folders: folders[start:end], NextCursor: nextCursor}, nil
}

func EncodeNextCursor(endIndex int) string {
	return base64.StdEncoding.EncodeToString([]byte("next_cursor:" + strconv.Itoa(endIndex)))
}

func DecodeNextCursor(encodedCursor string) (int, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return 0, err
	}

	index, err := strconv.Atoi(strings.Split(string(decodedCursor), ":")[1])
	if err != nil {
		return 0, err
	}

	return index, nil
}

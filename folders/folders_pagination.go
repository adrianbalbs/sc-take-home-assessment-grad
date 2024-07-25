package folders

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
)

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
	if req.Limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
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

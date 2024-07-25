package folders

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
)

type FetchFoldersPaginatedRequest struct {
	Token string
	Limit int
	OrgID uuid.UUID
}

type FetchFoldersPaginatedResponse struct {
	NextToken string
	Folders   []*Folder
}

func GetAllFoldersPaginated(req *FetchFoldersPaginatedRequest) (*FetchFoldersPaginatedResponse, error) {
	if req.Limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	start := 0
	if req.Token != "" {
		index, err := DecodeNextIndex(req.Token)
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

	nextToken := ""
	if end != len(folders) {
		nextToken = EncodeNextIndex(end)
	}

	return &FetchFoldersPaginatedResponse{Folders: folders[start:end], NextToken: nextToken}, nil
}

func EncodeNextIndex(nextIndex int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(nextIndex)))
}

func DecodeNextIndex(token string) (int, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}
	index, err := strconv.Atoi(string(decodedToken))
	if err != nil {
		return 0, err
	}

	return index, nil
}

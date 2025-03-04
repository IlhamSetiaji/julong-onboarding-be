package response

import (
	"github.com/google/uuid"
)

type SendFindOrganizationByIDMessageResponse struct {
	OrganizationID       string `json:"organization_id"`
	Name                 string `json:"name"`
	OrganizationCategory string `json:"organization_category"`
	OrganizationType     string `json:"organization_type"`
}

type SendFindOrganizationLocationByIDMessageResponse struct {
	OrganizationLocationID string `json:"organization_location_id"`
	Name                   string `json:"name"`
}

type SendFindOrganizationStructureByIDMessageResponse struct {
	OrganizationStructureID string `json:"organization_structure_id"`
	Name                    string `json:"name"`
}

type OrganizationResponse struct {
	ID                 uuid.UUID `json:"id"`
	OrganizationTypeID uuid.UUID `json:"organization_type_id"`
	Name               string    `json:"name"`
}

type OrganizationLocationResponse struct {
	ID               uuid.UUID `json:"id"`
	OrganizationID   uuid.UUID `json:"organization_id"`
	OrganizationName string    `json:"organization_name"`
	Name             string    `json:"name"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}

type OrganizationLocationPaginatedResponse struct {
	OrganizationLocations []OrganizationLocationResponse `json:"organization_locations"`
	Total                 int64                          `json:"total"`
	TotalNull             int64                          `json:"total_null"`
}

package handlers

import (
	"time"
)

// type Collection struct {
// 	ID          int       `json:"id"`
// 	Name        string    `json:"name"`
// 	Count       int       `json:"count"`
// 	ResourceURL string    `json:"resource_url"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// Release represents a vinyl release in the collection
type Release struct {
	ID          int       `json:"id"           db:"id"` // Release ID from Discogs
	InstanceID  int       `json:"instance_id"  db:"instance_id"`
	FolderID    int       `json:"folder_id"    db:"folder_id"`
	Rating      int       `json:"rating"       db:"rating"`
	Title       string    `json:"title"        db:"title"`
	Year        *int      `json:"year"         db:"year"` // Pointer since it can be null
	ResourceURL string    `json:"resource_url" db:"resource_url"`
	Thumb       string    `json:"thumb"        db:"thumb"`
	CoverImage  string    `json:"cover_image"  db:"cover_image"`
	CreatedAt   time.Time `json:"created_at"   db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"   db:"updated_at"`
	LastSynced  time.Time `json:"last_synced"  db:"last_synced"`

	// Related entities (for JSON marshaling/unmarshaling)
	Labels  []ReleaseLabel  `json:"labels,omitempty"`
	Artists []ReleaseArtist `json:"artists,omitempty"`
	Formats []Format        `json:"formats,omitempty"`
	Genres  []Genre         `json:"genres,omitempty"`
	Styles  []Style         `json:"styles,omitempty"`
	Notes   []ReleaseNote   `json:"notes,omitempty"`
}

// Label represents a record label
type Label struct {
	ID          int    `json:"id"           db:"id"` // Label ID from Discogs
	Name        string `json:"name"         db:"name"`
	ResourceURL string `json:"resource_url" db:"resource_url"`
	EntityType  string `json:"entity_type"  db:"entity_type"`
}

// ReleaseLabel represents the many-to-many relationship between releases and labels
type ReleaseLabel struct {
	ReleaseID int    `json:"release_id" db:"release_id"`
	LabelID   int    `json:"label_id"   db:"label_id"`
	CatNo     string `json:"catno"      db:"catno"` // Catalog number

	// For JSON marshaling/unmarshaling (embedded label info)
	Label *Label `json:"label,omitempty" db:"-"`
}

// Artist represents a music artist
type Artist struct {
	ID          int    `json:"id"           db:"id"` // Artist ID from Discogs
	Name        string `json:"name"         db:"name"`
	ResourceURL string `json:"resource_url" db:"resource_url"`
}

// ReleaseArtist represents the many-to-many relationship between releases and artists
type ReleaseArtist struct {
	ReleaseID    int    `json:"release_id"    db:"release_id"`
	ArtistID     int    `json:"artist_id"     db:"artist_id"`
	JoinRelation string `json:"join_relation" db:"join_relation"`
	ANV          string `json:"anv"           db:"anv"` // Artist name variation
	Tracks       string `json:"tracks"        db:"tracks"`
	Role         string `json:"role"          db:"role"`

	// For JSON marshaling/unmarshaling (embedded artist info)
	Artist *Artist `json:"artist,omitempty" db:"-"`
}

// Format represents the format of a release (e.g., "Vinyl", "CD")
type Format struct {
	ID        int    `json:"id"         db:"id"`
	ReleaseID int    `json:"release_id" db:"release_id"`
	Name      string `json:"name"       db:"name"` // e.g., "Vinyl", "CD"
	Qty       int    `json:"qty"        db:"qty"`

	// For JSON marshaling/unmarshaling
	Descriptions []string `json:"descriptions,omitempty" db:"-"`
}

// FormatDescription represents a description of a format
type FormatDescription struct {
	FormatID    int    `json:"format_id"   db:"format_id"`
	Description string `json:"description" db:"description"`
}

// Genre represents a music genre
type Genre struct {
	ID   int    `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
}

// ReleaseGenre represents the many-to-many relationship between releases and genres
type ReleaseGenre struct {
	ReleaseID int `json:"release_id" db:"release_id"`
	GenreID   int `json:"genre_id"   db:"genre_id"`
}

// Style represents a music style (sub-genre)
type Style struct {
	ID   int    `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
}

// ReleaseStyle represents the many-to-many relationship between releases and styles
type ReleaseStyle struct {
	ReleaseID int `json:"release_id" db:"release_id"`
	StyleID   int `json:"style_id"   db:"style_id"`
}

// ReleaseNote represents a note about a release
type ReleaseNote struct {
	ReleaseID int    `json:"release_id" db:"release_id"`
	FieldID   int    `json:"field_id"   db:"field_id"`
	Value     string `json:"value"      db:"value"`
}

// Folder represents a collection folder from Discogs
type Folder struct {
	ID          int       `json:"id"           db:"id"`
	Name        string    `json:"name"         db:"name"`
	Count       int       `json:"count"        db:"count"`
	ResourceURL string    `json:"resource_url" db:"resource_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastSynced  time.Time `json:"last_synced"`
}

// FoldersResponse represents the response from the Discogs API folders endpoint
type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}

// DiscogsResponse represents the paginated response from the Discogs API
type DiscogsResponse struct {
	Pagination struct {
		PerPage int `json:"per_page"`
		Pages   int `json:"pages"`
		Page    int `json:"page"`
		Items   int `json:"items"`
		URLs    struct {
			Next string `json:"next"`
			Last string `json:"last"`
		} `json:"urls"`
	} `json:"pagination"`
	Releases []DiscogsRelease `json:"releases"`
}

type DiscogsRelease struct {
	ID         int `json:"id"`
	InstanceID int `json:"instance_id"`
	FolderID   int `json:"folder_id"`
	Rating     int `json:"rating"`
	BasicInfo  struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Year        int    `json:"year"`
		ResourceURL string `json:"resource_url"`
		Thumb       string `json:"thumb"`
		CoverImage  string `json:"cover_image"`
		Formats     []struct {
			Qty          string   `json:"qty"`
			Descriptions []string `json:"descriptions"`
			Name         string   `json:"name"`
		} `json:"formats"`
		Labels []struct {
			ResourceURL string `json:"resource_url"`
			EntityType  string `json:"entity_type"`
			CatNo       string `json:"catno"`
			ID          int    `json:"id"`
			Name        string `json:"name"`
		} `json:"labels"`
		Artists []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Join        string `json:"join"`
			ResourceURL string `json:"resource_url"`
			ANV         string `json:"anv"`
			Tracks      string `json:"tracks"`
			Role        string `json:"role"`
		} `json:"artists"`
		Genres []string `json:"genres"`
		Styles []string `json:"styles"`
	} `json:"basic_information"`
	Notes []struct {
		FieldID int    `json:"field_id"`
		Value   string `json:"value"`
	} `json:"notes"`
}

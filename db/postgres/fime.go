package postgres

// UserStore Interface
type UserStore interface {
	User(id int64) (User, error)
	UserByEmail(email string) (User, error)
	Users(args ListParams) ([]User, error)
	CreateUser(args CreateUserParams) (User, error)
	UpdateUser(args UpdateUserParams) (User, error)
	DeleteUser(id int64) error
}

// ImageStore Interface
type ImageStore interface {
	Image(id int64) (Image, error)
	Images(args ListParams) ([]Image, error)
	ImagesByUser(args ListUserImagesParams) ([]Image, error)
	CreateImage(args CreateImageParams) (Image, error)
	UpdateImage(args UpdateImageParams) (Image, error)
	DeleteImage(id int64) error
}

// TagStore Interface
type TagStore interface {
	Tag(id int64) (Tag, error)
	Tags(args ListParams) ([]Tag, error)
	CreateTag(args CreateTagParams) (Tag, error)
	GetUserTags(arg ListUserTagsParams) ([]Tag, error)
	DeleteTag(id int64) error
}

// ImageTagStore Interface
type ImageTagStore interface {
	CreateImageTag(it ImageTag) error
	DeleteImageTag(it ImageTag) error
	DeleteAllImageTags(imgID int64) error
	GetTagsByImageID(imgID int64) ([]Tag, error)
	GetImagesByTagID(tagID int64) ([]Image, error)
}

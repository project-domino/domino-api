package errors

// 400 errors
var (
	BadParameters     = &Error{400, "Bad Parameters"}
	MissingParameters = &Error{400, "Missing Parameters"}

	InvalidAuthHeader = &Error{400, "Invalid Authorization header"}

	PasswordsDoNotMatch = &Error{400, "Passwords do not match"}
	UserExists          = &Error{400, "User already exists"}

	TagExists    = &Error{400, "Tag already exists"}
	InvalidPage  = &Error{400, "Page number is not valid"}
	InvalidItems = &Error{400, "Item count is not valid"}
)

// 401 errors
var (
	AuthRequired       = &Error{401, "Authorization required"}
	InvalidCredentials = &Error{401, "Invalid credentials"}
)

// 403 errors
var (
	UnknownAuthMethod = &Error{403, "Unknown Authorization method"}
	NoPermission      = &Error{403, "You do not have permission to access this resource"}

	NotNoteOwner       = &Error{403, "You are not the owner of this note"}
	NotCollectionOwner = &Error{403, "You are not the owner of this collection"}
	NotTextbookOwner   = &Error{403, "You are not the owner of this textbook"}
)

// 404 errors
var (
	NotFound           = &Error{404, "Page not found"}
	NoteNotFound       = &Error{404, "Note not found"}
	CollectionNotFound = &Error{404, "Collection not found"}
	UserNotFound       = &Error{404, "User not found"}
)

// 5xx errors
var (
	InternalError = &Error{500, "Server Error"}
	Debug         = &Error{500, "teh internets are asplode"}
	JSON          = &Error{500, "Could not convert to JSON"}
	UnknownValue  = &Error{500, "Unknown internal value"}
)

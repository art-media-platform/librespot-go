package respot

import (
	"github.com/art-media-platform/amp.SDK/stdlib/media"
	"github.com/art-media-platform/amp.SDK/stdlib/task"
	"github.com/art-media-platform/librespot-go/librespot/asset"
	"github.com/art-media-platform/librespot-go/librespot/core/crypto"
	"github.com/art-media-platform/librespot-go/librespot/mercury"
)

// Forward declartion to create a respot.Session.
var StartNewSession func(ctx *SessionContext) (Session, error)

// DefaultSessionContext creates a SessionContext with the given device label.
func DefaultSessionContext(deviceLabel string) *SessionContext {
	ctx := &SessionContext{
		DeviceName: deviceLabel,
	}
	return ctx
}

type SessionContext struct {
	task.Context              // logging & shutdown
	Login        SessionLogin // means for the session to login
	Info         SessionInfo  // filled in during Session.Login()
	Keys         crypto.Keys  // If nil, will be auto-generated
	DeviceName   string       // Label of the device being used
	DeviceUID    string       // if nil, auto-generated from DeviceName
}

type SessionLogin struct {
	Username  string
	Password  string // AUTHENTICATION_USER_PASS
	AuthData  []byte // AUTHENTICATION_STORED_SPOTIFY_CREDENTIALS
	AuthToken string // AUTHENTICATION_SPOTIFY_TOKEN
}

type SessionInfo struct {
	Username string // authenticated canonical username
	AuthBlob []byte // reusable authentication blob for Spotify Connect devices
	Country  string // user country returned by Spotify
}

type Session interface {
	Close() error

	// Returns the SessionContext current in use by this session
	Context() *SessionContext

	// Initiates login with params contained in Ctx.Login
	Login() error

	Search(query string, limit int) (*mercury.SearchResponse, error)
	Mercury() *mercury.Client

	// Initiates downloading ("pinning") with the given spotify track ID or URI.
	// The track data is accessible via media.Asset, supporting io.ReadSeekCloser and is essential for streaming or serving.
	PinTrack(trackID string, opts PinOpts) (media.Asset, error)
}

type PinOpts struct {

	// If set, media.Asset.OnStart(Ctx().Context) will be called on the returned media.Asset.
	// This is for convenience but not desirable when the asset is in a time-to-live cache, for example.
	StartInternally bool

	// Expresses preferences for the requested asset's codec and format.
	Format asset.AssetFormat
}

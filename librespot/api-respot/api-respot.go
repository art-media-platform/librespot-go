package respot

import (
	"github.com/arcspace/go-arcspace/arc/assets"
	"github.com/arcspace/go-cedar/process"
	"github.com/arcspace/go-librespot/librespot/core/crypto"
	"github.com/arcspace/go-librespot/librespot/mercury"
)

func DefaultSessionCtx(deviceLabel string) *SessionCtx {
	ctx := &SessionCtx{
		DeviceName: deviceLabel,
	}
	return ctx
}

type SessionLogin struct {
	Username   string
	Password   string
	AuthData   []byte
	OAuthToken string
}

type SessionCtx struct {
	process.Context              // logging & shutdown
	Info            SessionInfo  // filled in during Session.Login()
	Login           SessionLogin // means for the session to login
	Keys            crypto.Keys  // If left nil, will be auto-generated
	DeviceName      string       // Label of the device being used
	DeviceUID       string       // if nil, auto-generated from DeviceName
}

type SessionInfo struct {
	Username string //  authenticated canonical username
	AuthBlob []byte // reusable authentication blob for Spotify Connect devices
	Country  string // user country returned by Spotify
}

type PinOpts struct {

	// If set, MediaAsset.OnStart(Ctx().Context) will be called on the returned MediaAsset.
	// This is for convenience but not desirable when the asset is in a time-to-live cache, for example.
	StartInternally bool
}

type Session interface {
	// Returns the SessionCtx current in use by this session
	Ctx() *SessionCtx

	// Initiates login with params contained in Ctx.Login
	Login() error

	Search(query string, limit int) (*mercury.SearchResponse, error)
	Mercury() *mercury.Client

	// Initiates access ("pinning") with the given spotify track ID or URI
	PinTrack(trackID string, opts PinOpts) (assets.MediaAsset, error)
}

// Forward declared method to create a new Spotify session
var StartNewSession func(ctx *SessionCtx) (Session, error)

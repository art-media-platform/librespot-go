package core

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/arcspace/go-librespot/Spotify"
	respot "github.com/arcspace/go-librespot/librespot/api-respot"
	"github.com/arcspace/go-librespot/librespot/core/connection"
	"github.com/arcspace/go-librespot/librespot/utils"
)

var Version = "master"
var BuildID = "dev"

func (sess *Session) Login() error {
	var packet []byte
	login := &sess.ctx.Login

	if len(login.AuthData) > 0 && login.Username != "" {
		packet = sess.makeLoginBlobPacket(
			login.Username,
			login.AuthData,
			Spotify.AuthenticationType_AUTHENTICATION_STORED_SPOTIFY_CREDENTIALS.Enum(),
		)
	} else if login.OAuthToken != "" {
		packet = sess.makeLoginBlobPacket(
			"",
			[]byte(sess.ctx.Login.OAuthToken),
			Spotify.AuthenticationType_AUTHENTICATION_SPOTIFY_TOKEN.Enum(),
		)
	} else if login.Username != "" && login.Password != "" {
		packet = sess.makeLoginBlobPacket(
			login.Username,
			[]byte(login.Password),
			Spotify.AuthenticationType_AUTHENTICATION_USER_PASS.Enum(),
		)
	} else {
		return errors.New("no login method provided")
	}

	return sess.doLogin(packet, sess.ctx.Login.Username)
}

func (s *Session) doLogin(packet []byte, username string) error {
	s.ctx.Info = respot.SessionInfo{}

	err := s.stream.SendPacket(connection.PacketLogin, packet)
	if err != nil {
		log.Fatal("bad shannon write", err)
	}

	// Pll once for authentication response
	welcome, err := s.handleLogin()
	if err != nil {
		return err
	}

	// Store the few interesting values
	user := welcome.GetCanonicalUsername()
	if user == "" {
		user = s.ctx.Login.Username
	}
	s.ctx.Info.Username = user
	s.ctx.Info.AuthBlob = welcome.GetReusableAuthCredentials()

	// Poll for acknowledge before loading - needed for gopherjs
	// s.poll()
	go s.runPollLoop() // TODO: add context.Context exit!

	return nil
}

func (s *Session) handleLogin() (*Spotify.APWelcome, error) {
	cmd, data, err := s.stream.RecvPacket()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	if cmd == connection.PacketAuthFailure {
		errCode := Spotify.ErrorCode(data[1])
		return nil, fmt.Errorf("authentication failed: %v", errCode)
	} else if cmd == connection.PacketAPWelcome {
		welcome := &Spotify.APWelcome{}
		err := proto.Unmarshal(data, welcome)
		if err != nil {
			return nil, fmt.Errorf("authentication failed: %v", err)
		}
		return welcome, nil
	} else {
		return nil, fmt.Errorf("authentication failed: unexpected cmd %v", cmd)
	}
}

func (s *Session) getLoginBlobPacket(blob utils.BlobInfo) ([]byte, error) {
	data, _ := base64.StdEncoding.DecodeString(blob.DecodedBlob)
	buffer := bytes.NewBuffer(data)
	if _, err := buffer.ReadByte(); err != nil {
		return nil, fmt.Errorf("could not read byte: %+v", err)
	}
	_, err := readBytes(buffer)
	if err != nil {
		return nil, fmt.Errorf("could not read bytes: %+v", err)
	}
	if _, err := buffer.ReadByte(); err != nil {
		return nil, fmt.Errorf("could not read byte: %+v", err)
	}
	authNum := readInt(buffer)
	authType := Spotify.AuthenticationType(authNum)
	if _, err := buffer.ReadByte(); err != nil {
		return nil, fmt.Errorf("could not read byte: %+v", err)
	}
	authData, err := readBytes(buffer)
	if err != nil {
		return nil, fmt.Errorf("could not read bytes: %+v", err)
	}
	return s.makeLoginBlobPacket(blob.Username, authData, &authType), nil
}

func (s *Session) makeLoginBlobPacket(
	username string,
	authData []byte,
	authType *Spotify.AuthenticationType,
) []byte {
	versionString := "librespot_" + Version + "_" + BuildID
	packet := &Spotify.ClientResponseEncrypted{
		LoginCredentials: &Spotify.LoginCredentials{
			Username: proto.String(username),
			Typ:      authType,
			AuthData: authData,
		},
		SystemInfo: &Spotify.SystemInfo{
			CpuFamily:               Spotify.CpuFamily_CPU_UNKNOWN.Enum(),
			Os:                      Spotify.Os_OS_UNKNOWN.Enum(),
			SystemInformationString: proto.String("librespot-golang"),
			DeviceId:                proto.String(s.ctx.DeviceUID),
		},
		VersionString: proto.String(versionString),
	}
	buf, _ := proto.Marshal(packet)
	return buf
}

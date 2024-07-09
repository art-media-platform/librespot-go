package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"io"

	"github.com/amp-3d/amp-librespot-go/Spotify"
	respot "github.com/amp-3d/amp-librespot-go/librespot/api-respot"
	_ "github.com/amp-3d/amp-librespot-go/librespot/core" // bootstrapping
	"github.com/amp-3d/amp-librespot-go/librespot/core/oauth"
	"github.com/amp-3d/amp-sdk-go/stdlib/task"
)

const (
	// The device name that is registered to Spotify servers
	defaultDeviceName = "librespot"
)

func main() {
	err := mainStart()
	if err != nil {
		log.Fatal(err)
	}
}

func mainStart() error {

	// Read flags from commandline
	username := flag.String("username", "", "spotify username")
	password := flag.String("password", "", "spotify password")
	blobPath := flag.String("blob", "", "spotify auth blob")
	devicename := flag.String("devicename", defaultDeviceName, "name of device")
	flag.Parse()

	ctx := respot.DefaultSessionContext(*devicename)
	ctx.Context, _ = task.Start(&task.Task{
		Info: task.Info{
			Label: "main",
		},
	})

	sess, err := respot.StartNewSession(ctx)
	if err != nil {
		return err
	}

	{
		login := &ctx.Login
		login.Username = *username
		login.Password = *password
		login.AuthToken = ""

		// Authenticate reusing an existing blob
		if len(*blobPath) > 0 {
			login.AuthData, err = os.ReadFile(*blobPath)
			if err != nil {
				return fmt.Errorf("unable to read auth blob from %s: %s", *blobPath, err)
			}
		}

		if login.Password == "" && login.AuthToken == "" {
			var err error
			login.AuthToken, err = oauth.LoginOAuth(os.Getenv("client_id"), os.Getenv("client_secret"), os.Getenv("redirect_uri"))
			if err != nil {
				return err
			}
		}

		err = sess.Login()
		if err != nil {
			return err
		}
	}

	// Command loop
	reader := bufio.NewReader(os.Stdin)

	printHelp()

	for {
		var err error

		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		cmds := strings.Split(strings.TrimSpace(text), " ")

		switch cmds[0] {
		case "help":
			printHelp()

		case "track":
			if len(cmds) < 2 {
				fmt.Println("You must specify the Base62 Spotify ID of the track")
			} else {
				funcTrack(sess, cmds[1])
			}

		case "artist":
			if len(cmds) < 2 {
				fmt.Println("You must specify the Base62 Spotify ID of the artist")
			} else {
				funcArtist(sess, cmds[1])
			}

		case "album":
			if len(cmds) < 2 {
				fmt.Println("You must specify the Base62 Spotify ID of the album")
			} else {
				funcAlbum(sess, cmds[1])
			}

		case "playlists":
			funcPlaylists(sess)

		case "search":
			funcSearch(sess, cmds[1])

		case "save":
			if len(cmds) < 2 {
				err = errors.New("missing Base62 Spotify ID of the track")
			} else {
				err = funcSave(sess, cmds[1])
			}

		default:
			fmt.Println("Unknown command")
		}

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("save <track>:                   downloads the specified track by spotify base62 id")
	fmt.Println("track <track>:                  show details on specified track by spotify base62 id")
	fmt.Println("album <album>:                  show details on specified album by spotify base62 id")
	fmt.Println("artist <artist>:                show details on specified artist by spotify base62 id")
	fmt.Println("search <keyword>:               start a search on the specified keyword")
	fmt.Println("playlists:                      show your playlists")
	fmt.Println("help:                           show this help")
}

func funcTrack(session respot.Session, trackURI string) {
	fmt.Println("Loading track: ", trackURI)

	_, track, err := session.Mercury().GetTrack(trackURI)
	if err != nil {
		fmt.Printf("Error loading track %q - %v", trackURI, err)
		return
	}

	fmt.Println("Track title: ", track.GetName())
}

func funcArtist(session respot.Session, artistURI string) {
	_, artist, err := session.Mercury().GetArtist(artistURI)
	if err != nil {
		fmt.Printf("Error loading artist %q - %v", artistURI, err)
		return
	}

	fmt.Printf("Artist: %s\n", artist.GetName())
	fmt.Printf("Popularity: %.0f\n", artist.GetPopularity())
	fmt.Printf("Genre: %s\n", artist.GetGenre())

	if artist.GetTopTrack() != nil && len(artist.GetTopTrack()) > 0 {
		// Spotify returns top tracks in multiple countries. We take the first
		// one as example, but we should use the country data returned by the
		// Spotify server (session.Country())
		tt := artist.GetTopTrack()[0]
		fmt.Printf("\nTop tracks (country %s):\n", tt.GetCountry())

		for _, t := range tt.GetTrack() {
			// To save bandwidth, only track IDs are returned. If you want
			// the track name, you need to fetch it.
			fmt.Printf(" => %s\n", Spotify.ConvertTo62(t.GetGid()))
		}
	}

	fmt.Printf("\nAlbums:\n")
	for _, ag := range artist.GetAlbumGroup() {
		for _, a := range ag.GetAlbum() {
			fmt.Printf(" => %s\n", Spotify.ConvertTo62(a.GetGid()))
		}
	}

}

func funcAlbum(session respot.Session, albumURI string) {
	_, album, err := session.Mercury().GetAlbum(albumURI)
	if err != nil {
		fmt.Printf("Error loading album %q - %v", albumURI, err)
		return
	}

	fmt.Printf("Album: %s\n", album.GetName())
	fmt.Printf("Popularity: %.0f\n", album.GetPopularity())
	fmt.Printf("Genre: %s\n", album.GetGenre())
	fmt.Printf("Date: %d-%d-%d\n", album.GetDate().GetYear(), album.GetDate().GetMonth(), album.GetDate().GetDay())
	fmt.Printf("Label: %s\n", album.GetLabel())
	fmt.Printf("Type: %s\n", album.GetTyp())

	fmt.Printf("Artists: ")
	for _, artist := range album.GetArtist() {
		fmt.Printf("%s ", Spotify.ConvertTo62(artist.GetGid()))
	}
	fmt.Printf("\n")

	for _, disc := range album.GetDisc() {
		fmt.Printf("\nDisc %d (%s): \n", disc.GetNumber(), disc.GetName())

		for _, track := range disc.GetTrack() {
			fmt.Printf(" => %s\n", Spotify.ConvertTo62(track.GetGid()))
		}
	}

}

func funcPlaylists(session respot.Session) {
	fmt.Println("Listing playlists")

	playlist, err := session.Mercury().GetRootPlaylist(session.Context().Info.Username)

	if err != nil || playlist.Contents == nil {
		fmt.Println("Error getting root list: ", err)
		return
	}

	items := playlist.Contents.Items
	for i := 0; i < len(items); i++ {
		id := strings.TrimPrefix(items[i].GetUri(), "spotify:")
		id = strings.Replace(id, ":", "/", -1)
		list, _ := session.Mercury().GetPlaylist(id)
		fmt.Println(list.Attributes.GetName(), id)

		if list.Contents != nil {
			for j := 0; j < len(list.Contents.Items); j++ {
				item := list.Contents.Items[j]
				fmt.Println(" ==> ", *item.Uri)
			}
		}
	}
}

func funcSearch(session respot.Session, keyword string) {
	resp, err := session.Search(keyword, 12)

	if err != nil {
		fmt.Println("Failed to search:", err)
		return
	}

	res := resp.Results

	fmt.Println("Search results for ", keyword)
	fmt.Println("=============================")

	if res.Error != nil {
		fmt.Println("Search result error:", res.Error)
	}

	fmt.Printf("Albums: %d (total %d)\n", len(res.Albums.Hits), res.Albums.Total)

	for _, album := range res.Albums.Hits {
		fmt.Printf(" => %s (%s)\n", album.Name, album.Uri)
	}

	fmt.Printf("\nArtists: %d (total %d)\n", len(res.Artists.Hits), res.Artists.Total)

	for _, artist := range res.Artists.Hits {
		fmt.Printf(" => %s (%s)\n", artist.Name, artist.Uri)
	}

	fmt.Printf("\nTracks: %d (total %d)\n", len(res.Tracks.Hits), res.Tracks.Total)

	for _, track := range res.Tracks.Hits {
		fmt.Printf(" => %s (%s)\n", track.Name, track.Uri)
	}
}

func funcSave(sess respot.Session, trackID string) error {
	fmt.Println("Loading track for play: ", trackID)

	asset, err := sess.PinTrack(trackID, respot.PinOpts{
		StartInternally: true,
	})
	if err != nil {
		return fmt.Errorf("failed to pin track: %s", err)
	}
	r, err := asset.NewAssetReader()
	if err != nil {
		return err
	}
	defer r.Close()

	buffer, err := io.ReadAll(io.Reader(r))
	if err != nil {
		return err
	}

	err = os.WriteFile(asset.Label(), buffer, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}
	return nil
}

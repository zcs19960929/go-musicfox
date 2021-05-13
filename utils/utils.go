package utils

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "github.com/buger/jsonparser"
    "go-musicfox/constants"
    "go-musicfox/ds"
    "os"
    "os/exec"
    "os/user"
    "runtime"
    "strings"
)

// GetLocalDataDir 获取本地数据存储目录
func GetLocalDataDir() string {
    // Home目录
    homeDir, err := Home()
    if nil != err {
        panic("未获取到用户Home目录: " + err.Error())
    }

    projectDir := fmt.Sprintf("%s/%s", homeDir, constants.AppLocalDataDir)

    if _, err := os.Stat(projectDir); os.IsNotExist(err) {
        _ = os.Mkdir(projectDir, os.ModePerm)
    }

    return projectDir
}

// Home 获取当前用户的Home目录
func Home() (string, error) {
    curUser, err := user.Current()
    if nil == err {
        return curUser.HomeDir, nil
    }

    // cross compile support
    if "windows" == runtime.GOOS {
        return homeWindows()
    }

    // Unix-like system, so just assume Unix
    return homeUnix()
}

func homeUnix() (string, error) {
    // First prefer the HOME environmental variable
    if home := os.Getenv("HOME"); home != "" {
        return home, nil
    }

    // If that fails, try the shell
    var stdout bytes.Buffer
    cmd := exec.Command("sh", "-c", "eval echo ~$USER")
    cmd.Stdout = &stdout
    if err := cmd.Run(); err != nil {
        return "", err
    }

    result := strings.TrimSpace(stdout.String())
    if result == "" {
        return "", errors.New("blank output when reading home directory")
    }

    return result, nil
}

func homeWindows() (string, error) {
    drive := os.Getenv("HOMEDRIVE")
    path := os.Getenv("HOMEPATH")
    home := drive + path
    if drive == "" || path == "" {
        home = os.Getenv("USERPROFILE")
    }
    if home == "" {
        return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
    }

    return home, nil
}

// IDToBin convert autoincrement ID to []byte
func IDToBin(ID uint64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, ID)
    return b
}

// BinToID convert []byte to autoincrement ID
func BinToID(bin []byte) uint64 {
    ID := binary.BigEndian.Uint64(bin)

    return ID
}

type ResCode uint8
const (
    Success ResCode = iota
    UnknownError
    NetworkError
    NeedLogin
    PasswordError
)

// CheckCode 验证响应码
func CheckCode(code float64) ResCode {
    switch code {
    case 301, 302:
        return NeedLogin
    case 520:
        return NetworkError
    case 200:
        return Success
    }

    return PasswordError
}

// CheckUserInfo 验证用户信息
func CheckUserInfo(user *ds.User) ResCode {
    if user == nil || user.UserId == 0 {
        return NeedLogin
    }

    return Success
}

// ReplaceSpecialStr 替换特殊字符
func ReplaceSpecialStr(str string) string {
    replaceStr := map[string]string{
        "“": "\"",
        "”": "\"",
        "·": ".",
    }
    for oldStr, newStr := range replaceStr {
        str = strings.ReplaceAll(str, oldStr, newStr)
    }

    return str
}

// GetDailySongs 获取每日歌曲列表
func GetDailySongs(data []byte) (list []ds.Song) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
       if song, err := ds.NewSongFromDailySongsJson(value); err == nil {
           list = append(list, song)
       }

    }, "data", "dailySongs")

    return
}

// GetDailyPlaylists 获取播放列表
func GetDailyPlaylists(data []byte) (list []ds.Playlist) {

    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if playlist, err := ds.NewPlaylistFromJson(value); err == nil {
            list = append(list, playlist)
        }
    }, "recommend")

    return
}

// GetSongsOfPlaylist 获取播放列表的歌曲
func GetSongsOfPlaylist(data []byte) (list []ds.Song) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if song, err := ds.NewSongFromPlaylistSongsJson(value); err == nil {
            list = append(list, song)
        }

    }, "playlist", "tracks")

    return
}

// GetSongsOfAlbum 获取专辑的歌曲
func GetSongsOfAlbum(data []byte) (list []ds.Song) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if song, err := ds.NewSongFromAlbumSongsJson(value); err == nil {
            list = append(list, song)
        }

    }, "songs")

    return
}

// GetPlaylists 获取播放列表
func GetPlaylists(data []byte) (list []ds.Playlist) {

    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if playlist, err := ds.NewPlaylistFromJson(value); err == nil {
            list = append(list, playlist)
        }
    }, "playlist")

    return
}

// GetFmSongs 获取每日歌曲列表
func GetFmSongs(data []byte) (list []ds.Song) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if song, err := ds.NewSongFromFmJson(value); err == nil {
            list = append(list, song)
        }

    }, "data")

    return
}

// GetIntelligenceSongs 获取心动模式歌曲列表
func GetIntelligenceSongs(data []byte) (list []ds.Song) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if song, err := ds.NewSongFromIntelligenceJson(value); err == nil {
            list = append(list, song)
        }

    }, "data")

    return
}

// GetNewAlbums 获取最新专辑列表
func GetNewAlbums(data []byte) (albums []ds.Album) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

        if album, err := ds.NewAlbumFromAlbumJson(value); err == nil {
            albums = append(albums, album)
        }

    }, "albums")

    return
}

// GetTopAlbums 获取专辑列表
func GetTopAlbums(data []byte) (albums []ds.Album) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

        if album, err := ds.NewAlbumFromAlbumJson(value); err == nil {
            albums = append(albums, album)
        }

    }, "monthData")

    return
}

// GetArtistHotAlbums 获取歌手热门专辑列表
func GetArtistHotAlbums(data []byte) (albums []ds.Album) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

        if album, err := ds.NewAlbumFromAlbumJson(value); err == nil {
            albums = append(albums, album)
        }

    }, "hotAlbums")

    return
}

// GetSongsOfSearchResult 获取搜索结果的歌曲
func GetSongsOfSearchResult(data []byte) (list []ds.Song) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if song, err := ds.NewSongFromSearchResultJson(value); err == nil {
            list = append(list, song)
        }

    }, "result", "songs")

    return
}


// GetAlbumsOfSearchResult 获取搜索结果的专辑
func GetAlbumsOfSearchResult(data []byte) (list []ds.Album) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if album, err := ds.NewAlbumFromAlbumJson(value); err == nil {
            list = append(list, album)
        }

    }, "result", "albums")

    return
}

// GetPlaylistsOfSearchResult 获取搜索结果的歌单
func GetPlaylistsOfSearchResult(data []byte) (list []ds.Playlist) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if playlist, err := ds.NewPlaylistFromJson(value); err == nil {
            list = append(list, playlist)
        }

    }, "result", "playlists")

    return
}

// GetArtistsOfSearchResult 获取搜索结果的歌手
func GetArtistsOfSearchResult(data []byte) (list []ds.Artist) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if artist, err := ds.NewArtist(value); err == nil {
            list = append(list, artist)
        }

    }, "result", "artists")

    return
}

// GetSongsOfArtist 获取歌手的歌曲
func GetSongsOfArtist(data []byte) (list []ds.Song) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if song, err := ds.NewSongFromArtistSongsJson(value); err == nil {
            list = append(list, song)
        }

    }, "songs")

    return
}

// GetUsersOfSearchResult 从搜索结果中获取用户列表
func GetUsersOfSearchResult(data []byte) (list []ds.User) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if song, err := ds.NewUserFromSearchResultJson(value); err == nil {
            list = append(list, song)
        }

    }, "result", "userprofiles")

    return
}

// GetDjRadiosOfSearchResult 从搜索结果中获取电台列表
func GetDjRadiosOfSearchResult(data []byte) (list []ds.DjRadio) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if radio, err := ds.NewDjRadioFromJson(value); err == nil {
            list = append(list, radio)
        }

    }, "result", "djRadios")

    return
}

// GetSongsOfDjRadio 获取电台节目列表的歌曲
func GetSongsOfDjRadio(data []byte) (list []ds.Song) {
    _, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if song, err := ds.NewSongFromDjRadioProgramJson(value); err == nil {
           list = append(list, song)
        }
    }, "programs")

    return
}
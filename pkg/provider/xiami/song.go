package xiami

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) GetSong(songId string) (*provider.Song, error) {
	resp, err := a.GetSongDetailRaw(songId)
	if err != nil {
		return nil, err
	}

	_song := &resp.Data.Data.SongDetail
	a.patchSongLyric(_song)
	songs := resolve(_song)
	return songs[0], nil
}

// 获取歌曲详情
func (a *API) GetSongDetailRaw(songId string) (*SongDetailResponse, error) {
	token, err := a.getToken(APIGetSongDetail)
	if err != nil {
		return nil, err
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(songId)
	if err != nil {
		model["songStringId"] = songId
	} else {
		model["songId"] = songId
	}
	params := sreq.Params(signPayload(token, model))
	resp := new(SongDetailResponse)
	err = a.Request(sreq.MethodGet, APIGetSongDetail, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get song detail: %w", err)
	}

	return resp, nil
}

func (a *API) GetSongLyric(songId string) (string, error) {
	resp, err := a.GetSongLyricRaw(songId)
	if err != nil {
		return "", err
	}

	for _, i := range resp.Data.Data.Lyrics {
		if i.FlagOfficial == "1" && i.Type == "2" {
			return i.Content, nil
		}
	}

	return "", errors.New("official lyric not present")
}

// 获取歌词
func (a *API) GetSongLyricRaw(songId string) (*SongLyricResponse, error) {
	token, err := a.getToken(APIGetSongLyric)
	if err != nil {
		panic(err)
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(songId)
	if err != nil {
		model["songStringId"] = songId
	} else {
		model["songId"] = songId
	}
	params := signPayload(token, model)
	resp := new(SongLyricResponse)
	err = a.Request(sreq.MethodGet, APIGetSongLyric, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get song lyric: %w", err)
	}

	return resp, nil
}

// 批量获取歌曲，上限200首
func (a *API) GetSongsRaw(songIds ...string) (*SongsResponse, error) {
	token, err := a.getToken(APIGetSongs)
	if err != nil {
		return nil, err
	}

	if len(songIds) > SongRequestLimit {
		songIds = songIds[:SongRequestLimit]
	}
	model := map[string][]string{
		"songIds": songIds,
	}
	params := signPayload(token, model)
	resp := new(SongsResponse)
	err = a.Request(sreq.MethodGet, APIGetSongs, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get songs: %w", err)
	}

	return resp, nil
}

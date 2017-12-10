package provider

import (
	"github.com/ahmdrz/goinsta"
	"fmt"
)

type (
	Instagram interface {
		Login() error
		Logout() error
		LastPhotos(profile string, total int) ([]string, error)
		Follow(profile string) (error)
		Stories(profile string) ([]string, error)
	}

	instagram struct {
		provider *goinsta.Instagram
	}
)

func NewInstagram(login, password string) (Instagram, error) {
	i := &instagram{
		provider: goinsta.New(login, password),
	}
	return i, i.Login()
}

func (i *instagram) Login() error {
	return i.provider.Login()
}

func (i *instagram) Logout() error {
	return i.provider.Logout()
}

// TODO improve this
func (i *instagram) Stories(profile string) ([]string, error) {
	userResp, err := i.provider.GetUserByUsername(profile)
	if err != nil {
		return nil, err
	}

	fmt.Println(userResp)
	tray, err := i.provider.GetUserStories(userResp.User.ID)
	if err != nil {
		return nil, err
	}

	fmt.Println(tray)
	fmt.Println(tray.Media)
	var final []string
	for _, val := range tray.Media {
		imageURL := ""
		max := 0

		fmt.Println(val.ImageVersions2.Candidates)
		for _,image := range val.ImageVersions2.Candidates {
			if image.Height > max {
				max = image.Height
				imageURL = image.URL
			}
		}
		final = append(final, imageURL)

		videoURL := ""
		max = 0
		fmt.Println(val.VideoVersions)
		for _,video := range val.VideoVersions {
			if video.Height > max {
				max = video.Height
				videoURL = video.URL
			}
		}
		final = append(final, videoURL)
	}

	return final, nil
}

func (i *instagram) Follow(profile string) (error) {
	userResp, err := i.provider.GetUserByUsername(profile)
	if err != nil {
		return err
	}

	_, err = i.provider.Follow(userResp.User.ID)
	return err
}

func (i *instagram) LastPhotos(profile string, total int) ([]string, error) {
	user, err := i.provider.GetUserByUsername(profile)
	if err != nil {
		return nil, err
	}

	urlList, feedList := []string{}, []string{}
	maxID := ""
	hasMore := true
	for err == nil && len(urlList) < total && hasMore {
		feedList, maxID, hasMore, err = i.userFeed(user.User.ID, maxID)
		urlList = append(urlList, feedList...)
	}

	if len(urlList) > total {
		return urlList[:total], err
	}

	return urlList, err
}

func (i *instagram) userFeed(userID int64, maxID string) ([]string, string, bool, error) {
	userFeedResponse, err := i.provider.UserFeed(userID, maxID, "")
	if err != nil {
		return nil, "", false, err
	}

	urlList := []string{}

	for _, item := range userFeedResponse.Items {
		bestIndex := 0
		bestWidht := 0

		if len(item.ImageVersions2.Candidates) == 0 {
			continue
		}

		for index, candidate := range item.ImageVersions2.Candidates {
			if candidate.Width > bestWidht {
				bestIndex = index
				bestWidht = candidate.Width
			}
		}

		urlList = append(urlList, item.ImageVersions2.Candidates[bestIndex].URL)
	}

	return urlList, userFeedResponse.NextMaxID, userFeedResponse.MoreAvailable, nil
}
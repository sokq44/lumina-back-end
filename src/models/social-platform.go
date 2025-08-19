package models

import (
	"encoding/json"
	"strings"
)

type SocialType uint8

const (
	UnknownST   SocialType = 0
	GithubST    SocialType = 1
	TwitterST   SocialType = 2
	FacebookST  SocialType = 3
	LinkedInST  SocialType = 4
	InstagramST SocialType = 5
)

type UserSocialPlatform struct {
	UserId string
	Value  string
	Label  string
	Type   SocialType
}

type SocialPlatform struct {
	Type SocialType
	Name string
}

func (t *SocialType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = SocialPlatformGetType(s)
	return nil
}

func SocialPlatformGetName(t SocialType) string {
	switch t {
	case GithubST:
		return "Github"
	case TwitterST:
		return "Twitter"
	case FacebookST:
		return "Facebook"
	case LinkedInST:
		return "LinkedIn"
	case InstagramST:
		return "Instagram"
	default:
		return "Unknown"
	}
}

func SocialPlatformGetType(name string) SocialType {
	switch strings.ToLower(name) {
	case "github":
		return GithubST
	case "twitter":
		return TwitterST
	case "facebook":
		return FacebookST
	case "linkedin":
		return LinkedInST
	case "instagram":
		return InstagramST
	default:
		return UnknownST
	}
}

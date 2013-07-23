// Copyright 2013 The go-twitterstream AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twitterstream

type Tweet struct {
	Contributors         *TweetContributors  `json:"contributors,omitempty"`
	Coordinates          *TweetCoordinate    `json:"coordinates,omitempty"`
	CreatedAt            string              `json:"created_at,omitempty"`
	CurrentUserRetweet   *CurrentUserRetweet `json:"current_user_retweet,omitempty"`
	Entities             *TweetEntities      `json:"entities,omitempty"`
	FavoriteCount        int64               `json:"favorite_count,omitempty"`
	Favorited            bool                `json:"favorited,omitempty"`
	FilterLevel          string              `json:"filter_level,omitempty"`
	ID                   int64               `json:"id,omitempty"`
	IDStr                string              `json:"id_str,omitempty"`
	InReplyToScreenName  string              `json:"in_reply_to_screen_name,omitempty"`
	InReplyToStatusID    int64               `json:"in_reply_to_status_id,omitempty"`
	InReplyToStatusIDStr string              `json:"in_reply_to_status_id_str,omitempty"`
	InReplyToUserID      int64               `json:"in_reply_to_user_id,omitempty"`
	InReplyToUserIDStr   string              `json:"in_reply_to_user_id_str,omitempty"`
	Lang                 string              `json:"lang,omitempty"`
	Place                *Place              `json:"place,omitempty"`
	PossiblySensitive    bool                `json:"possibly_sensitive,omitempty"`
	Scopes               map[string]bool     `json:"scopes,omitempty"`
	RetweetCount         int                 `json:"retweet_count,omitempty"`
	Retweeted            bool                `json:"retweeted,omitempty"`
	Source               string              `json:"source,omitempty"`
	Text                 string              `json:"text,omitempty"`
	Truncated            bool                `json:"truncated,omitempty"`
	User                 *User               `json:"user,omitempty"`
	WithheldCopyright    bool                `json:"withheld_copyright,omitempty"`
	WithheldInCountries  []string            `json:"withheld_in_countries,omitempty"`
	WithheldScope        string              `json:"withheld_scope,omitempty"`
}

type TweetContributors struct {
	ID         int64  `json:"id,omitempty"`
	IDStr      string `json:"id_str,omitempty"`
	ScreenName string `json:"screen_name,omitempty"`
}

type TweetCoordinate struct {
	Coordinates [2]float64 `json:"coordinates,omitempty"`
	Type        string     `json:"type,omitempty"`
}

type CurrentUserRetweet struct {
	ID    int64  `json:"id,omitempty"`
	IDStr string `json:"id_str,omitempty"`
}

type Place struct {
	Attributes  *PlaceAttributes `json:"attributes,omitempty"`
	BoundingBox *BoundingBox     `json:"bounding_box,omitempty"`
	Country     string           `json:"country,omitempty"`
	CountryCode string           `json:"country_code,omitempty"`
	FullName    string           `json:"full_name,omitempty"`
	ID          string           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	PlaceType   string           `json:"place_type,omitempty"`
	URL         string           `json:"url,omitempty"`
}

type PlaceAttributes struct {
	StreetAddress string `json:"street_address,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Region        string `json:"region,omitempty"`
	ISO3          string `json:"iso3,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Twitter       string `json:"twitter,omitempty"`
	URL           string `json:"url,omitempty"`
}

type ll [2]float64
type llList []ll

type BoundingBox struct {
	Coordinates []llList `json:"coordinates,omitempty"`
	Type        string   `json:"type,omitempty"`
}

type TweetEntities struct {
	Hashtags     []HastagEntity      `json:"hashtags,omitempty"`
	URLs         []URLEntity         `json:"urls,omitempty"`
	UserMentions []UserMentionEntity `json:"user_mentions,omitempty"`
}

type UserMentionEntity struct {
	ID         int64    `json:"id,omitempty"`
	IDStr      string   `json:"id_str,omitempty"`
	Indices    [2]int64 `json:"indices,omitempty"`
	Name       string   `json:"name,omitempty"`
	ScreenName string   `json:"screen_name,omitempty"`
}

type HastagEntity struct {
	Indices [2]int64 `json:"indices,omitempty"`
	Text    string   `json:"text,omitempty"`
}

type MediaEntity struct {
	DisplayURL        string      `json:"display_url,omitempty"`
	ExpandedURL       string      `json:"expanded_url,omitempty"`
	ID                int64       `json:"id,omitempty"`
	IDStr             string      `json:"id_str,omitempty"`
	Indices           [2]int64    `json:"indices,omitempty"`
	MediaURL          string      `json:"media_url,omitempty"`
	MediaURLHTTPS     string      `json:"media_url_https,omitempty"`
	Sizes             *MediaSizes `json:"sizes,omitempty"`
	SourceStatusID    int64       `json:"source_status_id,omitempty"`
	SourceStatusIDStr string      `json:"source_status_id_str,omitempty"`
	Type              string      `json:"type,omitempty"`
	URL               string      `json:"url,omitempty"`
}

type MediaSizes struct {
	Thumb  *MediaSize `json:"thumb,omitempty"`
	Large  *MediaSize `json:"large,omitempty"`
	Medium *MediaSize `json:"medium,omitempty"`
	Small  *MediaSize `json:"small,omitempty"`
}

type MediaSize struct {
	Height int    `json:"h,omitempty"`
	Width  int    `json:"w,omitempty"`
	Resize string `json:"resize,omitempty"`
}

type URLEntity struct {
	DisplayURL  string   `json:"display_url,omitempty"`
	ExpandedURL string   `json:"expanded_url,omitempty"`
	Indices     [2]int64 `json:"indices,omitempty"`
	URL         string   `json:"url,omitempty"`
}

type User struct {
	ContributorsEnabled            bool     `json:"contributors_enabled,omitempty"`
	CreatedAt                      string   `json:"created_at,omitempty"`
	DefaultProfile                 bool     `json:"default_profile,omitempty"`
	DefaultProfileImage            bool     `json:"default_profile_image,omitempty"`
	Description                    string   `json:"description,omitempty"`
	FavouritesCount                int      `json:"favourites,omitempty"`
	FollowRequestSent              bool     `json:"follow_request_sent,omitempty"`
	Following                      bool     `json:"following,omitempty"`
	FollowersCount                 int      `json:"followers_count,omitempty"`
	FriendsCount                   int      `json:"friends_count,omitempty"`
	GeoEnabled                     bool     `json:"geo_enabled,omitempty"`
	ID                             int64    `json:"id,omitempty"`
	IDStr                          string   `json:"id_str,omitempty"`
	IsTranslator                   bool     `json:"is_translator,omitempty"`
	Lang                           string   `json:"lang,omitempty"`
	ListedCount                    int      `json:"listed_count,omitempty"`
	Location                       string   `json:"location,omitempty"`
	Name                           string   `json:"name,omitempty"`
	ProfileBackgroundColor         string   `json:"profile_background_color,omitempty"`
	ProfileBackgroundImageURL      string   `json:"profile_background_image_url,omitempty"`
	ProfileBackgroundImageURLHTTPS string   `json:"profile_background_image_url_https,omitempty"`
	ProfileBackgroundTile          bool     `json:"profile_background_tile,omitempty"`
	ProfileBannerURL               string   `json:"profile_banner_url,omitempty"`
	ProfileImageURL                string   `json:"profile_image_url,omitempty"`
	ProfileImageURLHTTPS           string   `json:"profile_image_url_https,omitempty"`
	ProfileLinkColor               string   `json:"profile_link_color,omitempty"`
	ProfileSidebarBorderColor      string   `json:"profile_sidebar_border_color,omitempty"`
	ProfileSidebarFillColor        string   `json:"profile_sidebar_fill_color,omitempty"`
	ProfileTextColor               string   `json:"profile_text_color,omitempty"`
	ProfileUseBackgroundImage      bool     `json:"profile_use_background_image,omitempty"`
	Protected                      bool     `json:"protected,omitempty"`
	ScreenName                     string   `json:"screen_name,omitempty"`
	ShowAllInlineMedia             bool     `json:"show_all_inline_media,omitempty"`
	Status                         *Tweet   `json:"status,omitempty"`
	StatusesCount                  int      `json:"statuses_count,omitempty"`
	TimeZone                       string   `json:"time_zone,omitempty"`
	URL                            string   `json:"url,omitempty"`
	UTCOffset                      int      `json:"utc_offset,omitempty"`
	Verified                       bool     `json:"verified,omitempty"`
	WithheldInCountries            []string `json:"withheld_in_countries,omitempty"`
	WithheldScope                  string   `json:"withheld_scope,omitempty"`
}

type DirectMessageNotice struct {
	DirectMessage *DirectMessage `json:"direct_message,omitempty"`
}

type DirectMessage struct {
	CreatedAt          string         `json:"created_at,omitempty"`
	Entities           *TweetEntities `json:"entities,omitempty"`
	ID                 int64          `json:"id,omitempty"`
	IDStr              string         `json:"id_str,omitempty"`
	Recipent           *User          `json:"recipent,omitempty"`
	RecipentID         int64          `json:"recipent_id,omitempty"`
	RecipentScreenName string         `json:"recipent_screen_name,omitempty"`
	Sender             *User          `json:"sender,omitempty"`
	SenderID           int64          `json:"sender_id,omitempty"`
	SenderScreenName   string         `json:"sender_screen_name,omitempty"`
	Text               string         `json:"text,omitempty"`
}

type TweetDeletionNotice struct {
	Delete *TweetDeletionNoticeStatus `json:"delete,omitempty"`
}

type TweetDeletionNoticeStatus struct {
	Status *DeletedStatus `json:"status,omitempty"`
}

type DeletedStatus struct {
	ID        int64  `json:"id,omitempty"`
	IDStr     string `json:"id_str,omitempty"`
	UserID    int64  `json:"user_id,omitempty"`
	UserIDStr string `json:"user_id_str,omitempty"`
}

type LocationDeletionNotice struct {
	ScrubGeo *ScrubGeo `json:"scrub_geo,omitempty"`
}

type ScrubGeo struct {
	UserID          int64  `json:"user_id,omitempty"`
	UserIDStr       string `json:"user_id_str,omitempty"`
	UpToStatusID    int64  `json:"up_to_status_id,omitempty"`
	UpToStatusIDStr string `json:"up_to_status_id_str,omitempty"`
}

type LimitNotice struct {
	Limit *Limit `json:"limit,omitempty"`
}

type Limit struct {
	Track int64 `json:"track,omitempty"`
}

type StatusWithheldNotice struct {
	StatusWithheld *StatusWithheld `json:"status_withheld,omitempty"`
}

type StatusWithheld struct {
	ID                  int64    `json:"id,omitempty"`
	UserID              int64    `json:"user_id,omitempty"`
	WithheldInCountries []string `json:"withheld_in_countries,omitempty"`
}

type UserWithheldNotice struct {
	UserWithheld *UserWithheld `json:"user_withheld,omitempty"`
}

type UserWithheld struct {
	ID                  int64    `json:"id,omitempty"`
	WithheldInCountries []string `json:"withheld_in_countries,omitempty"`
}

type WarningNotice struct {
	Warning *Warning `json:"warning,omitempty"`
}

type Warning struct {
	Code        string  `json:"code,omitempty"`
	Message     string  `json:"message,omitempty"`
	PercentFull float64 `json:"percent_full,omitempty"`
}

type DisconnectNotice struct {
	Disconnect *Disconnect `json:"disconnect,omitempty"`
}

type Disconnect struct {
	Code       int    `json:"code,omitempty"`
	StreamName string `json:"stream_name,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

func (d *Disconnect) ReasonByCode() string {
	if d.Code >= 1 && d.Code <= 12 {
		return DisconnectCode[d.Code]
	}
	return d.Reason
}

var DisconnectCode = map[int]string{
	1:  "Shutdown: The feed was shutdown (possibly a machine restart)",
	2:  "Duplicate stream: The same endpoint was connected too many times.",
	3:  "Control request: Control streams was used to close a stream (applies to sitestreams).",
	4:  "Stall: The client was reading too slowly and was disconnected by the server.",
	5:  "Normal: The client appeared to have initiated a disconnect.",
	6:  "Token revoked: An oauth token was revoked for a user (applies to site and userstreams).",
	7:  "Admin logout: The same credentials were used to connect a new stream and the oldest was disconnected.",
	8:  "Reserved for internal use. Will not be delivered to external clients.",
	9:  "Max message limit: The stream connected with a negative count parameter and was disconnected after all backfill was delivered.",
	10: "Stream exception: An internal issue disconnected the stream.",
	11: "Broker stall: An internal issue disconnected the stream.",
	12: "Shed load: The host the stream was connected to became overloaded and streams were disconnected to balance load. Reconnect as usual.",
}

type Event struct {
	Target       *User                  `json:"target,omitempty"`
	Source       *User                  `json:"source,omitempty"`
	Event        string                 `json:"event,omitempty"`
	TargetObject map[string]interface{} `json:"target_object,omitempty"`
	CreatedAt    string                 `json:"created_at,omitempty"`
}

type FriendsLists struct {
	Friends []int64 `json:"friends,omitempty"`
}

type TooManyFollow struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	UserID  int64  `json:"user_id,omitempty"`
}

type ForUser struct {
	ForUser string        `json:"for_user,omitempty"`
	Message *FriendsLists `json:"message,omitempty"`
}

type ControlNotice struct {
	Control *Control `json:"control,omitempty"`
}

type Control struct {
	ControlURI string `json:"control_uri,omitempty"`
}

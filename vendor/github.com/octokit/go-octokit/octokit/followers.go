package octokit

// URL templates for actions taken on the followers of users
//
// https://developer.github.com/v3/users/followers/
var (
	CurrentFollowerUrl  = Hyperlink("user/followers")
	FollowerUrl         = Hyperlink("users/{user}/followers")
	CurrentFollowingUrl = Hyperlink("user/following{/target}")
	FollowingUrl        = Hyperlink("users/{user}/following{/target}")
)

// Create a FollowersService
//
// https://developer.github.com/v3/users/followers/
func (c *Client) Followers() (followers *FollowersService) {
	followers = &FollowersService{client: c}
	return
}

// A service to return user followers
type FollowersService struct {
	client *Client
}

// Get a list of followers for the user
//
// https://developer.github.com/v3/users/followers/#list-followers-of-a-user
func (f *FollowersService) All(uri *Hyperlink, uriParams M) (followers []User, result *Result) {
	url, err := ExpandWithDefault(uri, &CurrentFollowerUrl, uriParams)
	if err != nil {
		return nil, &Result{Err: err}
	}

	result = f.client.get(url, &followers)
	return
}

// Checks if you are following a target user
//
// https://developer.github.com/v3/users/followers/#check-if-you-are-following-a-user
func (f *FollowersService) Check(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &CurrentFollowingUrl, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = f.client.get(url, nil)
	success = (result.Response.StatusCode == 204)
	return
}

// Follows a target user
//
// https://developer.github.com/v3/users/followers/#follow-a-user
func (f *FollowersService) Follow(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &CurrentFollowingUrl, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = f.client.put(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

// Unfollows a target user
//
// https://developer.github.com/v3/users/followers/#unfollow-a-user
func (f *FollowersService) Unfollow(uri *Hyperlink, uriParams M) (success bool, result *Result) {
	url, err := ExpandWithDefault(uri, &CurrentFollowingUrl, uriParams)
	if err != nil {
		return false, &Result{Err: err}
	}

	result = f.client.delete(url, nil, nil)
	success = (result.Response.StatusCode == 204)
	return
}

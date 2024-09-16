package response

type Final struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Errors  any    `json:"errors"`
	Status  Status `json:"status"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type FetchFacebookProfile struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Picture Picture `json:"picture"`
}

type Picture struct {
	Data PictureData `json:"data"`
}

type PictureData struct {
	Height       int    `json:"height"`
	IsSilhouette bool   `json:"is_silhouette"`
	URL          string `json:"url"`
	Width        int    `json:"width"`
}

type FetchGoogleProfile struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type OAuth struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
	IsRegistered bool   `json:"is_registered"`
}

type FinishRegister struct {
	//
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"`
}

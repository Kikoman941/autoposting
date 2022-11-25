package social_account

type SocialAccount struct {
	tableName  struct{} `pg:"public.social_accounts"`
	ID         int      `pg:"id,pk"`
	Network    string   `pg:"network"`
	Credential string   `pg:"credentials"`
}

type Group struct {
	tableName struct{}  `pg:"public.groups"`
	AccountID int       `pg:"account_id,fk"`
	Project   string    `pg:"project"`
	GroupInfo GroupInfo `pg:"group_info"`
}

type GroupInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

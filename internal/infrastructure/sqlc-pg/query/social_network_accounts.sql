-- name: CreateSocialNetworkAccount :one
INSERT INTO autoposting.social_network_accounts (social_network, credentials, access_token)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateSocialNetworkAccount :exec
UPDATE autoposting.social_network_accounts
SET social_network = $2, credentials = $3, access_token = $3
WHERE id = $1;

-- name: GetSocialNetworkAccounts :many
SELECT * FROM autoposting.social_network_accounts
WHERE (CASE WHEN @in_social_network_any_of::bool THEN social_network = ANY (@social_network_any_of::text[]) ELSE TRUE END);
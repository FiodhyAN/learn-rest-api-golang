-- name: GetUser :one
SELECT _id, username, email, password, email_verified, email_verification_token_expires_at 
FROM users 
WHERE email = $1 OR username = $1;

-- name: CreateUser :one
INSERT INTO users (name, username, email, password) 
VALUES ($1, $2, $3, $4) 
RETURNING _id, name, username, email, password, role, created_at;

-- name: UpdateUserVerificationExpired :exec
UPDATE users 
SET email_verification_token = $1, email_verification_token_expires_at = $2 
WHERE _id = $3;

-- name: GetUserById :one
SELECT _id, email_verified, email_verification_token, email_verification_token_expires_at 
FROM users 
WHERE _id = $1;

-- name: VerifyEmail :exec
UPDATE users 
SET email_verified = true, email_verification_token = NULL, email_verification_token_expires_at = NULL 
WHERE _id = $1;

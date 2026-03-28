package db

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                     string
	Email                  string
	Password               string
	GoogleID               string
	GoogleEmail            string
	GoogleName             string
	GooglePicture          string
	XAccessToken           string
	XRefreshToken          string
	XProfilePicture        string
	XUserID                string
	XName                  string
	XUsername              string
	LinkedInAccessToken    string
	LinkedInRefreshToken   string
	LinkedInProfilePicture string
	LinkedInUserID         string
	LinkedInName           string
}

const userSelectColumns = `
	SELECT
		id,
		email,
		COALESCE(password, ''),
		COALESCE(google_id, ''),
		COALESCE(google_email, ''),
		COALESCE(google_name, ''),
		COALESCE(google_picture, ''),
		COALESCE(x_access_token, ''),
		COALESCE(x_refresh_token, ''),
		COALESCE(x_profile_picture, ''),
		COALESCE(x_user_id, ''),
		COALESCE(x_name, ''),
		COALESCE(x_username, ''),
		COALESCE(linkedin_access_token, ''),
		COALESCE(linkedin_refresh_token, ''),
		COALESCE(linkedin_profile_picture, ''),
		COALESCE(linkedin_user_id, ''),
		COALESCE(linkedin_name, '')
	FROM app.users
`

func scanUser(scanner interface {
	Scan(dest ...any) error
}) (*User, error) {
	var user User

	err := scanner.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.GoogleID,
		&user.GoogleEmail,
		&user.GoogleName,
		&user.GooglePicture,
		&user.XAccessToken,
		&user.XRefreshToken,
		&user.XProfilePicture,
		&user.XUserID,
		&user.XName,
		&user.XUsername,
		&user.LinkedInAccessToken,
		&user.LinkedInRefreshToken,
		&user.LinkedInProfilePicture,
		&user.LinkedInUserID,
		&user.LinkedInName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}

		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	return scanUser(DB.QueryRow(userSelectColumns+` WHERE email = $1`, email))
}

func GetUserByGoogleID(googleID string) (*User, error) {
	return scanUser(DB.QueryRow(userSelectColumns+` WHERE google_id = $1`, googleID))
}

func GetUserByID(id string) (*User, error) {
	return scanUser(DB.QueryRow(userSelectColumns+` WHERE id = $1`, id))
}

func UpdateUserTwitterInfo(userID, accessToken, refreshToken, profilePicture, xUserID, xName, xUsername string) error {
	_, err := DB.Exec(
		`UPDATE app.users
		 SET x_access_token = $1,
		     x_refresh_token = $2,
		     x_profile_picture = $3,
		     x_user_id = $4,
		     x_name = $5,
		     x_username = $6
		 WHERE id = $7`,
		accessToken,
		refreshToken,
		profilePicture,
		xUserID,
		xName,
		xUsername,
		userID,
	)
	return err
}

func UpdateUserLinkedInInfo(userID, accessToken, refreshToken, profilePicture, linkedInUserID, linkedInName string) error {
	_, err := DB.Exec(
		`UPDATE app.users
		 SET linkedin_access_token = $1,
		     linkedin_refresh_token = $2,
		     linkedin_profile_picture = $3,
		     linkedin_user_id = $4,
		     linkedin_name = $5
		 WHERE id = $6`,
		accessToken,
		refreshToken,
		profilePicture,
		linkedInUserID,
		linkedInName,
		userID,
	)
	return err
}

func ClearUserTwitterInfo(userID string) error {
	_, err := DB.Exec(
		`UPDATE app.users
		 SET x_access_token = '',
		     x_refresh_token = '',
		     x_profile_picture = '',
		     x_user_id = '',
		     x_name = '',
		     x_username = ''
		 WHERE id = $1`,
		userID,
	)
	return err
}

func ClearUserLinkedInInfo(userID string) error {
	_, err := DB.Exec(
		`UPDATE app.users
		 SET linkedin_access_token = '',
		     linkedin_refresh_token = '',
		     linkedin_profile_picture = '',
		     linkedin_user_id = '',
		     linkedin_name = ''
		 WHERE id = $1`,
		userID,
	)
	return err
}

func UpsertGoogleUser(googleID, email, name, picture string) (*User, error) {
	user, err := GetUserByGoogleID(googleID)
	if err == nil {
		_, updateErr := DB.Exec(
			`UPDATE app.users
			 SET email = $1, google_email = $1, google_name = $2, google_picture = $3
			 WHERE id = $4`,
			email,
			name,
			picture,
			user.ID,
		)
		if updateErr != nil {
			return nil, updateErr
		}

		return GetUserByID(user.ID)
	}

	user, err = GetUserByEmail(email)
	if err == nil {
		_, updateErr := DB.Exec(
			`UPDATE app.users
			 SET google_id = $1,
			     google_email = $2,
			     google_name = $3,
			     google_picture = $4
			 WHERE id = $5`,
			googleID,
			email,
			name,
			picture,
			user.ID,
		)
		if updateErr != nil {
			return nil, updateErr
		}

		return GetUserByID(user.ID)
	}

	_, err = DB.Exec(
		`INSERT INTO app.users (email, password, google_id, google_email, google_name, google_picture)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		email,
		"",
		googleID,
		email,
		name,
		picture,
	)
	if err != nil {
		return nil, err
	}

	return GetUserByEmail(email)
}

func CreateUser(email, password string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	result, err := DB.Exec(
		`INSERT INTO app.users (email, password)
		 VALUES ($1, $2)`,
		email,
		string(hashedPassword),
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 1 {
		fmt.Println("User created successfully:", email)
	}

	return nil
}

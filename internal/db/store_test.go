package db

import (
	"context"
	"database/sql"
	"testing"
)

func TestLoginWithProvider_UserDoesNotExist(t *testing.T) {
	ctx := context.Background()

	email := generateUniqueEmail()
	provider := "google"
	providerID := generateUniqueProviderUserID(provider)

	user, err := sqlStore.LoginWithProvider(ctx, provider, providerID, email)

	mustOk(t, err)
	mustEqual(t, user.Email, email)
}

func TestLoginWithProvider_UserExistsLinkedProvider(t *testing.T) {
	ctx := context.Background()

	user := createTestUser(t, ctx, sqlStore)

	provider := "google"
	providerID := generateUniqueProviderUserID(provider)

	createTestOAuthProvider(t, ctx, sqlStore, user.ID, provider, providerID)

	result, err := sqlStore.LoginWithProvider(ctx, provider, providerID, user.Email)

	mustOk(t, err)
	mustEqual(t, result.ID, user.ID)
	mustEqual(t, result.Email, user.Email)

}

func TestLoginWithProvider_UserExistsNotLinkedProvider(t *testing.T) {
	ctx := context.Background()

	user := createTestUser(t, ctx, sqlStore)

	provider := "google"
	providerID := generateUniqueProviderUserID(provider)

	result, err := sqlStore.LoginWithProvider(ctx, provider, providerID, user.Email)

	mustOk(t, err)
	mustEqual(t, result.ID, user.ID)
	mustEqual(t, result.Email, user.Email)
}

func createTestUser(t *testing.T, ctx context.Context, q Querier) User {
	t.Helper()

	email := generateUniqueEmail()

	user, err := q.CreateUser(ctx, CreateUserParams{
		Email:        email,
		PasswordHash: sql.NullString{String: "test-hash", Valid: true},
	})

	mustOk(t, err)

	return user
}

func createTestOAuthProvider(t *testing.T, ctx context.Context, q Querier, userID int64, provider, providerUserID string) {
	t.Helper()

	_, err := q.CreateUserProvider(ctx, CreateUserProviderParams{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	mustOk(t, err)
}

func TestGetUserByEmail(t *testing.T) {

	user := createTestUser(t, ctx, sqlStore)
	dbUser, err := sqlStore.GetUserByEmail(ctx, user.Email)

	mustOk(t, err)
	mustEqual(t, dbUser.ID, user.ID)
	mustEqual(t, dbUser.Email, user.Email)
	mustEqual(t, dbUser.PasswordHash, user.PasswordHash)
}

func TestUpdateUserPassword(t *testing.T) {

	user := createTestUser(t, ctx, sqlStore)

	newPasswordHash := sql.NullString{String: "newpasswordhash", Valid: true}
	err := sqlStore.UpdateUserPassword(ctx, UpdateUserPasswordParams{
		PasswordHash: newPasswordHash,
		Email:        user.Email,
	})
	mustOk(t, err)

	dbUser, err := sqlStore.GetUserByEmail(ctx, user.Email)
	mustOk(t, err)
	mustEqual(t, dbUser.PasswordHash.String, newPasswordHash.String)
}

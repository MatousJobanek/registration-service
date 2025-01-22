package signup

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeUserID(t *testing.T) {
	t.Run("test valid user ID unchanged", func(t *testing.T) {
		userID := "abcde-12345"
		encoded := EncodeUserIdentifier(userID)
		require.Equal(t, userID, encoded)
	})
	t.Run("test user ID with invalid characters", func(t *testing.T) {
		userID := "abcde\\*-12345"
		encoded := EncodeUserIdentifier(userID)
		require.Equal(t, "c0177ca4-abcde-12345", encoded)
	})
	t.Run("test user ID with invalid prefix", func(t *testing.T) {
		userID := "-1234567"
		encoded := EncodeUserIdentifier(userID)
		require.Equal(t, "ca3e1e0f-1234567", encoded)
	})
	t.Run("test user ID that exceeds max length", func(t *testing.T) {
		userID := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-01234567890123456789"
		encoded := EncodeUserIdentifier(userID)
		require.Equal(t, "e3632025-0123456789abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqr", encoded)
	})
	t.Run("test user ID with colon separator", func(t *testing.T) {
		userID := "abc:xyz"
		encoded := EncodeUserIdentifier(userID)
		require.Equal(t, "a05a4053-abcxyz", encoded)
	})
	t.Run("test user ID with invalid end character", func(t *testing.T) {
		userID := "abc---"
		encoded := EncodeUserIdentifier(userID)
		require.Equal(t, "ed6bd2b5-abc", encoded)
	})
}

//go:build github_env

// -tags github_env を追加するとテストで実行される。
// GITHUB_TOKENの環境変数が必要なため、通常の実行では除外するようにしておく。

package impl

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCheckUserExist(t *testing.T) {
	g := NewGitHub("test")

	testCases := map[string]struct {
		userID   string
		expected bool
	}{
		"存在するユーザー": {
			userID:   "ikura-hamu",
			expected: true,
		},
		"存在しないユーザー": {
			userID:   uuid.New().String(),
			expected: false,
		},
		"存在するけどユーザーじゃない": {
			userID:   "traP-jp",
			expected: false,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			exist, err := g.CheckUserExist(ctx, test.userID)
			if err != nil {
				t.Fatal(err)
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expected, exist)
		})
	}
}

package valueobject_test

import (
	"strings"
	"testing"

	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    valueobject.Title
		wantErr error
	}{
		{
			name:  "正常: 通常のタイトル",
			input: "買い物リスト",
			want:  valueobject.Title("買い物リスト"),
		},
		{
			name:  "正常: 1文字",
			input: "a",
			want:  valueobject.Title("a"),
		},
		{
			name:  "正常: 100文字（上限）",
			input: strings.Repeat("a", 100),
			want:  valueobject.Title(strings.Repeat("a", 100)),
		},
		{
			name:    "異常: 空文字",
			input:   "",
			wantErr: valueobject.ErrTitleIsEmpty,
		},
		{
			name:    "異常: 101文字（上限超過）",
			input:   strings.Repeat("a", 101),
			wantErr: valueobject.ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewTitle(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.input, got.String())
		})
	}
}

package valueobject_test

import (
	"testing"

	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    valueobject.Status
		wantErr error
	}{
		{
			name:  "正常: pending",
			input: "pending",
			want:  valueobject.StatusPending,
		},
		{
			name:  "正常: in_progress",
			input: "in_progress",
			want:  valueobject.StatusInProgress,
		},
		{
			name:  "正常: completed",
			input: "completed",
			want:  valueobject.StatusCompleted,
		},
		{
			name:    "異常: 未定義のステータス",
			input:   "unknown",
			wantErr: valueobject.ErrInvalidStatus,
		},
		{
			name:    "異常: 大文字",
			input:   "Pending",
			wantErr: valueobject.ErrInvalidStatus,
		},
		{
			name:    "異常: 空文字",
			input:   "",
			wantErr: valueobject.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewStatus(tt.input)
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

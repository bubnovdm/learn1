package repo

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"sync"
	"testing"
)

func TestRepo_AddItem(t *testing.T) {
	type fields struct {
		userIDMap map[int][]*Item
		mu        sync.Mutex
	}
	type args struct {
		userID int
		items  []*Item
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repo{
				userIDMap: tt.fields.userIDMap,
				mu:        tt.fields.mu,
			}
			repo.AddItem(tt.args.userID, tt.args.items)
		})
	}
}

func TestRepo_ClearCard(t *testing.T) {
	repo := &Repo{userIDMap: map[int][]*Item{1: []*Item{{SkuID: 12345, Count: 2}}}}
	repo.ClearCard(1)
	require.Empty(t, repo.userIDMap)
}

func TestRepo_GetItems(t *testing.T) {
	type fields struct {
		userIDMap map[int][]*Item
		mu        sync.Mutex
	}
	type args struct {
		userID int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Item
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repo{
				userIDMap: tt.fields.userIDMap,
				mu:        tt.fields.mu,
			}
			if got := repo.GetItems(tt.args.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_RemoveItem(t *testing.T) {
	type fields struct {
		userIDMap map[int][]*Item
		mu        sync.Mutex
	}
	type args struct {
		userID int
		skuID  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repo{
				userIDMap: tt.fields.userIDMap,
				mu:        tt.fields.mu,
			}
			repo.RemoveItem(tt.args.userID, tt.args.skuID)
		})
	}
}

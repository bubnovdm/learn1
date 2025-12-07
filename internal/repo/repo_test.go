package repo

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepo_AddItem(t *testing.T) {
	repo := NewRepo()

	// Попробуем добавить 1 элемент
	repo.AddItem(1, []*Item{{SkuID: 1, Count: 5}})
	// Ожидаем длину 1
	require.Len(t, repo.GetItems(1), 1)

	// Попробуем добавить несколько элементов
	repo.AddItem(1, []*Item{{SkuID: 2, Count: 10}, {SkuID: 3, Count: 15}})
	// Ожидаем длину 3
	require.Len(t, repo.GetItems(1), 3)
	// Проверим ещё и значения
	expected := map[int]int{
		1: 5,
		2: 10,
		3: 15,
	}
	items := repo.GetItems(1)
	for _, item := range items {
		require.Equal(t, expected[item.SkuID], item.Count)
	}

	// Попробуем добавить Count для имеющегося элемента
	repo.AddItem(1, []*Item{{SkuID: 1, Count: 5}})
	//Ожидаем, что теперь Count для него 10
	items = repo.GetItems(1)
	// Но если поменяется логика добавления и нижележащий массив будет алоцироваться, надо будет переписать на цикл с поиском
	require.Equal(t, 10, items[0].Count)
}

func TestRepo_ClearCard(t *testing.T) {
	repo := &Repo{userIDMap: map[int][]*Item{1: []*Item{{SkuID: 12345, Count: 2}}}}
	repo.ClearCart(1)
	require.Empty(t, repo.userIDMap)
}

func TestRepo_GetItems(t *testing.T) {
	repo := &Repo{userIDMap: map[int][]*Item{1: []*Item{{SkuID: 12345, Count: 2}}}}

	items := repo.GetItems(1)

	require.Len(t, items, 1)
	require.Equal(t, 12345, items[0].SkuID)
	require.Equal(t, 2, items[0].Count)

	require.Nil(t, repo.GetItems(2))
}

func TestRepo_RemoveItem(t *testing.T) {
	repo := &Repo{userIDMap: map[int][]*Item{1: []*Item{
		{SkuID: 1, Count: 5},
		{SkuID: 2, Count: 10},
		{SkuID: 3, Count: 15},
		{SkuID: 4, Count: 20},
	}}}

	// Удаляем 1 элемент
	repo.RemoveItem(1, 3)
	items := repo.GetItems(1)

	// Ожидаем длину на 1 меньше, т.е. 3
	require.Len(t, items, 3)
	for _, item := range items {
		// Ожидаем, что элемента с SkuID 3 не будет
		require.NotEqual(t, 3, item.SkuID)
	}

	// Попробуем удалить элемент со SkuID 5 (ничего не произойдёт), ожидаем длину 3
	repo.RemoveItem(1, 5)
	require.Len(t, repo.GetItems(1), 3)

	//Удалим все остальные элементы поштучно
	repo.RemoveItem(1, 1)
	repo.RemoveItem(1, 2)
	repo.RemoveItem(1, 4)
	// Всё удалили, ожидаем длину 0
	require.Len(t, repo.GetItems(1), 0)
}

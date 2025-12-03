package repo

/*
 Наше in-memory хранилище, п.3 из Readme
*/

import "sync"

type Repo struct {
	userIDMap map[int][]*Item
	mu        sync.Mutex
}

type Item struct {
	SkuID int
	Count int
}

func NewRepo() Repo {
	return Repo{userIDMap: make(map[int][]*Item)}
}

func (repo *Repo) AddItem(userID int, items []*Item) {
	repo.mu.Lock()
	userItems := repo.userIDMap[userID]
	defer func() { repo.userIDMap[userID] = userItems }()
	for _, item := range items {
		isFound := false
		for _, userItem := range userItems {
			if userItem.SkuID == item.SkuID {
				isFound = true
				userItem.Count += item.Count
				break
			}

		}
		if !isFound {
			userItems = append(userItems, item)
		}
	}
	repo.mu.Unlock()
}

func (repo *Repo) ClearCard(userID int) {
	repo.mu.Lock()
	delete(repo.userIDMap, userID)
	repo.mu.Unlock()
}

func (repo *Repo) RemoveItem(userID int, skuID int) {
	repo.mu.Lock()
	userItems := repo.userIDMap[userID]
	defer func() { repo.userIDMap[userID] = userItems }()
	for index, userItem := range userItems {
		if userItem.SkuID == skuID {
			// Можно сделать без алокаций перемещая последний на место удаляемого
			userItems = append(userItems[:index], userItems[index+1:]...)
		}
	}
	repo.mu.Unlock()
}

func (repo *Repo) GetItems(userID int) []*Item {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.userIDMap[userID]
}

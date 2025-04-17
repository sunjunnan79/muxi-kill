package main

import (
	"gorm.io/gorm"
)

type DAO struct {
	db *gorm.DB
}

func NewDAO(db *gorm.DB) *DAO {
	return &DAO{db: db}
}

func (d *DAO) CreateItem(item *Item) error {
	return d.db.Create(item).Error
}

func (d *DAO) DeleteItem(item *Item) error {
	return d.db.Delete(item).Error
}

func (d *DAO) SaveItem(item *Item) error {
	return d.db.Save(item).Error
}

func (d *DAO) GetItemByID(id uint) (*Item, error) {
	var item Item
	if err := d.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *DAO) GetItems() ([]Item, error) {
	var items []Item
	if err := d.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

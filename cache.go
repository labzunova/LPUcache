package main

import (
	"container/list"
	"errors"
	"time"
)

type keyValueTTL struct {
	key      uint32
	value    string
	initTime time.Time
}

var initSize uint32 = 48    // каждый элемент изначальной структуры весит 8 байт
var mapElemSize uint32 = 16 // 8 байта ключ и 8 указатель на слайс
var elemSize uint32 = 128   // 128 байта строка и 24 Time

type Cache struct {
	size          uint32
	memory        uint32
	TTL           time.Duration
	values        *list.List               // для хранения значений в том порядке, в котором к ним обращались
	valuesHashes  map[uint32]*list.Element // для доступа к значению по ключу
	currentMemory uint32
}

func NewCache(size, TTL, memory uint32) Cache {
	return Cache{size, memory, time.Duration(TTL) * time.Second, list.New(), make(map[uint32]*list.Element),
		initSize}
}

func (c *Cache) Put(i uint32, value string) error {
	newVal := &keyValueTTL{
		value:    value,
		key:      i,
		initTime: time.Now(),
	}

	// если элемент с данным ключом уже есть в кэше
	elem, exists := c.valuesHashes[i]
	if exists {
		elem.Value = newVal
		c.values.MoveToFront(elem) // перемещаем в начало как самый "самый свежий"
		return nil
	}

	// если элемента с данным ключом нет
	if c.currentMemory+elemSize+mapElemSize > c.memory { // проверяем, не выйдет ли кэш за рамки заданного значения памяти
		return errors.New("cache overflow")
	}

	memPlace := c.values.PushFront(newVal)
	c.valuesHashes[i] = memPlace
	c.currentMemory += elemSize + mapElemSize // размер элемента слайса + мапки

	// если кэш заполнен - удаляем элемент, который последний в листе
	if c.values.Len() > int(c.size) {
		delete(c.valuesHashes, c.values.Back().Value.(*keyValueTTL).key) // из мапки
		c.values.Remove(c.values.Back())                                 // и из листа
	}

	return nil
}

func (c *Cache) Get(i uint32) string {
	// если элемент с данным ключом уже есть в кэше
	existingValue, exist := c.valuesHashes[i]
	if exist {
		if time.Now().Sub(existingValue.Value.(*keyValueTTL).initTime) > c.TTL { // проверка на TTL
			return "not found"
		}

		// если элемент еще не "протух", то перемещаем в начало слайса и возвращаем его
		existingValue.Value.(*keyValueTTL).initTime = time.Now()
		c.values.MoveToFront(existingValue)

		return existingValue.Value.(*keyValueTTL).value
	}

	return "not found"
}

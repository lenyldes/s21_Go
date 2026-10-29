package auth

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type familyRecord struct {
	current  uuid.UUID
	previous uuid.UUID
	userId   uuid.UUID
	role     string
}

type JwtCache struct {
	FamilyRecords map[uuid.UUID]*familyRecord
	JtiToFamily   map[uuid.UUID]uuid.UUID
	m             sync.Mutex
}

func NewJwtCache() *JwtCache {
	return &JwtCache{
		FamilyRecords: make(map[uuid.UUID]*familyRecord),
		JtiToFamily:   make(map[uuid.UUID]uuid.UUID),
		m:             sync.Mutex{},
	}
}

func (c *JwtCache) Rotate(cur uuid.UUID) (uuid.UUID, error) {
	c.m.Lock()
	defer c.m.Unlock()
	familyId, ok := c.JtiToFamily[cur]
	if !ok {
		return uuid.Nil, errors.New("family id not exists")
	}
	familyRecord, ok := c.FamilyRecords[familyId]
	if !ok {
		return uuid.Nil, errors.New("family record not exists")
	}

	if familyRecord.current != cur {
		delete(c.JtiToFamily, familyRecord.current)
		delete(c.JtiToFamily, familyRecord.previous)
		delete(c.FamilyRecords, familyId)
		return uuid.Nil, errors.New("cur uuid not match with last in jwt cache")
	}
	delete(c.JtiToFamily, familyRecord.previous)
	familyRecord.previous = familyRecord.current
	newCurrentTokenId := uuid.New()
	c.JtiToFamily[newCurrentTokenId] = familyId
	familyRecord.current = newCurrentTokenId
	return newCurrentTokenId, nil
}

func (c *JwtCache) CreateFamily(userId uuid.UUID, role string) uuid.UUID {
	c.m.Lock()
	defer c.m.Unlock()
	jti := uuid.New()
	familyId := uuid.New()
	c.JtiToFamily[jti] = familyId
	record := familyRecord{
		current:  jti,
		previous: uuid.Nil,
		userId:   userId,
		role:     role,
	}
	c.FamilyRecords[familyId] = &record
	return jti
}

func (c *JwtCache) LookUpUser(tokenId uuid.UUID) (uuid.UUID, string, error) {
	c.m.Lock()
	defer c.m.Unlock()
	familyId, ok := c.JtiToFamily[tokenId]
	if !ok {
		return uuid.Nil, "", errors.New("family id not found")
	}
	familyRecord, ok := c.FamilyRecords[familyId]
	if !ok {
		return uuid.Nil, "", errors.New("family record not found")
	}
	return familyRecord.userId, familyRecord.role, nil
}

package main

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// PrimaryKey : primary key schema struct for SQL objects
type PrimaryKey struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
}

// HealthcheckType : Types of healthchecks
type HealthcheckType struct {
	PrimaryKey
	Name string
}

// BeforeCreate : function for injecting primary keys into DB
func (pk *PrimaryKey) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New())
	return nil
}

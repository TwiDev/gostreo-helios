package db

import (
	"time"

	"gorm.io/gorm"
)

type Workflow struct {
	ID            string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID     string  `gorm:"type:uuid;not null"`
	CreatedByID   *string `gorm:"type:uuid"`
	Name          string  `gorm:"not null"`
	Description   *string
	IsActive      bool      `gorm:"default:true"`
	TriggerNodeID *string   `gorm:"type:uuid"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	Nodes      []WorkflowNode `gorm:"foreignKey:WorkflowID"`
	Executions []ExecutionLog `gorm:"foreignKey:WorkflowID"`
	Events     []Event        `gorm:"foreignKey:WorkflowID"`
}

type WorkflowNode struct {
	ID                string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkflowID        string `gorm:"type:uuid;not null;index"`
	Type              string `gorm:"not null"`
	Name              *string
	Position          map[string]any `gorm:"type:jsonb"`
	IsTrigger         bool           `gorm:"default:false"`
	NextNodeIDs       []string       `gorm:"type:text[]"`
	ConnectorActionID *string        `gorm:"type:uuid"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`

	Parameters []NodeParameter `gorm:"foreignKey:NodeID"`
	RunHistory []RunHistory    `gorm:"foreignKey:NodeID"`
}

type NodeParameter struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	NodeID    string    `gorm:"type:uuid;not null;index"`
	Key       string    `gorm:"not null"`
	Value     string    `gorm:"not null"`
	Type      string    `gorm:"default:'string'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type ExecutionLog struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkflowID    string    `gorm:"type:uuid;not null;index"`
	TriggerNodeID *string   `gorm:"type:uuid"`
	StartedAt     time.Time `gorm:"not null"`
	EndedAt       *time.Time
	Status        string `gorm:"not null"`
	ErrorMessage  *string
	Metadata      map[string]any `gorm:"type:jsonb;default:'{}'::jsonb"`

	RunHistory []RunHistory `gorm:"foreignKey:ExecutionID"`
	Events     []Event      `gorm:"foreignKey:ExecutionID"`
}

type RunHistory struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ExecutionID  string    `gorm:"type:uuid;not null;index"`
	NodeID       *string   `gorm:"type:uuid"`
	StartedAt    time.Time `gorm:"not null"`
	EndedAt      *time.Time
	Status       string         `gorm:"not null"`
	Input        map[string]any `gorm:"type:jsonb"`
	Output       map[string]any `gorm:"type:jsonb"`
	ErrorMessage *string
}

type Event struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkflowID  string    `gorm:"type:uuid;not null;index"`
	ReceivedAt  time.Time `gorm:"default:now()"`
	Source      *string
	Payload     map[string]any `gorm:"type:jsonb"`
	Processed   bool           `gorm:"default:false"`
	ExecutionID *string        `gorm:"type:uuid"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Workflow{},
		&WorkflowNode{},
		&NodeParameter{},
		&ExecutionLog{},
		&RunHistory{},
		&Event{},
	)
}

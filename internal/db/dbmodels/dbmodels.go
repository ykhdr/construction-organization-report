package dbmodels

import (
	"encoding/json"
	"time"
)

type ReportDB struct {
	ID                 int             `db:"id"`
	ProjectID          int             `db:"project_id"`
	ReportCreationDate time.Time       `db:"report_creation_date"`
	ReportFile         json.RawMessage `db:"report_file"`
}

type ScheduleDB struct {
	ID            int `db:"id"`
	WorkType      WorkTypeDB
	Team          TeamDB
	ProjectID     int       `db:"project_id"`
	PlanStartDate time.Time `db:"plan_start_date"`
	PlanEndDate   time.Time `db:"plan_end_date"`
	FactStartDate time.Time `db:"fact_start_date"`
	FactEndDate   time.Time `db:"fact_end_date"`
	PlanOrder     int       `db:"plan_order"`
	FactOrder     int       `db:"fact_order"`
}

type WorkTypeDB struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type TeamDB struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	ProjectID int    `db:"project_id"`
}

type EstimateDB struct {
	ID             int       `db:"id"`
	LastUpdateDate time.Time `db:"last_update_date"`
}

type MaterialUsageDB struct {
	EstimateID   int    `db:"estimate_id"`
	MaterialID   int    `db:"material_id"`
	Name         string `db:"name"`
	Cost         int    `db:"cost"`
	PlanQuantity int    `db:"plan_quantity"`
	FactQuantity int    `db:"fact_quantity"`
}

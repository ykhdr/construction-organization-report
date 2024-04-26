package model

import "time"

type Report struct {
	ID        int         `json:"id"`
	ProjectID int         `json:"project_id"`
	Schedules []*Schedule `json:"schedules"`
	Estimate  *Estimate   `json:"estimate"`
}

type Schedule struct {
	WorkType      *WorkType
	Team          *Team     `json:"construction_team"`
	PlanStartDate time.Time `json:"plan_start_date"`
	PlanEndDate   time.Time `json:"plan_end_date"`
	FactStartDate time.Time `json:"fact_start_date"`
	FactEndDate   time.Time `json:"fact_end_date"`
	PlanOrder     int       `json:"plan_order"`
	FactOrder     int       `json:"fact_order"`
}

type WorkType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Estimate struct {
	MaterialUsage  []*MaterialUsage `json:"material_usage"`
	LastUpdateDate time.Time        `json:"last_update_date"`
}

type MaterialUsage struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Cost         int    `json:"cost"`
	PlanQuantity int    `json:"plan_quantity"`
	FactQuantity int    `json:"fact_quantity"`
}

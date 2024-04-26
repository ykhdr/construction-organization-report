package db

import (
	"construction-organization-report/internal/db/dbmodels"
	"context"
	"database/sql"
)

func GetEstimate(ctx context.Context, db *sql.DB, projectId int) (*dbmodels.EstimateDB, error) {
	var entity dbmodels.EstimateDB
	err := db.QueryRowContext(ctx, `
	SELECT id, last_update_date 
	    FROM estimate
	WHERE id = $1
	`, projectId).Scan(&entity.ID, &entity.LastUpdateDate)

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func GetMaterialUsages(ctx context.Context, db *sql.DB, estimateId int) ([]*dbmodels.MaterialUsageDB, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT estimate_id, id, name, cost, plan_quantity, fact_quantity
	    FROM material_usage AS mu 
		JOIN material AS m ON mu.material_id = m.id
	WHERE estimate_id = $1
	`, estimateId)

	if err != nil {
		return nil, err
	}

	var entities []*dbmodels.MaterialUsageDB

	defer rows.Close()

	for rows.Next() {
		var entity dbmodels.MaterialUsageDB
		err = rows.Scan(&entity.EstimateID, &entity.MaterialID, &entity.Name, &entity.Cost, &entity.PlanQuantity, &entity.FactQuantity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func GetConstructionTeam(ctx context.Context, db *sql.DB, teamId int) (*dbmodels.TeamDB, error) {
	var entity dbmodels.TeamDB

	err := db.QueryRowContext(ctx, `
	SELECT id, name, project_id
	FROM construction_team AS ct
	WHERE ct.id = $1
	`, teamId).Scan(&entity.ID, &entity.Name, &entity.ProjectID)

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func GetSchedules(ctx context.Context, db *sql.DB, projectId int) ([]*dbmodels.ScheduleDB, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT ws.id, ct.id, ct.name, wt.id, wt.name, plan_start_date, plan_end_date, fact_start_date, fact_end_date, plan_order, fact_order, ws.project_id
	FROM work_schedule AS ws 
	JOIN construction_team AS ct ON ws.construction_team_id = ct.id
	JOIN work_type AS wt ON ws.work_type_id = wt.id
	WHERE ws.project_id = $1
	`, projectId)

	if err != nil {
		return nil, err
	}

	var entities []*dbmodels.ScheduleDB

	defer rows.Close()

	for rows.Next() {
		var entity dbmodels.ScheduleDB
		err := rows.Scan(&entity.ID, &entity.Team.ID, &entity.Team.Name, &entity.WorkType.ID, &entity.WorkType.Name, &entity.PlanStartDate, &entity.PlanEndDate, &entity.FactStartDate, &entity.FactEndDate, &entity.PlanOrder, &entity.FactOrder, &entity.ProjectID)
		if err != nil {
			return nil, err
		}

		entity.Team.ProjectID = projectId

		entities = append(entities, &entity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func GetWorkType(ctx context.Context, db *sql.DB, workTypeId int) (*dbmodels.WorkTypeDB, error) {
	var entity dbmodels.WorkTypeDB

	err := db.QueryRowContext(ctx, `
	SELECT id, name
	FROM work_type
	WHERE id = $1
	`, workTypeId).Scan(&entity.ID, &entity.Name)

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func GetReports(ctx context.Context, db *sql.DB, projectID int) ([]*dbmodels.ReportDB, error) {

	rows, err := db.QueryContext(ctx, `
	SELECT id, project_id, report_creation_date, report_file
	FROM report
	WHERE project_id = $1
	`, projectID)

	if err != nil {
		return nil, err
	}

	var entities []*dbmodels.ReportDB
	defer rows.Close()
	for rows.Next() {
		var entity dbmodels.ReportDB
		err := rows.Scan(&entity.ID, &entity.ProjectID, &entity.ReportCreationDate, &entity.ReportFile)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

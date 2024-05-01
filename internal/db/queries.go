package db

import (
	"construction-organization-report/internal/report"
	"context"
	"database/sql"
	"time"
)

func GetEstimate(ctx context.Context, db *sql.DB, projectId int) (*EstimateDB, error) {
	var entity EstimateDB
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

func GetMaterialUsages(ctx context.Context, db *sql.DB, estimateId int) ([]*MaterialUsageDB, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT estimate_id, id, name, cost, plan_quantity, fact_quantity
	    FROM material_usage AS mu 
		JOIN material AS m ON mu.material_id = m.id
	WHERE estimate_id = $1
	`, estimateId)

	if err != nil {
		return nil, err
	}

	var entities []*MaterialUsageDB

	defer rows.Close()

	for rows.Next() {
		var entity MaterialUsageDB
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

func GetSchedules(ctx context.Context, db *sql.DB, projectId int) ([]*ScheduleDB, error) {
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

	var entities []*ScheduleDB

	defer rows.Close()

	for rows.Next() {
		var entity ScheduleDB
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

func GetReports(ctx context.Context, db *sql.DB, projectID int) ([]*ReportDB, error) {

	rows, err := db.QueryContext(ctx, `
	SELECT id, project_id, report_creation_date, report_file
	FROM report
	WHERE project_id = $1
	`, projectID)

	if err != nil {
		return nil, err
	}

	var entities []*ReportDB
	defer rows.Close()
	for rows.Next() {
		var entity ReportDB
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

func SaveReport(ctx context.Context, db *sql.DB, report *ReportDB) (int, error) {

	err := db.QueryRowContext(ctx, `
	INSERT INTO report (project_id, report_creation_date, report_file)
	VALUES ($1, $2, $3)
	RETURNING id
	`, report.ProjectID, time.Now(), report.ReportFile).Scan(&report.ID)

	if err != nil {
		return -1, err
	}

	return report.ID, nil
}

func GetReport(ctx context.Context, db *sql.DB, reportID int) (*report.RawReport, error) {
	var rawReport report.RawReport
	err := db.QueryRowContext(ctx, `
	SELECT id, project_id, report_creation_date, report_file
	FROM report
	WHERE id = $1
	`, reportID).Scan(&rawReport.ID, &rawReport.ProjectID, &rawReport.ReportCreationDate, &rawReport.ReportFile)

	if err != nil {
		return nil, err
	}

	return &rawReport, nil
}

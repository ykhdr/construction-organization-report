package report

import (
	"construction-organization-report/internal/db"
	"construction-organization-report/internal/log"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

func CreateReport(projectID int, database *sql.DB) (*Report, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	errorChan := make(chan error, 3)
	estimateChan := make(chan *Estimate, 1)
	materialUsagesChan := make(chan []*MaterialUsage, 1)
	schedulesChan := make(chan []*Schedule, 1)

	go func() {
		defer wg.Done()
		estimate, err := db.GetEstimate(context.Background(), database, projectID)
		if err != nil {
			errorChan <- errors.Join(fmt.Errorf("error on get estimate"), err)
			return
		}

		estimateChan <- &Estimate{LastUpdateDate: estimate.LastUpdateDate}
	}()

	go func() {
		defer wg.Done()
		materialUsagesDB, err := db.GetMaterialUsages(context.Background(), database, projectID)
		if err != nil {
			errorChan <- errors.Join(fmt.Errorf("error on get material usages"), err)
			return
		}

		usages := make([]*MaterialUsage, 0)

		for _, usage := range materialUsagesDB {
			usages = append(usages, &MaterialUsage{
				ID:           usage.MaterialID,
				Name:         usage.Name,
				Cost:         usage.Cost,
				PlanQuantity: usage.PlanQuantity,
				FactQuantity: usage.FactQuantity,
			})
		}

		materialUsagesChan <- usages
	}()

	go func() {
		defer wg.Done()
		schedulesDB, err := db.GetSchedules(context.Background(), database, projectID)
		if err != nil {
			errorChan <- errors.Join(fmt.Errorf("error on get schedules"), err)
			return
		}

		schedules := make([]*Schedule, 0)
		for _, schedule := range schedulesDB {
			schedules = append(schedules, &Schedule{
				WorkType:      &WorkType{ID: schedule.WorkType.ID, Name: schedule.WorkType.Name},
				Team:          &Team{ID: schedule.Team.ID, Name: schedule.Team.Name},
				PlanStartDate: schedule.PlanStartDate,
				PlanEndDate:   schedule.PlanEndDate,
				FactStartDate: schedule.FactStartDate,
				FactEndDate:   schedule.FactEndDate,
				PlanOrder:     schedule.PlanOrder,
				FactOrder:     schedule.FactOrder,
			})
		}

		schedulesChan <- schedules
	}()

	wg.Wait()

	close(errorChan)
	close(estimateChan)
	close(materialUsagesChan)
	close(schedulesChan)

	report := &Report{ProjectID: projectID}

	var joinErrors error
	for err := range errorChan {
		joinErrors = errors.Join(joinErrors, err)
	}

	if joinErrors != nil {
		log.Logger.WithError(joinErrors).Error("errors while creating report")
		return nil, fmt.Errorf("errors while creating report : %v", joinErrors)
	}

	report.Estimate = <-estimateChan
	report.Schedules = <-schedulesChan
	report.Estimate.MaterialUsage = <-materialUsagesChan
	report.Estimate.LastUpdateDate = time.Now()

	log.Logger.Infoln("Saving report to database")
	reportDB, err := convertReportToReportDB(report)
	if err != nil {
		log.Logger.WithError(err).Error("error while converting report to raw report")
		return nil, fmt.Errorf("error while converting report to raw report")
	}

	_, err = db.SaveReport(context.Background(), database, reportDB)
	if err != nil {
		log.Logger.WithError(err).Error("error while saving report")
		return nil, fmt.Errorf("error while saving report")
	}
	log.Logger.Infoln("Report saved")

	return report, nil
}

func GetReports(projectID int, database *sql.DB) ([]*RawReport, error) {

	reports, err := db.GetReports(context.Background(), database, projectID)

	if err != nil {
		return nil, err
	}

	rawReports := make([]*RawReport, 0)

	for _, report := range reports {
		rawReports = append(rawReports, &RawReport{
			ID:                 report.ID,
			ProjectID:          report.ProjectID,
			ReportCreationDate: report.ReportCreationDate,
			ReportFile:         report.ReportFile,
		})
	}

	return rawReports, nil
}

func convertReportToReportDB(report *Report) (*db.ReportDB, error) {
	intermediate := struct {
		Schedules []*Schedule `json:"schedules"`
		Estimate  *Estimate   `json:"estimate"`
	}{
		Schedules: report.Schedules,
		Estimate:  report.Estimate,
	}

	reportData, err := json.Marshal(intermediate)
	if err != nil {
		return nil, err
	}

	reportDB := &db.ReportDB{
		ID:                 report.ID,
		ProjectID:          report.ProjectID,
		ReportCreationDate: time.Now(), // Пример заполнения даты создания отчета
		ReportFile:         json.RawMessage(reportData),
	}

	return reportDB, nil
}

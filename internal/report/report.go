package report

import (
	"construction-organization-report/internal/db"
	"construction-organization-report/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

func CreateReport(projectID int, database *sql.DB) (*model.Report, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	errorChan := make(chan error)
	estimateChan := make(chan *model.Estimate, 1)
	materialUsagesChan := make(chan []*model.MaterialUsage, 1)
	schedulesChan := make(chan []*model.Schedule, 1)

	go func() {
		defer wg.Done()
		estimate, err := db.GetEstimate(context.Background(), database, projectID)
		if err != nil {
			errorChan <- err
			return
		}

		estimateChan <- &model.Estimate{LastUpdateDate: estimate.LastUpdateDate}
	}()

	go func() {
		defer wg.Done()
		materialUsagesDB, err := db.GetMaterialUsages(context.Background(), database, projectID)
		if err != nil {
			errorChan <- err
			return
		}

		usages := make([]*model.MaterialUsage, 0)

		for _, usage := range materialUsagesDB {
			usages = append(usages, &model.MaterialUsage{
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
			errorChan <- err
			return
		}

		schedules := make([]*model.Schedule, 0)
		for _, schedule := range schedulesDB {
			schedules = append(schedules, &model.Schedule{
				WorkType:      &model.WorkType{ID: schedule.WorkType.ID, Name: schedule.WorkType.Name},
				Team:          &model.Team{ID: schedule.Team.ID, Name: schedule.Team.Name},
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

	report := &model.Report{ProjectID: projectID}

	var joinErrors error
	for err := range errorChan {
		joinErrors = errors.Join(joinErrors, err)
	}

	if joinErrors != nil {
		return nil, fmt.Errorf("errors while creating report : %v", joinErrors)
	}

	report.Estimate = <-estimateChan
	report.Schedules = <-schedulesChan
	report.Estimate.MaterialUsage = <-materialUsagesChan
	report.Estimate.LastUpdateDate = time.Now()

	return report, nil
}

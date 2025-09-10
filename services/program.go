package services

import (
	"context"
	"fmt"
	"linn221/Requester/requests"
)

// ProgramService handles program-related operations
type ProgramService struct {
	db Database
}

// NewProgramService creates a new ProgramService
func NewProgramService(db Database) *ProgramService {
	return &ProgramService{db: db}
}

// GetAllPrograms fetches all programs
func (s *ProgramService) GetAllPrograms(ctx context.Context) ([]requests.Program, error) {
	// Create transaction with context
	tx := s.db.WithContext(ctx)

	// Fetch all programs
	var programs []requests.Program
	if err := tx.Where("1 = 1").Order("created_at DESC").Find(&programs).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch programs: %v", err)
	}

	return programs, nil
}

// GetProgramByID fetches a single program by ID
func (s *ProgramService) GetProgramByID(ctx context.Context, id uint) (*requests.Program, error) {
	// Create transaction with context
	tx := s.db.WithContext(ctx)

	// Fetch program by ID
	var program requests.Program
	if err := tx.First(&program, id).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch program %d: %v", id, err)
	}

	return &program, nil
}

// CreateProgram creates a new program
func (s *ProgramService) CreateProgram(ctx context.Context, program *requests.Program) error {
	// Create transaction with context
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error() != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create the program
	if err := tx.Create(program).Error(); err != nil {
		return fmt.Errorf("failed to create program: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error(); err != nil {
		return fmt.Errorf("failed to commit program creation: %v", err)
	}

	return nil
}

// UpdateProgram updates an existing program
func (s *ProgramService) UpdateProgram(ctx context.Context, program *requests.Program) error {
	// Create transaction with context
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error() != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Update the program
	if err := tx.Model(program).Updates(program).Error(); err != nil {
		return fmt.Errorf("failed to update program: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error(); err != nil {
		return fmt.Errorf("failed to commit program update: %v", err)
	}

	return nil
}

// DeleteProgram deletes a program by ID
func (s *ProgramService) DeleteProgram(ctx context.Context, id uint) error {
	// Create transaction with context
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error() != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Delete the program
	if err := tx.Delete(&requests.Program{}, id).Error(); err != nil {
		return fmt.Errorf("failed to delete program: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error(); err != nil {
		return fmt.Errorf("failed to commit program deletion: %v", err)
	}

	return nil
}

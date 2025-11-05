package services

import (
    "context"
    "errors"
    "time"

    "github.com/superbkibbles/ecommerce/internal/domain/entities"
    "github.com/superbkibbles/ecommerce/internal/domain/ports"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type HomeSectionService struct {
    repo ports.HomeSectionRepository
}

func NewHomeSectionService(repo ports.HomeSectionRepository) ports.HomeSectionService {
    return &HomeSectionService{repo: repo}
}

func (s *HomeSectionService) Create(ctx context.Context, sectionType entities.HomeSectionType, title map[string]string, productIDs, categoryIDs []string, order int, active bool) (*entities.HomeSection, error) {
    // Convert IDs
    var pids []primitive.ObjectID
    for _, id := range productIDs {
        if id == "" { continue }
        oid, err := primitive.ObjectIDFromHex(id)
        if err != nil { return nil, errors.New("invalid product id") }
        pids = append(pids, oid)
    }
    var cids []primitive.ObjectID
    for _, id := range categoryIDs {
        if id == "" { continue }
        oid, err := primitive.ObjectIDFromHex(id)
        if err != nil { return nil, errors.New("invalid category id") }
        cids = append(cids, oid)
    }

    section, err := entities.NewHomeSection(sectionType, title, pids, cids, order, active)
    if err != nil { return nil, err }
    if err := s.repo.Create(ctx, section); err != nil { return nil, err }
    return section, nil
}

func (s *HomeSectionService) Get(ctx context.Context, id string) (*entities.HomeSection, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil { return nil, errors.New("invalid section id") }
    return s.repo.GetByID(ctx, oid)
}

func (s *HomeSectionService) Update(ctx context.Context, section *entities.HomeSection) error {
    if section.ID.IsZero() { return errors.New("id is required") }
    section.UpdatedAt = time.Now()
    return s.repo.Update(ctx, section)
}

func (s *HomeSectionService) Delete(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil { return errors.New("invalid section id") }
    return s.repo.Delete(ctx, oid)
}

func (s *HomeSectionService) List(ctx context.Context, onlyActive bool) ([]*entities.HomeSection, error) {
    return s.repo.List(ctx, onlyActive)
}



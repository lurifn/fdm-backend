package order

import (
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/pkg/repository"
)

// Create sends an email with the provided order according to the configurations loaded
func Create(order string, repo repository.Repository) error {
	log.Info.Printf("Order: %s\n", order)
	return repo.Save(order)
}

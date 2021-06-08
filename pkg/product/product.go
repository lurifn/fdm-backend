package product

import (
	"encoding/json"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/pkg/repository"
)

type Product struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price int `json:"price"` // divide per 100 to get the actual price
	Currency string `json:"currency"`
	ImageURL string `json:"image_url"`
}

func (p *Product) Save(repo repository.Repository) error {
	bin, err := json.Marshal(p)
	if err != nil {
		log.Error.Println("Error reading product: ", err.Error())
		return err
	}

	str := string(bin)
	log.Info.Printf("Save product: %s\n", str)

	id, err := repo.Save(str)
	if err != nil {
		log.Error.Println("Error saving product: ", err.Error())
		return err
	}

	p.ID = id

	return nil
}

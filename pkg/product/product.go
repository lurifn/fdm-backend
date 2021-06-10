package product

import (
	"encoding/json"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/pkg/repository"
)

type Product struct {
	ID          string `json:",omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"` // divide per 100 to get the actual price
	Currency    string `json:"currency,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

func (p *Product) Save(repo repository.Repository) error {
	if len(p.ID) > 0 {
		return p.Update(repo)
	}

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

func (p *Product) Update(repo repository.Repository) error {
	id := p.ID
	p.ID = ""
	bin, err := json.Marshal(map[string]Product{
		id: *p,
	})
	if err != nil {
		log.Error.Println("Error reading product: ", err.Error())
		return err
	}

	str := string(bin)
	log.Info.Printf("Save product: %s\n", str)

	res, err := repo.Update(str)
	if err != nil {
		log.Error.Println("Error saving product: ", err.Error())
		return err
	}

	log.Info.Println("Response: ", res)

	return nil
}

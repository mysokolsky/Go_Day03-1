package types


type Place struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type Page struct {
	Places   []Place
	Count    uint
	PageNum  uint
	Prev     uint
	Next     uint
	Last     uint
	NotFirst bool
	NotLast  bool
}

type ElasticClient struct {
	Es *elasticsearch.Client
	Index string
}
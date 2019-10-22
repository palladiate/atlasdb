package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var db *sqlx.DB

type Island struct {
	Name string `db:"name" json:"name"`
	Type string `db:"type" json:"type"`
	Resource string `db:"resource" json:"resource"`
}


const query = `
select
       i.name,
       r.type,
       r.resource
from
    atlas.resource_type r join
    atlas.goods g join
    atlas.islands i
on
    r.id = g.resourceId and
    i.id = g.islandId
;
`

func main() {
	var err error

	db, err = sqlx.Open("mysql", "atlas:pwd@tcp(10.0.0.22:3306)/atlas")
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("no database!")
	}

	db.Ping()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.Handle("GET", "/", handleIndex)

	ui := r.Group("/ui")
	{
		ui.Handle("GET", "/index", handleIndex)
	}

	rest := r.Group("/app")
	{
		rest.Handle("GET", "", handleRetrieve)
		rest.Handle("POST", "", handleNew)
		rest.Handle("PUT", "", handleUpdate)
	}

	if err = r.Run(":9000"); err != nil {
		panic(err)
	}
}

func handleNew(c *gin.Context) {
	c.JSON(200, nil)
}

func handleUpdate(c *gin.Context) {
	resp, err := db.Queryx(query, 1)
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, resp)
}

func handleRetrieve(c *gin.Context) {
	data, err := getIsland()
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, data)
}


func handleIndex(c *gin.Context) {
	body, err := getIsland()
	if err != nil {
		c.JSON( 500, nil )
	}
	c.HTML(200, "index", body )
}

func getIsland() (islands []Island, err error) {
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var island Island
		err = rows.StructScan(&island)
		if err != nil {
			log.Print(err)
		}
		islands = append(islands, island)
	}
	return
}
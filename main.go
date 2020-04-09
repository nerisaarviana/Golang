package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
	db, err =
		gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/tugas_go?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Gagal Conect Ke Database")
	}
	db.AutoMigrate(&person{})
}

type (
	person struct {
		gorm.Model
		Nama         string `json:"nama"`
		Alamat       string `json:"alamat"`
		JenisKelamin string `json:"jenis_kelamin"`
		NoHp         string `json:"no_hp"`
	}
	transformedPerson struct {
		ID           uint   `json:"id"`
		Nama         string `json:"nama"`
		Alamat       string `json:"alamat"`
		JenisKelamin string `json:"jenis_kelamin"`
		NoHp         string `json:"no_hp"`
	}
)

func cretedPerson(c *gin.Context) {
	var std transformedPerson
	var model person
	c.Bind(&std)
	validasi := validatorCreated(std)
	model = transferVoToModel(std)
	if validasi != "" {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": validasi})
	} else {
		db.Create(&model)
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": model})
	}
}

func fetchAllPersons(c *gin.Context) {
	var model []person
	var vo []transformedPerson

	db.Find(&model)

	if len(model) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}

	for _, item := range model {
		vo = append(vo, transferModelToVo(item))
	}
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": vo})
}

func fetchSinglePerson(c *gin.Context) {
	var model person
	var vo transformedPerson

	modelID := c.Param("id")
	db.Find(&model, modelID)

	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}
	vo = transferModelToVo(model)
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": vo})
}

func updatePerson(c *gin.Context) {
	var model person
	var vo transformedPerson
	modelID := c.Param("id")
	db.First(&model, modelID)

	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}
	c.Bind(&vo)

	validasi := validatorCreated(vo)
	if validasi != "" {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": validasi})
	} else {
		db.Model(&model).Update(transferVoToModel(vo))
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": model})
	}
}

func deletePerson(c *gin.Context) {
	var model person
	modelID := c.Param("id")

	db.First(&model, modelID)
	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak di Temukan"})
	}
	db.Delete(model)
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": "Data Telah berhasil di hapus"})
}

func transferModelToVo(model person) transformedPerson {
	var vo transformedPerson
	vo = transformedPerson{
		ID:           model.ID,
		Nama:         model.Nama,
		Alamat:       model.Alamat,
		JenisKelamin: model.JenisKelamin,
		NoHp:         model.NoHp,
	}
	return vo
}

func transferVoToModel(vo transformedPerson) person {
	var model person
	model = person{
		Nama:         vo.Nama,
		Alamat:       vo.Alamat,
		JenisKelamin: vo.JenisKelamin,
		NoHp:         vo.NoHp,
	}
	return model
}

func validatorCreated(vo transformedPerson) string {

	var kosong string = " Tidak Boleh Kosong"

	if vo.Nama == "" {
		return "Nama" + kosong
	}

	if vo.Alamat == "" {
		return "Alamat" + kosong
	}

	if vo.JenisKelamin == "" {
		return "Jenis Kelamin" + kosong
	}

	if vo.NoHp == "" {
		return "No Hp" + kosong
	}

	return ""
}

func main() {

	router := gin.Default()
	v1 := router.Group("/api/person")
	{
		v1.POST("", cretedPerson)
		v1.GET("", fetchAllPersons)
		v1.GET("/:id", fetchSinglePerson)
		v1.PUT("/:id", updatePerson)
		v1.DELETE("/:id", deletePerson)
	}
	router.Run(":5050")
}

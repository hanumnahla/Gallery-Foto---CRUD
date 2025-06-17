package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	dsn := "root:@tcp(127.0.0.1:3306)/upgris"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Tes koneksi database
	err = db.Ping()
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
}

// Handler untuk menampilkan halaman login

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user string
	err := db.QueryRow("SELECT username FROM users WHERE username = ? AND password = ?", username, password).Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
			return
		}
		log.Println("Error checking credentials:", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Internal server error"})
		return
	}

	c.Redirect(http.StatusFound, "/home")
}

// Handler untuk menampilakn halaman home beserta data foto

func HomePage(c *gin.Context) {
	rows, err := db.Query("SELECT id, judul, deskripsi, gambar, tanggal_upload FROM foto")
	if err != nil {
		log.Println("Error fetching foto:", err)
		c.HTML(http.StatusInternalServerError, "home.html", gin.H{"error": "Gagal mengambil data foto"})
		return
	}
	defer rows.Close()

	type Foto struct {
		ID            int
		Judul         string
		Deskripsi     string
		Gambar        string
		TanggalUpload string
	}

	var fotos []Foto

	for rows.Next() {
		var f Foto
		if err := rows.Scan(&f.ID, &f.Judul, &f.Deskripsi, &f.Gambar, &f.TanggalUpload); err != nil {
			log.Println("Error scanning foto:", err)
			continue
		}
		fotos = append(fotos, f)
	}

	var jumlahFoto int
	err = db.QueryRow("SELECT COUNT(*) FROM foto").Scan(&jumlahFoto)
	if err != nil {
		log.Println("Error menghitung jumlah foto:", err)
		jumlahFoto = 0
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"fotos":   fotos,
		"jumlah1": jumlahFoto,
	})
}

// Handler untuk menambahkan foto baru

func TambahFoto(c *gin.Context) {
	judul := c.PostForm("Judul")
	deskripsi := c.PostForm("deskripsi")

	file, err := c.FormFile("gambar")
	if err != nil {
		log.Println("Gagal mengambil file gambar:", err)
		c.String(http.StatusBadRequest, "Gagal upload gambar")
		return
	}
	// Simpan file ke folder "uploads"
	filename := filepath.Base(file.Filename)
	path := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		log.Println("Gagal menyimpan file:", err)
		c.String(http.StatusInternalServerError, "Gagal menyimpan gambar")
		return
	}

	// Tambahkan tanggal upload
	tanggalUpload := time.Now().Format("2006-01-02 15:04:05")
	// Simpan data ke database
	_, err = db.Exec("INSERT INTO foto (judul, deskripsi, gambar, tanggal_upload) VALUES (?, ?, ?, ?)", judul, deskripsi, filename, tanggalUpload)
	if err != nil {
		log.Println("Error inserting foto:", err)
	}

	c.Redirect(http.StatusFound, "/home")
}

// Handler untuk mengedit foto

func EditFoto(c *gin.Context) {
	id := c.PostForm("id")
	judul := c.PostForm("judul")
	deskripsi := c.PostForm("deskripsi")

	_, err := db.Exec("UPDATE foto SET judul = ?, deskripsi = ? WHERE id = ?", judul, deskripsi, id)
	if err != nil {
		log.Println("Error updating foto:", err)
	}

	c.Redirect(http.StatusFound, "/home")
}

// Handler untuk menghapuus foto berdasarkan ID

func HapusFoto(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM foto WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting foto:", err)
	}

	c.Redirect(http.StatusFound, "/home")
}

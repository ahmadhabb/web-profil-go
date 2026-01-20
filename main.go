package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Feature struct {
	Icon  string
	Title string
	Desc  string
}

type Testimonial struct {
	Name    string
	Company string
	Text    string
	Avatar  string
}

func main() {

	// =========================
	// DATABASE CONNECTION
	// =========================
	dsn := "postgres://cp_user:password123@localhost:5432/company_profile?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Gagal connect ke database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Gagal ping database:", err)
	}

	log.Println("PostgreSQL connected successfully!")

	// Run database migrations
	log.Println("Running database migrations...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migrate driver:", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Migrations completed successfully!")

	// Debug: Cek working directory
	dir, _ := os.Getwd()
	log.Println("Current working directory:", dir)

	// Tentukan command npm berdasarkan OS
	npmCmd := "npm"
	if runtime.GOOS == "windows" {
		npmCmd = "npm.cmd"
	}

	// Build CSS sekali di awal (Synchronous) untuk memastikan file style.css terisi
	log.Println("Building Tailwind CSS (Initial Build)...")
	if err := exec.Command(npmCmd, "run", "build").Run(); err != nil {
		log.Printf("Gagal build Tailwind: %v. Pastikan 'npm install' sudah dijalankan.", err)
	}

	// Jalankan Tailwind CSS watcher secara otomatis di background
	go runTailwind(npmCmd)

	// Inisialisasi template engine dengan path yang benar
	engine := html.New("./views", ".html")

	// Reload template untuk development
	engine.Reload(true)

	// Debug: Cek apakah template bisa dibaca
	engine.Debug(true)

	// Buat instance Fiber dengan template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware
	app.Use(logger.New())
	// Pastikan baris ini ada! Ini melayani file di folder static (CSS/JS)
	app.Static("/static", "./static", fiber.Static{
		CacheDuration: -1, // Matikan cache browser saat development agar perubahan CSS terbaca
	})

	// Data untuk halaman
	companyData := map[string]interface{}{
		"CompanyName":    "TechSolution Inc.",
		"CompanyTagline": "Solusi Digital untuk Masa Depan Bisnis Anda",
		"CompanyAddress": "Jl. Teknologi No. 123, Jakarta",
		"CompanyPhone":   "+62 21 1234 5678",
		"CompanyEmail":   "info@techsolution.com",
	}

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		// Debug: Cek path
		log.Println("Mencoba render index.html")

		features, err := getFeatures(db)
		if err != nil {
			log.Printf("Gagal get features: %v", err)
			features = []Feature{} // fallback empty
		}

		testimonials, err := getTestimonials(db)
		if err != nil {
			log.Printf("Gagal get testimonials: %v", err)
			testimonials = []Testimonial{} // fallback empty
		}

		data := fiber.Map{
			"Title":        "Beranda - " + companyData["CompanyName"].(string),
			"Company":      companyData,
			"Active":       "home",
			"Features":     features,
			"Testimonials": testimonials,
		}
		return c.Render("index", data, "main_layout")
	})

	app.Get("/about", func(c *fiber.Ctx) error {
		data := fiber.Map{
			"Title":   "Tentang Kami - " + companyData["CompanyName"].(string),
			"Company": companyData,
			"Active":  "about",
			"Team": []map[string]string{
				{
					"name":     "Wahyu Chaw",
					"position": "CEO & Founder",
					"bio":      "Berpengalaman 10 tahun di industri teknologi.",
					"avatar":   "👨‍💼",
				},
				{
					"name":     "Ahmad Hab",
					"position": "Lead Developer",
					"bio":      "Ahli dalam Go, Python, dan JavaScript.",
					"avatar":   "👨‍🔧",
				},
			},
		}
		return c.Render("about", data, "main_layout")
	})

	app.Get("/services", func(c *fiber.Ctx) error {
		data := fiber.Map{
			"Title":   "Layanan - " + companyData["CompanyName"].(string),
			"Company": companyData,
			"Active":  "services",
			"Services": []map[string]interface{}{
				{
					"name":        "Web Development",
					"description": "Kami membuat website yang responsif, cepat, dan SEO-friendly.",
					"price":       "Mulai dari Rp 5.000.000",
					"features":    []string{"Responsive Design", "SEO Optimization", "CMS Integration", "Maintenance"},
				},
				{
					"name":        "Mobile App Development",
					"description": "Aplikasi mobile untuk iOS dan Android dengan performa optimal.",
					"price":       "Mulai dari Rp 15.000.000",
					"features":    []string{"iOS & Android", "API Integration", "Push Notification", "App Store Submission"},
				},
				{
					"name":        "Cloud Solutions",
					"description": "Migrasi dan manajemen sistem cloud untuk bisnis Anda.",
					"price":       "Mulai dari Rp 10.000.000",
					"features":    []string{"Cloud Migration", "Server Management", "Backup Solutions", "24/7 Monitoring"},
				},
			},
		}
		return c.Render("services", data, "main_layout")
	})

	app.Get("/contact", func(c *fiber.Ctx) error {
		data := fiber.Map{
			"Title":   "Kontak - " + companyData["CompanyName"].(string),
			"Company": companyData,
			"Active":  "contact",
		}
		return c.Render("contact", data, "main_layout")
	})

	app.Post("/contact", func(c *fiber.Ctx) error {
		// Simulasi proses form submission
		name := c.FormValue("name")
		email := c.FormValue("email")
		message := c.FormValue("message")

		// Log pesan yang diterima
		log.Printf("Pesan baru diterima:\nNama: %s\nEmail: %s\nPesan: %s\n", name, email, message)

		return c.Render("contact", fiber.Map{
			"Title":      "Kontak - " + companyData["CompanyName"].(string),
			"Company":    companyData,
			"Active":     "contact",
			"Success":    true,
			"SuccessMsg": "Terima kasih " + name + ", pesan Anda telah dikirim! Kami akan membalas ke " + email + " segera.",
		}, "main_layout")
	})

	// Tambahkan 404.html terlebih dahulu
	app.Get("/404", func(c *fiber.Ctx) error {
		return c.Render("404", fiber.Map{
			"Title":   "404 - Halaman Tidak Ditemukan",
			"Company": companyData,
		}, "main_layout")
	})

	// Route untuk mengecek file static
	app.Get("/check-static", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"views_exists":    fileExists("./views"),
			"layout_exists":   fileExists("./views/main_layout.html"),
			"index_exists":    fileExists("./views/index.html"),
			"about_exists":    fileExists("./views/about.html"),
			"services_exists": fileExists("./views/services.html"),
			"contact_exists":  fileExists("./views/contact.html"),
		})
	})

	// Static files root (favicon.ico, etc) - ditempatkan setelah route agar tidak menimpa route "/"
	app.Static("/", "./static")

	// 404 Handler (Pindahkan ke sini, setelah semua route didefinisikan)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("404", fiber.Map{
			"Title":   "404 - Halaman Tidak Ditemukan",
			"Company": companyData,
		}, "main_layout")
	})

	log.Println("Server berjalan di http://localhost:3002")
	log.Println("Cek file template di: http://localhost:3002/check-static")
	log.Fatal(app.Listen(":3002"))
}

func getFeatures(db *sql.DB) ([]Feature, error) {
	rows, err := db.QueryContext(context.Background(), "SELECT icon, title, description FROM features ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var features []Feature
	for rows.Next() {
		var f Feature
		err := rows.Scan(&f.Icon, &f.Title, &f.Desc)
		if err != nil {
			return nil, err
		}
		features = append(features, f)
	}
	return features, rows.Err()
}

func getTestimonials(db *sql.DB) ([]Testimonial, error) {
	rows, err := db.QueryContext(context.Background(), "SELECT name, company, text, avatar FROM testimonials ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testimonials []Testimonial
	for rows.Next() {
		var t Testimonial
		err := rows.Scan(&t.Name, &t.Company, &t.Text, &t.Avatar)
		if err != nil {
			return nil, err
		}
		testimonials = append(testimonials, t)
	}
	return testimonials, rows.Err()
}

func runTailwind(npmCmd string) {
	cmd := exec.Command(npmCmd, "run", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Starting Tailwind CSS watcher...")
	if err := cmd.Run(); err != nil {
		log.Printf("Tailwind CSS error: %v", err)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

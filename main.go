package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
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

		data := fiber.Map{
			"Title":   "Beranda - " + companyData["CompanyName"].(string),
			"Company": companyData,
			"Active":  "home",
			"Features": []map[string]string{
				{
					"icon":  "💻",
					"title": "Web Development",
					"desc":  "Kami membuat website yang responsif dan modern.",
				},
				{
					"icon":  "📱",
					"title": "Mobile Apps",
					"desc":  "Aplikasi mobile untuk iOS dan Android.",
				},
				{
					"icon":  "☁️",
					"title": "Cloud Solutions",
					"desc":  "Solusi cloud untuk bisnis Anda.",
				},
				{
					"icon":  "🔒",
					"title": "Cybersecurity",
					"desc":  "Melindungi data dan sistem Anda.",
				},
			},
			"Testimonials": []map[string]string{
				{
					"name":    "Budi Santoso",
					"company": "ABC Corporation",
					"text":    "Pelayanan sangat memuaskan, website kami jadi lebih modern.",
					"avatar":  "👨",
				},
				{
					"name":    "Sari Dewi",
					"company": "XYZ Enterprises",
					"text":    "Tim yang profesional dan hasil kerja berkualitas tinggi.",
					"avatar":  "👩",
				},
			},
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

	log.Println("Server berjalan di http://localhost:3000")
	log.Println("Cek file template di: http://localhost:3000/check-static")
	log.Fatal(app.Listen(":3000"))
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

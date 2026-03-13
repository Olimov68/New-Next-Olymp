package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/config"
	"github.com/nextolympservice/go-backend/internal/middleware"

	// Existing modules
	"github.com/nextolympservice/go-backend/internal/modules/auth"
	"github.com/nextolympservice/go-backend/internal/modules/telegram"
	"github.com/nextolympservice/go-backend/internal/modules/user"

	// Panel auth
	panelauth "github.com/nextolympservice/go-backend/internal/panel/auth"

	// User modules
	userbalance "github.com/nextolympservice/go-backend/internal/user/balance"
	usercerts "github.com/nextolympservice/go-backend/internal/user/certificates"
	userexams "github.com/nextolympservice/go-backend/internal/user/exams"
	userfeedback "github.com/nextolympservice/go-backend/internal/user/feedback"
	usermocktests "github.com/nextolympservice/go-backend/internal/user/mocktests"
	usernews "github.com/nextolympservice/go-backend/internal/user/news"
	usernotifs "github.com/nextolympservice/go-backend/internal/user/notifications"
	userolympiads "github.com/nextolympservice/go-backend/internal/user/olympiads"
	userresults "github.com/nextolympservice/go-backend/internal/user/results"

	// Admin centralized routes
	adminroutes "github.com/nextolympservice/go-backend/internal/admin/routes"

	// Superadmin centralized routes
	saroutes "github.com/nextolympservice/go-backend/internal/superadmin/routes"

	"github.com/nextolympservice/go-backend/internal/utils"
	"github.com/nextolympservice/go-backend/pkg/response"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.Static("/uploads", cfg.Upload.Dir)

	// JWT managers
	jwtManager := utils.NewJWTManager(&cfg.JWT)
	panelJWT := utils.NewPanelJWTManager(&cfg.PanelJWT)

	// ─── Existing modules ─────────────────────────────────────────────
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, jwtManager)
	authHandler := auth.NewHandler(authService)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, &cfg.Upload)
	userHandler := user.NewHandler(userService)

	telegramRepo := telegram.NewRepository(db)
	telegramService := telegram.NewService(telegramRepo, &cfg.Telegram)
	telegramHandler := telegram.NewHandler(telegramService)

	// ─── Panel auth ───────────────────────────────────────────────────
	panelAuthRepo := panelauth.NewRepository(db)
	panelAuthService := panelauth.NewService(panelAuthRepo, panelJWT)
	panelAuthHandler := panelauth.NewHandler(panelAuthService)

	// ─── User modules ─────────────────────────────────────────────────
	olympiadsHandler := userolympiads.NewHandler(userolympiads.NewService(userolympiads.NewRepository(db)))
	mockTestsHandler := usermocktests.NewHandler(usermocktests.NewService(usermocktests.NewRepository(db)))
	newsHandler := usernews.NewHandler(usernews.NewService(usernews.NewRepository(db)))
	certsHandler := usercerts.NewHandler(usercerts.NewService(usercerts.NewRepository(db)))
	feedbackHandler := userfeedback.NewHandler(userfeedback.NewService(userfeedback.NewRepository(db)))
	examsHandler := userexams.NewHandler(db)
	balanceHandler := userbalance.NewHandler(db)
	notifsHandler := usernotifs.NewHandler(db)
	userResultsHandler := userresults.NewHandler(db)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, 200, "OK", gin.H{"status": "healthy"})
	})

	api := r.Group("/api/v1")

	// ============================================================
	// MAVJUD USER AUTH ROUTES (o'zgarmadi)
	// ============================================================
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthRequired(jwtManager, db))
	{
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/auth/me", authHandler.Me)

		profileGroup := protected.Group("/profile")
		{
			profileGroup.POST("/complete", userHandler.CompleteProfile)
			profileGroup.PUT("/me", userHandler.UpdateProfile)
			profileGroup.GET("/me", userHandler.GetProfile)
			profileGroup.POST("/photo", userHandler.UploadPhoto)
		}

		telegramGroup := protected.Group("/telegram")
		{
			telegramGroup.POST("/verify", telegramHandler.VerifyCode)
			telegramGroup.GET("/status", telegramHandler.CheckStatus)
		}
	}
	api.POST("/telegram/webhook", telegramHandler.Webhook)

	// ============================================================
	// USER ROUTES — /api/v1/user/...
	// ============================================================
	userAPI := api.Group("/user")
	userAPI.Use(middleware.AuthRequired(jwtManager, db))
	userAPI.Use(middleware.ProfileRequired())
	{
		// Olympiads
		og := userAPI.Group("/olympiads")
		{
			og.GET("", olympiadsHandler.List)
			og.GET("/my", olympiadsHandler.MyOlympiads)
			og.GET("/:id", olympiadsHandler.GetByID)
			og.POST("/:id/join", olympiadsHandler.Join)
		}

		// Mock tests
		mg := userAPI.Group("/mock-tests")
		{
			mg.GET("", mockTestsHandler.List)
			mg.GET("/my", mockTestsHandler.MyMockTests)
			mg.GET("/:id", mockTestsHandler.GetByID)
			mg.POST("/:id/join", mockTestsHandler.Join)
		}

		// News
		ng := userAPI.Group("/news")
		{
			ng.GET("", newsHandler.List)
			ng.GET("/:id", newsHandler.GetByID)
		}
		// Announcements
		ag := userAPI.Group("/announcements")
		{
			ag.GET("", newsHandler.List)
			ag.GET("/:id", newsHandler.GetByID)
		}

		// Certificates
		cg := userAPI.Group("/certificates")
		{
			cg.GET("", certsHandler.List)
			cg.GET("/:id", certsHandler.GetByID)
		}

		// Feedback
		fg := userAPI.Group("/feedback")
		{
			fg.POST("", feedbackHandler.Create)
			fg.GET("", feedbackHandler.List)
			fg.GET("/:id", feedbackHandler.GetByID)
		}

		// Exam — test topshirish
		eg := userAPI.Group("/exams")
		{
			// Mock test topshirish
			eg.POST("/mock-tests/:id/start", examsHandler.StartMockTest)
			eg.POST("/mock-tests/attempts/:attempt_id/answer", examsHandler.SubmitMockAnswer)
			eg.POST("/mock-tests/attempts/:attempt_id/finish", examsHandler.FinishMockTest)
			eg.GET("/mock-tests/attempts/:attempt_id/result", examsHandler.GetMockAttemptResult)
			eg.GET("/mock-tests/:id/my-attempts", examsHandler.GetMyMockAttempts)
			// Olympiad topshirish
			eg.POST("/olympiads/:id/start", examsHandler.StartOlympiad)
			eg.POST("/olympiads/attempts/:attempt_id/answer", examsHandler.SubmitOlympiadAnswer)
			eg.POST("/olympiads/attempts/:attempt_id/finish", examsHandler.FinishOlympiad)
			eg.GET("/olympiads/attempts/:attempt_id/result", examsHandler.GetOlympiadAttemptResult)
		}

		// Balance
		bg := userAPI.Group("/balance")
		{
			bg.GET("", balanceHandler.GetBalance)
			bg.GET("/transactions", balanceHandler.GetTransactions)
			bg.POST("/topup", balanceHandler.TopUp)
		}

		// Notifications
		ntfg := userAPI.Group("/notifications")
		{
			ntfg.GET("", notifsHandler.List)
			ntfg.GET("/unread-count", notifsHandler.UnreadCount)
			ntfg.PATCH("/:id/read", notifsHandler.MarkAsRead)
			ntfg.POST("/read-all", notifsHandler.MarkAllAsRead)
			ntfg.DELETE("/:id", notifsHandler.Delete)
		}

		// Results
		resg := userAPI.Group("/results")
		{
			resg.GET("", userResultsHandler.GetMyResults)
			resg.GET("/mock-tests", userResultsHandler.GetMockTestResults)
			resg.GET("/olympiads", userResultsHandler.GetOlympiadResults)
		}
	}

	// ============================================================
	// PANEL AUTH ROUTES — /api/v1/panel/auth/...
	// ============================================================
	panelPub := api.Group("/panel/auth")
	{
		panelPub.POST("/login", panelAuthHandler.Login)
		panelPub.POST("/refresh", panelAuthHandler.RefreshToken)
	}

	panelProt := api.Group("/panel/auth")
	panelProt.Use(middleware.PanelAuthRequired(panelJWT, db))
	{
		panelProt.GET("/me", panelAuthHandler.Me)
		panelProt.GET("/permissions", panelAuthHandler.Permissions)
		panelProt.POST("/logout", panelAuthHandler.Logout)
	}

	// ============================================================
	// ADMIN ROUTES — /api/v1/admin/...
	// admin + superadmin kira oladi (centralized)
	// ============================================================
	adminroutes.Register(api, panelJWT, db, cfg)

	// ============================================================
	// SUPERADMIN ROUTES — /api/v1/superadmin/...
	// faqat superadmin kiradi (centralized)
	// ============================================================
	saroutes.Register(api, panelJWT, db)

	return r
}

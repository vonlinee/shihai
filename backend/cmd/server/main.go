package main

import (
	"flag"
	"log"
	"shihai/internal/config"
	"shihai/internal/handlers"
	"shihai/internal/middleware"
	"shihai/internal/repository"
	"shihai/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// App 应用容器
type App struct {
	db                  *gorm.DB
	userHandler         *handlers.UserHandler
	poemHandler         *handlers.PoemHandler
	commentHandler      *handlers.CommentHandler
	announcementHandler *handlers.AnnouncementHandler
	rbacHandler         *handlers.RBACHandler
	rbacMiddleware      *middleware.RBACMiddleware
}

func main() {
	// Parse command-line flags
	configFile := flag.String("config", "", "path to JSON config file (default: config.json)")
	flag.Parse()

	// Load configuration
	cfg := config.Load(*configFile)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	db, err := config.InitDB(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize application
	app := initApp(db)

	// Initialize default roles and permissions
	if err := initRBACData(app); err != nil {
		log.Println("Warning: Failed to initialize RBAC data:", err)
	}

	// Create Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup routes
	setupRoutes(r, app)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initApp 初始化应用依赖
func initApp(db *gorm.DB) *App {
	// Repository layer
	userRepo := repository.NewUserRepository(db)
	poemRepo := repository.NewPoemRepository(db)
	dynastyRepo := repository.NewDynastyRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	rolePermissionRepo := repository.NewRolePermissionRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)

	// Service layer
	userService := services.NewUserService(userRepo)
	poemService := services.NewPoemService(poemRepo, dynastyRepo)
	commentService := services.NewCommentService(commentRepo)
	announcementService := services.NewAnnouncementService(announcementRepo)
	rbacService := services.NewRBACService(roleRepo, permissionRepo, rolePermissionRepo, userRoleRepo)

	// Handler layer
	userHandler := handlers.NewUserHandler(userService)
	poemHandler := handlers.NewPoemHandler(poemService)
	commentHandler := handlers.NewCommentHandler(commentService)
	announcementHandler := handlers.NewAnnouncementHandler(announcementService)
	rbacHandler := handlers.NewRBACHandler(rbacService)

	// Middleware
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService)

	return &App{
		db:                  db,
		userHandler:         userHandler,
		poemHandler:         poemHandler,
		commentHandler:      commentHandler,
		announcementHandler: announcementHandler,
		rbacHandler:         rbacHandler,
		rbacMiddleware:      rbacMiddleware,
	}
}

// initRBACData 初始化RBAC默认数据
func initRBACData(app *App) error {
	rbacService := services.NewRBACService(
		repository.NewRoleRepository(app.db),
		repository.NewPermissionRepository(app.db),
		repository.NewRolePermissionRepository(app.db),
		repository.NewUserRoleRepository(app.db),
	)
	return rbacService.InitDefaultRolesAndPermissions()
}

// setupRoutes 设置路由
func setupRoutes(r *gin.Engine, app *App) {
	api := r.Group("/api")
	{
		// ========== Public Routes ==========
		// Auth
		api.POST("/auth/register", app.userHandler.Register)
		api.POST("/auth/login", app.userHandler.Login)

		// Poems - Public
		api.GET("/poems", app.poemHandler.GetPoemList)
		api.GET("/poems/random", app.poemHandler.GetRandomPoems)
		api.GET("/poems/:id", app.poemHandler.GetPoemByID)
		api.POST("/poems/:id/like", app.poemHandler.LikePoem)
		api.GET("/dynasties", app.poemHandler.GetDynastyList)

		// Announcements - Public
		api.GET("/announcements", app.announcementHandler.GetAnnouncements)
		api.GET("/announcements/:id", app.announcementHandler.GetAnnouncementByID)

		// Comments - Public (list), Protected (create, delete)
		api.GET("/comments", app.commentHandler.GetComments)
		api.POST("/comments/vote", app.commentHandler.VoteComment)

		// ========== Protected Routes ==========
		authorized := api.Group("/")
		authorized.Use(middleware.Auth())
		{
			// User
			authorized.GET("/user/profile", app.userHandler.GetProfile)
			authorized.PUT("/user/profile", app.userHandler.UpdateProfile)
			authorized.PUT("/user/password", app.userHandler.ChangePassword)

			// Comments - Protected
			authorized.POST("/comments", app.commentHandler.CreateComment)
			authorized.DELETE("/comments/:id", app.commentHandler.DeleteComment)
		}

		// ========== RBAC Routes ==========
		// 角色管理（需要role:list权限）
		api.GET("/rbac/roles", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:list"), app.rbacHandler.GetRoleList)
		api.GET("/rbac/roles/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:read"), app.rbacHandler.GetRoleByID)
		api.POST("/rbac/roles", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:create"), app.rbacHandler.CreateRole)
		api.PUT("/rbac/roles/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:update"), app.rbacHandler.UpdateRole)
		api.DELETE("/rbac/roles/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:delete"), app.rbacHandler.DeleteRole)
		api.GET("/rbac/roles/:id/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:read"), app.rbacHandler.GetRolePermissions)
		api.PUT("/rbac/roles/:id/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:assign"), app.rbacHandler.AssignPermissionsToRole)

		// 权限管理（需要permission:list权限）
		api.GET("/rbac/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission("permission:list"), app.rbacHandler.GetPermissionList)
		api.GET("/rbac/permissions/all", middleware.Auth(), app.rbacMiddleware.RequirePermission("permission:list"), app.rbacHandler.GetAllPermissions)
		api.GET("/rbac/permissions/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission("permission:read"), app.rbacHandler.GetPermissionByID)
		api.POST("/rbac/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission("permission:create"), app.rbacHandler.CreatePermission)
		api.PUT("/rbac/permissions/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission("permission:update"), app.rbacHandler.UpdatePermission)
		api.DELETE("/rbac/permissions/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission("permission:delete"), app.rbacHandler.DeletePermission)

		// 用户角色管理（需要role:assign权限）
		api.GET("/rbac/users/:id/roles", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:read"), app.rbacHandler.GetUserRoles)
		api.PUT("/rbac/users/:id/roles", middleware.Auth(), app.rbacMiddleware.RequirePermission("role:assign"), app.rbacHandler.AssignRolesToUser)
		api.GET("/rbac/users/:id/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission("permission:read"), app.rbacHandler.GetUserPermissions)

		// 当前用户权限查询
		api.GET("/rbac/my/permissions", middleware.Auth(), app.rbacHandler.GetMyPermissions)
		api.POST("/rbac/check", middleware.Auth(), app.rbacHandler.CheckUserPermission)

		// ========== Admin Routes ==========
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(), app.rbacMiddleware.RequireAnyRole("admin", "editor"))
		{
			// Users
			admin.GET("/users", app.rbacMiddleware.RequirePermission("user:list"), app.userHandler.GetUserList)
			admin.GET("/users/:id", app.rbacMiddleware.RequirePermission("user:read"), app.userHandler.GetUserByID)
			admin.DELETE("/users/:id", app.rbacMiddleware.RequirePermission("user:delete"), app.userHandler.DeleteUser)

			// Poems
			admin.POST("/poems", app.rbacMiddleware.RequirePermission("poem:create"), app.poemHandler.CreatePoem)
			admin.PUT("/poems/:id", app.rbacMiddleware.RequirePermission("poem:update"), app.poemHandler.UpdatePoem)
			admin.DELETE("/poems/:id", app.rbacMiddleware.RequirePermission("poem:delete"), app.poemHandler.DeletePoem)

			// Announcements
			admin.POST("/announcements", app.rbacMiddleware.RequirePermission("announcement:create"), app.announcementHandler.CreateAnnouncement)
			admin.PUT("/announcements/:id", app.rbacMiddleware.RequirePermission("announcement:update"), app.announcementHandler.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", app.rbacMiddleware.RequirePermission("announcement:delete"), app.announcementHandler.DeleteAnnouncement)

			// Comments Admin
			admin.GET("/comments/all", app.rbacMiddleware.RequirePermission("comment:list"), app.commentHandler.GetComments)
		}
	}
}

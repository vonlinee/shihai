package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"shihai/internal/apidocs"
	"shihai/internal/config"
	"shihai/internal/database"
	"shihai/internal/handlers"
	"shihai/internal/middleware"
	"shihai/internal/models"
	"shihai/internal/repository"
	"shihai/internal/services"
	"shihai/internal/webui"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// App 应用容器
type App struct {
	db                    *gorm.DB
	userHandler           *handlers.UserHandler
	poemHandler           *handlers.PoemHandler
	commentHandler        *handlers.CommentHandler
	forumHandler          *handlers.ForumHandler
	correctionHandler     *handlers.CorrectionHandler
	announcementHandler   *handlers.AnnouncementHandler
	workCollectionHandler *handlers.WorkCollectionHandler
	textConversionHandler *handlers.TextConversionHandler
	rbacHandler           *handlers.RBACHandler
	rbacMiddleware        *middleware.RBACMiddleware
}

func main() {
	// Parse command-line flags
	configFile := flag.String("config", "", "path to JSON config file (default: config.json)")
	port := flag.String("port", "", "server listen port (overrides SERVER_PORT and config file)")
	flag.Parse()

	// Load configuration
	cfg := config.Load(*configFile)
	if err := applyServerPortOverride(cfg, *port); err != nil {
		log.Fatal(err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	db, err := config.InitDB(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate Database Model
	err = config.AutoMigrateDatabaseModel(db, err)
	if err != nil {
		log.Fatal("Failed to auto migrate database model:", err)
	}

	// Initialize application
	app := initApp(db)

	database.Init(*configFile)

	// Initialize default roles and permissions
	if err := initRBACData(app); err != nil {
		log.Println("Warning: Failed to initialize RBAC data:", err)
		return
	}

	// Create Gin router
	r := gin.Default()
	healthHandler := handlers.NewHealthHandler()

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
	r.GET("/health", healthHandler.Check)

	// Setup routes
	setupRoutes(r, app)
	apidocs.RegisterRoutes(r)
	if err := webui.RegisterRoutes(r); err != nil {
		log.Fatal("Failed to register frontend routes:", err)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func applyServerPortOverride(cfg *config.Config, port string) error {
	if port == "" {
		return nil
	}

	parsedPort, err := strconv.Atoi(port)
	if err != nil || parsedPort < 1 || parsedPort > 65535 {
		return fmt.Errorf("invalid server port %q: must be an integer between 1 and 65535", port)
	}

	cfg.Server.Port = strconv.Itoa(parsedPort)
	return nil
}

// initApp 初始化应用依赖
func initApp(db *gorm.DB) *App {
	// Repository layer
	userRepo := repository.NewUserRepository(db)
	poemRepo := repository.NewPoemRepository(db)
	poemAnnotationRepo := repository.NewPoemAnnotationRepository(db)
	dynastyRepo := repository.NewDynastyRepository(db)
	authorRepo := repository.NewAuthorRepository(db)
	poetRepo := repository.NewPoetRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	forumRepo := repository.NewForumRepository(db)
	correctionRepo := repository.NewCorrectionRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)
	workCollectionRepo := repository.NewWorkCollectionRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	rolePermissionRepo := repository.NewRolePermissionRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	// Service layer
	userService := services.NewUserService(userRepo, roleRepo, userRoleRepo)
	poemService := services.NewPoemService(poemRepo, dynastyRepo, authorRepo, poetRepo)
	poemAnnotationService := services.NewPoemAnnotationService(poemAnnotationRepo, poemRepo)
	poemService.SetPoemAnnotationRepository(poemAnnotationRepo)
	commentService := services.NewCommentService(commentRepo)
	forumService := services.NewForumService(forumRepo)
	correctionService := services.NewCorrectionService(correctionRepo)
	announcementService := services.NewAnnouncementService(announcementRepo)
	workCollectionService := services.NewWorkCollectionService(workCollectionRepo, poemRepo)
	textConversionService := services.NewTextConversionService()
	rbacService := services.NewRBACService(roleRepo, rolePermissionRepo, userRoleRepo, permRepo)

	// Handler layer
	userHandler := handlers.NewUserHandler(userService)
	poemHandler := handlers.NewPoemHandler(poemService, poemAnnotationService)
	commentHandler := handlers.NewCommentHandler(commentService)
	forumHandler := handlers.NewForumHandler(forumService)
	correctionHandler := handlers.NewCorrectionHandler(correctionService)
	announcementHandler := handlers.NewAnnouncementHandler(announcementService)
	workCollectionHandler := handlers.NewWorkCollectionHandler(workCollectionService)
	textConversionHandler := handlers.NewTextConversionHandler(textConversionService)
	rbacHandler := handlers.NewRBACHandler(rbacService)

	// Middleware
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService)

	return &App{
		db:                    db,
		userHandler:           userHandler,
		poemHandler:           poemHandler,
		commentHandler:        commentHandler,
		forumHandler:          forumHandler,
		correctionHandler:     correctionHandler,
		announcementHandler:   announcementHandler,
		workCollectionHandler: workCollectionHandler,
		textConversionHandler: textConversionHandler,
		rbacHandler:           rbacHandler,
		rbacMiddleware:        rbacMiddleware,
	}
}

// initRBACData 初始化RBAC默认数据和超级管理员
func initRBACData(app *App) error {
	rbacService := services.NewRBACService(
		repository.NewRoleRepository(app.db),
		repository.NewRolePermissionRepository(app.db),
		repository.NewUserRoleRepository(app.db),
		repository.NewPermissionRepository(app.db),
	)

	// 1. 初始化默认角色和权限
	if err := rbacService.InitDefaultRolesAndPermissions(); err != nil {
		log.Println("Warning: Failed to initialize RBAC roles/permissions:", err)
	}

	// 2. 创建默认超级管理员账号（如不存在）
	userRepo := repository.NewUserRepository(app.db)
	userRoleRepo := repository.NewUserRoleRepository(app.db)
	roleRepo := repository.NewRoleRepository(app.db)

	adminUsername := "root"
	adminPassword := "shihai2024"

	_, err := userRepo.GetByUsername(adminUsername)
	if err != nil {
		// root 用户不存在，创建
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		rootUser := &models.User{
			Username: adminUsername,
			Password: string(hashedPassword),
			Name:     "超级管理员",
			IsActive: true,
		}
		if err := userRepo.Create(rootUser); err != nil {
			return err
		}

		// 分配 admin 角色
		adminRole, err := roleRepo.GetByName("admin")
		if err == nil {
			err := userRoleRepo.Create(&models.UserRole{
				UserID: rootUser.ID,
				RoleID: adminRole.ID,
			})
			if err != nil {
				return err
			}
		}

		log.Printf("Default admin user created: %s / %s", adminUsername, adminPassword)
	}

	return nil
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
		api.GET("/poets", app.poemHandler.GetPoetList)
		api.GET("/genres", app.poemHandler.GetGenreList)

		// Announcements - Public
		api.GET("/announcements", app.announcementHandler.GetAnnouncements)
		api.GET("/announcements/:id", app.announcementHandler.GetAnnouncementByID)

		// Comments - Public (list), Protected (create, delete)
		api.GET("/comments", app.commentHandler.GetComments)
		api.POST("/comments/vote", app.commentHandler.VoteComment)

		// Forum - Public read
		api.GET("/forum/posts", app.forumHandler.ListPosts)
		api.GET("/forum/posts/:id", app.forumHandler.GetPostByID)
		api.GET("/forum/posts/:id/replies", app.forumHandler.ListReplies)

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

			// Forum - Protected write
			authorized.POST("/forum/posts", app.rbacMiddleware.RequirePermission(models.PermForumCreate), app.forumHandler.CreatePost)
			authorized.PUT("/forum/posts/:id", app.forumHandler.UpdatePost)
			authorized.DELETE("/forum/posts/:id", app.forumHandler.DeletePost)
			authorized.POST("/forum/posts/:id/replies", app.rbacMiddleware.RequirePermission(models.PermForumCreate), app.forumHandler.CreateReply)
			authorized.DELETE("/forum/replies/:id", app.forumHandler.DeleteReply)
		}

		// ========== RBAC Routes ==========
		// 角色管理（需要 role:* 系列权限，权限编码定义见 models/permission_codes.go）
		api.GET("/rbac/roles", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleList), app.rbacHandler.GetRoleList)
		api.GET("/rbac/roles/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleRead), app.rbacHandler.GetRoleByID)
		api.POST("/rbac/roles", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleCreate), app.rbacHandler.CreateRole)
		api.PUT("/rbac/roles/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleUpdate), app.rbacHandler.UpdateRole)
		api.DELETE("/rbac/roles/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleDelete), app.rbacHandler.DeleteRole)
		api.GET("/rbac/roles/:id/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleRead), app.rbacHandler.GetRolePermissions)
		api.PUT("/rbac/roles/:id/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleAssign), app.rbacHandler.AssignPermissionsToRole)

		// 用户角色管理（需要 role:assign 权限）
		api.GET("/rbac/users/:id/roles", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleRead), app.rbacHandler.GetUserRoles)
		api.PUT("/rbac/users/:id/roles", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermRoleAssign), app.rbacHandler.AssignRolesToUser)
		api.GET("/rbac/users/:id/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermPermissionRead), app.rbacHandler.GetUserPermissions)

		// 当前用户权限查询（仅需登录即可，不需额外权限）
		api.GET("/rbac/my/permissions", middleware.Auth(), app.rbacHandler.GetMyPermissions)
		api.POST("/rbac/check", middleware.Auth(), app.rbacHandler.CheckUserPermission)

		// 权限管理
		api.GET("/rbac/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermPermissionList), app.rbacHandler.GetPermissionList)
		api.GET("/rbac/permissions/all", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermPermissionList), app.rbacHandler.GetAllPermissions)
		api.GET("/rbac/permissions/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermPermissionRead), app.rbacHandler.GetPermissionByID)
		api.POST("/rbac/permissions", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermPermissionCreate), app.rbacHandler.CreatePermission)
		api.PUT("/rbac/permissions/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermPermissionUpdate), app.rbacHandler.UpdatePermission)
		api.DELETE("/rbac/permissions/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermPermissionDelete), app.rbacHandler.DeletePermission)

		// ========== Admin Routes ==========
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(), app.rbacMiddleware.RequireAnyRole(models.RoleAdmin, models.RoleEditor))
		{
			// Users
			admin.GET("/users", app.rbacMiddleware.RequirePermission(models.PermUserList), app.userHandler.GetUserList)
			admin.GET("/users/:id", app.rbacMiddleware.RequirePermission(models.PermUserRead), app.userHandler.GetUserByID)
			admin.POST("/users", app.rbacMiddleware.RequirePermission(models.PermUserCreate), app.userHandler.AdminCreateUser)
			admin.DELETE("/users/:id", app.rbacMiddleware.RequirePermission(models.PermUserDelete), app.userHandler.DeleteUser)

			// Poems
			admin.POST("/poems", app.rbacMiddleware.RequirePermission(models.PermPoemCreate), app.poemHandler.CreatePoem)
			admin.PUT("/poems/:id", app.rbacMiddleware.RequirePermission(models.PermPoemUpdate), app.poemHandler.UpdatePoem)
			admin.DELETE("/poems", app.rbacMiddleware.RequirePermission(models.PermPoemDelete), app.poemHandler.BatchDeletePoems)
			admin.DELETE("/poems/:id", app.rbacMiddleware.RequirePermission(models.PermPoemDelete), app.poemHandler.DeletePoem)
			admin.GET("/poems/:id/annotations", app.rbacMiddleware.RequirePermission(models.PermPoemRead), app.poemHandler.GetPoemAnnotations)
			admin.POST("/poems/:id/annotations", app.rbacMiddleware.RequirePermission(models.PermPoemUpdate), app.poemHandler.CreatePoemAnnotation)
			admin.PUT("/poem-annotations/:id", app.rbacMiddleware.RequirePermission(models.PermPoemUpdate), app.poemHandler.UpdatePoemAnnotation)
			admin.DELETE("/poem-annotations/:id", app.rbacMiddleware.RequirePermission(models.PermPoemUpdate), app.poemHandler.DeletePoemAnnotation)
			admin.POST("/text-conversion", app.rbacMiddleware.RequirePermission(models.PermPoemUpdate), app.textConversionHandler.ConvertTexts)

			// Dynasties & Poets（复用 poem:* 权限，不独立划分权限点）
			admin.POST("/dynasties", app.rbacMiddleware.RequirePermission(models.PermPoemCreate), app.poemHandler.CreateDynasty)
			admin.PUT("/dynasties/:id", app.rbacMiddleware.RequirePermission(models.PermPoemUpdate), app.poemHandler.UpdateDynasty)
			admin.DELETE("/dynasties", app.rbacMiddleware.RequirePermission(models.PermPoemDelete), app.poemHandler.BatchDeleteDynasties)
			admin.DELETE("/dynasties/:id", app.rbacMiddleware.RequirePermission(models.PermPoemDelete), app.poemHandler.DeleteDynasty)
			admin.POST("/poets", app.rbacMiddleware.RequirePermission(models.PermPoemCreate), app.poemHandler.CreatePoet)
			admin.PUT("/poets/:id", app.rbacMiddleware.RequirePermission(models.PermPoemUpdate), app.poemHandler.UpdatePoet)
			admin.DELETE("/poets", app.rbacMiddleware.RequirePermission(models.PermPoemDelete), app.poemHandler.BatchDeletePoets)
			admin.DELETE("/poets/:id", app.rbacMiddleware.RequirePermission(models.PermPoemDelete), app.poemHandler.DeletePoet)

			// Work Collections
			admin.GET("/work-collections", app.rbacMiddleware.RequirePermission(models.PermWorkCollectionList), app.workCollectionHandler.GetWorkCollections)
			admin.GET("/work-collections/:id", app.rbacMiddleware.RequirePermission(models.PermWorkCollectionRead), app.workCollectionHandler.GetWorkCollectionByID)
			admin.POST("/work-collections", app.rbacMiddleware.RequirePermission(models.PermWorkCollectionCreate), app.workCollectionHandler.CreateWorkCollection)
			admin.PUT("/work-collections/:id", app.rbacMiddleware.RequirePermission(models.PermWorkCollectionUpdate), app.workCollectionHandler.UpdateWorkCollection)
			admin.DELETE("/work-collections/:id", app.rbacMiddleware.RequirePermission(models.PermWorkCollectionDelete), app.workCollectionHandler.DeleteWorkCollection)
			admin.POST("/work-collections/:id/items", app.rbacMiddleware.RequirePermission(models.PermWorkCollectionItemManage), app.workCollectionHandler.AddWorkCollectionItem)
			admin.PUT("/work-collections/:id/items/:itemId", app.rbacMiddleware.RequirePermission(models.PermWorkCollectionItemManage), app.workCollectionHandler.UpdateWorkCollectionItem)
			admin.DELETE("/work-collections/:id/items/:itemId", app.rbacMiddleware.RequirePermission(models.PermWorkCollectionItemManage), app.workCollectionHandler.DeleteWorkCollectionItem)

			// Announcements
			admin.POST("/announcements", app.rbacMiddleware.RequirePermission(models.PermAnnouncementCreate), app.announcementHandler.CreateAnnouncement)
			admin.PUT("/announcements/:id", app.rbacMiddleware.RequirePermission(models.PermAnnouncementUpdate), app.announcementHandler.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", app.rbacMiddleware.RequirePermission(models.PermAnnouncementDelete), app.announcementHandler.DeleteAnnouncement)

			// Comments Admin
			admin.GET("/comments/all", app.rbacMiddleware.RequirePermission(models.PermCommentList), app.commentHandler.GetAllComments)

			// Corrections Admin
			admin.GET("/corrections", app.rbacMiddleware.RequirePermission(models.PermCorrectionList), app.correctionHandler.ListCorrections)
		}

		// Forum Admin（forum:moderate 权限可独立授予 reviewer，不套用 admin/editor 角色门槛）
		api.GET("/admin/forum/posts", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermForumModerate), app.forumHandler.ListAllPosts)
		api.PUT("/admin/forum/posts/:id/pin", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermForumModerate), app.forumHandler.SetPostPinned)
		api.DELETE("/admin/forum/posts/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermForumModerate), app.forumHandler.AdminDeletePost)
		api.DELETE("/admin/forum/replies/:id", middleware.Auth(), app.rbacMiddleware.RequirePermission(models.PermForumModerate), app.forumHandler.AdminDeleteReply)
	}
}

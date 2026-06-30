package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/topmentor/backend/internal/config"
	"github.com/topmentor/backend/internal/handler"
	"github.com/topmentor/backend/internal/middleware"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/internal/service"
	"github.com/topmentor/backend/pkg/agora"
	"github.com/topmentor/backend/pkg/database"
	jwtpkg "github.com/topmentor/backend/pkg/jwt"
	"github.com/topmentor/backend/pkg/payment"
	"github.com/topmentor/backend/pkg/roomsession"
	redisclient "github.com/topmentor/backend/pkg/redis"
	"github.com/topmentor/backend/pkg/slotlock"
	smspkg "github.com/topmentor/backend/pkg/sms"
	"github.com/topmentor/backend/pkg/upload"
	wxpkg "github.com/topmentor/backend/pkg/wx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	pg, err := database.NewPostgres(cfg.DSN())
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pg.Close()

	rdb := redisclient.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer func() { _ = rdb.Close() }()

	jwtMgr := jwtpkg.NewManager(cfg.JWTSecret, cfg.JWTExpireHours)
	wxClient := wxpkg.NewClient(cfg.WxAppID, cfg.WxAppSecret, cfg.WxMockMode)
	smsStore := smspkg.NewOTPStore(rdb.Raw(), cfg.SMSMockMode, cfg.SMSMockCode)
	payFactory := payment.NewFactory(cfg.PaymentMerchantConfig())

	userRepo := repository.NewUserRepository(pg.DB())
	mentorRepo := repository.NewMentorRepository(pg.DB())
	packageRepo := repository.NewPackageRepository(pg.DB())
	rechargeRepo := repository.NewRechargeRepository(pg.DB())

	slotRepo := repository.NewSlotRepository(pg.DB())
	bookingRepo := repository.NewBookingRepository(pg.DB())
	slotLocker := slotlock.NewLocker(rdb.Raw())

	orderRepo := repository.NewOrderRepository(pg.DB())
	roomSession := roomsession.NewStore(rdb.Raw())
	agoraSvc := agora.NewTokenService(agora.Config{
		AppID:          cfg.AgoraAppID,
		AppCertificate: cfg.AgoraAppCertificate,
		MockMode:       cfg.AgoraMockMode,
	})

	reportRepo := repository.NewReportRepository(pg.DB())
	walletRepo := repository.NewWalletRepository(pg.DB())

	authSvc := service.NewAuthService(userRepo, mentorRepo, wxClient, jwtMgr, smsStore, cfg.SMSMockMode)
	userSvc := service.NewUserService(userRepo)
	mentorSvc := service.NewMentorService(mentorRepo)
	rechargeSvc := service.NewRechargeService(packageRepo, rechargeRepo, userRepo, payFactory)
	bookingSvc := service.NewBookingService(slotRepo, bookingRepo, slotLocker)
	roomSvc := service.NewRoomService(orderRepo, roomSession, agoraSvc, cfg.LessonDurationMin)
	mentorPortalSvc := service.NewMentorPortalService(
		mentorRepo, orderRepo, slotRepo, reportRepo, walletRepo,
		cfg.LessonEarnYuan, cfg.WithdrawMockMode,
	)
	mentorApplySvc := service.NewMentorApplyService(mentorRepo, userRepo)
	healthSvc := service.NewHealthService(pg, rdb)

	uploadStore, err := upload.NewStore(cfg.UploadDir)
	if err != nil {
		log.Fatalf("upload store: %v", err)
	}

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	mentorHandler := handler.NewMentorHandler(mentorSvc)
	rechargeHandler := handler.NewRechargeHandler(rechargeSvc)
	bookingHandler := handler.NewBookingHandler(bookingSvc)
	roomHandler := handler.NewRoomHandler(roomSvc)
	mentorPortalHandler := handler.NewMentorPortalHandler(mentorPortalSvc, uploadStore)
	mentorApplyHandler := handler.NewMentorApplyHandler(mentorApplySvc, uploadStore)
	adminRepo := repository.NewAdminRepository(pg.DB())
	adminSvc := service.NewAdminService(adminRepo, jwtMgr, cfg.AdminUsername, cfg.AdminPassword)
	adminHandler := handler.NewAdminHandler(adminSvc)
	healthHandler := handler.NewHealthHandler(healthSvc)

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger(), middleware.CORS(cfg.CORSOrigins))
	router.Static("/uploads", cfg.UploadDir)

	api := router.Group("/api/v1")
	api.GET("/health", healthHandler.Check)

	api.POST("/auth/sms/send", authHandler.SendSMS)
	api.POST("/auth/sms/login", authHandler.SmsLogin)
	api.POST("/auth/wx-login", authHandler.WxLogin)

	api.POST("/admin/auth/login", adminHandler.Login)

	api.GET("/mentors", mentorHandler.List)
	api.GET("/mentors/:id/slots", bookingHandler.ListMentorSlots)
	api.GET("/mentors/:id", mentorHandler.GetByID)

	api.GET("/packages", rechargeHandler.ListPackages)
	api.GET("/payment/channels", rechargeHandler.PaymentChannels)
	api.POST("/payment/notify/wechat", rechargeHandler.NotifyWechat)

	authed := api.Group("")
	authed.Use(middleware.Auth(jwtMgr))
	authed.PUT("/users/profile", userHandler.UpdateProfile)
	authed.GET("/users/lessons", rechargeHandler.GetLessonBalance)
	authed.POST("/recharge", rechargeHandler.CreateRecharge)
	authed.GET("/recharge/:id", rechargeHandler.GetRechargeOrder)
	authed.POST("/bookings", bookingHandler.CreateBooking)
	authed.GET("/users/orders", roomHandler.ListUserOrders)
	authed.GET("/orders/:id", roomHandler.GetOrder)
	authed.POST("/rooms/:orderId/join", roomHandler.Join)
	authed.POST("/rooms/:orderId/heartbeat", roomHandler.Heartbeat)
	authed.POST("/rooms/:orderId/complete", roomHandler.Complete)
	authed.GET("/reports/:orderId", mentorPortalHandler.GetReport)
	authed.GET("/mentor/orders", mentorPortalHandler.ListOrders)
	authed.GET("/mentor/slots", mentorPortalHandler.ListSlots)
	authed.PUT("/mentor/slots", mentorPortalHandler.SetSlots)
	authed.POST("/mentor/reports", mentorPortalHandler.SubmitReport)
	authed.GET("/mentor/wallet", mentorPortalHandler.Wallet)
	authed.POST("/mentor/withdraw", mentorPortalHandler.Withdraw)
	authed.GET("/mentor/profile", mentorPortalHandler.GetProfile)
	authed.PUT("/mentor/profile", mentorPortalHandler.UpdateProfile)
	authed.POST("/mentor/upload", mentorPortalHandler.Upload)
	authed.GET("/mentor/apply/status", mentorApplyHandler.GetStatus)
	authed.POST("/mentor/apply", mentorApplyHandler.Submit)
	authed.POST("/mentor/apply/upload", mentorApplyHandler.Upload)

	admin := api.Group("/admin")
	admin.Use(middleware.AdminAuth(jwtMgr))
	admin.GET("/mentors/pending", adminHandler.ListPendingMentors)
	admin.PUT("/mentors/:id/review", adminHandler.ReviewMentor)
	admin.GET("/courseware", adminHandler.ListCourseware)
	admin.POST("/courseware", adminHandler.CreateCourseware)
	admin.PUT("/courseware/:id", adminHandler.UpdateCourseware)
	admin.DELETE("/courseware/:id", adminHandler.DeleteCourseware)
	admin.GET("/finance/summary", adminHandler.FinanceSummary)

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	srv := &http.Server{Addr: addr, Handler: router}

	go func() {
		log.Printf("server listening on %s (payment=%s agora_mock=%v)", addr, cfg.PaymentMode, cfg.AgoraMockMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}

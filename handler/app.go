package handler

import (
	"finalProject2/infrastructure/config"
	"finalProject2/infrastructure/database"
	middlewares "finalProject2/pkg/middleware"
	"finalProject2/repository/comment_repository/comment_pg"
	"finalProject2/repository/photo_repository/photo_pg"
	social_media_pg "finalProject2/repository/socialMedia_repository/socialMedia_pg"
	"finalProject2/repository/user_repository/user_pg"
	"finalProject2/service"

	"github.com/gin-gonic/gin"
)

func StartApp() {
	config.LoadAppConfig()
	database.InitiliazeDatabase()
	db := database.GetDatabaseInstance()

	// user
	userRepo := user_pg.NewOrderPG(db)
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	// photo
	photoRepo := photo_pg.NewPhotoPG(db)
	photoService := service.NewPhotoService(photoRepo)
	photoHandler := NewPhotoHandler(photoService)

	// comment
	commentRepo := comment_pg.NewCommentPG(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := NewCommentHandler(commentService)

	//social media
	socialMediaRepo := social_media_pg.NewSocialMediaPG(db)
	socialMediaService := service.NewSocialMediaService(socialMediaRepo)
	socialMediaHandler := NewSocialMediaHandler(socialMediaService)

	r := gin.Default()

	//user
	r.POST("/users/register", userHandler.register)
	r.POST("/users/login", userHandler.Login)

	users := r.Group("/users")
	photos := r.Group("/photos")
	comments := r.Group("/comments")
	socialMedias := r.Group("/socialmedias")

	users.Use(middlewares.Authentication())
	{
		users.PUT("", userHandler.EditUser)
		users.DELETE("", userHandler.DeleteUser)
	}

	photos.Use(middlewares.Authentication())
	{
		photos.POST("", photoHandler.CreatePhoto)
		photos.GET("", photoHandler.GetPhotos)
		photos.PUT("/:photoId", photoHandler.EditPhoto)
		photos.DELETE("/:photoId", photoHandler.DeletePhoto)
	}

	comments.Use(middlewares.Authentication())
	{
		comments.POST("", commentHandler.CreateComment)
		comments.GET("", commentHandler.GetComments)
		comments.PUT("/:commentId", commentHandler.UpdateComment)
		comments.DELETE("/:commentId", commentHandler.DeleteComment)
	}

	socialMedias.Use(middlewares.Authentication())
	{
		socialMedias.POST("", socialMediaHandler.CreateSocialMedia)
		socialMedias.GET("", socialMediaHandler.GetSocialMedias)
		socialMedias.PUT("/:socialMediaId", socialMediaHandler.UpdateSocialMedia)
		socialMedias.DELETE("/:socialMediaId", socialMediaHandler.DeleteSocialMedia)
	}

	r.Run(":" + config.GetAppConfig().Port)

}

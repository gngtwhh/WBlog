package handler

// App contains all handlers
type App struct {
	Index   *IndexHandler
	Article *ArticleHandler
	User    *UserHandler
}

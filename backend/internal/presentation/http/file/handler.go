package file

// Handler はファイル関連のハンドラーをまとめた構造体です。
type Handler struct {
	GetDownloadURL *GetDownloadURLHandler
}

func NewHandler(getDownloadURL *GetDownloadURLHandler) *Handler {
	return &Handler{
		GetDownloadURL: getDownloadURL,
	}
}

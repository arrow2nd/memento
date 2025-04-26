package watcher

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/arrow2nd/memento/config"
	"github.com/arrow2nd/memento/logparser"
	"github.com/arrow2nd/memento/picture"
	"github.com/fsnotify/fsnotify"
	"github.com/sqweek/dialog"
)

type Watcher struct {
	watcher            *fsnotify.Watcher
	watcherMutex       sync.Mutex
	watchingSubDirName string
	config             *config.Config
}

// New: 作成
func New(config *config.Config) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher:            watcher,
		watchingSubDirName: getCurrentDate(),
		config:             config,
	}, nil
}

// Start: 監視を開始
func (w *Watcher) Start() {
	w.watcherMutex.Lock()
	defer w.watcherMutex.Unlock()
	defer w.watcher.Close()

	// 監視対象の設定
	w.Setup()

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// 作成イベントでないならスキップ
			if !event.Has(fsnotify.Create) {
				continue
			}

			fi, err := os.Stat(event.Name)
			if err != nil {
				continue
			}

			// 新しいディレクトリが作成された
			if fi.IsDir() {
				if newDirName, ok := isCurrentMonthDir(event.Name); ok {
					log.Println("新しいディレクトリを検出:", event.Name)

					if err := w.updateWatchingSubDir(newDirName); err != nil {
						dialog.Message("%s\n%s",
							"監視するフォルダの更新に失敗しました。",
							"アプリを再起動してください。",
						).Title("エラー").Error()
						log.Println("監視対象の更新に失敗:", err)
						continue
					}
				}

				continue
			}

			// 新しい写真が作成された
			if w.isVRCPicture(event.Name) {
				log.Println("新しい写真を検出:", event.Name)

				latestWorldVisit, err := logparser.FindLatestWorldVisit(w.config.VRCLogDirPath)
				if err != nil {
					log.Println("ログの解析に失敗:", err)
					continue
				}

				log.Println("ワールド訪問履歴を取得:", latestWorldVisit)

				// 写真をワールド名のディレクトリに移動
				err = picture.MoveToWorldNameDir(picture.MoveToWorldNameDirOpts{
					PicturePath:   event.Name,
					TargetDirPath: filepath.Join(w.config.PictureDirPath, w.watchingSubDirName),
					WorldVisit:    latestWorldVisit,
				}, w.config.ConvertToJpeg)

				if err != nil {
					log.Println("写真の移動に失敗:", err)
					continue
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

			log.Println("監視エラー:", err)
		}
	}
}

package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/arrow2nd/memento/config"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher            *fsnotify.Watcher
	watcherMutex       sync.Mutex
	watchingSubDirName string
	config             *config.Config
}

// New: 新しいWatcherを作成
func New(config *config.Config) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("watcherの作成に失敗: ", err)
	}

	return &Watcher{
		watcher:            watcher,
		watchingSubDirName: getCurrentDate(),
		config:             config,
	}
}

// Start: 監視を開始
func (w *Watcher) Start() {
	w.watcherMutex.Lock()
	defer w.watcherMutex.Unlock()
	defer w.watcher.Close()

	subDirPath := filepath.Join(w.config.RootDirPath, w.watchingSubDirName)

	// 監視対象に追加
	err := w.addWatchDir(w.config.RootDirPath, subDirPath)
	if err != nil {
		log.Fatal("監視対象の追加に失敗: ", err)
	}

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
					log.Println("新しいディレクトリを検出: ", event.Name)

					if err := w.updateWatchingSubDir(newDirName); err != nil {
						log.Println("監視対象の更新に失敗: ", err)
						continue
					}
				}

				continue
			}

			// 新しい写真が作成された
			if w.isVRCPicture(event.Name) {
				fmt.Println("新しい画像を検出: ", event.Name)
				// TODO: ここでログの読み取りと画像の移動を行う
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

			log.Println("監視エラー: ", err)
		}
	}
}

package watcher

import (
	"log"
	"os"
	"path/filepath"
)

// addWatchDir: 監視対象のディレクトリを追加
func (w *Watcher) addWatchDir(paths ...string) error {
	for _, p := range paths {
		// 存在しないならスキップ
		if _, err := os.Stat(p); os.IsNotExist(err) {
			log.Println("ディレクトリが存在しないため監視対象への追加をスキッ: ", p)
			continue
		}

		if err := w.watcher.Add(p); err != nil {
			return err
		}

		log.Println("監視対象に追加: ", p)
	}

	return nil
}

// updateWatchingSubDir: 監視対象のサブディレクトリを更新
func (w *Watcher) updateWatchingSubDir(newDirName string) error {
	newDirPath := filepath.Join(w.config.PictureDirPath, newDirName)

	// 監視対象に追加
	if err := w.addWatchDir(newDirPath); err != nil {
		return err
	}

	oldWatchingSubDirPath := filepath.Join(w.config.PictureDirPath, w.watchingSubDirName)

	// 監視対象のディレクトリが変更されていないなら、削除処理をスキップ
	if oldWatchingSubDirPath == newDirPath {
		// NOTE:
		// 次のような時にここに引っかかるはず
		// - 監視対象のディレクトリが無くてmementoが作成したとき
		// - 監視対象のディレクトリが削除されて、再作成されたとき
		log.Println("監視対象のディレクトリは変更されていません: ", oldWatchingSubDirPath)
		return nil
	}

	// 前の監視対象を削除
	if err := w.watcher.Remove(oldWatchingSubDirPath); err != nil {
		return err
	}

	log.Println("監視対象から削除: ", w.watchingSubDirName)

	// 新しい監視対象のサブディレクトリ名で更新
	w.watchingSubDirName = newDirName

	return nil
}

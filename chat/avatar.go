package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

//ErrNoAvatarURL はAvatarインスタンスがアバターのURLを返すことができないエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

//Avatar はゆーざーのプロフィール画像を表す型
type Avatar interface {
	//問題が発生したらエラーを返す
	AvatarURL(ChatUser) (string, error)
}

//TryAvatars Avatarオブジェクトのスライス
type TryAvatars []Avatar

//AvatarURL URLが呼び出し元に返る
func (a TryAvatars) AvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.AvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

//AuthAvatar AuthAvatar型
type AuthAvatar struct{}

//UseAuthAvatar 変数
var UseAuthAvatar AuthAvatar

//AvatarURL avatar_url値を探しそれを返す
func (AuthAvatar) AvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

//GravatarAvatar GravatarAvatar型
type GravatarAvatar struct{}

//UseGravatar 変数
var UseGravatar GravatarAvatar

//AvatarURL メールアドレス処理、ハッシュ値算出、URL埋め込み
func (GravatarAvatar) AvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

//FileSystemAvatar 型
type FileSystemAvatar struct{}

//UseFileSystemAvatar 変数
var UseFileSystemAvatar FileSystemAvatar

//AvatarURL urlidを元にして画像のURL文字列生成
func (FileSystemAvatar) AvatarURL(u ChatUser) (string, error) {

	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}

	return "", ErrNoAvatarURL
}

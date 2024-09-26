package telegram

import (
	tele "gopkg.in/telebot.v3"
	"io"
	"os"
	"time"
)

type memeData struct {
	PhotoPath  string
	TopText    string
	BottomText string
}

var storage = map[int64]memeData{}

type ImageInterface interface {
	DrawText(inputFileName, topText, bottomText string) (string, error)
}

type Service struct {
	b              *tele.Bot
	imageInterface ImageInterface
}

func NewTelegramService(token string, imageInterface ImageInterface) *Service {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		photo := c.Message().Photo
		rc, err := b.File(&tele.File{FileID: photo.FileID})
		if err != nil {
			return err
		}

		f, err := os.CreateTemp("", "telegramupload")
		if err != nil {
			return err
		}

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}

		storage[c.Sender().ID] = memeData{
			PhotoPath: f.Name(),
		}

		return c.Send("Теперь скинь верхний текст")
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		if memedata, ok := storage[c.Sender().ID]; ok {
			if memedata.TopText == "" {
				memedata.TopText = c.Message().Text
				storage[c.Sender().ID] = memedata
				return c.Send("Отправь нижний текст")
			} else {
				result, err := imageInterface.DrawText(memedata.PhotoPath, memedata.TopText, c.Message().Text)
				if err != nil {
					return err
				}
				return c.Send(&tele.Photo{File: tele.FromDisk(result)})
			}
		}

		return nil
	})
	return &Service{b: b, imageInterface: imageInterface}
}

func (s *Service) Start() {
	s.b.Start()
}

func (s *Service) Stop() {
	s.b.Stop()
}

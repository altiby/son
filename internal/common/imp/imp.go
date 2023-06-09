package imp

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"github.com/altiby/son/internal/config"
	"github.com/altiby/son/internal/domain"
	"github.com/altiby/son/internal/hasher"
	"github.com/altiby/son/internal/storage"
	"github.com/altiby/son/internal/user"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

type UsersImporter struct {
	fileName string
	cfg      *config.Config
}

func NewUsersImporter(fileName string, config *config.Config) *UsersImporter {
	return &UsersImporter{fileName: fileName, cfg: config}
}

func (i UsersImporter) Do() error {
	postgres, err := storage.InitPostgresql(i.cfg.Postgresql)
	if err != nil {
		return fmt.Errorf("postgresql inint failed, %w", err)
	}
	h := hasher.New()

	userStorage := storage.NewUserStorage(postgres)
	userService := user.NewService(userStorage, h)
	counter := 0
	err = i.fileRead(func(line string) error {
		counter++
		if counter%100000 == 0 {
			log.Info().Msgf("imported %d", counter)
		}

		splitLine := strings.Split(line, ",")
		if len(splitLine) != 3 {
			return nil
		}
		names := strings.Split(splitLine[0], " ")
		if len(names) != 2 {
			return nil
		}

		years, _ := strconv.Atoi(splitLine[1])

		_, err := userService.RegisterUser(context.TODO(), domain.User{
			Role:       "user",
			FirstName:  names[0],
			SecondName: names[1],
			Birthdate:  time.Now().AddDate(-years, 0, 0).Format(time.DateOnly),
			Biography:  "biography",
			City:       splitLine[2],
		}, "NoP/2ssv0rd")
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("import file process failed, %w", err)
	}

	return err
}

func (i UsersImporter) fileRead(lineProcessor func(line string) error) error {
	zf, err := zip.OpenReader(i.fileName)
	if err != nil {
		return fmt.Errorf("failed to read input zip, %w", err)
	}
	defer func() {
		err := zf.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to close input zip")
			return
		}
	}()

	for _, file := range zf.File {
		fileUnpacked, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to unzip file, %w", err)
		}
		txtFile := bufio.NewScanner(fileUnpacked)
		txtFile.Split(bufio.ScanLines)

		for txtFile.Scan() {
			if err := lineProcessor(txtFile.Text()); err != nil {
				return fmt.Errorf("failed to process file line, %w", err)
			}
		}

		err = fileUnpacked.Close()
		if err != nil {
			return fmt.Errorf("failed to close unzip file, %w", err)
		}
	}

	return nil
}

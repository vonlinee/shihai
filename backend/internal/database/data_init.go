package database

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"shihai/internal/config"
	"shihai/internal/models"
)

func Init(configFile string) {
	cfg := config.Load(configFile)
	db, err := config.InitDB(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := config.AutoMigrateDatabaseModel(db, nil); err != nil {
		log.Fatal("Failed to auto migrate database model:", err)
	}

	if err := seedPoetryReferenceData(db); err != nil {
		log.Fatal("Failed to initialize poetry reference data:", err)
	}

	log.Println("Poetry reference data initialized successfully")
}

// seedPoetryReferenceData initializes dynasty and poem_type reference data.
func seedPoetryReferenceData(db *gorm.DB) error {
	if err := seedDynasties(db, initialDynasties()); err != nil {
		return err
	}
	if err := seedPoemTypes(db, InitialPoemTypes()); err != nil {
		return err
	}
	return nil
}

// initialDynasties returns the built-in dynasty seed data.
func initialDynasties() []models.Dynasty {
	return []models.Dynasty{
		newDynasty("先秦", "Pre-Qin", intPtr(-2070), intPtr(-221)),
		newDynasty("两汉", "Han", intPtr(-206), intPtr(220)),
		newDynasty("魏晋", "Wei-Jin", intPtr(220), intPtr(420)),
		newDynasty("南北朝", "Northern and Southern", intPtr(420), intPtr(589)),
		newDynasty("隋", "Sui", intPtr(581), intPtr(618)),
		newDynasty("唐", "Tang", intPtr(618), intPtr(907)),
		newDynasty("五代", "Five Dynasties", intPtr(907), intPtr(960)),
		newDynasty("宋", "Song", intPtr(960), intPtr(1279)),
		newDynasty("元", "Yuan", intPtr(1271), intPtr(1368)),
		newDynasty("清", "Qing", intPtr(1644), intPtr(1912)),
		newDynasty("其他", "Other", nil, nil),
	}
}

// InitialPoemTypes returns the built-in poem type seed data.
func InitialPoemTypes() []models.PoemType {
	return []models.PoemType{
		newPoemType(10, "唐诗", "唐诗", nil, nil, "诗"),
		newPoemType(11, "五言绝句", "唐诗", intPtr(4), intPtr(5), "四句，每句五字"),
		newPoemType(12, "七言绝句", "唐诗", intPtr(4), intPtr(7), "四句，每句七字"),
		newPoemType(13, "五言律诗", "唐诗", intPtr(8), intPtr(5), "八句，每句五字"),
		newPoemType(14, "七言律诗", "唐诗", intPtr(8), intPtr(7), "八句，每句七字"),
		newPoemType(15, "五言古诗", "唐诗", nil, intPtr(5), "不限句数，每句五字"),
		newPoemType(16, "七言古诗", "唐诗", nil, intPtr(7), "不限句数，每句七字"),
		newPoemType(17, "乐府诗", "唐诗", nil, nil, "不限句数，不限字数"),
		newPoemType(20, "宋词", "宋词", nil, nil, "长短句"),
		newPoemType(21, "五代词", "词", nil, nil, "长短句"),
		newPoemType(30, "元曲", "曲", nil, nil, "散曲"),
		newPoemType(40, "蒙学", "蒙学", nil, nil, "蒙学"),
		newPoemType(50, "诗经", "诗经", nil, nil, "诗经"),
		newPoemType(60, "论语", "论语", nil, nil, "论语"),
		newPoemType(70, "楚辞", "楚辞", nil, nil, "楚辞"),
		newPoemType(80, "四书五经", "四书五经", nil, nil, "四书五经"),
		newPoemType(99, "其他", "其他", nil, nil, "不规则或其他形式"),
	}
}

func seedDynasties(db *gorm.DB, dynasties []models.Dynasty) error {
	for _, dynasty := range dynasties {
		var existing models.Dynasty
		if err := db.Where("name = ?", dynasty.Name).Attrs(dynasty).FirstOrCreate(&existing).Error; err != nil {
			return fmt.Errorf("seed dynasty %s: %w", dynasty.Name, err)
		}
	}
	return nil
}

func seedPoemTypes(db *gorm.DB, poemTypes []models.PoemType) error {
	if len(poemTypes) == 0 {
		return nil
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&poemTypes).Error; err != nil {
		return fmt.Errorf("seed poem types: %w", err)
	}
	return nil
}

func newDynasty(name string, nameEn string, startYear *int, endYear *int) models.Dynasty {
	return models.Dynasty{
		Name:      name,
		NameEn:    nameEn,
		StartYear: startYear,
		EndYear:   endYear,
		Period:    formatPeriod(startYear, endYear),
	}
}

func newPoemType(id uint64, name string, category string, lines *int, charsPerLine *int, description string) models.PoemType {
	return models.PoemType{
		BaseModel:    models.BaseModel{ID: id},
		Name:         name,
		Category:     category,
		Lines:        lines,
		CharsPerLine: charsPerLine,
		Description:  description,
	}
}

func formatPeriod(startYear *int, endYear *int) string {
	if startYear == nil || endYear == nil {
		return ""
	}
	return fmt.Sprintf("%d-%d", *startYear, *endYear)
}

func intPtr(value int) *int {
	return &value
}

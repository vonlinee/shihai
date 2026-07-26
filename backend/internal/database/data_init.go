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
	if err := normalizeLegacyPoemTypeCategories(db); err != nil {
		return err
	}
	if err := backfillPoemGenreCategories(db); err != nil {
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
		newPoemType(10, "唐诗", "诗", nil, nil, "唐代诗歌作品"),
		newPoemType(11, "五言绝句", "诗", intPtr(4), intPtr(5), "近体诗，四句，每句五字"),
		newPoemType(12, "七言绝句", "诗", intPtr(4), intPtr(7), "近体诗，四句，每句七字"),
		newPoemType(13, "五言律诗", "诗", intPtr(8), intPtr(5), "近体诗，八句，每句五字"),
		newPoemType(14, "七言律诗", "诗", intPtr(8), intPtr(7), "近体诗，八句，每句七字"),
		newPoemType(15, "五言古诗", "诗", nil, intPtr(5), "古体诗，不限句数，每句五字"),
		newPoemType(16, "七言古诗", "诗", nil, intPtr(7), "古体诗，不限句数，每句七字"),
		newPoemType(17, "乐府诗", "诗", nil, nil, "配乐入歌的诗体，后也指仿乐府作品"),
		newPoemType(18, "四言诗", "诗", nil, intPtr(4), "古体诗，每句四字"),
		newPoemType(19, "杂言古诗", "诗", nil, nil, "古体诗，句式长短较自由"),
		newPoemType(22, "歌行", "诗", nil, nil, "古体诗一类，篇幅和句式较自由"),
		newPoemType(23, "绝句", "诗", intPtr(4), nil, "近体诗，四句成篇"),
		newPoemType(24, "律诗", "诗", intPtr(8), nil, "近体诗，八句成篇，中间两联通常对仗"),
		newPoemType(25, "排律", "诗", nil, nil, "律诗的扩展形式，通常超过八句"),
		newPoemType(50, "诗经", "诗", nil, nil, "先秦诗歌总集"),
		newPoemType(70, "楚辞", "诗", nil, nil, "楚地辞赋诗歌体式"),
		newPoemType(20, "宋词", "词", nil, nil, "宋代词体作品"),
		newPoemType(21, "五代词", "词", nil, nil, "五代时期词体作品"),
		newPoemType(26, "小令", "词", nil, nil, "篇幅较短的词调"),
		newPoemType(27, "中调", "词", nil, nil, "篇幅居中的词调"),
		newPoemType(28, "长调", "词", nil, nil, "篇幅较长的词调"),
		newPoemType(30, "元曲", "曲", nil, nil, "元代曲体作品"),
		newPoemType(31, "散曲", "曲", nil, nil, "不以舞台演出为主的曲体"),
		newPoemType(32, "戏曲", "曲", nil, nil, "用于舞台演出的曲体作品"),
		newPoemType(35, "赋", "文", nil, nil, "铺陈描写、兼具韵散特征的文体"),
		newPoemType(36, "骈文", "文", nil, nil, "讲究对偶、声律和辞藻的文体"),
		newPoemType(37, "散文", "文", nil, nil, "不受骈偶和韵律严格约束的文体"),
		newPoemType(38, "序跋", "文", nil, nil, "置于作品前后的说明或评论文字"),
		newPoemType(39, "记", "文", nil, nil, "记事、写景或抒怀的文章体裁"),
		newPoemType(40, "蒙学", "文", nil, nil, "传统启蒙读物"),
		newPoemType(41, "铭", "文", nil, nil, "刻器记功或自警的短篇韵文"),
		newPoemType(42, "说", "文", nil, nil, "议论说明类文章"),
		newPoemType(43, "传", "文", nil, nil, "记述人物事迹的文章"),
		newPoemType(60, "论语", "文", nil, nil, "儒家经典语录体文献"),
		newPoemType(80, "四书五经", "文", nil, nil, "儒家经典文献"),
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

func normalizeLegacyPoemTypeCategories(db *gorm.DB) error {
	type legacyCategoryFix struct {
		name        string
		oldCategory string
		newCategory string
	}
	fixes := []legacyCategoryFix{
		{name: "唐诗", oldCategory: "唐诗", newCategory: "诗"},
		{name: "五言绝句", oldCategory: "唐诗", newCategory: "诗"},
		{name: "七言绝句", oldCategory: "唐诗", newCategory: "诗"},
		{name: "五言律诗", oldCategory: "唐诗", newCategory: "诗"},
		{name: "七言律诗", oldCategory: "唐诗", newCategory: "诗"},
		{name: "五言古诗", oldCategory: "唐诗", newCategory: "诗"},
		{name: "七言古诗", oldCategory: "唐诗", newCategory: "诗"},
		{name: "乐府诗", oldCategory: "唐诗", newCategory: "诗"},
		{name: "诗经", oldCategory: "诗经", newCategory: "诗"},
		{name: "楚辞", oldCategory: "楚辞", newCategory: "诗"},
		{name: "宋词", oldCategory: "宋词", newCategory: "词"},
		{name: "五代词", oldCategory: "词", newCategory: "词"},
		{name: "元曲", oldCategory: "曲", newCategory: "曲"},
		{name: "蒙学", oldCategory: "蒙学", newCategory: "文"},
		{name: "论语", oldCategory: "论语", newCategory: "文"},
		{name: "四书五经", oldCategory: "四书五经", newCategory: "文"},
	}
	for _, fix := range fixes {
		if err := db.Model(&models.PoemType{}).
			Where("name = ? AND category = ?", fix.name, fix.oldCategory).
			Update("category", fix.newCategory).Error; err != nil {
			return fmt.Errorf("normalize legacy poem type category %s: %w", fix.name, err)
		}
	}
	return nil
}

func backfillPoemGenreCategories(db *gorm.DB) error {
	if err := db.Exec(`
UPDATE poem AS p
SET genre_category = COALESCE(NULLIF(pt.category, ''), '其他')
FROM poem_type AS pt
WHERE p.genre = pt.name
  AND (p.genre_category IS NULL OR p.genre_category = '');`).Error; err != nil {
		return fmt.Errorf("backfill poem genre categories from poem types: %w", err)
	}
	if err := db.Exec(`
UPDATE poem
SET genre_category = '其他'
WHERE genre <> ''
  AND (genre_category IS NULL OR genre_category = '');`).Error; err != nil {
		return fmt.Errorf("backfill poem genre categories for custom genres: %w", err)
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

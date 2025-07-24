package service

import (
	"go.uber.org/zap"
)

/*
 Main Goal for this data processor is to process data in a way so that it can be stored and searched in much more efficient manner,
 Data will be same but the data structure will be different.
*/

type Cache struct {
	log       *zap.SugaredLogger
	runTimeDB *RunTimeDB
	lit       *LineItemService
}

func NewDataProcessorService(log *zap.SugaredLogger, db *RunTimeDB, lit *LineItemService) *Cache {
	return &Cache{
		log:       log,
		runTimeDB: db,
		lit:       lit,
	}
}

func (d *Cache) PopulateCache() {
	for id, item := range d.lit.items {
		d.runTimeDB.AddKeyWords(item.Keywords, id)
		d.runTimeDB.AddCategory(item.Categories, id)
		d.runTimeDB.AddPlacements(item.Placement, id)
		totalParam := len(item.Categories) + len(item.Keywords)
		d.runTimeDB.AddParameterCount(id, totalParam)
		if len(item.Keywords) == 0 && len(item.Categories) == 0 {
			d.runTimeDB.AddTargetFree(id)
		}
	}
}

// This function will populate variable and then swap the maps
func (d *Cache) RePopulateCache() {
	// swap variable
}

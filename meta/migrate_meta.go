package meta

type Migrating struct {
	From string
	To   string
	Slot int
}

type MigrateTask struct {
	From  string
	To    string
	Tasks []Range
}

type MigrateMeta struct {
	MigrateKeysEach    int
	MigrateTimeout     int
	MigrateConcurrency int
	MigrateDoing       []Migrating
	MigrateTasks       []MigrateTask
}

const (
	DEFAULT_MIGRATE_KEYS_EACH    = 1000
	DEFAULT_MIGRATE_TIMEOUT      = 100
	DEFAULT_MINGRATE_CONCURRENCY = 1
)

func (migrateMeta *MigrateMeta) FetchMigrateMeta() error {
	err := FetchMetaDB(".MigrateMeta", migrateMeta)
	if err != nil {
		return err
	}

	if migrateMeta.MigrateKeysEach == CONFIG_NIL {
		migrateMeta.MigrateKeysEach = DEFAULT_MIGRATE_KEYS_EACH
	}
	if migrateMeta.MigrateTimeout == CONFIG_NIL {
		migrateMeta.MigrateTimeout = DEFAULT_MIGRATE_TIMEOUT
	}
	if migrateMeta.MigrateConcurrency == CONFIG_NIL {
		migrateMeta.MigrateConcurrency = DEFAULT_MINGRATE_CONCURRENCY
	}

	return nil
}

func (migrateMeta *MigrateMeta) FetchMigrateDoing() error {
	err := FetchMetaDB(".MigrateDoing", &migrateMeta.MigrateDoing)
	if err != nil {
		return err
	}

	return nil
}

func (migrateMeta *MigrateMeta) FetchMigrateTasks() error {
	err := FetchMetaDB(".MigrateTasks", &migrateMeta.MigrateTasks)
	if err != nil {
		return err
	}

	return nil
}

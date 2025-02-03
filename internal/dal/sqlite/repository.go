package sqlite

import (
	"github.com/jcocozza/jbf/internal/metadata"
	"database/sql"
	"time"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) CreateTag(metadataID int, name string) error {
	_, err := r.db.Exec("insert into tag (metadata_id, tag_name) values (?,?)", metadataID, name)
	return err
}

func (r *SQLiteRepository) ReadTagExists(tagName string) bool {
	row := r.db.QueryRow("select tag_name from tag where tag_name = ?;", tagName)
	var fpath string
	err := row.Scan(&fpath)
	return err != sql.ErrNoRows
}

func (r *SQLiteRepository) ReadTags(metadataID int) ([]string, error) {
	rows, err := r.db.Query("select tag_name from tags where metadata_id = ?", metadataID)
	if err == sql.ErrNoRows {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := []string{}
	for rows.Next() {
		var tagName string
		err := rows.Scan(&tagName)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tagName)
	}
	return tags, nil
}

func (r *SQLiteRepository) DeleteTag(tagName string) error {
	_, err := r.db.Exec("delete from tag where tag_name = ?", tagName)
	return err
}

func (r *SQLiteRepository) CreateMetadata(m metadata.Metadata) (int, error) {
	q := "insert into metadata (filepath, title, author, created, last_updated) values (?,?,?,?,?)"
	result, err := r.db.Exec(q, m.Filepath, m.Title, m.Author, time.Time(m.Created), time.Time(m.LastUpdated))
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (r *SQLiteRepository) ReadMetadataExists(filepath string) bool {
	row := r.db.QueryRow("select filepath from metadata where filepath = ?;", filepath)
	var fpath string
	err := row.Scan(&fpath)
	return err != sql.ErrNoRows
}

func (r *SQLiteRepository) ReadMetadata(filepath string) (metadata.Metadata, error) {
	row := r.db.QueryRow("select id, filepath, title, author, created, last_updated, is_home from metadata where filepath = ?", filepath)
	var m metadata.Metadata
	err := row.Scan(&m.ID, &m.Filepath, &m.Title, &m.Author, &m.Created, &m.LastUpdated, &m.IsHome)
	if err != nil {
		return metadata.Metadata{}, err
	}
	m.Filepath = filepath
	return m, nil
}

func (r *SQLiteRepository) ReadMetadataFiles() ([]string, error) {
	rows, err := r.db.Query("select filepath from metadata")
	if err == sql.ErrNoRows {
		return []string{}, nil
	}
	fileList := []string{}
	for rows.Next() {
		var fname string
		err := rows.Scan(&fname)
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, fname)
	}
	return fileList, nil
}

func (r *SQLiteRepository) ReadAllMetadata() ([]metadata.Metadata, error) {
	rows, err := r.db.Query("select id, filepath, title, author, created, last_updated, is_home from metadata ORDER BY created desc")
	if err == sql.ErrNoRows {
		return []metadata.Metadata{}, nil
	}
	if err != nil {
		return nil, err
	}
	mLst := []metadata.Metadata{}
	for rows.Next() {
		m := metadata.Metadata{}
		err := rows.Scan(&m.ID, &m.Filepath, &m.Title, &m.Author, &m.Created, &m.LastUpdated, &m.IsHome)
		if err != nil {
			return nil, err
		}
		mLst = append(mLst, m)
	}
	return mLst, nil
}

func (r *SQLiteRepository) UpdateMetadata(m metadata.Metadata) error {
	q := "update metadata set title = ?, author = ?, created = ?, last_updated = ?, is_home = ? where filepath = ?"
	_, err := r.db.Exec(q, m.Title, m.Author, time.Time(m.Created), time.Time(m.LastUpdated), m.IsHome, m.Filepath)
	return err
}

func (r SQLiteRepository) DeleteMetadata(filepath string) error {
	_, err := r.db.Exec("delete from metadata where filepath = ?;", filepath)
	return err
}

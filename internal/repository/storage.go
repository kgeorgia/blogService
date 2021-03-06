package repository

import (
    "blog_service2/internal/model"
    "database/sql"
    "fmt"
)

type PostgresDB struct {
    DB *sql.DB
}

func New(dbUser, dbPassword, dbName string) (*PostgresDB, error) {
    str := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s ",
        "database", "5432", dbName, "disable", dbUser, dbPassword)

    db, err := sql.Open("postgres", str)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
        `posts("id" SERIAL PRIMARY KEY,` +
        `"title" varchar(50), "text" varchar(1024),` +
        `"created_at" TIMESTAMP NOT NULL DEFAULT NOW())`)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
        `tags("id" SERIAL PRIMARY KEY,` +
        `"tag" varchar(20), "post_id" int)`)
    if err != nil {
        return nil, err
    }

    return &PostgresDB{
        DB: db,
    }, nil
}

func (d *PostgresDB) Insert(post model.Record) error {
    strm, err := d.DB.Prepare("INSERT INTO posts VALUES (default, $1, $2) RETURNING id")
    if err != nil {
        return err
    }

    var postId int
    err = strm.QueryRow(post.Title, post.Text).Scan(&postId)
    if err != nil {
        return err
    }

    for _, tag := range post.Tags {
        _, err = d.DB.Exec("INSERT INTO tags VALUES (default, $1, $2)", tag, postId)
        if err != nil {
            return err
        }
    }

    return nil
}

func (d *PostgresDB) Remove(id int) error {
    _, err := d.DB.Exec("DELETE FROM posts WHERE id=$1", id)
    if err != nil {
        return err
    }
    _, err = d.DB.Exec("DELETE  FROM tags WHERE post_id=$1", id)

    return err
}

func (d *PostgresDB) Update(post model.Record) (int64, error) {
    res, err := d.DB.Exec("UPDATE posts SET title = $1, text = $2 WHERE id=$3",
        post.Title, post.Text, post.Id)
    if err != nil {
        return 0, err
    }

    n, err := res.RowsAffected()
    if err != nil {
        return 0, err
    }

    return n, nil
}

func (d *PostgresDB) ReadOne(id int) (model.Record, error) {
    var rec model.Record
    row := d.DB.QueryRow("SELECT * FROM posts WHERE id=$1 ORDER BY id", id)
    tags, err := d.DB.Query("SELECT tag FROM tags WHERE post_id=$1 ORDER BY id", id)
    if err != nil {
        return rec, err
    }

    if err = row.Scan(&rec.Id, &rec.Title, &rec.Text, &rec.CreatedAt); err != nil {
        return rec, err
    }

    var tmp string
    for tags.Next() {
        if err = tags.Scan(&tmp); err != nil {
            return rec, err
        }
        rec.Tags = append(rec.Tags, tmp)
    }

    return rec, nil
}

func (d *PostgresDB) Read(str string) ([]model.Record, error) {
    var rows *sql.Rows
    var err error
    if str != "" {
        rows, err = d.DB.Query("SELECT id FROM posts WHERE title LIKE $1 ORDER BY id", "%"+str+"%")
    } else {
        rows, err = d.DB.Query("SELECT id FROM posts ORDER BY id")
    }
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var rs = make([]model.Record, 0)
    var rec model.Record
    var tmpId int
    for rows.Next() {
        if err = rows.Scan(&tmpId); err != nil {
            return nil, err
        }
        rec, err = d.ReadOne(tmpId)
        if err != nil {
            return nil, err
        }
        rs = append(rs, rec)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return rs, nil
}

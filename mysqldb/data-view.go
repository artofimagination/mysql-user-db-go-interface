package mysqldb

import (
	"database/sql"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrNoUserViewerAdded = errors.New("No data views have been added for this user")
var ErrNoViewForThisUser = errors.New("No data views associated with this user")

var DeleteViewerIDsByUserIDQuery = "DELETE FROM users_viewers where users_id = UUID_TO_BIN(?)"

func DeleteViewerIDsByUserID(userID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteViewerIDsByUserIDQuery, userID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoViewForThisUser
	}

	return nil
}

var GetViewerIDsByUserIDQuery = "SELECT BIN_TO_UUID(viewer_id), is_owner, BIN_TO_UUID(projects_id) from users_viewers where users_id = UUID_TO_BIN(?)"

func GetViewerIDsByUserID(userID *uuid.UUID) (models.ViewersList, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(GetViewerIDsByUserIDQuery, userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	defer rows.Close()
	usersViewers := make(models.ViewersList, 0)
	for rows.Next() {
		viewer := models.Viewer{}
		err := rows.Scan(&viewer.ViewerID, &viewer.IsOwner, &viewer.ProjectID)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		usersViewers = append(usersViewers, viewer)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}
	return usersViewers, nil
}

var GetUserIDsByViewerIDQuery = "SELECT BIN_TO_UUID(users_id), is_owner from users_viewers where viewer_id = UUID_TO_BIN(?)"

func GetUserIDsByViewerID(viewerID *uuid.UUID) (*models.ViewUsers, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(GetUserIDsByViewerIDQuery, viewerID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	defer rows.Close()
	viewUsers := &models.ViewUsers{
		UsersList: make([]uuid.UUID, 0),
	}
	for rows.Next() {
		isOwner := false
		userID := uuid.UUID{}
		err := rows.Scan(&userID, &isOwner)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		if isOwner {
			viewUsers.OwnerID = userID
		}
		viewUsers.UsersList = append(viewUsers.UsersList, userID)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}
	return viewUsers, nil
}

var AddUsersViewersQuery = "INSERT INTO users_viewers (users_id, projects_id, viewer_id, is_owner) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?), ?)"

func AddUsersViewers(userID *uuid.UUID, userViewers models.ViewersList, tx *sql.Tx) error {
	for _, viewer := range userViewers {
		result, err := tx.Exec(AddUsersViewersQuery, userID, viewer.ProjectID, viewer.ViewerID, viewer.IsOwner)
		if err != nil {
			return RollbackWithErrorStack(tx, err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return RollbackWithErrorStack(tx, err)
		}

		if affected == 0 {
			if errRb := tx.Rollback(); errRb != nil {
				return err
			}
			return ErrNoUserViewerAdded
		}
	}
	return nil
}

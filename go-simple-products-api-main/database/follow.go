package database

import (
	"Capstone/models"
	"context"
)

func CreateFollow(ctx context.Context, Follow models.Follow) (models.Follow, error) {
	err := DB.WithContext(ctx).Create(&Follow).Error
	if err != nil {
		return models.Follow{}, err
	}

	// Preload user data for the created Follow
	err = DB.WithContext(ctx).Preload("User").First(&Follow).Error
	if err != nil {
		return models.Follow{}, err
	}

	return Follow, nil
}

func DeleteFollows(ctx context.Context, id int) error {
	var Follow models.Follow

	result := DB.WithContext(ctx).Where("id = ?", id).Delete(&Follow)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrIDNotFound
	}

	return nil
}


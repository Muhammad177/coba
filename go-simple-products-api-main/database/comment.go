package database

import (
	"Capstone/models"
	"context"
)

func CreateComment(ctx context.Context, Comment models.Comment) (models.Comment, error) {
	err := DB.WithContext(ctx).Create(&Comment).Error
	if err != nil {
		return models.Comment{}, err
	}

	// Preload user data for the created Comment
	err = DB.WithContext(ctx).Preload("Thread").First(&Comment).Error
	if err != nil {
		return models.Comment{}, err
	}

	return Comment, nil
}

func DeleteComments(ctx context.Context, id int) error {
	var Comment models.Comment

	result := DB.WithContext(ctx).Where("id = ?", id).Delete(&Comment)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrIDNotFound
	}

	return nil
}

func UpdateComments(ctx context.Context, id int, Comment models.Comment) (interface{}, error) {
	result := DB.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", id).Updates(Comment)
	if result.Error != nil {
		return Comment, result.Error
	}

	if result.RowsAffected == 0 {
		return Comment, ErrIDNotFound
	}

	return Comment, nil
}

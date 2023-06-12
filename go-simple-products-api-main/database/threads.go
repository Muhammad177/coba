package database

import (
	"Capstone/models"
	"context"
)

func GetThreads(ctx context.Context) ([]models.Thread, error) {

	var thread []models.Thread

	err := DB.WithContext(ctx).Preload("User").Find(&thread).Error
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func GetThreadsByID(ctx context.Context, id int) (models.Thread, error) {
	var thread models.Thread

	err := DB.WithContext(ctx).Preload("Comments").Where("id = ?", id).First(&thread).Error
	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}
func GetThreadByTitle(ctx context.Context, title string) (thread []models.Thread, err error) {
	title = "%" + title + "%"
	err = DB.WithContext(ctx).Where("title LIKE ? OR topic LIKE ?", title, title).Find(&thread).Error
	if err != nil {
		return thread, err
	}

	return thread, nil
}

func CreateThreads(ctx context.Context, thread models.Thread) (models.Thread, error) {
	err := DB.WithContext(ctx).Create(&thread).Error
	if err != nil {
		return models.Thread{}, err
	}

	// Preload user data for the created thread
	err = DB.WithContext(ctx).Preload("User").First(&thread).Error
	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func DeleteThreads(ctx context.Context, id int) error {
	var thread models.Thread

	result := DB.WithContext(ctx).Where("id = ?", id).Delete(&thread)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrIDNotFound
	}

	return nil
}

func UpdateThreads(ctx context.Context, id int, thread models.Thread) (interface{}, error) {
	result := DB.WithContext(ctx).Model(&models.Thread{}).Where("id = ?", id).Updates(thread)
	if result.Error != nil {
		return thread, result.Error
	}

	if result.RowsAffected == 0 {
		return thread, ErrIDNotFound
	}

	return thread, nil
}

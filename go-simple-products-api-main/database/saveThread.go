package database

import (
	"Capstone/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

func CreateSaveThreads(ctx context.Context, user_id int, thread_id int) (models.User, error) {
	var thread models.Thread

	fmt.Println(thread_id, user_id)
	err := DB.WithContext(ctx).Where("id = ?", thread_id).First(&thread).Error
	if err != nil {
		return models.User{}, err
	}

	fmt.Print(thread)
	err = DB.WithContext(ctx).Model(&models.User{Model: gorm.Model{ID: uint(user_id)}}).Association("Bookmarked").Append(&thread)
	if err != nil {
		return models.User{}, err
	}

	var svThread models.User
	err = DB.WithContext(ctx).Preload("Bookmarked").Find(&svThread).Error
	if err != nil {
		return models.User{}, err
	}

	return svThread, nil
}

func DeleteSaveThreads(ctx context.Context, user_id int, thread_id int) error {
	var thread models.Thread

	fmt.Println(thread_id, user_id)
	err := DB.WithContext(ctx).Where("id = ?", thread_id).First(&thread).Error
	if err != nil {
		return err
	}

	err = DB.WithContext(ctx).Model(&models.User{Model: gorm.Model{ID: uint(user_id)}}).Association("Bookmarked").Delete(&thread)
	if err != nil {
		return err
	}

	return nil
}

func GetSaveThreads(ctx context.Context, user_id int) (models.User, error) {
	var svThread models.User

	err := DB.WithContext(ctx).Where(" id = ? ", user_id).Preload("Bookmarked").Find(&svThread).Error
	if err != nil {
		return models.User{}, err
	}

	return svThread, nil
}

package controllers

import (
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NewCategoryInput struct {
	Name        string     `json:"name"`
	ParentId    *uuid.UUID `json:"parent_id"`
	Thumbnail   string     `json:"thumbnail"`
	Description string     `json:"description"`
	IsActive    *bool      `json:"is_active"`
}

func CreateCategory(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusForbidden, "Error Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}

	var (
		input       NewCategoryInput
		newCategory models.Category
	)

	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "InvalidInput", consts.InvalidReqInput)
	}

	if !utils.IsValidStrLen(input.Name, 100) {
		return ResponseError(c, fiber.StatusBadRequest, "InvalidInput", consts.InvalidReqInput)
	}

	if input.ParentId != nil {
		parentCategory, err := repo.GetCategoryByIdAndCenterId(*input.ParentId, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, "InvalidInput", consts.DataNotFound)
		}
		if parentCategory.ParentId != nil {
			return ResponseError(c, fiber.StatusBadRequest, "InvalidInput", consts.InvalidReqInput)
		}
		newCategory.ParentId = input.ParentId
	}

	if _, err := repo.GetCategoryByNameAndCenterId(input.Name, *user.CenterId); err == nil {
		return ResponseError(c, fiber.StatusBadRequest, "Tên đã tồn tại", consts.ERROR_CATEGORY_EXISTS)
	}
	newCategory.Name = input.Name
	newCategory.Description = input.Description
	newCategory.IsActive = input.IsActive
	newCategory.CenterId = *user.CenterId
	newCategory.CreatedBy = user.ID
	newCategory.Thumbnail = input.Thumbnail
	category, err := repo.CreateCategory(&newCategory)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed", consts.CreateFailed)
	}
	return ResponseSuccess(c, fiber.StatusCreated, "Success", category)
}

func ReadCategory(c *fiber.Ctx) error {
	// Lấy thông tin người dùng từ context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	// Kiểm tra xem người dùng có thuộc trung tâm không
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}

	// Lấy thông tin categoryId từ URL param
	categoryId := c.Params("id")
	if categoryId == "" {
		return ResponseError(c, fiber.StatusBadRequest, "Category ID is required", consts.InvalidReqInput)
	}

	// Chuyển categoryId thành uuid
	categoryUUID, err := uuid.Parse(categoryId)
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "ID không hợp lệ", consts.InvalidReqInput)
	}

	// Truy vấn danh mục theo id và centerId
	category, err := repo.GetCategoryByIdAndCenterId(categoryUUID, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, "Category not found", consts.DataNotFound)
	}

	// Trả về kết quả
	return ResponseSuccess(c, fiber.StatusOK, "Success", category)
}

func ReadListCategory(c *fiber.Ctx) error {
	// Lấy thông tin người dùng từ context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	// Kiểm tra xem người dùng có thuộc trung tâm không
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}

	// Lấy các query parameters từ request
	isActiveQuery := c.Query("is_active") // Ví dụ: /categories?is_active=true

	var isActive *bool
	if isActiveQuery != "" {
		active := (isActiveQuery == "true")
		isActive = &active
	}

	// Truy vấn danh sách category của centerId, có thể lọc theo trạng thái active
	categories, err := repo.GetCategoriesByCenterIdAndActive(*user.CenterId, isActive)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed to retrieve categories", consts.GetFailed)
	}

	// Trả về kết quả
	return ResponseSuccess(c, fiber.StatusOK, "Success", categories)
}

func DeleteCategory(c *fiber.Ctx) error {
	// Lấy thông tin người dùng từ context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	// Kiểm tra xem người dùng có thuộc trung tâm không
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}

	// Lấy ID từ request
	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "ID không hợp lệ", consts.InvalidReqInput)
	}

	// Kiểm tra xem danh mục có tồn tại và thuộc trung tâm của người dùng không
	_, err = repo.GetCategoryByIdAndCenterId(categoryID, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, "Không tìm thấy dữ liệu", consts.ERROR_CATEGORY_NOT_FOUND)
	}

	// Kiểm tra xem danh mục có danh mục con không
	count, err := repo.CountChildCategories(categoryID)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed to check dependencies", consts.GetFailed)
	}
	if count > 0 {
		return ResponseError(c, fiber.StatusBadRequest, "Cannot delete category with child categories", consts.ERROR_CATEGORY_HAS_CHILDREN)
	}

	// Kiểm tra xem danh mục có ràng buộc với `Curriculums` hoặc `Subjects`
	hasDependencies, err := repo.HasCategoryDependencies(categoryID)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed to check dependencies", consts.GetFailed)
	}
	if hasDependencies {
		return ResponseError(c, fiber.StatusBadRequest, "Cannot delete category with linked curriculums or subjects", consts.ERROR_CATEGORY_HAS_DEPENDENCIES)
	}

	// Xóa danh mục
	if err := repo.DeleteCategoryById(categoryID); err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Xóa không thành công", consts.DeletedFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, "Xóa thành công", nil)
}

func UpdateCategory(c *fiber.Ctx) error {
	// Lấy thông tin người dùng từ context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied", consts.ERROR_PERMISSION_DENIED)
	}

	// Kiểm tra xem người dùng có thuộc trung tâm không
	if user.CenterId == nil {
		return ResponseError(c, fiber.StatusForbidden, "Permission denied - center", consts.ERROR_PERMISSION_DENIED)
	}

	// Lấy ID từ request
	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid category ID", consts.InvalidReqInput)
	}

	// Lấy danh mục cần cập nhật
	category, err := repo.GetCategoryByIdAndCenterId(categoryID, *user.CenterId)
	if err != nil {
		return ResponseError(c, fiber.StatusNotFound, "Category not found", consts.DataNotFound)
	}

	// Parse dữ liệu từ request body
	var input NewCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid input", consts.InvalidReqInput)
	}

	// Kiểm tra độ dài của tên danh mục
	if !utils.IsValidStrLen(input.Name, 100) {
		return ResponseError(c, fiber.StatusBadRequest, "Invalid input", consts.InvalidReqInput)
	}

	// Kiểm tra danh mục cha (nếu có)
	if input.ParentId != nil {
		parentCategory, err := repo.GetCategoryByIdAndCenterId(*input.ParentId, *user.CenterId)
		if err != nil {
			return ResponseError(c, fiber.StatusBadRequest, "Invalid parent category", consts.DataNotFound)
		}
		if parentCategory.ParentId != nil {
			return ResponseError(c, fiber.StatusBadRequest, "A subcategory cannot have another subcategory", consts.InvalidReqInput)
		}
		category.ParentId = input.ParentId
	}

	// Kiểm tra xem tên danh mục có bị trùng không (trừ chính nó)
	existingCategory, err := repo.GetCategoryByNameAndCenterId(input.Name, *user.CenterId)
	if err == nil && existingCategory.ID != category.ID {
		return ResponseError(c, fiber.StatusBadRequest, "Category name already exists", consts.ERROR_CATEGORY_EXISTS)
	}

	// Cập nhật thông tin danh mục
	category.Name = input.Name
	category.Description = input.Description
	category.IsActive = input.IsActive
	category.Thumbnail = input.Thumbnail

	// Thực hiện cập nhật trong database
	updatedCategory, err := repo.UpdateCategory(&category)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Failed to update category", consts.UpdateFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, "Category updated successfully", updatedCategory)
}

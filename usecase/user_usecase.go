package usecase

import (
	"fmt"
	"go-clean-arch/pkg/hash"
	"go-clean-arch/pkg/jwt"
	"go-clean-arch/repository"
	"go-clean-arch/structure/entity"
	"go-clean-arch/structure/request"
	"go-clean-arch/structure/response"
)

type UserUsecase interface {
	CreateTeacher(req *request.RegisterTeacher) error
	CreateStudent(req *request.RegisterStudent) error
	GetUserByEmail(email string, password string) (string, string, error)

	GetUserByClaims(claims map[string]interface{}) (interface{}, error)
	GetAllTeacher() ([]response.TeacherResponse, error)
	GetAllStudent() ([]response.StudentResponse, error)
	SendEvent(year uint, claims map[string]interface{}) error
	GetStudentsAndYearsByCertifier(claims map[string]interface{}) ([]response.StudentYear, error)

	UpdateTeacherByID(req *request.RegisterTeacher, claims map[string]interface{}) error
	UpdateStudentByID(req *request.RegisterStudent, claims map[string]interface{}) error
	UpdateRoleByID(userID uint, role string) error
	UpdateStatusDones(certifierID uint, userID uint, status bool, comment string) error

	MarkNewsAsRead(newsID uint) error
	GetNews(claims map[string]interface{}) ([]response.NewsResponse, error)

	GetDoneFiltered(year uint, status bool, facultyID *uint) ([]response.ListResponse, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
	jwt      jwt.JWTService
}

func NewUserUsecase(userRepo repository.UserRepository, jwt jwt.JWTService) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (u *userUsecase) CreateTeacher(req *request.RegisterTeacher) error {
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	req.Password = hashedPassword
	if req.Role == "" {
		return fmt.Errorf("failed to input role")
	}

	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
	
	teacher := &entity.Teacher{
		TitleName: req.TitleName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Code:      req.Code,
	}

	if err := u.userRepo.CreateTeacher(user, teacher); err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	return nil
}

func (u *userUsecase) CreateStudent(req *request.RegisterStudent) error {
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	req.Password = hashedPassword

	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     "student",
	}
	student := &entity.Student{
		TitleName: req.TitleName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Code:      req.Code,
		Year:      req.Year,
		BranchId:  req.BranchId,
	}

	if err := u.userRepo.CreateStudent(user, student); err != nil {
		return fmt.Errorf("failed to create student: %w", err)
	}

	return nil
}

func (u *userUsecase) GetUserByEmail(email string, password string) (string, string, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if user == nil || err != nil {
		return "", "", fmt.Errorf("invalid email or password")
	}
	if !hash.CheckPasswordHash(password, user.Password) {
		return "", "", fmt.Errorf("invalid email or password")
	}
	token, err := u.jwt.GenerateJWT(user.UserID, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user.Role, nil
}

func (u *userUsecase) getTeacherByUserID(userID uint) (*response.TeacherByClaimsResponse, error) {
	result, superUser, err := u.userRepo.GetTeacherByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher by user id %d: %w", userID, err)
	}
	teacher := response.TeacherByClaimsResponse{
		UserID:    result.UserID,
		TitleName: result.TitleName,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Phone:     result.Phone,
		Code:      result.Code,
		SuperUser: superUser,
	}
	return &teacher, nil
}

func (u *userUsecase) MapStudentResponse(s *entity.Student) response.StudentResponse {
	return response.StudentResponse{
		UserID:      s.UserID,
		TitleName:   s.TitleName,
		FirstName:   s.FirstName,
		LastName:    s.LastName,
		Phone:       s.Phone,
		Code:        s.Code,
		Year:        s.Year,
		BranchID:    s.BranchId,
		BranchName:  s.Branch.BranchName,
		FacultyID:   s.Branch.Faculty.FacultyID,
		FacultyName: s.Branch.Faculty.FacultyName,
	}
}

func (u *userUsecase) getStudentByUserID(userID uint) (*response.StudentResponse, error) {
	student, err := u.userRepo.GetStudentByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student by user id %d: %w", userID, err)
	}
	res := u.MapStudentResponse(student)

	return &res, nil
}

func (u *userUsecase) GetUserByClaims(claims map[string]interface{}) (interface{}, error) {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	// ตรวจสอบ role จาก claims
	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role in claims")
	}

	switch role {
	case "student":
		result, err := u.getStudentByUserID(userID)
		if err != nil {
			return nil, err
		}
		result.Role = role
		return result, nil
	case "teacher", "admin":
		result, err := u.getTeacherByUserID(userID)
		if err != nil {
			return nil, err
		}
		result.Role = role
		return result, nil
	case "superadmin":
		return map[string]interface{}{
			"message": "welcome superadmin",
			"role":    "superadmin",
		}, nil
	default:
		return nil, fmt.Errorf("unauthorized access")
	}
}

func (u *userUsecase) GetAllTeacher() ([]response.TeacherResponse, error) {
	allTeacher, err := u.userRepo.GetAllTeacher()
	if err != nil {
		return nil, fmt.Errorf("usecase: %w", err)
	}
	return allTeacher, nil
}

func (u *userUsecase) GetAllStudent() ([]response.StudentResponse, error) {
	allStudents, err := u.userRepo.GetAllStudent()
	if err != nil {
		return nil, fmt.Errorf("usecase: %w", err)
	}
	var res []response.StudentResponse
	for _, s := range allStudents {
		studentResponse := u.MapStudentResponse(&s)
		res = append(res, studentResponse)
	}
	return res, nil
}

func (u *userUsecase) UpdateTeacherByID(req *request.RegisterTeacher, claims map[string]interface{}) error {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	teacher := entity.Teacher{
		UserID:    userID,
		TitleName: req.TitleName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Code:      req.Code,
	}
	return u.userRepo.UpdateTeacherByID(&teacher)
}

func (u *userUsecase) UpdateStudentByID(req *request.RegisterStudent, claims map[string]interface{}) error {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)
	role, ok := claims["role"].(string)
	if !ok {
		return fmt.Errorf("invalid role in claims")
	}
	var student entity.Student
	if role == "admin" {
		student = entity.Student{
			UserID:    req.UserID, // Admin can update any user
			TitleName: req.TitleName,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Phone:     req.Phone,
			Code:      req.Code,
			Year:      req.Year,
			BranchId:  req.BranchId,
		}
	} else {
		student = entity.Student{
			UserID:    userID, // Regular users can only update themselves
			TitleName: req.TitleName,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Phone:     req.Phone,
			Code:      req.Code,
			Year:      req.Year,
			BranchId:  req.BranchId,
		}
	}

	return u.userRepo.UpdateStudentByID(&student)
}

func (u *userUsecase) UpdateRoleByID(userID uint, role string) error {
	if role == "admin" || role == "student" || role == "teacher" {
		return u.userRepo.UpdateRoleByID(userID, role)
	}
	return fmt.Errorf("incorrect role")
}

func (u *userUsecase) SendEvent(year uint, claims map[string]interface{}) error {
	// ดึง user_id จาก claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	// คำนวณชั่วโมงจาก EventInside และ EventOutside
	outsideHour, insideHour, err := u.userRepo.GetTotalWorkingHours(userID, year)
	if err != nil {
		return fmt.Errorf("failed to get total working hours: %v", err)
	}

	// ตรวจสอบชั่วโมง
	if insideHour >= 18 && (outsideHour+insideHour) >= 36 {
		// ดึง superUserID ของ student
		superUserID, err := u.userRepo.GetSuperUserForStudent(userID)
		if err != nil {
			return fmt.Errorf("failed to get super user: %v", err)
		}

		// สร้าง Dones
		err = u.userRepo.CreateDones(userID, year, *superUserID)
		if err != nil {
			return fmt.Errorf("failed to create dones: %v", err)
		}
		return nil
	}

	return fmt.Errorf("insufficient working hours: inside=%d, outside=%d", insideHour, outsideHour)
}

func (u *userUsecase) GetStudentsAndYearsByCertifier(claims map[string]interface{}) ([]response.StudentYear, error) {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	result, err := u.userRepo.GetStudentsAndYearsByCertifier(userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *userUsecase) UpdateStatusDones(certifierID uint, userID uint, status bool, comment string) error {
	return u.userRepo.UpdateStatusDones(certifierID, userID, status, comment)
}
func (u *userUsecase) GetNews(claims map[string]interface{}) ([]response.NewsResponse, error) {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	result, err := u.userRepo.GetNews(userID)
	if err != nil {
		return nil, err
	}
	var news []response.NewsResponse
	for _, r := range result {
		n := response.NewsResponse{
			NewsID:  r.NewsID,
			Title:   r.Title,
			UserID:  r.UserID,
			Message: r.Message,
			IsRead:  r.IsRead,
			// CreatedAt: r.CreatedAt,
		}
		news = append(news, n)
	}

	return news, nil
}

func (u *userUsecase) MarkNewsAsRead(newsID uint) error {
	err := u.userRepo.MarkNewsAsRead(newsID)
	if err != nil {
		return fmt.Errorf("failed to update read: %v", err)
	}
	return nil
}

func (u *userUsecase) GetDoneFiltered(year uint, status bool, facultyID *uint) ([]response.ListResponse, error) {
	result, err := u.userRepo.GetDoneFiltered(year, status, facultyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get: %v", err)
	}
	var list []response.ListResponse
	for _, r := range result {
		l := response.ListResponse{
			UserID: r.User,
			Student: struct {
				TitleName   string `json:"title_name"`
				FirstName   string `json:"first_name"`
				LastName    string `json:"last_name"`
				Phone       string `json:"phone"`
				Code        string `json:"code"`
				Year        uint   `json:"year"`
				BranchName  string `json:"branch_name"`
				FacultyName string `json:"faculty_name"`
			}{
				TitleName:   r.Student.TitleName,
				FirstName:   r.Student.FirstName,
				LastName:    r.Student.LastName,
				Phone:       r.Student.Phone,
				Code:        r.Student.Code,
				Year:        r.Student.Year,
				BranchName:  r.Student.Branch.BranchName,
				FacultyName: r.Student.Branch.Faculty.FacultyName,
			},
			Status:  r.Status,
			Comment: r.Comment,
		}
		list = append(list, l)
	}

	return list, nil
}

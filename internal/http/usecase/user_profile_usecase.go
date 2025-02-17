package usecase

import (
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IUserProfileUseCase interface {
	FillUserProfile(req *request.FillUserProfileRequest, userID uuid.UUID) (*response.UserProfileResponse, error)
	FillUserProfileMessage(req *request.FillUserProfileRequest, userID uuid.UUID) (string, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.UserProfileResponse, int64, error)
	FindByID(id uuid.UUID) (*response.UserProfileResponse, error)
	FindByUserID(userID uuid.UUID) (*response.UserProfileResponse, error)
	UpdateStatusUserProfile(req *request.UpdateStatusUserProfileRequest) (*response.UserProfileResponse, error)
	DeleteUserProfile(id uuid.UUID) error
	UpdateAvatar(id uuid.UUID, avatarPath string) (*response.UserProfileResponse, error)
}

type UserProfileUseCase struct {
	Log                      *logrus.Logger
	Repository               repository.IUserProfileRepository
	DTO                      dto.IUserProfileDTO
	WorkExperienceRepository repository.IWorkExperienceRepository
	EducationRepository      repository.IEducationRepository
	SkillRepository          repository.ISkillRepository
	Viper                    *viper.Viper
}

func NewUserProfileUseCase(
	log *logrus.Logger,
	repo repository.IUserProfileRepository,
	uDTO dto.IUserProfileDTO,
	weRepository repository.IWorkExperienceRepository,
	edRepository repository.IEducationRepository,
	sRepository repository.ISkillRepository,
	viper *viper.Viper,
) IUserProfileUseCase {
	return &UserProfileUseCase{
		Log:                      log,
		Repository:               repo,
		DTO:                      uDTO,
		WorkExperienceRepository: weRepository,
		EducationRepository:      edRepository,
		SkillRepository:          sRepository,
		Viper:                    viper,
	}
}

func UserProfileUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IUserProfileUseCase {
	repo := repository.UserProfileRepositoryFactory(log)
	uDTO := dto.UserProfileDTOFactory(log, viper)
	weRepository := repository.WorkExperienceRepositoryFactory(log)
	edRepository := repository.EducationRepositoryFactory(log)
	sRepository := repository.SkillRepositoryFactory(log)
	return NewUserProfileUseCase(log, repo, uDTO, weRepository, edRepository, sRepository, viper)
}

func (uc *UserProfileUseCase) FillUserProfile(req *request.FillUserProfileRequest, userID uuid.UUID) (*response.UserProfileResponse, error) {
	parsedBirthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing birth date: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing birth date: " + err.Error())
	}
	if req.ID == "" || req.ID == uuid.Nil.String() {
		exist, err := uc.Repository.FindByUserID(userID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: " + err.Error())
		}
		if exist != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] user profile already exist")
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] user profile already exist")
		}
		createdProfile, err := uc.Repository.CreateUserProfile(&entity.UserProfile{
			UserID:          &userID,
			Name:            req.Name,
			MaritalStatus:   entity.MaritalStatusEnum(req.MaritalStatus),
			Gender:          entity.UserGender(req.Gender),
			PhoneNumber:     req.PhoneNumber,
			Age:             req.Age,
			BirthDate:       parsedBirthDate,
			BirthPlace:      req.BirthPlace,
			Ktp:             req.KtpPath,
			CurriculumVitae: req.CvPath,
			Address:         req.Address,
			Bilingual:       req.Bilingual,
			Status:          entity.USER_INACTIVE,
		})
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating user profile: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating user profile: " + err.Error())
		}
		if len(req.WorkExperiences) > 0 {
			for _, we := range req.WorkExperiences {
				_, err := uc.WorkExperienceRepository.CreateWorkExperience(&entity.WorkExperience{
					UserProfileID:  createdProfile.ID,
					CompanyName:    we.CompanyName,
					Name:           we.Name,
					YearExperience: we.YearExperience,
					JobDescription: we.JobDescription,
					Certificate:    we.CertificatePath,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating work experience: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating work experience: " + err.Error())
				}
			}
		}
		if len(req.Educations) > 0 {
			for _, ed := range req.Educations {
				parsedEndDate, err := time.Parse("2006-01-02", ed.EndDate)
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing end date: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing end date: " + err.Error())
				}
				var gpa *float64
				if ed.Gpa != nil {
					gpa = ed.Gpa
				}
				_, err = uc.EducationRepository.CreateEducation(&entity.Education{
					UserProfileID:  createdProfile.ID,
					SchoolName:     ed.SchoolName,
					Major:          ed.Major,
					GraduateYear:   ed.GraduateYear,
					Certificate:    ed.CertificatePath,
					EndDate:        parsedEndDate,
					Gpa:            gpa,
					EducationLevel: entity.EducationLevelEnum(ed.EducationLevel),
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating education: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating education: " + err.Error())
				}
			}
		}
		if len(req.Skills) > 0 {
			for _, s := range req.Skills {
				var level *int
				if s.Level != nil {
					level = s.Level
				}
				_, err := uc.SkillRepository.CreateSkill(&entity.Skill{
					UserProfileID: createdProfile.ID,
					Name:          s.Name,
					Description:   s.Description,
					Certificate:   s.CertificatePath,
					Level:         level,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating skill: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating skill: " + err.Error())
				}
			}
		}
	} else {
		parsedID, err := uuid.Parse(req.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing ID: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing ID: " + err.Error())
		}

		exist, err := uc.Repository.FindByID(parsedID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when finding user profile by ID: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when finding user profile by ID: " + err.Error())
		}
		if exist == nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] user profile not found")
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] user profile not found")
		}

		updatedProfile, err := uc.Repository.UpdateUserProfile(&entity.UserProfile{
			ID:              parsedID,
			Name:            req.Name,
			UserID:          &userID,
			MaritalStatus:   entity.MaritalStatusEnum(req.MaritalStatus),
			Gender:          entity.UserGender(req.Gender),
			PhoneNumber:     req.PhoneNumber,
			Age:             req.Age,
			BirthDate:       parsedBirthDate,
			BirthPlace:      req.BirthPlace,
			Ktp:             req.KtpPath,
			CurriculumVitae: req.CvPath,
			Address:         req.Address,
			Bilingual:       req.Bilingual,
			Status:          entity.USER_INACTIVE,
		})
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when updating user profile: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when updating user profile: " + err.Error())
		}

		var workExpIDs []uuid.UUID
		if len(req.WorkExperiences) > 0 {
			for _, we := range req.WorkExperiences {
				if we.ID != nil {
					woExpID, err := uuid.Parse(*we.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing work experience ID: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing work experience ID: " + err.Error())
					}
					workExpIDs = append(workExpIDs, woExpID)
				}
			}
		}
		err = uc.WorkExperienceRepository.DeleteNotInIDAndUserProfileID(workExpIDs, updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting work experiences: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when deleting work experiences: " + err.Error())
		}
		if len(req.WorkExperiences) > 0 {
			for _, we := range req.WorkExperiences {
				if we.ID != nil {
					woExpID, err := uuid.Parse(*we.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing work experience ID: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing work experience ID: " + err.Error())
					}

					_, err = uc.WorkExperienceRepository.UpdateWorkExperience(&entity.WorkExperience{
						ID:             woExpID,
						UserProfileID:  updatedProfile.ID,
						CompanyName:    we.CompanyName,
						Name:           we.Name,
						YearExperience: we.YearExperience,
						JobDescription: we.JobDescription,
						Certificate:    we.CertificatePath,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when updating work experience: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when updating work experience: " + err.Error())
					}
				} else {
					_, err = uc.WorkExperienceRepository.CreateWorkExperience(&entity.WorkExperience{
						UserProfileID:  updatedProfile.ID,
						CompanyName:    we.CompanyName,
						Name:           we.Name,
						YearExperience: we.YearExperience,
						JobDescription: we.JobDescription,
						Certificate:    we.CertificatePath,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating work experience: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating work experience: " + err.Error())
					}
				}
			}
		}

		var eduIDs []uuid.UUID
		if len(req.Educations) > 0 {
			for _, ed := range req.Educations {
				if ed.ID != nil {
					eduID, err := uuid.Parse(*ed.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing education ID: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing education ID: " + err.Error())
					}
					eduIDs = append(eduIDs, eduID)
				}
			}
		}

		err = uc.EducationRepository.DeleteNotInIDAndUserProfileID(eduIDs, updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting educations: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when deleting educations: " + err.Error())
		}

		if len(req.Educations) > 0 {
			for _, ed := range req.Educations {
				if ed.ID == nil {
					parsedEndDate, err := time.Parse("2006-01-02", ed.EndDate)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing end date: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing end date: " + err.Error())
					}
					var gpa *float64
					if ed.Gpa != nil {
						gpa = ed.Gpa
					}
					_, err = uc.EducationRepository.CreateEducation(&entity.Education{
						UserProfileID:  updatedProfile.ID,
						SchoolName:     ed.SchoolName,
						Major:          ed.Major,
						GraduateYear:   ed.GraduateYear,
						Certificate:    ed.CertificatePath,
						EndDate:        parsedEndDate,
						EducationLevel: entity.EducationLevelEnum(ed.EducationLevel),
						Gpa:            gpa,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating education: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating education: " + err.Error())
					}
				} else {
					eduID, err := uuid.Parse(*ed.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing education ID: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing education ID: " + err.Error())
					}
					parsedEndDate, err := time.Parse("2006-01-02", ed.EndDate)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing end date: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing end date: " + err.Error())
					}
					var gpa *float64
					if ed.Gpa != nil {
						gpa = ed.Gpa
					}
					_, err = uc.EducationRepository.UpdateEducation(&entity.Education{
						ID:             eduID,
						UserProfileID:  updatedProfile.ID,
						SchoolName:     ed.SchoolName,
						Major:          ed.Major,
						GraduateYear:   ed.GraduateYear,
						Certificate:    ed.CertificatePath,
						EndDate:        parsedEndDate,
						EducationLevel: entity.EducationLevelEnum(ed.EducationLevel),
						Gpa:            gpa,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating education: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating education: " + err.Error())
					}
				}
			}
		}

		var skillIDs []uuid.UUID
		if len(req.Skills) > 0 {
			for _, s := range req.Skills {
				if s.ID != nil {
					skillID, err := uuid.Parse(*s.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing skill ID: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing skill ID: " + err.Error())
					}
					skillIDs = append(skillIDs, skillID)
				}
			}
		}

		err = uc.SkillRepository.DeleteNotInIDAndUserProfileID(skillIDs, updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting skills: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when deleting skills: " + err.Error())
		}

		if len(req.Skills) > 0 {
			for _, s := range req.Skills {
				if s.ID == nil {
					var level *int
					if s.Level != nil {
						level = s.Level
					}
					_, err := uc.SkillRepository.CreateSkill(&entity.Skill{
						UserProfileID: updatedProfile.ID,
						Name:          s.Name,
						Description:   s.Description,
						Certificate:   s.CertificatePath,
						Level:         level,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating skill: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating skill: " + err.Error())
					}
				} else {
					skillID, err := uuid.Parse(*s.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing skill ID: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing skill ID: " + err.Error())
					}
					var level *int
					if s.Level != nil {
						level = s.Level
					}
					_, err = uc.SkillRepository.UpdateSkill(&entity.Skill{
						ID:            skillID,
						UserProfileID: updatedProfile.ID,
						Name:          s.Name,
						Description:   s.Description,
						Certificate:   s.CertificatePath,
						Level:         level,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when updating skill: %s", err.Error())
						return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when updating skill: " + err.Error())
					}
				}
			}
		}
	}
	// find user profile by user ID
	profile, err := uc.Repository.FindByUserID(userID)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: " + err.Error())
	}

	if profile == nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] user profile not found")
		return nil, errors.New("[UserProfileUseCase.FillUserProfile] user profile not found")
	}

	resp, err := uc.DTO.ConvertEntityToResponse(profile)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when converting entity to response: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when converting entity to response: " + err.Error())
	}
	return resp, nil
}

func (uc *UserProfileUseCase) FillUserProfileMessage(req *request.FillUserProfileRequest, userID uuid.UUID) (string, error) {
	parsedBirthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing birth date: %s", err.Error())
		return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing birth date: " + err.Error())
	}
	if req.ID == "" || req.ID == uuid.Nil.String() {
		exist, err := uc.Repository.FindByUserID(userID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: %s", err.Error())
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: " + err.Error())
		}
		if exist != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] user profile already exist")
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] user profile already exist")
		}
		createdProfile, err := uc.Repository.CreateUserProfile(&entity.UserProfile{
			UserID:          &userID,
			Name:            req.Name,
			MaritalStatus:   entity.MaritalStatusEnum(req.MaritalStatus),
			Gender:          entity.UserGender(req.Gender),
			PhoneNumber:     req.PhoneNumber,
			Age:             req.Age,
			BirthDate:       parsedBirthDate,
			BirthPlace:      req.BirthPlace,
			Ktp:             req.KtpPath,
			CurriculumVitae: req.CvPath,
			Address:         req.Address,
			Bilingual:       req.Bilingual,
			Status:          entity.USER_INACTIVE,
		})
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating user profile: %s", err.Error())
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when creating user profile: " + err.Error())
		}
		if len(req.WorkExperiences) > 0 {
			for _, we := range req.WorkExperiences {
				_, err := uc.WorkExperienceRepository.CreateWorkExperience(&entity.WorkExperience{
					UserProfileID:  createdProfile.ID,
					CompanyName:    we.CompanyName,
					Name:           we.Name,
					YearExperience: we.YearExperience,
					JobDescription: we.JobDescription,
					Certificate:    we.CertificatePath,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating work experience: %s", err.Error())
					return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when creating work experience: " + err.Error())
				}
			}
		}
		if len(req.Educations) > 0 {
			for _, ed := range req.Educations {
				parsedEndDate, err := time.Parse("2006-01-02", ed.EndDate)
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing end date: %s", err.Error())
					return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing end date: " + err.Error())
				}
				var gpa *float64
				if ed.Gpa != nil {
					gpa = ed.Gpa
				}
				_, err = uc.EducationRepository.CreateEducation(&entity.Education{
					UserProfileID:  createdProfile.ID,
					SchoolName:     ed.SchoolName,
					Major:          ed.Major,
					GraduateYear:   ed.GraduateYear,
					Certificate:    ed.CertificatePath,
					EndDate:        parsedEndDate,
					Gpa:            gpa,
					EducationLevel: entity.EducationLevelEnum(ed.EducationLevel),
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating education: %s", err.Error())
					return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when creating education: " + err.Error())
				}
			}
		}
		if len(req.Skills) > 0 {
			for _, s := range req.Skills {
				var level *int
				if s.Level != nil {
					level = s.Level
				}
				_, err := uc.SkillRepository.CreateSkill(&entity.Skill{
					UserProfileID: createdProfile.ID,
					Name:          s.Name,
					Description:   s.Description,
					Certificate:   s.CertificatePath,
					Level:         level,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating skill: %s", err.Error())
					return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when creating skill: " + err.Error())
				}
			}
		}
	} else {
		parsedID, err := uuid.Parse(req.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing ID: %s", err.Error())
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing ID: " + err.Error())
		}

		exist, err := uc.Repository.FindByID(parsedID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when finding user profile by ID: %s", err.Error())
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when finding user profile by ID: " + err.Error())
		}
		if exist == nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] user profile not found")
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] user profile not found")
		}

		updatedProfile, err := uc.Repository.UpdateUserProfile(&entity.UserProfile{
			ID:              parsedID,
			Name:            req.Name,
			UserID:          &userID,
			MaritalStatus:   entity.MaritalStatusEnum(req.MaritalStatus),
			Gender:          entity.UserGender(req.Gender),
			PhoneNumber:     req.PhoneNumber,
			Age:             req.Age,
			BirthDate:       parsedBirthDate,
			BirthPlace:      req.BirthPlace,
			Ktp:             req.KtpPath,
			CurriculumVitae: req.CvPath,
			Address:         req.Address,
			Status:          entity.USER_INACTIVE,
		})
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when updating user profile: %s", err.Error())
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when updating user profile: " + err.Error())
		}

		var workExpIDs []uuid.UUID
		if len(req.WorkExperiences) > 0 {
			for _, we := range req.WorkExperiences {
				if we.ID != nil {
					woExpID, err := uuid.Parse(*we.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing work experience ID: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing work experience ID: " + err.Error())
					}
					workExpIDs = append(workExpIDs, woExpID)
				}
			}
		}
		err = uc.WorkExperienceRepository.DeleteNotInIDAndUserProfileID(workExpIDs, updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting work experiences: %s", err.Error())
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when deleting work experiences: " + err.Error())
		}
		if len(req.WorkExperiences) > 0 {
			for _, we := range req.WorkExperiences {
				if we.ID != nil {
					woExpID, err := uuid.Parse(*we.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing work experience ID: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing work experience ID: " + err.Error())
					}

					_, err = uc.WorkExperienceRepository.UpdateWorkExperience(&entity.WorkExperience{
						ID:             woExpID,
						UserProfileID:  updatedProfile.ID,
						CompanyName:    we.CompanyName,
						Name:           we.Name,
						YearExperience: we.YearExperience,
						JobDescription: we.JobDescription,
						Certificate:    we.CertificatePath,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when updating work experience: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when updating work experience: " + err.Error())
					}
				} else {
					_, err = uc.WorkExperienceRepository.CreateWorkExperience(&entity.WorkExperience{
						UserProfileID:  updatedProfile.ID,
						CompanyName:    we.CompanyName,
						Name:           we.Name,
						YearExperience: we.YearExperience,
						JobDescription: we.JobDescription,
						Certificate:    we.CertificatePath,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating work experience: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when creating work experience: " + err.Error())
					}
				}
			}
		}

		var eduIDs []uuid.UUID
		if len(req.Educations) > 0 {
			for _, ed := range req.Educations {
				if ed.ID != nil {
					eduID, err := uuid.Parse(*ed.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing education ID: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing education ID: " + err.Error())
					}
					eduIDs = append(eduIDs, eduID)
				}
			}
		}

		err = uc.EducationRepository.DeleteNotInIDAndUserProfileID(eduIDs, updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting educations: %s", err.Error())
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when deleting educations: " + err.Error())
		}

		if len(req.Educations) > 0 {
			for _, ed := range req.Educations {
				if ed.ID == nil {
					parsedEndDate, err := time.Parse("2006-01-02", ed.EndDate)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing end date: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing end date: " + err.Error())
					}
					var gpa *float64
					if ed.Gpa != nil {
						gpa = ed.Gpa
					}
					_, err = uc.EducationRepository.CreateEducation(&entity.Education{
						UserProfileID:  updatedProfile.ID,
						SchoolName:     ed.SchoolName,
						Major:          ed.Major,
						GraduateYear:   ed.GraduateYear,
						Certificate:    ed.CertificatePath,
						EndDate:        parsedEndDate,
						EducationLevel: entity.EducationLevelEnum(ed.EducationLevel),
						Gpa:            gpa,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating education: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when creating education: " + err.Error())
					}
				} else {
					eduID, err := uuid.Parse(*ed.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing education ID: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing education ID: " + err.Error())
					}
					parsedEndDate, err := time.Parse("2006-01-02", ed.EndDate)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing end date: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing end date: " + err.Error())
					}
					var gpa *float64
					if ed.Gpa != nil {
						gpa = ed.Gpa
					}
					_, err = uc.EducationRepository.UpdateEducation(&entity.Education{
						ID:             eduID,
						UserProfileID:  updatedProfile.ID,
						SchoolName:     ed.SchoolName,
						Major:          ed.Major,
						GraduateYear:   ed.GraduateYear,
						Certificate:    ed.CertificatePath,
						EndDate:        parsedEndDate,
						EducationLevel: entity.EducationLevelEnum(ed.EducationLevel),
						Gpa:            gpa,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating education: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when creating education: " + err.Error())
					}
				}
			}
		}

		var skillIDs []uuid.UUID
		if len(req.Skills) > 0 {
			for _, s := range req.Skills {
				if s.ID != nil {
					skillID, err := uuid.Parse(*s.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing skill ID: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing skill ID: " + err.Error())
					}
					skillIDs = append(skillIDs, skillID)
				}
			}
		}

		err = uc.SkillRepository.DeleteNotInIDAndUserProfileID(skillIDs, updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting skills: %s", err.Error())
			return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when deleting skills: " + err.Error())
		}

		if len(req.Skills) > 0 {
			for _, s := range req.Skills {
				if s.ID == nil {
					var level *int
					if s.Level != nil {
						level = s.Level
					}
					_, err := uc.SkillRepository.CreateSkill(&entity.Skill{
						UserProfileID: updatedProfile.ID,
						Name:          s.Name,
						Description:   s.Description,
						Certificate:   s.CertificatePath,
						Level:         level,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating skill: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when creating skill: " + err.Error())
					}
				} else {
					skillID, err := uuid.Parse(*s.ID)
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing skill ID: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when parsing skill ID: " + err.Error())
					}
					var level *int
					if s.Level != nil {
						level = s.Level
					}
					_, err = uc.SkillRepository.UpdateSkill(&entity.Skill{
						ID:            skillID,
						UserProfileID: updatedProfile.ID,
						Name:          s.Name,
						Description:   s.Description,
						Certificate:   s.CertificatePath,
						Level:         level,
					})
					if err != nil {
						uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when updating skill: %s", err.Error())
						return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when updating skill: " + err.Error())
					}
				}
			}
		}
	}
	// find user profile by user ID
	profile, err := uc.Repository.FindByUserID(userID)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: %s", err.Error())
		return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: " + err.Error())
	}

	if profile == nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] user profile not found")
		return "Failed", errors.New("[UserProfileUseCase.FillUserProfile] user profile not found")
	}

	return "Success", nil
}

func (uc *UserProfileUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.UserProfileResponse, int64, error) {
	profiles, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FindAllPaginated] error when finding all paginated: %s", err.Error())
		return nil, 0, errors.New("[UserProfileUseCase.FindAllPaginated] error when finding all paginated: " + err.Error())
	}

	if profiles == nil {
		uc.Log.Errorf("[UserProfileUseCase.FindAllPaginated] user profiles not found")
		return nil, 0, errors.New("[UserProfileUseCase.FindAllPaginated] user profiles not found")
	}

	resp := make([]response.UserProfileResponse, 0)
	for _, profile := range *profiles {
		converted, err := uc.DTO.ConvertEntityToResponse(&profile)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FindAllPaginated] error when converting entity to response: %s", err.Error())
			return nil, 0, errors.New("[UserProfileUseCase.FindAllPaginated] error when converting entity to response: " + err.Error())
		}
		resp = append(resp, *converted)
	}
	return &resp, total, nil
}

func (uc *UserProfileUseCase) FindByID(id uuid.UUID) (*response.UserProfileResponse, error) {
	profile, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FindByID] error when finding user profile by ID: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.FindByID] error when finding user profile by ID: " + err.Error())
	}

	if profile == nil {
		uc.Log.Errorf("[UserProfileUseCase.FindByID] user profile not found")
		return nil, errors.New("[UserProfileUseCase.FindByID] user profile not found")
	}

	return uc.DTO.ConvertEntityToResponse(profile)
}

func (uc *UserProfileUseCase) FindByUserID(userID uuid.UUID) (*response.UserProfileResponse, error) {
	profile, err := uc.Repository.FindByUserID(userID)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FindByUserID] error when finding user profile by user ID: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.FindByUserID] error when finding user profile by user ID: " + err.Error())
	}

	if profile == nil {
		uc.Log.Errorf("[UserProfileUseCase.FindByUserID] user profile not found")
		return nil, errors.New("[UserProfileUseCase.FindByUserID] user profile not found")
	}

	return uc.DTO.ConvertEntityToResponse(profile)
}

func (uc *UserProfileUseCase) UpdateStatusUserProfile(req *request.UpdateStatusUserProfileRequest) (*response.UserProfileResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.UpdateStatusUserProfile] error when parsing ID: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.UpdateStatusUserProfile] error when parsing ID: " + err.Error())
	}

	profile, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.UpdateStatusUserProfile] error when finding user profile by ID: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.UpdateStatusUserProfile] error when finding user profile by ID: " + err.Error())
	}
	if profile == nil {
		uc.Log.Errorf("[UserProfileUseCase.UpdateStatusUserProfile] user profile not found")
		return nil, errors.New("[UserProfileUseCase.UpdateStatusUserProfile] user profile not found")
	}

	updatedProfile, err := uc.Repository.UpdateUserProfile(&entity.UserProfile{
		ID:     parsedID,
		Status: entity.UserStatus(req.Status),
	})
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.UpdateStatusUserProfile] error when updating user profile: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.UpdateStatusUserProfile] error when updating user profile: " + err.Error())
	}

	return uc.DTO.ConvertEntityToResponse(updatedProfile)
}

func (uc *UserProfileUseCase) DeleteUserProfile(id uuid.UUID) error {
	err := uc.Repository.DeleteUserProfile(id)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.DeleteUserProfile] error when deleting user profile: %s", err.Error())
		return errors.New("[UserProfileUseCase.DeleteUserProfile] error when deleting user profile: " + err.Error())
	}
	return nil
}

func (uc *UserProfileUseCase) UpdateAvatar(id uuid.UUID, avatarPath string) (*response.UserProfileResponse, error) {
	profile, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.UpdateAvatar] error when finding user profile by ID: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.UpdateAvatar] error when finding user profile by ID: " + err.Error())
	}
	if profile == nil {
		uc.Log.Errorf("[UserProfileUseCase.UpdateAvatar] user profile not found")
		return nil, errors.New("[UserProfileUseCase.UpdateAvatar] user profile not found")
	}

	updatedProfile, err := uc.Repository.UpdateUserProfile(&entity.UserProfile{
		ID:     id,
		Avatar: avatarPath,
	})
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.UpdateAvatar] error when updating user profile: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.UpdateAvatar] error when updating user profile: " + err.Error())
	}

	return uc.DTO.ConvertEntityToResponse(updatedProfile)
}

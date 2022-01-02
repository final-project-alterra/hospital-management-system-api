package business

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/config"
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
	"github.com/golang-jwt/jwt"
)

type authBusiness struct {
	adminBusiness  admins.IBusiness
	doctorBusiness doctors.IBusiness
	nurseBusiness  nurses.IBusiness
}

func (a *authBusiness) Login(email string, password string) (string, error) {
	const op errors.Op = "auth.business.Login"
	var errMessage errors.ErrClientMessage = "Wrong email or password"

	// Check admin
	admin, err := a.adminBusiness.FindAdminByEmail(email)
	if err == nil {
		doesMatch := hash.Validate(admin.Password, password)
		if !doesMatch {
			err = errors.New("Wrong email or password")
			return "", errors.E(err, op, errMessage, errors.KindUnauthorized)
		}

		token, err := a.createToken(admin.ID, "admin")
		if err != nil {
			return "", errors.E(err, op)
		}
		return token, nil
	}

	if errors.Kind(err) != errors.KindNotFound {
		return "", errors.E(err, op)
	}

	// Check doctor
	doctor, err := a.doctorBusiness.FindDoctorByEmail(email)
	if err == nil {
		doesMatch := hash.Validate(doctor.Password, password)
		if !doesMatch {
			err = errors.New("Wring email or password")
			return "", errors.E(err, op, errMessage, errors.KindUnauthorized)
		}
		token, err := a.createToken(doctor.ID, "doctor")
		if err != nil {
			return "", errors.E(err, op)
		}
		return token, nil
	}

	if errors.Kind(err) != errors.KindNotFound {
		return "", errors.E(err, op)
	}

	// Check nurse
	nurse, err := a.nurseBusiness.FindNurseByEmail(email)
	if err == nil {
		doesMatch := hash.Validate(nurse.Password, password)
		if !doesMatch {
			err = errors.New("Wring email or password")
			return "", errors.E(err, op, errMessage, errors.KindUnauthorized)
		}
		token, err := a.createToken(nurse.ID, "nurse")
		if err != nil {
			return "", errors.E(err, op)
		}
		return token, nil
	}

	if errors.Kind(err) != errors.KindNotFound {
		return "", errors.E(err, op)
	}

	err = errors.New("Account does not exsist")
	errMessage = "Account does not exsist"
	return "", errors.E(err, op, errMessage, errors.KindNotFound)
}

func (a *authBusiness) createToken(id int, role string) (string, error) {
	const op errors.Op = "auth.business.createToken"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,
		"role":   role,
		"exp":    time.Now().Add(6 * time.Hour).Unix(),
	})

	token, err := claims.SignedString([]byte(config.ENV.JWT_SECRET))
	if err != nil {
		return "", errors.E(err, op, errMessage, errors.KindServerError)
	}
	return token, nil
}

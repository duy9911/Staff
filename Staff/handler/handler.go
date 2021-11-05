package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/duy9911/Staff/handler/logger"
	"github.com/duy9911/Staff/handler/redis"
	"github.com/duy9911/Staff/models"
)

type Id struct {
	Domain string
	Id     int
}

func CreateStaff(staff models.Staff) {
	validatedStaff, err := prepareStaff(staff)
	if err != nil {
		logger.Logger("error prepare data staff ", err)
		return
	}
	redis.Set(validatedStaff.ID, validatedStaff)
	logger.Logger("Create successfull staff's information ", validatedStaff)
}

func ReturnStaffs(domain string) {
	staffs, err := redis.GetAll(domain)
	if err != nil {
		logger.Logger("error get all ", err)
		return
	}
	for _, v := range staffs {
		logger.Logger("Staff: ", v)
	}
}

func UpdateStaff(key string, staff models.Staff) {
	updateStaff, err := prepareUpdate(key, staff)
	if err != nil {
		logger.Logger("error prepare update staff", err)
		return
	}

	if err := redis.Set(key, updateStaff); err != nil {
		logger.Logger("error update", err)
		return
	}
	logger.Logger("Updated ", updateStaff)
}

func Deletestaff(key string) {
	_, err := redis.Get(key)
	if err != nil {
		logger.Logger("error key", errors.New("doesn't match any key"))
		return
	}

	if err := redis.Delete(key); err != nil {
		logger.Logger("error delete", err)
		return
	}
	logger.Logger("deleted ", key)
}

func validateStaff(staff models.Staff) error {
	switch {
	case staff.Name == " ":
		return errors.New("name can not empty")
	case staff.Gender == " ":
		return errors.New("gender can not empty")
	case staff.Salary == 0:
		return errors.New("salary can not empty")
	}

	const layout = "2006-01-02"
	now := time.Now()
	workingAge := now.AddDate(-18, 0, 0)
	dob, err := time.Parse(layout, staff.Dob)

	if err != nil {
		return errors.New("format day of birth is wrong")
	}
	if !dob.Before(workingAge) {
		return errors.New("opp! your staff too young, must be more than 18 years old ")
	}
	return nil
}

func prepareUpdate(key string, staff models.Staff) (models.Staff, error) {
	updateStaff := models.Staff{}

	if _, err := redis.Get(key); err != nil {
		return updateStaff, errors.New("doesn't match any key ")
	}
	err := validateStaff(staff)
	if err != nil {
		return updateStaff, err
	}
	updateStaff = models.Staff{
		ID:     key,
		Name:   staff.Name,
		Dob:    staff.Dob,
		Gender: staff.Gender,
		Salary: staff.Salary,
	}
	return updateStaff, err
}

func prepareStaff(s models.Staff) (models.Staff, error) {
	staff := models.Staff{}
	err := validateStaff(s)
	if err != nil {
		return s, err
	}
	domain := "staff"
	nextKey, err := GenerateId(domain)
	if err != nil {
		return staff, err
	}
	staff = models.Staff{
		ID:     nextKey,
		Name:   s.Name,
		Gender: s.Gender,
		Salary: s.Salary,
		Dob:    s.Dob,
	}
	return staff, nil
}

// staff-id: N+1
// team-id: N+1

// staffs: Hashes
// teams: Hashes

func GenerateId(domain string) (string, error) {
	keyLasest := "lastId"
	idLatest, err := redis.Get(keyLasest)

	// check lastest key is empty or not
	if err == nil {
		idStruct := Id{}
		err := json.Unmarshal([]byte(idLatest), &idStruct)
		if err != nil {
			return domain, err
		}
		nextId := &Id{
			Domain: idStruct.Domain,
			Id:     idStruct.Id + 1,
		}
		redis.Set(keyLasest, nextId)
		concatenated := nextId.Domain + strconv.Itoa(nextId.Id)
		return concatenated, nil
	}

	newId := Id{
		Domain: domain,
		Id:     1,
	}
	redis.Set(keyLasest, newId)
	concatenated := newId.Domain + strconv.Itoa(newId.Id)
	return concatenated, nil
}

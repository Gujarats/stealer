package stealer

import "errors"

type Stealer struct {
	Datas map[string][]string
}

// Steal some variable data from path phpfile
func Steal(path string) (error, *Stealer) {
	var stealer Stealer
	var err error
	err, stealer.Datas = ReadFile(path)
	return err, &stealer
}

// save all the variables and its values to new path
// TODO : write test case where the path is not exist and need to create the path first
// lets say path/to/specific/file.go need to create the folder first
func (s *Stealer) Save(path, packageName string) error {
	var err error
	if s.Datas == nil {
		err = errors.New("Stealer Datas are empty")
		return err
	}
	err = WriteFile(path, packageName, s.Datas)
	if err != nil {
		return err
	}
	return nil
}

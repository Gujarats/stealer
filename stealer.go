package stealer

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
func (s *Stealer) Save(path, packageName string) error {
	var err error
	err = WriteFile(path, packageName, s.Datas)
	return err
}

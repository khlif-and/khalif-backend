package enums

// ShahihStatus represents the authenticity status of a hadist
type ShahihStatus string

const (
	ShahihStatusShahih  ShahihStatus = "shahih"
	ShahihStatusHasan   ShahihStatus = "hasan"
	ShahihStatusDhaif   ShahihStatus = "dhaif"
	ShahihStatusMaudhu  ShahihStatus = "maudhu"
)

func (s ShahihStatus) IsValid() bool {
	switch s {
	case ShahihStatusShahih, ShahihStatusHasan, ShahihStatusDhaif, ShahihStatusMaudhu:
		return true
	}
	return false
}

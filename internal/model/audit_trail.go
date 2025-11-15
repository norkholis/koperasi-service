package model

import (
	"time"

	"gorm.io/gorm"
)

// AuditTrail represents a comprehensive audit log entry for all system changes
type AuditTrail struct {
	gorm.Model
	UserID      uint      `gorm:"not null;index" json:"user_id"`                       // Who made the change
	Action      string    `gorm:"type:varchar(50);not null;index" json:"action"`       // CREATE, UPDATE, DELETE, LOGIN, etc.
	EntityTable string    `gorm:"type:varchar(50);not null;index" json:"entity_table"` // Which table was affected
	RecordID    uint      `gorm:"index" json:"record_id"`                              // ID of the affected record
	OldValues   string    `gorm:"type:text" json:"old_values"`                         // JSON of old values (for UPDATE/DELETE)
	NewValues   string    `gorm:"type:text" json:"new_values"`                         // JSON of new values (for CREATE/UPDATE)
	IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address"`                  // Client IP address
	UserAgent   string    `gorm:"type:text" json:"user_agent"`                         // Client user agent
	Description string    `gorm:"type:text" json:"description"`                        // Human-readable description
	Timestamp   time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"timestamp"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for AuditTrail model
func (AuditTrail) TableName() string {
	return "audit_trails"
}

// TransactionHistory represents a comprehensive financial transaction log
type TransactionHistory struct {
	gorm.Model
	UserID          uint       `gorm:"not null;index" json:"user_id"`                           // User involved in transaction
	TransactionType string     `gorm:"type:varchar(50);not null;index" json:"transaction_type"` // SIMPANAN, PINJAMAN, ANGSURAN, SHU
	ReferenceTable  string     `gorm:"type:varchar(50);not null" json:"reference_table"`        // Source table (simpanan, pinjaman, angsuran, etc.)
	ReferenceID     uint       `gorm:"not null;index" json:"reference_id"`                      // ID of the source record
	Amount          float64    `gorm:"type:decimal(15,2);not null" json:"amount"`               // Transaction amount
	BalanceBefore   float64    `gorm:"type:decimal(15,2)" json:"balance_before"`                // Balance before transaction
	BalanceAfter    float64    `gorm:"type:decimal(15,2)" json:"balance_after"`                 // Balance after transaction
	Status          string     `gorm:"type:varchar(20);not null;index" json:"status"`           // PENDING, COMPLETED, CANCELLED, VERIFIED
	TransactionDate time.Time  `gorm:"default:CURRENT_TIMESTAMP;index" json:"transaction_date"`
	VerifiedBy      uint       `gorm:"index" json:"verified_by"`     // Admin who verified (if applicable)
	VerifiedAt      *time.Time `json:"verified_at"`                  // When it was verified
	Description     string     `gorm:"type:text" json:"description"` // Transaction description
	Metadata        string     `gorm:"type:text" json:"metadata"`    // Additional JSON metadata
	User            User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	VerifiedByUser  User       `gorm:"foreignKey:VerifiedBy" json:"verified_by_user,omitempty"`
}

// TableName specifies the table name for TransactionHistory model
func (TransactionHistory) TableName() string {
	return "transaction_histories"
}

// SystemReport represents comprehensive system reports for admin analysis
type SystemReport struct {
	gorm.Model
	ReportType       string    `gorm:"type:varchar(50);not null;index" json:"report_type"` // DAILY, WEEKLY, MONTHLY, YEARLY, CUSTOM
	StartDate        time.Time `gorm:"not null;index" json:"start_date"`                   // Report period start
	EndDate          time.Time `gorm:"not null;index" json:"end_date"`                     // Report period end
	GeneratedBy      uint      `gorm:"not null" json:"generated_by"`                       // Admin who generated
	ReportData       string    `gorm:"type:longtext" json:"report_data"`                   // JSON report data
	TotalUsers       int       `json:"total_users"`                                        // Summary statistics
	TotalSimpanan    float64   `gorm:"type:decimal(15,2)" json:"total_simpanan"`           // Total savings
	TotalPinjaman    float64   `gorm:"type:decimal(15,2)" json:"total_pinjaman"`           // Total loans
	TotalAngsuran    float64   `gorm:"type:decimal(15,2)" json:"total_angsuran"`           // Total installments
	TotalSHU         float64   `gorm:"type:decimal(15,2)" json:"total_shu"`                // Total SHU distributed
	Status           string    `gorm:"type:varchar(20);default:'GENERATED'" json:"status"` // GENERATED, ARCHIVED
	GeneratedBy_User User      `gorm:"foreignKey:GeneratedBy" json:"generated_by_user,omitempty"`
}

// TableName specifies the table name for SystemReport model
func (SystemReport) TableName() string {
	return "system_reports"
}

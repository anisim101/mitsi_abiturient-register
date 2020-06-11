package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"mitsoСhat/internal/app/model"
)

type Store struct {
	db             *sql.DB
}

func New(db *sql.DB) *Store {

	sqlStore :=  &Store{
		db: db,
	}
	return  sqlStore
}


func (s *Store)AddPersson(uni *model.Uni) error {


	_,err := s.db.Query("insert into univercity (first_name, second_name, third_name, latin_first_name, latin_second_name, gender, date_of_birth, is_foreigner, nationality, family_status, region, post_index, birth_place, register_address, home_address, live_address, phone_number, home_phone_number, document_type, serial_and_passport_number, issued_by, when_issud, valid_until, personal_nomer, form_studing, faculty, specialization, type_education, educational_institution, education_specialization, graduated_institution_year, subject_1_name, subject_2_name, subject_3_name, subject_1_rating, subject_2_rating, subject_3_rating, subject_1_setificate, subject_2_setificate, subject_3_setificate, forigen_lang_name, forigen_lang_rating, document_type_education, is_selskaya_shkola, is_gorodskaya, a2019, medal, ssuz, ptu, document_education_number_nomer, document_education_at, average_mark, olimp, spec_fond_rb, olimp_sport, perviy_razrad, lgotniy_kredit, hostel, disabled, chernobil, additional_information, email, privileges, info_privileges, has_father, second_name_father, third_name_father, adress_father, work_adress_father, contact_father, has_mother, second_mother_name, mother_first_name, third_mother_name, work_mother_address, mother_contact_info, owner_passport, owner_passport_number, owner_passport_id, owner_passport_kto_vidal, owner_passport_kogda_vidal, owner_passport_srok, pasport_photo, ct_1_photo, ct_2_photo, ct_3_photo, ct_4_photo, owner_pasport_photo, education_document_photo, first_name_father, adress_mother, sum_marks)  values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92 )",
		&uni.FirstName,
		&uni.SecondName,
		&uni.ThirdName,
		&uni.LatinFirstName,
		&uni.LatinSecondName,
		&uni.Gender,
		&uni.DateOfBirth,
		&uni.IsForeigner,
		&uni.Nationality,
		&uni.FamilyStatus,
		&uni.Region,
		&uni.PostIndex,
		&uni.BirthPlace,
		&uni.RegisterAddress,
		&uni.HomeAddress,
		&uni.LiveAddress,
		&uni.PhoneNumber,
		&uni.HomePhoneNumber,
		&uni.DocumentType,
		&uni.SerialAndPassportNumber,
		&uni.IssuedBy,
		&uni.WhenIssud,
		&uni.ValidUntil,
		&uni.PersonalNomer,
		&uni.FormStuding,
		&uni.Faculty,
		&uni.Specialization,
		&uni.TypeEducation,
		&uni.EducationalInstitution,
		&uni.EducationSpecialization,
		&uni.GraduatedInstitutionYear,
		&uni.Subject1_Name,
		&uni.Subject2_Name,
		&uni.Subject3_Name,
		&uni.Subject1_Rating,
		&uni.Subject2_Rating,
		&uni.Subject3_Rating,
		&uni.Subject1_Setificate,
		&uni.Subject2_Setificate,
		&uni.Subject3_Setificate,
		&uni.ForigenLangName,
		&uni.ForigenLangRating,
		&uni.DocumentTypeEducation,
		&uni.IsSelskayaShkola,
		&uni.IsGorodskaya,
		&uni.A2019,
		&uni.Medal,
		&uni.Ssuz,
		&uni.Ptu,
		&uni.DocumentEducationNumberNomer,
		&uni.DocumentEducationAt,
		&uni.AverageMark,
		&uni.Olimp,
		&uni.SpecFondRb,
		&uni.OlimpSport,
		&uni.PerviyRazrad,
		&uni.LgotniyKredit,
		&uni.Hostel,
		&uni.Disabled,
		&uni.Chernobil,
		&uni.AdditionalInformation,
		&uni.Email,
		&uni.Privileges,
		&uni.InfoPrivileges,
		&uni.HasFather,
		&uni.SecondNameFather,
		&uni.ThirdNameFather,
		&uni.AdressFather,
		&uni.WorkAdressFather,
		&uni.ContactFather,
		&uni.HasMother,
		&uni.SecondMotherName,
		&uni.MotherFirstName,
		&uni.ThirdMotherName,
		&uni.WorkMotherAddress,
		&uni.MotherContactInfo,
		&uni.OwnerPassport,
		&uni.OwnerPassportNumber,
		&uni.OwnerPassportID,
		&uni.OwnerPassportKtoVidal,
		&uni.OwnerPassportKogdaVidal,
		&uni.OwnerPassportSrok,
		&uni.PasportPhoto,
		&uni.CT1_Photo,
		&uni.CT2_Photo,
		&uni.CT3_Photo,
		&uni.CT4_Photo,
		&uni.OwnerPasportPhoto,
		&uni.EducationDocumentPhoto,
		&uni.FirstNameFather,
		&uni.AdressMother,

		&uni.SumMarks,
	)

	return err
}






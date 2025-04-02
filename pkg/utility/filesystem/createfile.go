package filesystem

import (
	// "RESTAPI/usecase"
	"bytes"
	"fmt"
	"go-clean-arch/structure/response"
	"log"
	"os"
	// "time"

	"github.com/signintech/gopdf"
)

func CreatePDF(data response.OutsideResponse) ([]byte, string, error) {
    pdf := gopdf.GoPdf{}
	// กำหนดการตั้งค่า PDF เช่น ขนาดหน้ากระดาษ A4
	landscapeSize := gopdf.Rect{W: 841.89, H: 595.28}
	pdf.Start(gopdf.Config{PageSize: landscapeSize})
	pdf.AddPage()

    watermark(&pdf,landscapeSize)

	err := pdf.AddTTFFont("THSarabunNew", "./pkg/utility/filesystem/assets/THSarabunNew/THSarabunNew.ttf")
	if err != nil {
        log.Fatalf("Error adding font: %s", err)
	}
    
	err = pdf.AddTTFFont("THSarabunNewBold", "./pkg/utility/filesystem/assets/THSarabunNew/THSarabunNew Bold.ttf")
	if err != nil {
        log.Fatalf("Error adding font: %s", err)
	}
    
    header(&pdf,data)

	table(&pdf,data)

	drawImage(&pdf)

   
   // สร้างไฟล์ PDF โดยการเขียนลงในไฟล์ชั่วคราว (temporary file)
	tempFile, err := os.CreateTemp("", "tempfile-*.pdf")
	if err != nil {
		return nil, "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer tempFile.Close()

	// เขียน PDF ลงในไฟล์ชั่วคราว
	err = pdf.WritePdf(tempFile.Name())
	if err != nil {
		return nil, "", fmt.Errorf("error creating PDF: %v", err)
	}

	// อ่านข้อมูล PDF จากไฟล์ชั่วคราวไปยัง memory buffer
	var buf bytes.Buffer
	_, err = buf.ReadFrom(tempFile)
	if err != nil {
		return nil, "", fmt.Errorf("error reading PDF file to memory buffer: %v", err)
	}

	// คืนค่าเป็น []byte และชื่อไฟล์
	// fileName := fmt.Sprintf("แบบฟอร์มบันทึกกิจกรรม.pdf")
	fileName := "แบบฟอร์มบันทึกกิจกรรม.pdf"
	return buf.Bytes(), fileName, nil

}

func watermark(pdf *gopdf.GoPdf,landscapeSize gopdf.Rect){
    imageWidth := 247.8
	imageHeight := 464.1

	pageWidth, pageHeight := landscapeSize.W, landscapeSize.H
	x := (pageWidth - imageWidth) / 2
	y := (pageHeight - imageHeight) / 2

	err := pdf.Image("./pkg/utility/filesystem/assets/image/bg2.png", x, y, &gopdf.Rect{W: imageWidth, H: imageHeight})
	if err != nil {
		log.Fatal(err)
	}
}

func header(pdf *gopdf.GoPdf,data response.OutsideResponse){

	pdf.SetFont("THSarabunNewBold", "", 20)
	pdf.SetXY(130, 50)
	pdf.Cell(nil, "แบบบันทึกการเข้าร่วมกิจกรรม/โครงการจิตอาสา ประจำปีการศึกษา............... มหาวิทยาลัยเทคโนโลยีราชมงคลอีสาน")

	pdf.SetFont("THSarabunNewBold", "", 18)
	pdf.SetXY(150, 85)
	pdf.Cell(nil, "ชื่อ-สกุล................................................................... หมายเลขโทรศัพท์............................ รหัสนักศึกษา...............................")

    
	pdf.SetFont("THSarabunNewBold", "", 16)
	pdf.SetXY(510, 49)
	pdf.Cell(nil, string(data.SchoolYear))

	pdf.SetXY(196, 82)
	pdf.Cell(nil, data.Student.TitleName+data.Student.FirstName+" "+data.Student.LastName)
    
	pdf.SetXY(500, 82)
	pdf.Cell(nil,data.Student.Phone)

	pdf.SetXY(655, 82)
	pdf.Cell(nil, data.Student.Code)

	pdf.SetFont("THSarabunNewBold", "", 18)
	pdf.SetXY(180, 120)
	pdf.Cell(nil, "สาขา..................................................................... คณะ..................................................................................")

	pdf.SetFont("THSarabunNewBold", "", 16)
	pdf.SetXY(216, 117)
	pdf.Cell(nil, data.Student.BranchName)

	pdf.SetXY(456, 117)
	pdf.Cell(nil, data.Student.FacultyName)

    err := pdf.Image("./pkg/utility/filesystem/assets/image/logo.png", 40, 30, &gopdf.Rect{W: 53.1, H: 99.45})
	if err != nil {
		log.Fatal(err)
	}
	err = pdf.Image("./pkg/utility/filesystem/assets/image/logo2.png", 700, 80, &gopdf.Rect{W: 100, H: 100})
	if err != nil {
		log.Fatal(err)
	}
}


func table(pdf *gopdf.GoPdf,data response.OutsideResponse) {
	tableStartY := 180.0
	marginLeft := 51.0
	table := pdf.NewTableLayout(marginLeft, tableStartY, 30, 1)


	pdf.SetFont("THSarabunNewBold", "", 16)
	table.AddColumn("โครงการ/กิจกรรมจิตอาสา", 260, "left")
	table.AddColumn("วันเดือนปี ที่เข้าร่วม", 100, "center")
	table.AddColumn("สถานที่", 150, "center")
	table.AddColumn("เวลามา-เวลากลับ", 150, "center")
	table.AddColumn("จำนวนชั่งโมง", 80, "center")
	table.AddRow([]string{data.EventName, "" ,data.Location, "10:00-16:00", fmt.Sprint(data.WorkingHour)})
	table.DrawTable()

    pdf.SetFont("THSarabunNewBold", "", 18)

    lineText := "................................................................"
    topText := "(  "+data.Intendant+"  )"
    bottomText := "ผู้รับรองการเข้าร่วมโครงการ"

    lineTextWidth,_ := pdf.MeasureTextWidth(lineText)
    topTextWidth, _ := pdf.MeasureTextWidth(topText)
    bottomTextWidth, _ := pdf.MeasureTextWidth(bottomText)

	pdf.SetXY(570, 360)
	pdf.Cell(nil,lineText)

    newX := 570 + (lineTextWidth-topTextWidth)/2

    pdf.SetXY(newX, 390)
    pdf.Cell(nil, topText)

    // คำนวณตำแหน่ง X ใหม่สำหรับข้อความด้านล่าง
    newX = 570 + (lineTextWidth-bottomTextWidth)/2

    // ตั้งค่าตำแหน่งใหม่สำหรับข้อความด้านล่าง
    pdf.SetXY(newX, 420)
    pdf.Cell(nil, bottomText)

}

func drawImage(pdf *gopdf.GoPdf) {
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)
	pdf.RectFromUpperLeftWithStyle(50, 260, 464, 290, "D")

	pdf.SetFont("THSarabunNewBold", "", 16)
	pdf.SetXY(265, 390)
	pdf.Cell(nil, "ใส่รูปภาพ")

}

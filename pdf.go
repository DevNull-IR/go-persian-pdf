package main

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
	"os"
	"strings"
	"strconv"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jung-kurt/gofpdf"
	"go-census-report/farsi"
)

func reverseString(s string) string {
	fg := farsi.NewFarsi()
    return fg.PersiaText(string(s), "fa", "pChars", false)
}

func truncateText(text string, maxWidth float64, font, pdf *gofpdf.Tpl) string {
    runes := []rune(text)
    for len(runes) > 0 {
        sub := string(runes)
        width := pdf.GetStringWidth(sub)
        if width <= maxWidth {
            return sub
        }
        runes = runes[:len(runes)-1]
    }
    return ""
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func reverseSlice(s []string) []string {
    reversed := make([]string, len(s))
    for i, v := range s {
        reversed[len(s)-1-i] = v
    }
    return reversed
}


func main() {

    // خواندن اطلاعات اتصال
    dsnData, err := ioutil.ReadFile("connection.db")
    if err != nil {
        log.Fatal("خطا در خواندن connection.db:", err)
    }
    dsn := string(dsnData)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("خطا در اتصال به دیتابیس:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("اتصال به دیتابیس برقرار نیست:", err)
    }
    query := "select * from users;"

    rows, err := db.Query(query)
    if err != nil {
        log.Fatal("خطا در اجرای کوئری:", err)
    }
    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        log.Fatal("خطا در گرفتن نام ستون‌ها:", err)
    }
	// columns = reverseSlice(columns);

    values := make([]interface{}, len(columns))
    valuePtrs := make([]interface{}, len(columns))
    var data [][]string

    for rows.Next() {
        for i := range columns {
            valuePtrs[i] = &values[i]
        }
        if err := rows.Scan(valuePtrs...); err != nil {
            log.Fatal("خطا در اسکن ردیف:", err)
        }

        entry := make([]string, len(columns))
        for i, raw := range values {
            switch v := raw.(type) {
            case nil:
                entry[i] = "—"
            case []byte:
                entry[i] = string(v)
            default:
                entry[i] = fmt.Sprintf("%v", v)
            }
        }
        data = append(data, entry)
    }

    if err = rows.Err(); err != nil {
        log.Fatal("خطا در حین خواندن ردیف‌ها:", err)
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetMargins(10, 10, 10)
    pdf.SetAutoPageBreak(true, 10)

    pdf.AddUTF8Font("vazir", "", "vazir.ttf") 
    pdf.SetFont("vazir", "", 10)

    pageWidth := 190.0
    colCount := float64(len(columns))
    colWidth := pageWidth / colCount

    rowHeight := 5.0

    pdf.AddPage()

    cell := func(text string, width float64, border int, fill bool) {
        if fill {
            pdf.SetFillColor(240, 240, 240)
        } else {
            pdf.SetFillColor(255, 255, 255)
        }
        text = reverseString(text)

        pdf.CellFormat(width, rowHeight, text, strconv.Itoa(border), 0, "R", fill, 0, "")
    }

drawHeader := func() {
    pdf.SetFont("vazir", "", 5)
    pdf.SetFillColor(240, 240, 240)
    for _, col := range columns {
        cell(col, colWidth, 1, true)
    }
    pdf.Ln(rowHeight)
}

pdf.SetFont("vazir", "", 5)

  const rowsPerPage = 50

for i := 0; i < len(data); i += rowsPerPage {
    if i > 0 {
        pdf.AddPage()
    }

    drawHeader()

    end := i + rowsPerPage
    if end > len(data) {
        end = len(data)
    }
    for _, row := range data[i:end] {
        for _, cellText := range row {
            cell(cellText, colWidth, 1, false)
        }
        pdf.Ln(rowHeight)
    }
}

    pdf.SetFont("vazir", "", 5)
    for _, row := range data {
        for _, cellText := range row {
            cell(cellText, colWidth, 1, false)
        }
        pdf.Ln(rowHeight)
    }

    err = pdf.OutputFileAndClose("output.pdf")
    if err != nil {
        log.Fatal("خطا در ذخیره PDF:", err)
    }

    fmt.Println("✅ فایل PDF با موفقیت ایجاد شد: output.pdf")
}

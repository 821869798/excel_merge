package merge

import (
	"github.com/xuri/excelize/v2"
	"testing"
)

func TestInsertSheet(t *testing.T) {
	excel, err := excelize.OpenFile("D:/test.xlsx")
	if err != nil {
		t.Fatal(err)
	}

	defer excel.Close()

	err = ExcelInsertSheet(excel, 0, "test_new")
	if err != nil {
		t.Fatal(err)
	}

	err = excel.SaveAs("D:/test_new.xlsx")
	if err != nil {
		t.Fatal(err)
	}

}

func TestSwapSheet(t *testing.T) {
	excel, err := excelize.OpenFile("D:/test.xlsx")
	if err != nil {
		t.Fatal(err)
	}

	defer excel.Close()

	err = ExcelSwapSheetByName(excel, "Sheet2", "Sheet3")
	if err != nil {
		t.Fatal(err)
	}

	err = excel.SaveAs("D:/test_new2.xlsx")
	if err != nil {
		t.Fatal(err)
	}

}

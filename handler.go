package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"no.name/qldt"
)

const l02_01_2006 = "02/01/2006"
const l2006_01_02T04_05_06 = "2006-01-02T04:05:06"

func ensureAuthenticated(w http.ResponseWriter, r *http.Request) (*qldt.TokenResponse, error) {
	res, err := qldt.FetchToken(r)

	if err != nil {
        // For browser to request basic credentials from user.
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)

		return nil, err
	}

    return res, nil
}

func examHandler(w http.ResponseWriter, r *http.Request) {
    tokenResp, err := ensureAuthenticated(w, r)
    if err != nil {
        fmt.Fprintf(w, "[ERROR]: %v\n", err)
        return
    }

    lichThiResp, err := qldt.FetchLichThi(tokenResp.AccessToken, tokenResp.Name)
    if err != nil {
        fmt.Fprintf(w, "[ERROR]: %v\n", err)
        return
    }

    for _, lichThi := range lichThiResp.Data.DsLichThi {
        examBytes, err := generateExamSchedule(lichThi)
        if err != nil {
            fmt.Fprintf(w, "[ERROR]: %v\n", err)
            return
        }

        fmt.Fprintf(w, "%v", string(examBytes))
    }
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
    tokenResp, err := ensureAuthenticated(w, r)
    if err != nil {
        fmt.Fprintf(w, "[ERROR]: %v\n", err)
        return
    }

	scheduleResp, err := qldt.FetchDSTKB(tokenResp.AccessToken, tokenResp.Name)
	if err != nil {
        fmt.Fprintf(w, "[ERROR]: %v\n", err)
		return
	}

	t := r.PathValue("t")

	now := time.Now()
	if len(t) > 0 {
		rt, err := strconv.Atoi(t)
		if err != nil {
            fmt.Fprintf(w, "[ERROR]: %v\n", err)
			return
		}

		now = now.Add(time.Duration(rt) * 7 * 24 * time.Hour)
	}

	var curTuanTKB qldt.TuanTKB
	for _, value := range scheduleResp.Data.DSTuanTKB {
		startTime, err := time.Parse(l02_01_2006, value.NgayBatDau)
		if err != nil {
            fmt.Fprintf(w, "[ERROR]: %v\n", err)
			continue
		}

		endTime, err := time.Parse(l02_01_2006, value.NgayKetThuc)
		if err != nil {
            fmt.Fprintf(w, "[ERROR]: %v\n", err)
			continue
		}

		if now.After(startTime) && now.Before(endTime.Add(24*time.Hour-time.Nanosecond)) {
			curTuanTKB = value
			break
		}
	}

	tempTKB := make([]qldt.ThoiKhoaBieu, 0)
	for _, value := range curTuanTKB.DSThoiKhoaBieu {
		if len(tempTKB) == 0 {
			tempTKB = append(tempTKB, value)
			continue
		}

		ngayHoc, err := time.Parse(l2006_01_02T04_05_06, value.NgayHoc)
		if err != nil {
            fmt.Fprintf(w, "[ERROR]: %v\n", err)
			continue
		}

		currNgayHoc, err := time.Parse(l2006_01_02T04_05_06, tempTKB[0].NgayHoc)
		if err != nil {
            fmt.Fprintf(w, "[ERROR]: %v\n", err)
			continue
		}

		if currNgayHoc.Equal(ngayHoc) {
			tempTKB = append(tempTKB, value)
			continue
		}

		scheduleBytes, err := generateSchedule(tempTKB, scheduleResp)
        if err != nil {
            fmt.Fprintf(w, "[ERROR]: %v\n", err)
            return
        }

		tempTKB = append([]qldt.ThoiKhoaBieu{}, value)
		fmt.Fprintf(w, "%v", string(scheduleBytes))
	}

    scheduleBytes, err := generateSchedule(tempTKB, scheduleResp)
    if err != nil {
        fmt.Fprintf(w, "[ERROR]: %v\n", err)
        return
    }

	fmt.Fprintf(w, "%v", string(scheduleBytes))
}

func generateExamSchedule(lichThi qldt.LichThi) ([]rune, error) {
    ngayThi, err := time.Parse(l02_01_2006, lichThi.NgayThi);
    if err != nil {
        return nil, err
    }

    examBytes := []rune(RefExam)

    var tmp bytes.Buffer
    fmt.Fprintf(&tmp, "%v %.02d %v", (ngayThi.Weekday()).String()[:3], ngayThi.Day(), ngayThi.Month())

    copy(examBytes[(ExamColums*ExamNgaySpot.Row)+ExamNgaySpot.StartIndex:(ExamColums*ExamNgaySpot.Row)+ExamNgaySpot.StartIndex+ExamNgaySpot.Length], []rune(tmp.String()))

    tmp.Reset()
    fmt.Fprintf(&tmp, "%v", lichThi.GioBatDau)
    copy(examBytes[(ExamColums*ExamThoiGianSpot.Row)+ExamThoiGianSpot.StartIndex:(ExamColums*ExamThoiGianSpot.Row)+ExamThoiGianSpot.StartIndex+ExamThoiGianSpot.Length], []rune(tmp.String()))

    tmp.Reset()
    fmt.Fprintf(&tmp, "%v", lichThi.SoPhut)
    copy(examBytes[(ExamColums*ExamSoPhutSpot.Row)+ExamSoPhutSpot.StartIndex:(ExamColums*ExamSoPhutSpot.Row)+ExamSoPhutSpot.StartIndex+ExamSoPhutSpot.Length], []rune(tmp.String()))

    tmp.Reset()
    fmt.Fprintf(&tmp, "%v", lichThi.TenMon)
    copy(examBytes[(ExamColums*ExamMonSpot.Row)+ExamMonSpot.StartIndex:(ExamColums*ExamMonSpot.Row)+ExamMonSpot.StartIndex+ExamMonSpot.Length], []rune(tmp.String()))

    tmp.Reset()
    fmt.Fprintf(&tmp, "%v", lichThi.MaPhong)
    copy(examBytes[(ExamColums*ExamPhongSpot.Row)+ExamPhongSpot.StartIndex:(ExamColums*ExamPhongSpot.Row)+ExamPhongSpot.StartIndex+ExamPhongSpot.Length], []rune(tmp.String()))

    tmp.Reset()
    fmt.Fprintf(&tmp, "%v", lichThi.HinhThucThi)
    copy(examBytes[(ExamColums*ExamHinhThucThiSpot.Row)+ExamHinhThucThiSpot.StartIndex:(ExamColums*ExamHinhThucThiSpot.Row)+ExamHinhThucThiSpot.StartIndex+ExamHinhThucThiSpot.Length], []rune(tmp.String()))

    return examBytes, nil
}

func generateSchedule(tkb []qldt.ThoiKhoaBieu, scheduleResp *qldt.ScheduleResponse) ([]rune, error) {
	ngayHoc, err := time.Parse(l2006_01_02T04_05_06, tkb[0].NgayHoc)
	if err != nil {
		return nil, err
	}

	scheduleBytes := []rune(RefTop)

	var tmp bytes.Buffer
	fmt.Fprintf(&tmp, "%v %.02d %v", (ngayHoc.Weekday() + 3).String()[:3], ngayHoc.Day(), ngayHoc.Month())

	copy(scheduleBytes[(Colums*NgaySpot.Row)+NgaySpot.StartIndex:(Colums*NgaySpot.Row)+NgaySpot.StartIndex+NgaySpot.Length], []rune(tmp.String()))

	for i := 0; i < len(tkb); i++ {
		var ref []rune

		tiet := tkb[i]

		thoiGianBatDau := scheduleResp.Data.DSTietTrongNgay[tiet.TietBatDau-1].GioBatDau
		thoiGianKetThuc := scheduleResp.Data.DSTietTrongNgay[tiet.TietBatDau+tiet.SoTiet-1].GioKetThuc

		if i != len(tkb)-1 {
			ref = []rune(RefMiddle)

			tmp.Reset()
			fmt.Fprintf(&tmp, "%v - %v", thoiGianBatDau, thoiGianKetThuc)
			copy(ref[(Colums*ThoiGianSpot.Row)+ThoiGianSpot.StartIndex:(Colums*ThoiGianSpot.Row)+ThoiGianSpot.StartIndex+ThoiGianSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.SoTiet)
			copy(ref[(Colums*TietSpot.Row)+TietSpot.StartIndex:(Colums*TietSpot.Row)+TietSpot.StartIndex+TietSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.TenMon)
			copy(ref[(Colums*MonSpot.Row)+MonSpot.StartIndex:(Colums*MonSpot.Row)+MonSpot.StartIndex+MonSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.MaCoSo)
			copy(ref[(Colums*CoSoSpot.Row)+CoSoSpot.StartIndex:(Colums*CoSoSpot.Row)+CoSoSpot.StartIndex+CoSoSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.MaPhong)
			copy(ref[(Colums*PhongSpot.Row)+PhongSpot.StartIndex:(Colums*PhongSpot.Row)+PhongSpot.StartIndex+PhongSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.TenGiangVien)
			copy(ref[(Colums*GiangVienSpot.Row)+GiangVienSpot.StartIndex:(Colums*GiangVienSpot.Row)+GiangVienSpot.StartIndex+GiangVienSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.MaNhom)
			copy(ref[(Colums*NhomSpot.Row)+NhomSpot.StartIndex:(Colums*NhomSpot.Row)+NhomSpot.StartIndex+NhomSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.MaToTh)
			copy(ref[(Colums*NhomTHSpot.Row)+NhomTHSpot.StartIndex:(Colums*NhomTHSpot.Row)+NhomTHSpot.StartIndex+NhomTHSpot.Length], []rune(tmp.String()))

		} else {
			ref = []rune(RefBottom)

			tmp.Reset()
			fmt.Fprintf(&tmp, "%v - %v", thoiGianBatDau, thoiGianKetThuc)
			copy(ref[(Colums*ThoiGianSpot.Row)+ThoiGianSpot.StartIndex:(Colums*ThoiGianSpot.Row)+ThoiGianSpot.StartIndex+ThoiGianSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.SoTiet)
			copy(ref[(Colums*TietSpot.Row)+TietSpot.StartIndex:(Colums*TietSpot.Row)+TietSpot.StartIndex+TietSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.TenMon)
			copy(ref[(Colums*MonSpot.Row)+MonSpot.StartIndex:(Colums*MonSpot.Row)+MonSpot.StartIndex+MonSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.MaCoSo)
			copy(ref[(Colums*CoSoSpot.Row)+CoSoSpot.StartIndex:(Colums*CoSoSpot.Row)+CoSoSpot.StartIndex+CoSoSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.MaPhong)
			copy(ref[(Colums*PhongSpot.Row)+PhongSpot.StartIndex:(Colums*PhongSpot.Row)+PhongSpot.StartIndex+PhongSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.TenGiangVien)
			copy(ref[(Colums*GiangVienSpot.Row)+GiangVienSpot.StartIndex:(Colums*GiangVienSpot.Row)+GiangVienSpot.StartIndex+GiangVienSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.MaNhom)
			copy(ref[(Colums*NhomSpot.Row)+NhomSpot.StartIndex:(Colums*NhomSpot.Row)+NhomSpot.StartIndex+NhomSpot.Length], []rune(tmp.String()))
			tmp.Reset()
			fmt.Fprintf(&tmp, "%v", tiet.MaToTh)
			copy(ref[(Colums*NhomTHSpot.Row)+NhomTHSpot.StartIndex:(Colums*NhomTHSpot.Row)+NhomTHSpot.StartIndex+NhomTHSpot.Length], []rune(tmp.String()))
		}

		scheduleBytes = append(scheduleBytes[0:len(scheduleBytes)-Colums-1], []rune(ref)...)
	}

	return scheduleBytes, nil
}


package qldt

import "encoding/json"

type LichThi struct {
    SoThuTu             uint8 `json:"so_thu_tu"`
    KyThi               string `json:"ky_thi"`
    DotThi              string `json:"dot_thi"`
    MaMon               string `json:"ma_mon"`
    TenMon              string `json:"ten_mon"`
    MaPhong             string `json:"ma_phong"`
    MaCoSo              string `json:"ma_co_so"`
    NgayThi             string `json:"ngay_thi"`
    TietBatDau          string `json:"tiet_bat_dau"`
    SoTiet              string `json:"so_tiet"`
    GioBatDau           string `json:"gio_bat_dau"`
    SoPhut              string `json:"so_phut"`
    HinhThucThi         string `json:"hinh_thuc_thi"`
    GhiChuSV            string `json:"ghi_chu"`
    ToThi               string `json:"to_thi"`
    GhiChuHtt           string `json:"ghi_chu_htt"`
    GhiChuDuThi         string `json:"ghi_chu_du_thi"`
    SiSo                uint `json:"si_so"`
    NhomThi             string `json:"nhom_thi"`
}

type LichThiData struct {
    TotalItems          uint `json:"total_items"`
    TotalPages          uint `json:"total_pages"`
    DsLichThi           []LichThi `json:"ds_lich_thi"`
    ThongBaoNoHocPhi    string `json:"thong_bao_no_hoc_phi"`
}

type LichThiResponse struct {
    Data                LichThiData `json:"data"`
    ThongBaoGhiChu      string `json:"thong_bao_ghi_chu"`
    Response
}

type LichThiRequestBody struct {
    Filter struct {
        HocKy           string `json:"hoc_ky"`
        IsGiuaHocKy     bool `json:"ten_hoc_ky"`
    } `json:"filter"`
    Additional struct {
        Paging struct {
            Page        uint `json:"page"`
            Limit       uint `json:"limit"`
        } `json:"paging"`
        Ordering []Ordering `json:"ordering"`
    }
}

type Ordering struct {
    Name        string `json:"name"` 
    OrderType   string `json:"order"`
}

type TietTrongNgay struct {
    Tiet                uint8 `json:"tiet"`
    GioBatDau           string `json:"gio_bat_dau"`
    GioKetThuc          string `json:"gio_ket_thuc"`
    SoPhut              uint8 `json:"so_phut"`
}

type TuanTKB struct {
    TuanHocKy           uint `json:"tuan_hoc_ky"`
    TuanTuyetDoi        uint `json:"tuan_tuyet_doi"`
    ThongTinTuan        string `json:"thong_tin_tuan"`
    NgayBatDau          string `json:"ngay_bat_dau"`
    NgayKetThuc         string `json:"ngay_ket_thuc"`
    DSThoiKhoaBieu      []ThoiKhoaBieu `json:"ds_thoi_khoa_bieu"`
}

type ThoiKhoaBieu struct {
    IsHkLienTruoc       uint8 `json:"is_hk_lien_truoc"`
    ThuKieuSo           uint8 `json:"thu_kieu_so"`
    TietBatDau          uint8 `json:"tiet_bat_dau"`
    SoTiet              uint8 `json:"so_tiet"`
    MaMon               string `json:"ma_mon"`
    TenMon              string `json:"ten_mon"`
    SoTinChi            string `json:"so_tin_chi"`
    MaNhom              string `json:"ma_nhom"`
    MaToTh              string `json:"ma_to_th"`
    MaToHopPhan         string `json:"ma_to_hop"`
    MaGiangVien         string `json:"ma_giang_vien"`
    TenGiangVien        string `json:"ten_giang_vien"`
    MaLop               string `json:"ma_lop"`
    TenLop              string `json:"ten_lop"`
    MaPhong             string `json:"ma_phong"`
    MaCoSo              string `json:"ma_co_so"`
    IsDayBu             bool `json:"is_day_bu"`
    NgayHoc             string `json:"ngay_hoc"`
    IsNghiDay           bool `json:"is_nghi_day"`
}

type ScheduleData struct {
    TotalItems          uint `json:"total_items"`
    TotalPages          uint `json:"total_pages"`
    DSTietTrongNgay     []TietTrongNgay `json:"ds_tiet_trong_ngay"`
    DSTuanTKB           []TuanTKB `json:"ds_tuan_tkb"`
    ThongBao            string `json:"thong_bao"`
}

type ScheduleRequestBody struct {
    Filter struct {
        HocKy           string `json:"hoc_ky"`
        TenHocKy        string `json:"ten_hoc_ky"`
    } `json:"filter"`
}

type ScheduleResponse struct {
    Data                ScheduleData `json:"data"`
    Response
}

type TokenResponse struct {
    AccessToken         string `json:"access_token"`
    TokenType           string `json:"token_type"`
    RefreshToken        string `json:"refresh_token"`
    Username            string `json:"username"`
    Name                string `json:"name"`
    Principal           string `json:"principal"`
    Role                string `json:"roles"`
    ExpiresAt           string `json:".expires"`
    IssuedAt            string `json:".issued"`
    Response
}

type TokenErrorResponse struct {
    Message             string `json:"message"`
    ValidatedMessage    string `json:"validated_message"`
    Response
}

type Code string

type Response struct {
    Code                Code `json:"code"`
}

func (s *Code) UnmarshalJSON(bytes []byte) error {
    if bytes[0] == '"' {
        var str string
        if err := json.Unmarshal(bytes, &str); err != nil {
            return err
        }

        *s = Code(str)
    } else {
        var num json.Number
        if err := json.Unmarshal(bytes, &num); err != nil {
            return err
        }

        *s = Code(num.String())
    }   

    return nil
}

package model

type ZLMRecordMp4Data struct {
	APP       string `json:"app"`
	Stream    string `json:"stream"`
	FileName  string `json:"file_name"`
	FilePath  string `json:"file_path"`
	FileSize  int    `json:"file_size"`
	Folder    string `json:"folder"`
	StartTime int64  `json:"start_time"`
	TimeLen   int    `json:"time_len"`
	URL       string `json:"url"`
}

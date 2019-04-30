package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// globalSupportedTimeZones represents the instance that will be initialized once on runtime
// and load all TimeZones form Microsoft, correlate them to IANA and set proper time.Location
var globalSupportedTimeZones supportedTimeZones

// supportedTimeZones represents multiple instances grabbed by https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/outlookuser_supportedtimezones
type supportedTimeZones struct {
	Value []supportedTimeZone
}

// GetTimeZoneByAlias searches in the given set of supportedTimeZones for the TimeZone with the given alias. Returns
// either the time.Location or an error if it can not be found.
func (s supportedTimeZones) GetTimeZoneByAlias(alias string) (*time.Location, error) {
	for _, searchItem := range s.Value {
		if searchItem.Alias == alias {
			return searchItem.TimeLoc, nil
		}
	}
	return nil, fmt.Errorf("Could not find given time.Location for Alias %v", alias)
}

// GetTimeZoneByAlias searches in the given set of supportedTimeZones for the TimeZone with the given alias. Returns
// either the time.Location or an error if it can not be found.
func (s supportedTimeZones) GetTimeZoneByDisplayName(displayName string) (*time.Location, error) {
	for _, searchItem := range s.Value {
		if searchItem.DisplayName == displayName {
			return searchItem.TimeLoc, nil
		}
	}
	return nil, fmt.Errorf("Could not find given time.Location for DisplayName %v", displayName)
}

// supportedTimeZone represents one instance grabbed by https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/outlookuser_supportedtimezones
type supportedTimeZone struct {
	Alias       string
	DisplayName string
	TimeLoc     *time.Location
}

func (s *supportedTimeZone) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Alias       string `json:"alias"`
		DisplayName string `json:"displayName"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.Alias = tmp.Alias
	s.DisplayName = tmp.DisplayName

	ianaName, ok := WinIANA[s.DisplayName]
	if !ok {
		return fmt.Errorf("Can not map %v to IANA", s.DisplayName)
	}

	loc, err := time.LoadLocation(ianaName)
	if err != nil {
		return fmt.Errorf("Can not time.LoadLocation for original \"%v\" mapped to IANA \"%v\"", s.DisplayName, ianaName)
	}
	s.TimeLoc = loc

	return nil
}

// WinIANA contains a mapping for all Windows Time Zones to IANA time zones usable for time.LoadLocation.
// This list was initially copied from https://github.com/thinkovation/windowsiana/blob/master/windowsiana.go
// on 30th of August 2018, 14:00 and then extended on the same day.
//
// The full list of time zones that have been added and are now supported come from an an API-Call described here:
// https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/outlookuser_supportedtimezones
var WinIANA = map[string]string{
	"(UTC-12:00) International Date Line West":                      "Etc/GMT+12",
	"(UTC-11:00) Co-ordinated Universal Time-11":                    "Etc/GMT+11",
	"(UTC-11:00) Coordinated Universal Time-11":                     "Etc/GMT+11",
	"(UTC-10:00) Aleutian Islands":                                  "US/Aleutian",
	"(UTC-10:00) Hawaii":                                            "Pacific/Honolulu",
	"(UTC-09:30) Marquesas Islands":                                 "Pacific/Marquesas",
	"(UTC-09:00) Alaska":                                            "America/Anchorage",
	"(UTC-09:00) Co-ordinated Universal Time-09":                    "Etc/GMT+9",
	"(UTC-09:00) Coordinated Universal Time-09":                     "Etc/GMT+9",
	"(UTC-08:00) Baja California":                                   "America/Tijuana",
	"(UTC-08:00) Co-ordinated Universal Time-08":                    "Etc/GMT+8",
	"(UTC-08:00) Coordinated Universal Time-08":                     "Etc/GMT+8",
	"(UTC-08:00) Pacific Time (US & Canada)":                        "America/Los_Angeles",
	"(UTC-07:00) Arizona":                                           "America/Phoenix",
	"(UTC-07:00) Chihuahua, La Paz, Mazatlan":                       "America/Chihuahua",
	"(UTC-07:00) Mountain Time (US & Canada)":                       "America/Denver",
	"(UTC-06:00) Central America":                                   "America/Guatemala",
	"(UTC-06:00) Central Time (US & Canada)":                        "America/Chicago",
	"(UTC-06:00) Easter Island":                                     "Pacific/Easter",
	"(UTC-06:00) Guadalajara, Mexico City, Monterrey":               "America/Mexico_City",
	"(UTC-06:00) Saskatchewan":                                      "America/Regina",
	"(UTC-05:00) Bogota, Lima, Quito, Rio Branco":                   "America/Bogota",
	"(UTC-05:00) Chetumal":                                          "America/Cancun",
	"(UTC-05:00) Eastern Time (US & Canada)":                        "America/New_York",
	"(UTC-05:00) Haiti":                                             "America/Port-au-Prince",
	"(UTC-05:00) Havana":                                            "America/Havana",
	"(UTC-05:00) Indiana (East)":                                    "America/Indianapolis",
	"(UTC-05:00) Turks and Caicos":                                  "Etc/GMT+5",
	"(UTC-04:00) Asuncion":                                          "America/Asuncion",
	"(UTC-04:00) Atlantic Time (Canada)":                            "America/Halifax",
	"(UTC-04:00) Caracas":                                           "America/Caracas",
	"(UTC-04:00) Cuiaba":                                            "America/Cuiaba",
	"(UTC-04:00) Georgetown, La Paz, Manaus, San Juan":              "America/La_Paz",
	"(UTC-04:00) Santiago":                                          "America/Santiago",
	"(UTC-04:00) Turks and Caicos":                                  "America/Grand_Turk",
	"(UTC-03:30) Newfoundland":                                      "America/St_Johns",
	"(UTC-03:00) Araguaina":                                         "America/Araguaina",
	"(UTC-03:00) Brasilia":                                          "America/Sao_Paulo",
	"(UTC-03:00) Cayenne, Fortaleza":                                "America/Cayenne",
	"(UTC-03:00) City of Buenos Aires":                              "America/Buenos_Aires",
	"(UTC-03:00) Greenland":                                         "America/Godthab",
	"(UTC-03:00) Montevideo":                                        "America/Montevideo",
	"(UTC-03:00) Punta Arenas":                                      "America/Punta_Arenas",
	"(UTC-03:00) Saint Pierre and Miquelon":                         "America/Miquelon",
	"(UTC-03:00) Salvador":                                          "America/Bahia",
	"(UTC-02:00) Co-ordinated Universal Time-02":                    "Etc/GMT+2",
	"(UTC-02:00) Coordinated Universal Time-02":                     "Etc/GMT+2",
	"(UTC-02:00) Mid-Atlantic - Old":                                "Etc/GMT+2",
	"(UTC-01:00) Azores":                                            "Atlantic/Azores",
	"(UTC-01:00) Cabo Verde Is.":                                    "Atlantic/Cape_Verde",
	"(UTC) Co-ordinated Universal Time":                             "Etc/GMT",
	"(UTC) Coordinated Universal Time":                              "Etc/GMT",
	"(UTC+00:00) Casablanca":                                        "Africa/Casablanca",
	"(UTC+00:00) Dublin, Edinburgh, Lisbon, London":                 "Europe/London",
	"(UTC+00:00) Monrovia, Reykjavik":                               "Atlantic/Reykjavik",
	"(UTC+01:00) Amsterdam, Berlin, Bern, Rome, Stockholm, Vienna":  "Europe/Berlin",
	"(UTC+01:00) Belgrade, Bratislava, Budapest, Ljubljana, Prague": "Europe/Budapest",
	"(UTC+01:00) Brussels, Copenhagen, Madrid, Paris":               "Europe/Paris",
	"(UTC+01:00) Casablanca":                                        "Africa/Casablanca",
	"(UTC+01:00) Sarajevo, Skopje, Warsaw, Zagreb":                  "Europe/Warsaw",
	"(UTC+01:00) West Central Africa":                               "Africa/Lagos",
	"(UTC+01:00) Windhoek":                                          "Africa/Windhoek",
	"(UTC+02:00) Amman":                                             "Asia/Amman",
	"(UTC+02:00) Athens, Bucharest":                                 "Europe/Bucharest",
	"(UTC+02:00) Beirut":                                            "Asia/Beirut",
	"(UTC+02:00) Cairo":                                             "Africa/Cairo",
	"(UTC+02:00) Chisinau":                                          "Europe/Chisinau",
	"(UTC+02:00) Damascus":                                          "Asia/Damascus",
	"(UTC+02:00) Gaza, Hebron":                                      "Asia/Gaza",
	"(UTC+02:00) Harare, Pretoria":                                  "Africa/Johannesburg",
	"(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius":     "Europe/Kiev",
	"(UTC+02:00) Istanbul":                                          "Europe/Istanbul",
	"(UTC+03:00) Istanbul":                                          "Europe/Istanbul",
	"(UTC+02:00) Jerusalem":                                         "Asia/Jerusalem",
	"(UTC+02:00) Kaliningrad":                                       "Europe/Kaliningrad",
	"(UTC+02:00) Windhoek":                                          "Africa/Windhoek",
	"(UTC+02:00) Khartoum":                                          "Africa/Khartoum",
	"(UTC+02:00) Tripoli":                                           "Africa/Tripoli",
	"(UTC+03:00) Baghdad":                                           "Asia/Baghdad",
	"(UTC+03:00) Kuwait, Riyadh":                                    "Asia/Riyadh",
	"(UTC+03:00) Minsk":                                             "Europe/Minsk",
	"(UTC+03:00) Moscow, St. Petersburg":                            "Europe/Moscow",
	"(UTC+03:00) Moscow, St. Petersburg, Volgograd":                 "Europe/Moscow",
	"(UTC+03:00) Nairobi":                                           "Africa/Nairobi",
	"(UTC+03:30) Tehran":                                            "Asia/Tehran",
	"(UTC+04:00) Abu Dhabi, Muscat":                                 "Asia/Dubai",
	"(UTC+04:00) Astrakhan, Ulyanovsk":                              "Europe/Samara",
	"(UTC+04:00) Baku":                                              "Asia/Baku",
	"(UTC+04:00) Izhevsk, Samara":                                   "Europe/Samara",
	"(UTC+04:00) Port Louis":                                        "Indian/Mauritius",
	"(UTC+04:00) Saratov":                                           "Europe/Saratov",
	"(UTC+04:00) Tbilisi":                                           "Asia/Tbilisi",
	"(UTC+04:00) Volgograd":                                         "Europe/Volgograd",
	"(UTC+04:00) Yerevan":                                           "Asia/Yerevan",
	"(UTC+04:30) Kabul":                                             "Asia/Kabul",
	"(UTC+05:00) Ashgabat, Tashkent":                                "Asia/Tashkent",
	"(UTC+05:00) Ekaterinburg":                                      "Asia/Yekaterinburg",
	"(UTC+05:00) Islamabad, Karachi":                                "Asia/Karachi",
	"(UTC+05:00) Qyzylorda":                                         "Asia/Qyzylorda",
	"(UTC+05:30) Chennai, Kolkata, Mumbai, New Delhi":               "Asia/Calcutta",
	"(UTC+05:30) Sri Jayawardenepura":                               "Asia/Colombo",
	"(UTC+05:45) Kathmandu":                                         "Asia/Kathmandu",
	"(UTC+06:00) Astana":                                            "Asia/Almaty",
	"(UTC+06:00) Dhaka":                                             "Asia/Dhaka",
	"(UTC+06:00) Omsk":                                              "Asia/Omsk",
	"(UTC+06:00) Novosibirsk":                                       "Asia/Novosibirsk",
	"(UTC+06:30) Yangon (Rangoon)":                                  "Asia/Rangoon",
	"(UTC+07:00) Bangkok, Hanoi, Jakarta":                           "Asia/Bangkok",
	"(UTC+07:00) Barnaul, Gorno-Altaysk":                            "Asia/Krasnoyarsk",
	"(UTC+07:00) Hovd":                                              "Asia/Hovd",
	"(UTC+07:00) Krasnoyarsk":                                       "Asia/Krasnoyarsk",
	"(UTC+07:00) Novosibirsk":                                       "Asia/Novosibirsk",
	"(UTC+07:00) Tomsk":                                             "Asia/Tomsk",
	"(UTC+08:00) Beijing, Chongqing, Hong Kong SAR, Urumqi":         "Asia/Shanghai",
	"(UTC+08:00) Beijing, Chongqing, Hong Kong, Urumqi":             "Asia/Shanghai",
	"(UTC+08:00) Irkutsk":                                           "Asia/Irkutsk",
	"(UTC+08:00) Kuala Lumpur, Singapore":                           "Asia/Singapore",
	"(UTC+08:00) Perth":                                             "Australia/Perth",
	"(UTC+08:00) Taipei":                                            "Asia/Taipei",
	"(UTC+08:00) Ulaanbaatar":                                       "Asia/Ulaanbaatar",
	"(UTC+08:30) Pyongyang":                                         "Asia/Pyongyang",
	"(UTC+09:00) Pyongyang":                                         "Asia/Pyongyang",
	"(UTC+08:45) Eucla":                                             "Australia/Eucla",
	"(UTC+09:00) Chita":                                             "Asia/Chita",
	"(UTC+09:00) Osaka, Sapporo, Tokyo":                             "Asia/Tokyo",
	"(UTC+09:00) Seoul":                                             "Asia/Seoul",
	"(UTC+09:00) Yakutsk":                                           "Asia/Yakutsk",
	"(UTC+09:30) Adelaide":                                          "Australia/Adelaide",
	"(UTC+09:30) Darwin":                                            "Australia/Darwin",
	"(UTC+10:00) Brisbane":                                          "Australia/Brisbane",
	"(UTC+10:00) Canberra, Melbourne, Sydney":                       "Australia/Sydney",
	"(UTC+10:00) Guam, Port Moresby":                                "Pacific/Port_Moresby",
	"(UTC+10:00) Hobart":                                            "Australia/Hobart",
	"(UTC+10:00) Vladivostok":                                       "Asia/Vladivostok",
	"(UTC+10:30) Lord Howe Island":                                  "Australia/Lord_Howe",
	"(UTC+11:00) Bougainville Island":                               "Pacific/Bougainville",
	"(UTC+11:00) Chokurdakh":                                        "Asia/Srednekolymsk",
	"(UTC+11:00) Magadan":                                           "Asia/Magadan",
	"(UTC+11:00) Norfolk Island":                                    "Pacific/Norfolk",
	"(UTC+11:00) Sakhalin":                                          "Asia/Sakhalin",
	"(UTC+11:00) Solomon Is., New Caledonia":                        "Pacific/Guadalcanal",
	"(UTC+12:00) Anadyr, Petropavlovsk-Kamchatsky":                  "Asia/Kamchatka",
	"(UTC+12:00) Auckland, Wellington":                              "Pacific/Auckland",
	"(UTC+12:00) Co-ordinated Universal Time+12":                    "Etc/GMT-12",
	"(UTC+12:00) Coordinated Universal Time+12":                     "Etc/GMT-12",
	"(UTC+12:00) Petropavlovsk-Kamchatsky - Old":                    "Etc/GMT-12",
	"(UTC+12:00) Fiji":                                              "Pacific/Fiji",
	"(UTC+12:45) Chatham Islands":                                   "Pacific/Chatham",
	"(UTC+13:00) Nuku'alofa":                                        "Pacific/Tongatapu",
	"(UTC+13:00) Co-ordinated Universal Time+13":                    "Etc/GMT-13",
	"(UTC+13:00) Coordinated Universal Time+13":                     "Etc/GMT-13",
	"(UTC+13:00) Samoa":                                             "Pacific/Apia",
	"(UTC+14:00) Kiritimati Island":                                 "Pacific/Kiritimati"}

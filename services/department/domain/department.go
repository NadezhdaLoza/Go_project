package domain

type Department struct {
    ID                  int    `json:"id"`
    Name                string `json:"name"`
    WorkConditionsClass string `json:"work_conditions_class,omitempty"`
    WorkDayStart        string `json:"work_day_start,omitempty"`
    WorkDayEnd          string `json:"work_day_end,omitempty"`
}
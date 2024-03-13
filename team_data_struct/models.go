package team_data_struct

// SensorData structure model for storing sensor data

type SensorDataInterface interface {
    SetPower(power float32)
    GetVoltage() *float32
    GetCurrent() *float32
}

type SensorDataInterfaceDualPower interface {
    SetPower(power float32)
    GetVoltage1() *float32
    GetCurrent1() *float32
    GetVoltage2() *float32
    GetCurrent2() *float32
}

type SensorDataInterfaceQuadPower interface {
    SetPower(power float32)
    GetVoltage1() *float32
    GetCurrent1() *float32
    GetVoltage2() *float32
    GetCurrent2() *float32
    GetVoltage3() *float32
    GetCurrent3() *float32
    GetVoltage4() *float32
    GetCurrent4() *float32
}

type GenericSensorData struct {
    Temp1   *float32 `json:"generic1_temp1"`
    Temp2   *float32 `json:"generic1_temp2"`
    Temp3   *float32 `json:"generic1_temp3"`
    Temp4   *float32 `json:"generic1_temp4"`
    Temp5   *float32 `json:"generic1_temp5"`
    Voltage *float32 `json:"generic1_voltage"`
    Current *float32 `json:"generic1_current"`
    Power   *float32 `json:"generic1_power"`
    Lat     *float64 `json:"generic1_lat"`
    Lon     *float64 `json:"generic1_lon"`
}
// Assurez-vous que toutes vos structures de données implémentent cette interface.
func (g *GenericSensorData) SetPower(power float32)     { g.Power = &power }
func (g *GenericSensorData) GetVoltage() *float32        { return g.Voltage }
func (g *GenericSensorData) GetCurrent() *float32        { return g.Current }

type GenericSensorDataDualPower struct {
  Temp1   		*float32 `json:"generic2_temp1"`
  Temp2   		*float32 `json:"generic2_temp2"`
  Temp3           *float32 `json:"generic2_temp3"`
  Voltage1        *float32 `json:"generic2_voltage1"`
  Voltage2        *float32 `json:"generic2_voltage2"`
  Voltage3        *float32 `json:"generic2_voltage3"`
  Current1        *float32 `json:"generic2_current1"`
  Current2        *float32 `json:"generic2_current2"`
  Current3        *float32 `json:"generic2_current3"`
  Power           *float32 `json:"generic2_power"`			// compute automaticaly
}
func (g *GenericSensorDataDualPower) SetPower(power float32)     { g.Power = &power }
func (g *GenericSensorDataDualPower) GetVoltage1() *float32        { return g.Voltage1 }
func (g *GenericSensorDataDualPower) GetCurrent1() *float32        { return g.Current1 }
func (g *GenericSensorDataDualPower) GetVoltage2() *float32        { return g.Voltage2 }
func (g *GenericSensorDataDualPower) GetCurrent2() *float32        { return g.Current2 }
func (g *GenericSensorDataDualPower) GetVoltage() *float32        { return g.Voltage3 }
func (g *GenericSensorDataDualPower) GetCurrent() *float32        { return g.Current3 }

type GenericSensorDataQuadPower struct {
  Temp1           *float32 `json:"generic3_temp1"`
  Temp2           *float32 `json:"generic3_temp2"`
  Temp3           *float32 `json:"generic3_temp3"`
  Temp4           *float32 `json:"generic3_temp4"`
  Voltage1        *float32 `json:"generic3_voltage1"`
  Voltage2        *float32 `json:"generic3_voltage2"`
  Voltage3        *float32 `json:"generic3_voltage3"`
  Voltage4        *float32 `json:"generic3_voltage4"`
  Current1        *float32 `json:"generic3_current1"`
  Current2        *float32 `json:"generic3_current2"`
  Current3        *float32 `json:"generic3_current3"`
  Current4        *float32 `json:"generic3_current3"`
  Power           *float32 `json:"generic3_power"`			// compute automaticaly
}
func (g *GenericSensorDataQuadPower) SetPower(power float32)     { g.Power = &power }
func (g *GenericSensorDataQuadPower) GetVoltage1() *float32        { return g.Voltage1 }
func (g *GenericSensorDataQuadPower) GetCurrent1() *float32        { return g.Current1 }
func (g *GenericSensorDataQuadPower) GetVoltage2() *float32        { return g.Voltage2 }
func (g *GenericSensorDataQuadPower) GetCurrent2() *float32        { return g.Current2 }
func (g *GenericSensorDataQuadPower) GetVoltage3() *float32        { return g.Voltage3 }
func (g *GenericSensorDataQuadPower) GetCurrent3() *float32        { return g.Current3 }
func (g *GenericSensorDataQuadPower) GetVoltage4() *float32        { return g.Voltage4 }
func (g *GenericSensorDataQuadPower) GetCurrent4() *float32        { return g.Current4 }
func (g *GenericSensorDataQuadPower) GetVoltage()  *float32        { return g.Voltage1 }
func (g *GenericSensorDataQuadPower) GetCurrent()  *float32        { return g.Current1 }

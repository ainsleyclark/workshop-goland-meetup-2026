package openmeteo

// WeatherCode represents the present weather conditions as reported from
// a manned weather station. These codes correspond to different weather
// conditions (e.g., clear sky, partly cloudy, rain).
//
// See WMO weather codes:
// https://www.nodc.noaa.gov/archive/arc0021/0002199/1.1/data/0-data/HTML/WMO-CODE/WMO4677.HTM
type WeatherCode int

const (
	ClearSky WeatherCode = iota
	CloudsDissolving
	SkyUnchanged
	CloudsForming
	VisibilityReducedBySmoke
	Haze
	WidespreadDust
	DustOrSandRaisedByWind
	DustOrSandWhirls
	DustStormOrSandstorm
	Mist
	ShallowFogPatches
	ShallowFogContinuous
	LightningNoThunder
	PrecipitationNotReachingGround
	DistantPrecipitation
	NearbyPrecipitation
	ThunderstormNoPrecipitation
	Squalls
	FunnelClouds
	RecentDrizzleOrSnowGrains
	RecentRain
	RecentSnow
	RecentMixedPrecipitation
	RecentFreezingPrecipitation
	RecentRainShower
	RecentSnowOrMixedShower
	RecentHailShower
	RecentFog
	RecentThunderstorm
	DuststormDecreasedLastHour
	DuststormNoChangeLastHour
	DuststormIncreasedLastHour
	SevereDuststormDecreasedLastHour
	SevereDuststormNoChangeLastHour
	SevereDuststormIncreasedLastHour
	SlightBlowingSnowLow
	HeavyDriftingSnowLow
	SlightBlowingSnowHigh
	HeavyDriftingSnowHigh
	FogAtDistance
	FogInPatches
	FogSkyVisibleThinning
	FogSkyInvisibleThinning
	FogSkyVisibleNoChange
	FogSkyInvisibleNoChange
	FogSkyVisibleThickening
	FogSkyInvisibleThickening
	FogDepositingRimeSkyVisible
	FogDepositingRimeSkyInvisible
	SlightIntermittentDrizzle
	SlightContinuousDrizzle
	ModerateIntermittentDrizzle
	ModerateContinuousDrizzle
	HeavyIntermittentDrizzle
	HeavyContinuousDrizzle
	SlightFreezingDrizzle
	ModerateFreezingDrizzle
	SlightDrizzleAndRain
	ModerateDrizzleAndRain
	SlightIntermittentRain
	SlightContinuousRain
	ModerateIntermittentRain
	ModerateContinuousRain
	HeavyIntermittentRain
	HeavyContinuousRain
	SlightFreezingRain
	ModerateFreezingRain
	SlightRainAndSnow
	ModerateRainAndSnow
	SlightIntermittentSnow
	SlightContinuousSnow
	ModerateIntermittentSnow
	ModerateContinuousSnow
	HeavyIntermittentSnow
	HeavyContinuousSnow
	DiamondDust
	SnowGrains
	IsolatedSnowCrystals
	IcePellets
	SlightRainShowers
	ModerateRainShowers
	ViolentRainShowers
	SlightRainAndSnowShowers
	ModerateRainAndSnowShowers
	SlightSnowShowers
	ModerateSnowShowers
	SlightSnowPelletShowers
	ModerateSnowPelletShowers
	SlightHailShowers
	ModerateHailShowers
	SlightRainThunderstormLastHour
	ModerateRainThunderstormLastHour
	SlightSnowThunderstormLastHour
	ModerateSnowThunderstormLastHour
	ThunderstormSlightRain
	ThunderstormModerateRain
	ThunderstormSlightHail
	ThunderstormHeavyHail
	ThunderstormWithDuststorm
	ThunderstormWithHeavyHail
)

// WeatherCodeDescriptions maps WeatherCode to its string description
var WeatherCodeDescriptions = map[WeatherCode]string{
	ClearSky:                         "Clear sky",
	CloudsDissolving:                 "Clouds dissolving or becoming less developed",
	SkyUnchanged:                     "State of sky on the whole unchanged",
	CloudsForming:                    "Clouds generally forming or developing",
	VisibilityReducedBySmoke:         "Visibility reduced by smoke",
	Haze:                             "Haze",
	WidespreadDust:                   "Widespread dust in suspension in the air",
	DustOrSandRaisedByWind:           "Dust or sand raised by wind",
	DustOrSandWhirls:                 "Well developed dust or sand whirls",
	DustStormOrSandstorm:             "Duststorm or sandstorm within sight",
	Mist:                             "Mist",
	ShallowFogPatches:                "Shallow fog patches",
	ShallowFogContinuous:             "Shallow fog continuous",
	LightningNoThunder:               "Lightning visible, no thunder heard",
	PrecipitationNotReachingGround:   "Precipitation within sight, not reaching ground",
	DistantPrecipitation:             "Distant precipitation reaching ground",
	NearbyPrecipitation:              "Nearby precipitation reaching ground",
	ThunderstormNoPrecipitation:      "Thunderstorm without precipitation",
	Squalls:                          "Squalls",
	FunnelClouds:                     "Funnel cloud(s)",
	RecentDrizzleOrSnowGrains:        "Recent drizzle or snow grains",
	RecentRain:                       "Recent rain",
	RecentSnow:                       "Recent snow",
	RecentMixedPrecipitation:         "Recent rain and snow or ice pellets",
	RecentFreezingPrecipitation:      "Recent freezing drizzle or freezing rain",
	RecentRainShower:                 "Recent rain shower",
	RecentSnowOrMixedShower:          "Recent snow shower or mixed rain and snow shower",
	RecentHailShower:                 "Recent shower of hail or rain and hail",
	RecentFog:                        "Recent fog",
	RecentThunderstorm:               "Recent thunderstorm",
	DuststormDecreasedLastHour:       "Slight/moderate duststorm or sandstorm, decreased in last hour",
	DuststormNoChangeLastHour:        "Slight/moderate duststorm or sandstorm, no change in last hour",
	DuststormIncreasedLastHour:       "Slight/moderate duststorm or sandstorm, increased in last hour",
	SevereDuststormDecreasedLastHour: "Severe duststorm or sandstorm, decreased in last hour",
	SevereDuststormNoChangeLastHour:  "Severe duststorm or sandstorm, no change in last hour",
	SevereDuststormIncreasedLastHour: "Severe duststorm or sandstorm, increased in last hour",
	SlightBlowingSnowLow:             "Slight/moderate blowing snow, generally low",
	HeavyDriftingSnowLow:             "Heavy drifting snow, generally low",
	SlightBlowingSnowHigh:            "Slight/moderate blowing snow, generally high",
	HeavyDriftingSnowHigh:            "Heavy drifting snow, generally high",
	FogAtDistance:                    "Fog at distance",
	FogInPatches:                     "Fog in patches",
	FogSkyVisibleThinning:            "Fog sky visible, thinning",
	FogSkyInvisibleThinning:          "Fog sky invisible, thinning",
	FogSkyVisibleNoChange:            "Fog sky visible, no change",
	FogSkyInvisibleNoChange:          "Fog sky invisible, no change",
	FogSkyVisibleThickening:          "Fog sky visible, becoming thicker",
	FogSkyInvisibleThickening:        "Fog sky invisible, becoming thicker",
	FogDepositingRimeSkyVisible:      "Fog depositing rime, sky visible",
	FogDepositingRimeSkyInvisible:    "Fog depositing rime, sky invisible",
	SlightIntermittentDrizzle:        "Slight intermittent drizzle",
	SlightContinuousDrizzle:          "Slight continuous drizzle",
	ModerateIntermittentDrizzle:      "Moderate intermittent drizzle",
	ModerateContinuousDrizzle:        "Moderate continuous drizzle",
	HeavyIntermittentDrizzle:         "Heavy intermittent drizzle",
	HeavyContinuousDrizzle:           "Heavy continuous drizzle",
	SlightFreezingDrizzle:            "Slight freezing drizzle",
	ModerateFreezingDrizzle:          "Moderate/heavy freezing drizzle",
	SlightDrizzleAndRain:             "Slight drizzle and rain",
	ModerateDrizzleAndRain:           "Moderate/heavy drizzle and rain",
	SlightIntermittentRain:           "Slight intermittent rain",
	SlightContinuousRain:             "Slight continuous rain",
	ModerateIntermittentRain:         "Moderate intermittent rain",
	ModerateContinuousRain:           "Moderate continuous rain",
	HeavyIntermittentRain:            "Heavy intermittent rain",
	HeavyContinuousRain:              "Heavy continuous rain",
	SlightFreezingRain:               "Slight freezing rain",
	ModerateFreezingRain:             "Moderate/heavy freezing rain",
	SlightRainAndSnow:                "Slight rain and snow",
	ModerateRainAndSnow:              "Moderate/heavy rain and snow",
	SlightIntermittentSnow:           "Slight intermittent snow",
	SlightContinuousSnow:             "Slight continuous snow",
	ModerateIntermittentSnow:         "Moderate intermittent snow",
	ModerateContinuousSnow:           "Moderate continuous snow",
	HeavyIntermittentSnow:            "Heavy intermittent snow",
	HeavyContinuousSnow:              "Heavy continuous snow",
	DiamondDust:                      "Diamond dust",
	SnowGrains:                       "Snow grains",
	IsolatedSnowCrystals:             "Isolated star-like snow crystals",
	IcePellets:                       "Ice pellets",
	SlightRainShowers:                "Slight rain showers",
	ModerateRainShowers:              "Moderate/heavy rain showers",
	ViolentRainShowers:               "Violent rain showers",
	SlightRainAndSnowShowers:         "Slight rain and snow showers",
	ModerateRainAndSnowShowers:       "Moderate/heavy rain and snow showers",
	SlightSnowShowers:                "Slight snow showers",
	ModerateSnowShowers:              "Moderate/heavy snow showers",
	SlightSnowPelletShowers:          "Slight showers of snow pellets/small hail",
	ModerateSnowPelletShowers:        "Moderate/heavy showers of snow pellets/small hail",
	SlightHailShowers:                "Slight hail showers",
	ModerateHailShowers:              "Moderate/heavy hail showers",
	SlightRainThunderstormLastHour:   "Slight rain, thunderstorm in last hour",
	ModerateRainThunderstormLastHour: "Moderate/heavy rain, thunderstorm in last hour",
	SlightSnowThunderstormLastHour:   "Slight snow/mixed precipitation, thunderstorm in last hour",
	ModerateSnowThunderstormLastHour: "Moderate/heavy snow/mixed precipitation, thunderstorm in last hour",
	ThunderstormSlightRain:           "Thunderstorm with slight rain",
	ThunderstormModerateRain:         "Thunderstorm with moderate/heavy rain",
	ThunderstormSlightHail:           "Thunderstorm with slight hail",
	ThunderstormHeavyHail:            "Thunderstorm with heavy hail",
	ThunderstormWithDuststorm:        "Thunderstorm with duststorm",
	ThunderstormWithHeavyHail:        "Thunderstorm with heavy hail",
}

// String returns the string description of the WeatherCode
func (w WeatherCode) String() string {
	if desc, ok := WeatherCodeDescriptions[w]; ok {
		return desc
	}
	return "Unknown weather code"
}

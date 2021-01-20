package main

import "testing"

func TestPostText(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedOutput string
	}{
		"multiple texts": {
			input:          "\"authorEndpoint\":{\"clickTrackingParams\":\"CIIBEPS8AiITCMfK99Xnqu4CFTNATAgdVoYJvA==\",\"commandMetadata\":{\"webCommandMetadata\":{\"url\":\"/channel/UCLF4C8ikMjOElQrNhTrO8yg\",\"webPageType\":\"WEB_PAGE_TYPE_CHANNEL\",\"rootVe\":3611,\"apiUrl\":\"/youtubei/v1/browse\"}},\"browseEndpoint\":{\"browseId\":\"UCLF4C8ikMjOElQrNhTrO8yg\"}},\"contentText\":{\"runs\":[{\"text\":\"99c members: At \"},{\"text\":\"10:30\"},{\"text\":\" ET I will be releasing a video on AIRG \u0026 PCTI. I like these stocks \u0026 I am invested in both. These are 5G but do not compete with CRNT (the vid I did yesterday), they're in different areas of the 5G Space. I think this is a good sector to focus on because it's not overexposed like a lot of other hot sectors. I found this very interesting: 4G wavelengths have a range of about 10 miles. 5G wavelengths have a range of about 1000 feet, not even 2% of 4g's range. There needs to be reliable signal everywhere (lampposts, traffic lights, etc because even trees can block the signal) and I think this is why these 2 companies are both actively buying back stock. They provide the antennas needed to make 5G possible. The reason I am talking about both is because they both are similar but a little different. AIRG might be the leader in this space but PCTI has a dividend and is profitable, where as AIRG is close to being profitable but have a partnership with AT\u0026T and also have news announcements that might send the stock higher in \"early 2021\".\"}]},\"expandButton\":{\"buttonRenderer\":{\"style\":\"STYLE_TEXT\",\"size\":\"SIZE_DEFAULT\",\"text\":{\"accessibility\":{\"accessibilityData\":{\"label\":\"Read more\"}},\"simpleText\":\"Read more\"},\"accessibility\":{\"label\":\"Read more\"},\"trackingParams\":\"CIgBEK_YAiITCMfK99Xnqu4CFTNATAgdVoYJvA==\"}},\"collapseButton\":{\"buttonRenderer\":{\"style\":\"STYLE_TEXT\",\"size\":\"SIZE_DEFAULT\",\"text\":{\"accessibility\":{\"accessibilityData\":{\"label\":\"Show less\"}},\"simpleText\":\"Show less\"},\"accessibility\":{\"label\":\"Show less\"},\"trackingParams\":\"CIcBELDYAiITCMfK99Xnqu4CFTNATAgdVoYJvA==\"}},\"publishedTimeText\":{\"runs\":[{\"text\":\"18 minutes ago\",\"navigationEndpoint\":{\"clickTrackingParams\":\"CIIBEPS8AiITCMfK99Xnqu4CFTNATAgdVoYJvA==\",\"commandMetadata\":{\"webCommandMetadata\":{\"url\":\"/post/UgzoegY3IawhR6QmYNN4AaABCQ\",\"webPageType\":\"WEB_PAGE_TYPE_CHANNEL\",\"rootVe\":3611,\"apiUrl\":\"/youtubei/v1/browse\"}},\"browseEndpoint\":{\"browseId\":\"UCLF4C8ikMjOElQrNhTrO8yg\",\"params\":\"Egljb21tdW5pdHnKAR2yARpVZ3pvZWdZM0lhd2hSNlFtWU5ONEFhQUJDUeoCBBABGAE%3D\",\"canonicalBaseUrl\":\"/post/UgzoegY3IawhR6QmYNN4AaABCQ\"}}}]}",
			expectedOutput: "99c members: At 10:30 ET I will be releasing a video on AIRG & PCTI.  I like these stocks & I am invested in both.  These are 5G but do not compete with CRNT (the vid I did yesterday), they're in different areas of the 5G Space. I think this is a good sector to focus on because it's not overexposed like a lot of other hot sectors.  I found this very interesting: 4G wavelengths have a range of about 10 miles.  5G wavelengths have a range of about 1000 feet, not even 2% of 4g's range.  There needs to be reliable signal everywhere (lampposts, traffic lights, etc because even trees can block the signal) and I think this is why these 2 companies are both actively buying back stock.  They provide the antennas needed to make 5G possible.  The reason I am talking about both is because they both are similar but a little different.  AIRG might be the leader in this space but PCTI has a dividend and is profitable, where as AIRG is close to being profitable but have a partnership with AT&T and also have news announcements that might send the stock higher in \"early 2021\".",
		},
	}
	for name, test := range tests {
		out, err := postText(test.input)
		if err != nil {
			t.Fatalf("test \"%s\" failed, expected no error, got error %v", name, err)
		} else if out != test.expectedOutput {
			t.Fatalf("test \"%s\" failed, expected output %s, got %s", name, test.expectedOutput, out)
		}
	}
}

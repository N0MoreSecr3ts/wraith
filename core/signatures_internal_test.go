package core

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetEntropyIntBasicBehaviour(t *testing.T) {
	Convey("Given strings with different character distributions", t, func() {
		low := "aaaaaaaaaa"
		high := "a1b2c3d4e5"

		Convey("When getEntropyInt is called", func() {
			lowEntropy := getEntropyInt(low)
			highEntropy := getEntropyInt(high)

			Convey("The mixed string should have higher entropy than the repeated one", func() {
				So(highEntropy, ShouldBeGreaterThan, lowEntropy)
			})
		})
	})
}

func TestIsSafeTextUsesSafeFunctionSignatures(t *testing.T) {
	Convey("Given a match string and SafeFunctionSignatures", t, func() {
		// reset any global state from other tests
		SafeFunctionSignatures = nil

		safePattern := regexp.MustCompile(`(?i)not_a_secret`)
		SafeFunctionSignatures = []SafeFunctionSignature{
			{
				match: safePattern,
			},
		}

		Convey("When the text matches a safe function signature", func() {
			text := "NOT_A_SECRET"
			result := IsSafeText(&text)

			Convey("The text should be considered safe", func() {
				So(result, ShouldBeTrue)
			})
		})

		Convey("When the text does not match a safe function signature", func() {
			text := "really-secret-value"
			result := IsSafeText(&text)

			Convey("The text should not be considered safe", func() {
				So(result, ShouldBeFalse)
			})
		})
	})
}

func TestConfirmEntropyRespectsThresholdAndSafeText(t *testing.T) {
	Convey("Given a match string and entropy threshold", t, func() {
		// ensure no safe signatures interfere
		SafeFunctionSignatures = nil

		match := "abc123XYZ"

		Convey("When the threshold is zero, entropy alone should allow the match", func() {
			ok := confirmEntropy(match, 0)
			So(ok, ShouldBeTrue)
		})

		Convey("When the threshold is higher than the string entropy, the match should be rejected", func() {
			entropy := getEntropyInt(match)
			ok := confirmEntropy(match, entropy+1.0)
			So(ok, ShouldBeFalse)
		})

		Convey("When the text is marked as safe, the match should be rejected even with low threshold", func() {
			safePattern := regexp.MustCompile(regexp.QuoteMeta(match))
			SafeFunctionSignatures = []SafeFunctionSignature{
				{
					match: safePattern,
				},
			}

			ok := confirmEntropy(match, 0)
			So(ok, ShouldBeFalse)
		})
	})
}

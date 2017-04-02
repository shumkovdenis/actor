package account

import (
	"encoding/json"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"go.uber.org/zap"
)

func authorize(username, password string) Message {
	conf := config.AccountAPI()

	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   conf.Type + "_CLIENT_AUTH",
			"auth_username": username,
			"auth_password": password,
		}).
		Post(conf.URL)
	if err != nil {
		logger.L().Error("authorization failed",
			zap.Error(err),
		)
		return &AuthorizationFailed{}
	}

	res := &struct {
		Result string `json:"result"`
		Code   int    `json:"code"`
		Groups []struct {
			Title string `json:"title"`
			Games []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			}
		} `json:"groups"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		logger.L().Error("authorization failed",
			zap.Error(err),
		)
		return &AuthorizationFailed{}
	}

	if res.Result == "Error" {
		logger.L().Error("authorization failed",
			zap.String("result", res.Result),
			zap.Int("code", res.Code),
		)
		return &AuthorizationFailed{}
	}

	categories := make([]Category, len(res.Groups))
	for i, group := range res.Groups {
		games := make([]Game, len(group.Games))
		for j, game := range group.Games {
			thumb := ""
			if _, ok := thumbs[game.ID]; ok {
				thumb = "./games/Assets/thumbs/180_130/" + thumbs[game.ID]
			}

			games[j] = Game{
				ID:    game.ID,
				Title: game.Title,
				Thumb: thumb,
			}
		}
		categories[i] = Category{
			Title: group.Title,
			Games: games,
		}
	}

	return &Authorized{categories}
}

var (
	thumbs = map[string]string{
		"59":  "book_of_ra_deluxe.png",
		"60":  "dolphins_pearl_deluxe.png",
		"61":  "lucky_ladys_deluxe.png",
		"62":  "columbus_deluxe.png",
		"63":  "xtra_hot.png",
		"64":  "power_stars.png",
		"65":  "just_jewels_deluxe.png",
		"66":  "beetle_mania_deluxe.png",
		"67":  "lord_of_ocean.png",
		"68":  "classic/columbus.png",
		"69":  "book_of_ra.png",
		"70":  "katana.png",
		"72":  "plenty_on_twenty.png",
		"73":  "fruits_and_royals.png",
		"74":  "hollywood_star.png",
		"75":  "golden_ark.png",
		"76":  "dazzling_diamonds.png",
		"77":  "diamond_seven.png",
		"81":  "flame_dancer.png",
		"82":  "ultra_hot_deluxe.png",
		"83":  "sizzling_hot_deluxe.png",
		"84":  "queen_of_hearts_deluxe.png",
		"85":  "cinderella.png",
		"86":  "fairy_queen.png",
		"87":  "aztec_power.png",
		"88":  "marilyn_red_carpet.png",
		"89":  "sea_sirens.png",
		"90":  "chicago.png",
		"91":  "captain_venture.png",
		"92":  "secret_elixir.png",
		"93":  "mystic_secrets.png",
		"94":  "lucky_rose.png",
		"95":  "fruitilicious.png",
		"96":  "fruit_sensation.png",
		"97":  "golden_cobras_deluxe.png",
		"98":  "royal_dynasty.png",
		"99":  "rumpel_wildspins.png",
		"100": "own/fifa_world_cup.png",
		"105": "mayan_moons.png",
		"106": "golden_7.png",
		"107": "mystic_secrets.png",
		"108": "indian_spirit.png",
		"109": "the_real_king_aloha_hawai.png",
		"110": "two_mayans.png",
		"111": "fruit_farm.png",
		"112": "reel_king.png",
		"113": "the_royals.png",
		"114": "classic/emperors_china.png",
		"115": "easy_peasy_lemon_squeezy.png",
		"116": "rex.png",
		"117": "spinning_stars.png",
		"118": "classic/threee.png",
		"119": "mega_joker.png",
		"120": "rainbow_reels.png",
		"121": "rainbow_king.png",
		"122": "roaring_forties.png",
		"123": "pharaohs_tomb.png",
		"124": "gorilla.png",
		"125": "hoffmeister.png",
		"126": "sharky.png",
		"127": "classic/marco_polo.png",
		"128": "anubix.png",
		"129": "bloody_love.png",
		"130": "bullion_bars.png",
		"131": "clockwork_oranges.png",
		"132": "dolphins_pearl.png",
		"133": "flamenco_roses.png",
		"134": "gemstone_jackpot.png",
		"135": "happy_fruits.png",
		"136": "hot_chance.png",
		"137": "orca.png",
		"138": "showgirls.png",
		"139": "spinderella.png",
		"140": "classic/always_hot.png",
		"141": "classic/attila.png",
		"142": "classic/banana_splash.png",
		"143": "classic/bananas_go_bahamas.png",
		"144": "classic/dynasty_of_ming.png",
		"145": "classic/golden_planet.png",
		"146": "classic/hot_target.png",
		"147": "classic/illusionist.png",
		"148": "classic/king_of_cards.png",
		"149": "classic/olivers_bar.png",
		"150": "classic/panter_moon.png",
		"151": "classic/pharaons_gold_2.png",
		"152": "classic/pharaons_gold_3.png",
		"153": "classic/polar_fox.png",
		"154": "classic/riches_of_india.png",
		"155": "classic/roller_coaster.png",
		"156": "classic/royal_treasures.png",
		"157": "classic/safari_heat.png",
		"158": "classic/sparta.png",
		"159": "classic/the_money_game.png",
		"160": "classic/unicorn_magic.png",
		"161": "classic/wonderful_flute.png",
		"162": "own/halloween_nightmare.png",
		"163": "own/happy_easter.png",
		"164": "own/pobeda.png",
		"165": "own/victory.png",
	}
)

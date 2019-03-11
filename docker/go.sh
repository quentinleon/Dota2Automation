sh ./dota/game/dota.sh \
	-dedicated -high -fill_with_bots +exec autoexec.cfg
	2> /dev/null | awk -f scripts/gameToJSON.awk > gamedata.json 
cat gamedata.json

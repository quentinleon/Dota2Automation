sh prepDota.sh ~/.local/share/Steam/steamapps/common/dota\ 2\ beta/
echo "Playing Dota Game..."
sh ~/.local/share/Steam/steamapps/common/dota\ 2\ beta/game/dota.sh \
	-dedicated -high -fill_with_bots +exec autoexec.cfg \
	2> /dev/null | awk -f gameToJSON.awk > gamedata.json 
echo "Done!"

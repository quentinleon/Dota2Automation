BEGIN{
	items = 1;
	lastHit = 2;
	deny = 3;
	claimedMiss = 4;
	claimedDeny = 5;
	miss = 6;
	level = 7;
	gold = 8;
	kill = 9;
	death = 10;
	assist = 11;
	actualPlayerDamage = 12;
	actualBuildingDamage = 13;
	actualHealing = 14;
	scaledPlayerDamage = 15;
	scaledBuildingDamage = 16;
	scaledHealing = 17;
	xpPerMin = 18;
	goldPerMin = 19;
	claimedFarm = 20;
	supportGold = 21;
}
{
	# save radiant heros to goodguys[]
	if(mode == "goodguys"){
		if(++goodguysCount == 5){
			mode = 0;
		};
		goodguys[i] = $3;
		i++;
	}
	if($0 ~ /^selecting radiant$/){mode = "goodguys"; i = 1};

	# save dire heros to badguys[]
	if(mode == "badguys"){
		if(++badguysCount == 5){
			mode = 0;
		};
		badguys[i] = $3;
		i++;
	}

	# other
	if($0 ~ /^selecting dire$/){mode = "badguys"; i = 1};
	if($0 ~ /Match signout:  duration[0-9]*/){gameDuration = $5};
	if($0 ~ /good guys win = [01]/){
		if($11 == 1){winner = "goodguys"}
		else{winner = "badguys"}
	}
	if($0 ~ /Match start date: /){matchDate = $4 " " $5 " " $6 " " $7 " " $8}
	if($0 ~ /Num players on teams/){numPlayersGood = $6; numPlayersBad = $8}
	
	# hero stats
	if($0 ~ /Team 0 Player [0-4]/){
		team = $2;
		player = $4 + 1;
		badguys[$4+1,items] = sprintf("%s %s %s %s %s %s %s %s %s", $9, $10, $11, $12, $13, $14, $15, $16, $17);
	}
	if($0 ~ /Team 1 Player [0-4]/){
		team = $2;
		player = $4 + 1;
		goodguys[$4+1,items] = sprintf("%s %s %s %s %s %s %s %s %s", $9, $10, $11, $12, $13, $14, $15, $16, $17);
	}
	if($0 ~ /LastHit = [0-9]* *Deny = [0-9]*/){
		if(team){
			goodguys[player, lastHit] = $3;
			goodguys[player, deny] = $6;
			goodguys[player, claimedMiss] = $9;
			goodguys[player, claimedDeny] = $12;
			goodguys[player, miss] = $15;			
		}else{
			badguys[player, lastHit] = $3;
			badguys[player, deny] = $6;
			badguys[player, claimedMiss] = $9;
			badguys[player, claimedDeny] = $12;
			badguys[player, miss] = $15;	
		}
	}else if($0 ~ /Level: [0-25]/){
		if(team){
			goodguys[player,level] = $2;
			goodguys[player,gold] = $4;
			goodguys[player,kill] = $6;
			goodguys[player,death] = $8;
			goodguys[player,assist] = $10;
		}else{
			badguys[player,level] = $2;
			badguys[player,gold] = $4;
			badguys[player,kill] = $6;
			badguys[player,death] = $8;
			badguys[player,assist] = $10;
		}
	}else if($0 ~ /Actual Player Damage: [0-9]* Actual Building Damage: [0-9]*/){
		if(team){
			goodguys[player,actualPlayerDamage] = $4;
			goodguys[player,actualBuildingDamage] = $8;
			goodguys[player,actualHealing] = $11;
		}else{
			badguys[player,actualPlayerDamage] = $4;
			badguys[player,actualBuildingDamage] = $8;
			badguys[player,actualHealing] = $11;
		}
	}else if($0 ~ /Scaled Player Damage: [0-9]* Scaled Building Damage: [0-9]*/){
		if(team){
			goodguys[player,scaledPlayerDamage] = $4;
			goodguys[player,scaledBuildingDamage] = $8;
			goodguys[player,scaledHealing] = $11;
		}else{
			badguys[player,scaledPlayerDamage] = $4;
			badguys[player,scaledBuildingDamage] = $8;
			badguys[player,scaledHealing] = $11;
		}
	}else if($0 ~ /XP per min: [0-9]/){
		if(team){
			goodguys[player,xpPerMin] = $4;
			goodguys[player,goldPerMin] = $8;
			goodguys[player,claimedFarm] = $11;
			goodguys[player,supportGold] = $14;
		}else{
			badguys[player,xpPerMin] = $4;
			badguys[player,goldPerMin] = $8;
			badguys[player,claimedFarm] = $11;
			badguys[player,supportGold] = $14;
		}
	}
}
END{
	print "{";
	printf "\"matchDate\" : \"%s\",\n", matchDate
	printf "\"gameDuration\": %d,\n", gameDuration
	printf "\"winner\" : \"%s\",\n", winner
	print "\"numPlayersGood\" : " numPlayersGood ","
	print "\"numPlayersBad\" : " numPlayersBad ","
	print "\"goodguys\" : {"
		for(i = 1; i <= 5; i++){
			printf "\t\"%s\" : {\n", goodguys[i];
			printf "\t\t\"items\": [%s],\n", goodguys[i,items];
			printf "\t\t\"lastHit\": %s,\n", goodguys[i,lastHit];
			printf "\t\t\"deny\": %s,\n", goodguys[i,deny];
			printf "\t\t\"claimedMiss\": %s,\n", goodguys[i,claimedMiss];
			printf "\t\t\"claimedDeny\": %s,\n", goodguys[i,claimedDeny];
			printf "\t\t\"miss\": %s,\n", goodguys[i,miss];
			printf "\t\t\"level\": %s,\n", goodguys[i,level];
			printf "\t\t\"gold\": %s,\n", goodguys[i,gold];
			printf "\t\t\"kill\": %s,\n", goodguys[i,kill];
			printf "\t\t\"death\": %s,\n", goodguys[i,death];
			printf "\t\t\"assist\": %s,\n", goodguys[i,assist];
			printf "\t\t\"actualPlayerDamage\": %s,\n", goodguys[i,actualPlayerDamage];
			printf "\t\t\"actualBuildingDamage\": %s,\n", goodguys[i,actualBuildingDamage];
			printf "\t\t\"actualHealing\": %s,\n", goodguys[i,actualHealing];
			printf "\t\t\"scaledPlayerDamage\": %s,\n", goodguys[i,scaledPlayerDamage];
			printf "\t\t\"scaledBuildingDamage\": %s,\n", goodguys[i,scaledBuildingDamage];
			printf "\t\t\"xpPerMin\": %s,\n", goodguys[i,xpPerMin];
			printf "\t\t\"goldPerMin\": %s,\n", goodguys[i,goldPerMin];
			printf "\t\t\"claimedFarm\": %s,\n", goodguys[i,claimedFarm];
			printf "\t\t\"supportGold\": %s\n", goodguys[i,supportGold];
			printf "\t}";
			if (i < 5){printf ",\n"};
		}
	print "},"
	print "\"badguys\" : {"
		for(i = 1; i <= 5; i++){
			printf "\t\"%s\" : {\n", badguys[i];
			printf "\t\t\"items\": [%s],\n", badguys[i,items];
			printf "\t\t\"lastHit\": %s,\n", badguys[i,lastHit];
			printf "\t\t\"deny\": %s,\n", badguys[i,deny];
			printf "\t\t\"claimedMiss\": %s,\n", badguys[i,claimedMiss];
			printf "\t\t\"claimedDeny\": %s,\n", badguys[i,claimedDeny];
			printf "\t\t\"miss\": %s,\n", badguys[i,miss];
			printf "\t\t\"level\": %s,\n", badguys[i,level];
			printf "\t\t\"gold\": %s,\n", badguys[i,gold];
			printf "\t\t\"kill\": %s,\n", badguys[i,kill];
			printf "\t\t\"death\": %s,\n", badguys[i,death];
			printf "\t\t\"assist\": %s,\n", badguys[i,assist];
			printf "\t\t\"actualPlayerDamage\": %s,\n", badguys[i,actualPlayerDamage];
			printf "\t\t\"actualBuildingDamage\": %s,\n", badguys[i,actualBuildingDamage];
			printf "\t\t\"actualHealing\": %s,\n", badguys[i,actualHealing];
			printf "\t\t\"scaledPlayerDamage\": %s,\n", badguys[i,scaledPlayerDamage];
			printf "\t\t\"scaledBuildingDamage\": %s,\n", badguys[i,scaledBuildingDamage];
			printf "\t\t\"xpPerMin\": %s,\n", badguys[i,xpPerMin];
			printf "\t\t\"goldPerMin\": %s,\n", badguys[i,goldPerMin];
			printf "\t\t\"claimedFarm\": %s,\n", badguys[i,claimedFarm];
			printf "\t\t\"supportGold\": %s\n", badguys[i,supportGold];
			printf "\t}";
			if (i < 5){printf ","};
			printf "\n";
		}
	printf "}}";
}
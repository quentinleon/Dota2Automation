package main

//Directory Tree
/*
/genetic
	/data
		/G1
			/best_genes
				gene_chaos_knight.lua
				gene_bane.lua
				gene_ogre_magi.lua
				gene_juggernaut.lua
				gene_lich.lua
			/I1
				/genes
					gene_chaos_knight.lua
					gene_bane.lua
					gene_ogre_magi.lua
					gene_juggernaut.lua
					gene_lich.lua
				/gamedata
					game1.json
					game2.json
					.
					gameX.json
			/I2
			.
			/IX
		/G2
		.
		/GX


*/

type Individual struct {
	path  string
	score int
	gene  []float64
}

type Generation struct {
	path       string
	population []Individual
}

package main

func SeedPokemon() {
	pokemons := []Pokemon{
		{Name: "Bulbasaur", Type: "Grass", Sprite: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png"},
		{Name: "Charmander", Type: "Fire", Sprite: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/4.png"},
		{Name: "Squirtle", Type: "Water", Sprite: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/7.png"},
		{Name: "Pikachu", Type: "Electric", Sprite: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png"},
		{Name: "Mewtwo", Type: "Psychic", Sprite: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/150.png"},
	}

	for _, p := range pokemons {
		DB.FirstOrCreate(&p, Pokemon{Name: p.Name})
	}
}

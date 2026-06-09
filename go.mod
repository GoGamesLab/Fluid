module github.com/GoGamesLab/Fluid

require github.com/GoGamesLab/Materials  v0.0.0

replace (
    github.com/GoGamesLab/Energy => ../Energy
    github.com/GoGamesLab/Inventory => ../Inventory
    github.com/GoGamesLab/Materials => ../Materials
)

go 1.26.1

package main

type FoodGroup int

const ( 
   VEGETAIS FoodGroup = 1 + iota 
   LEGUMES
   FRUTAS
   LATICINIOS
   CARNES
)
module Interface exposing (main)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)


-- STYLES 


inputContainerStyle =
    [ style "display" "flex" -- utilise flexbox pour la mise en page
    , style "flex-direction" "column" -- disposition en colonne
    , style "gap" "15px" -- espace entre les éléments
    , style "background" "white"
    , style "padding" "20px" -- espace intérieur pour que le contenu ne touche pas les bords
    , style "border-radius" "8px" -- coins arrondis
    , style "box-shadow" "0 4px 6px rgba(0,0,0,0.1)" -- ombre
    , style "width" "200px"
    ]

buttonStyle =
    [ style "background-color" "#4a90e2"
    , style "color" "white"
    , style "border" "none"
    , style "padding" "12px"
    , style "border-radius" "4px"
    , style "cursor" "pointer" -- change le curseur au survol
    , style "font-weight" "bold" -- texte en gras
    ]


-- MAIN


main =
  Browser.sandbox { init = init, update = update, view = view }



-- MODEL


type alias Model =
  { lat : String
  , long : String
  , d : String
  , validate : Bool
  , typeError : Bool
  }


init : Model
init =
  Model "" "" "" False False



-- UPDATE


type Msg
  = Lat String
  | Long String
  | Deniv String
  | Validate


update : Msg -> Model -> Model
update msg model =
  case msg of
    Lat lat ->
      { model | lat = lat, validate = False, typeError = False } -- Dès que le champ est modifié, réinitialisation de la validation et de l'erreur de type

    Long long ->
      { model | long = long, validate = False, typeError = False }

    Deniv d ->
      { model | d = d, validate = False, typeError = False }
      
    Validate ->
      if model.lat /= "" && model.long /= "" && model.d /= "" then
        if String.toFloat model.d == Nothing || String.toFloat model.long == Nothing || String.toFloat model.lat == Nothing then -- Si l'une des valeurs rentrées ne correspond pas à un float
          { model | validate = False, typeError = True }
      else 
          { model | validate = True }
      else
        { model | validate = False } -- Reste en false si un champ manque


-- VIEW


view : Model -> Html Msg
view model =
  div []
  [ div inputContainerStyle
    [ h1  
      [ style "color" "#2c3e50"
      , style "font-family" "Segoe UI, sans-serif"
      , style "font-size" "20px"
      , style "text-align" "center"  -- Centre le texte sur ses lignes
      , style "margin-top" "0"       -- Enlève l'espace inutile en haut
      ]
      [ text "Calculateur de Zone Atteignable" ]
    , viewInput "text" "Latitude" model.lat Lat
    , viewInput "text" "Longitude" model.long Long
    , viewInput "text" "Dénivelé" model.d Deniv
    , button 
            ([ onClick Validate ] ++ buttonStyle ++ 
                -- On ajoute une couleur différente si c'est cliqué
                [ style "background-color" (if model.validate then "#27ae60" else "#4a90e2") ]
            ) 
            [ text (if model.validate then "Calcul en cours..." else "Calculer la zone") ]
    , viewValidation model
    ]
  ]


viewInput : String -> String -> String -> (String -> msg) -> Html msg
viewInput t p v toMsg =
  input [ type_ t, placeholder p, value v, onInput toMsg ] []

viewValidation : Model -> Html msg
viewValidation model =
  if model.typeError then
    div [ style "color" "red", style "font-size" "13px", style "text-align" "center" ] [ text "Veuillez entrer des données valides" ]
  else
    text ""

module Interface exposing (mainView, Model, Msg(..), ButtonStatus(..))

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)


-- STYLES 


inputContainerStyle : List (Attribute msg)
inputContainerStyle =

    [ style "position" "fixed"
    , style "left" "50px"
    , style "bottom" "50px"
    , style "z-index" "1000"
    , style "display" "flex" -- utilise flexbox pour la mise en page
    , style "flex-direction" "column" -- disposition en colonne
    , style "gap" "15px" -- espace entre les éléments
    , style "background" "white"
    , style "padding" "20px" -- espace intérieur pour que le contenu ne touche pas les bords
    , style "border-radius" "8px" -- coins arrondis
    , style "box-shadow" "rgba(50, 50, 50, 0.3) 10px 10px 20px" -- ombre
    , style "width" "200px"
    ]

buttonStyle : List (Attribute msg)
buttonStyle =
    [ style "background-color" "#4a90e2"
    , style "color" "white"
    , style "border" "none"
    , style "padding" "12px"
    , style "border-radius" "4px"
    , style "cursor" "pointer" -- change le curseur au survol
    , style "font-weight" "bold" -- texte en gras
    ]

-- MODEL

type alias Model =
    { lat : String
    , long : String
    , d : String
    , accuracy : String
    , status : ButtonStatus
    , typeError : Bool
    }

type ButtonStatus
    = Idle          -- Au repos (pas encore cliqué)
    | Loading       -- En cours de chargement
    | Success       -- Requête réussie ("C'est bon")
    | Error         -- Erreur serveur ("Pas bon")

type Msg
  = Lat String
  | Long String
  | Deniv String
  | Accuracy String
  | Validate

-- VIEW

mainView : Model -> Html Msg 
mainView model =
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
      , viewDropdown model.accuracy
      , let
          ( btnText, btnColor ) =
            case model.status of
              Idle ->
                ( "Calculer la zone", "#4a90e2" )

              Loading ->
                ( "Calcul en cours...", "#f39c12" )

              Success ->
                ( "C'est bon !", "#27ae60" )

              Error ->
                ( "Erreur serveur", "#e74c3c" )
        in
        button
          (onClick Validate
              :: buttonStyle
              ++ [ style "background-color" btnColor ]
          )
          [ text btnText ]

      , viewValidation model
      ]
  ]

viewDropdown : String -> Html Msg
viewDropdown currentAccuracy =
    let
        renderOption val =
            option 
                [ value val, selected (val == currentAccuracy) ] 
                [ text val ]
    in
    select [ onInput Accuracy ]
        [ renderOption "1m"
        , renderOption "5m"
        , renderOption "25m"
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
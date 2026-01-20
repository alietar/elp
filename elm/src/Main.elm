module Main exposing (main)

import Browser
import Html exposing (Html)
import Carte
import Draw_square


-- MODEL

type alias Model =
    { carte : Carte.Model }


-- DONNÉES DE TEST : PLUSIEURS CARRÉS

squaresParams : List Draw_square.Params
squaresParams =
    [ { size = 20, centerLng = 2.35, centerLat = 48.85 }   -- Paris
    , { size = 30, centerLng = 5.37, centerLat = 43.29 }   -- Marseille
    , { size = 40, centerLng = -1.55, centerLat = 44.83 }  -- Bordeaux
    ] 

-- Il faut implanter ici le raisonemment sur nos carrés


-- INIT

init : () -> ( Model, Cmd msg )
init _ =
    let
        ( carteModel, carteCmd ) =
            Carte.init

        boundsList =
            List.map Draw_square.computeBounds squaresParams

        drawCmds =
            List.map Carte.drawSquare boundsList
    in
    ( { carte = carteModel }
    , Cmd.batch (carteCmd :: drawCmds)
    )


-- VIEW

view : Model -> Html msg
view model =
    Carte.view model.carte


-- MAIN

main : Program () Model msg
main =
    Browser.element
        { init = init
        , view = view
        , update = \_ model -> ( model, Cmd.none )
        , subscriptions = \_ -> Sub.none
        }
